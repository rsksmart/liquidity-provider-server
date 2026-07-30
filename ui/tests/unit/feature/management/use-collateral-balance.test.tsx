import { useCollateralBalance } from '@feature/management/hooks/use-collateral-balance'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { loggedInFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { useCallback } from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { apiFetchMock } = vi.hoisted(() => {
  const fn = vi.fn()
  return { apiFetchMock: Object.assign(fn, { get: fn, post: fn }) }
})

vi.mock('@api/management/utils/api-fetch', () => ({
  apiFetch: apiFetchMock,
}))

function mockBalance(wei: string) {
  return new Response(JSON.stringify({ collateral: wei }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })
}

function BalanceProbe({ endpoint }: { endpoint: string }) {
  const { balance, loading, refresh } = useCollateralBalance(endpoint)
  const handleRefresh = useCallback(() => {
    void refresh()
  }, [refresh])
  return (
    <div>
      <span data-testid="balance">{loading ? 'loading' : (balance ?? 'none')}</span>
      <button type="button" onClick={handleRefresh}>
        refresh
      </button>
    </div>
  )
}

describe('useCollateralBalance', () => {
  beforeEach(() => {
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
  })

  it('ignores stale responses when refresh is called again', async () => {
    let resolveFirst: (value: Response) => void = () => undefined
    const first = new Promise<Response>((resolve) => {
      resolveFirst = resolve
    })

    apiFetchMock.mockReturnValueOnce(first).mockResolvedValueOnce(mockBalance('2000000000000000000'))

    render(<BalanceProbe endpoint="/pegin/collateral" />)

    await waitFor(() => {
      expect(screen.getByTestId('balance')).toHaveTextContent('loading')
    })

    const user = userEvent.setup()
    await user.click(screen.getByRole('button', { name: 'refresh' }))

    resolveFirst(mockBalance('1000000000000000000'))

    await waitFor(() => {
      expect(screen.getByTestId('balance')).toHaveTextContent('2 rBTC')
    })
  })

  it('returns null balance when fetch fails', async () => {
    apiFetchMock.mockRejectedValueOnce(new Error('network'))

    render(<BalanceProbe endpoint="/pegin/collateral" />)

    await waitFor(() => {
      expect(screen.getByTestId('balance')).toHaveTextContent('none')
    })
  })
})
