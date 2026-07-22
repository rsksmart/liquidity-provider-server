import { useTrustedAccounts } from '@feature/management/hooks/use-trusted-accounts'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { loggedInFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { useCallback } from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { apiFetchMock, toastErrorMock } = vi.hoisted(() => {
  const fn = vi.fn()
  return {
    apiFetchMock: Object.assign(fn, { get: fn, post: fn }),
    toastErrorMock: vi.fn(),
  }
})

vi.mock('@api/management/utils/api-fetch', () => ({
  apiFetch: apiFetchMock,
}))

vi.mock('sonner', () => ({
  toast: {
    error: toastErrorMock,
    success: vi.fn(),
  },
}))

function mockAccounts(accounts: unknown[]) {
  return new Response(JSON.stringify({ accounts }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })
}

function HookProbe() {
  const { accounts, loading, error, refresh } = useTrustedAccounts()
  const onRefresh = useCallback(() => {
    void refresh()
  }, [refresh])
  return (
    <div>
      <span data-testid="loading">{loading ? 'yes' : 'no'}</span>
      <span data-testid="error">{error ?? ''}</span>
      <ul data-testid="accounts">
        {accounts.map((account) => (
          <li key={account.address}>{account.name || 'Unknown'}</li>
        ))}
      </ul>
      <button type="button" onClick={onRefresh}>
        refresh
      </button>
    </div>
  )
}

describe('useTrustedAccounts', () => {
  beforeEach(() => {
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
    toastErrorMock.mockReset()
  })

  it('loads accounts from GET /management/trusted-accounts', async () => {
    apiFetchMock.mockResolvedValueOnce(
      mockAccounts([
        {
          name: 'Alice',
          address: '0xabc',
          btcLockingCap: '1',
          rbtcLockingCap: '2',
        },
      ]),
    )

    render(<HookProbe />)

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('no')
    })
    expect(screen.getByTestId('accounts')).toHaveTextContent('Alice')
    expect(apiFetchMock).toHaveBeenCalledWith('/management/trusted-accounts')
  })

  it('surfaces fetch errors and toasts', async () => {
    apiFetchMock.mockRejectedValueOnce(new Error('network down'))

    render(<HookProbe />)

    await waitFor(() => {
      expect(screen.getByTestId('error')).toHaveTextContent('network down')
    })
    expect(toastErrorMock).toHaveBeenCalledWith(
      'Failed to load trusted accounts: network down',
    )
  })

  it('ignores stale responses when refresh is called again', async () => {
    let resolveFirst: (value: Response) => void = () => undefined
    const first = new Promise<Response>((resolve) => {
      resolveFirst = resolve
    })

    apiFetchMock.mockReturnValueOnce(first).mockResolvedValueOnce(
      mockAccounts([
        {
          name: 'Second',
          address: '0x2',
          btcLockingCap: '1',
          rbtcLockingCap: '1',
        },
      ]),
    )

    render(<HookProbe />)

    await waitFor(() => {
      expect(screen.getByTestId('loading')).toHaveTextContent('yes')
    })

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'refresh' }))

    resolveFirst(
      mockAccounts([
        {
          name: 'First',
          address: '0x1',
          btcLockingCap: '1',
          rbtcLockingCap: '1',
        },
      ]),
    )

    await waitFor(() => {
      expect(screen.getByTestId('accounts')).toHaveTextContent('Second')
    })
    expect(screen.getByTestId('accounts')).not.toHaveTextContent('First')
  })
})
