import { ApiFetchError } from '@api/management/types/errors'
import { TrustedAccountsCard } from '@feature/management/components'
import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { loggedInFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { apiFetchMock, toastErrorMock, toastSuccessMock } = vi.hoisted(() => {
  const fn = vi.fn()
  return {
    apiFetchMock: Object.assign(fn, { get: fn, post: fn }),
    toastErrorMock: vi.fn(),
    toastSuccessMock: vi.fn(),
  }
})

vi.mock('@api/management/utils/api-fetch', () => ({
  apiFetch: apiFetchMock,
}))

vi.mock('sonner', () => ({
  toast: {
    error: toastErrorMock,
    success: toastSuccessMock,
  },
}))

const SAMPLE_ACCOUNTS = [
  {
    name: 'Alice LP',
    address: '0x1234567890123456789012345678901234567890',
    btcLockingCap: '1500000000000000000',
    rbtcLockingCap: '2500000000000000000',
  },
  {
    address: '0xabcdefabcdefabcdefabcdefabcdefabcdefabcd',
    btcLockingCap: '100000000000000000',
    rbtcLockingCap: '50000000000000000',
  },
]

function mockAccountsResponse(accounts: unknown[], once = true) {
  const response = new Response(JSON.stringify({ accounts }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })
  if (once) {
    apiFetchMock.mockResolvedValueOnce(response)
  } else {
    apiFetchMock.mockResolvedValue(response)
  }
}

describe('TrustedAccountsCard', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
    toastErrorMock.mockReset()
    toastSuccessMock.mockReset()
  })

  it('renders legacy columns, formatted caps, and Unknown name fallback', async () => {
    mockAccountsResponse(SAMPLE_ACCOUNTS)

    render(<TrustedAccountsCard />)

    expect(
      await screen.findByRole('columnheader', { name: 'Name' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('columnheader', { name: 'Address' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('columnheader', { name: 'BTC Cap' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('columnheader', { name: 'rBTC Cap' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('columnheader', { name: 'Actions' }),
    ).toBeInTheDocument()

    expect(screen.getByText('Alice LP')).toBeInTheDocument()
    expect(screen.getByText('Unknown')).toBeInTheDocument()
    expect(screen.getByText('1.5 BTC')).toBeInTheDocument()
    expect(screen.getByText('2.5 rBTC')).toBeInTheDocument()
    expect(screen.getByText('0.1 BTC')).toBeInTheDocument()
    expect(screen.getByText('0.05 rBTC')).toBeInTheDocument()
  })

  it('shows loading indicator while the initial GET is in flight', async () => {
    let resolveGet: (value: Response) => void = () => undefined
    const getPromise = new Promise<Response>((resolve) => {
      resolveGet = resolve
    })
    apiFetchMock.mockReturnValueOnce(getPromise)

    render(<TrustedAccountsCard />)

    expect(screen.getByTestId('trusted-accounts-loading-bar')).toHaveClass(
      'is-visible',
    )

    resolveGet(
      new Response(JSON.stringify({ accounts: [] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    expect(
      await screen.findByText('No trusted accounts found.'),
    ).toBeInTheDocument()
    await waitFor(() => {
      expect(
        screen.getByTestId('trusted-accounts-loading-bar'),
      ).not.toHaveClass('is-visible')
    })
  })

  it('renders empty state when there are no accounts', async () => {
    mockAccountsResponse([])
    render(<TrustedAccountsCard />)
    expect(
      await screen.findByText('No trusted accounts found.'),
    ).toBeInTheDocument()
  })

  it('renders table error state and toasts on load failure', async () => {
    apiFetchMock.mockRejectedValueOnce(new Error('network down'))
    render(<TrustedAccountsCard />)

    expect(await screen.findByText('Error: network down')).toBeInTheDocument()
    await waitFor(() => {
      expect(toastErrorMock).toHaveBeenCalledWith(
        'Failed to load trusted accounts: network down',
      )
    })
  })

  it('cancels remove without calling DELETE', async () => {
    mockAccountsResponse(SAMPLE_ACCOUNTS)
    const user = userEvent.setup()
    render(<TrustedAccountsCard />)

    expect(await screen.findByText('Alice LP')).toBeInTheDocument()
    await user.click(
      screen.getByTestId(
        'remove-trusted-account-0x1234567890123456789012345678901234567890',
      ),
    )

    const confirmDialog = await screen.findByRole('alertdialog')
    expect(
      within(confirmDialog).getByText(
        /Are you sure you want to remove the trusted account with address 0x1234567890123456789012345678901234567890\?/,
      ),
    ).toBeInTheDocument()

    await user.click(
      within(confirmDialog).getByRole('button', { name: 'Cancel' }),
    )

    await waitFor(() => {
      expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()
    })
    expect(apiFetchMock).toHaveBeenCalledTimes(1) // initial GET only
  })

  it('confirms remove with encoded address, refreshes, and toasts', async () => {
    mockAccountsResponse(SAMPLE_ACCOUNTS)
    const user = userEvent.setup()
    render(<TrustedAccountsCard />)

    expect(await screen.findByText('Alice LP')).toBeInTheDocument()
    await user.click(
      screen.getByTestId(
        'remove-trusted-account-0x1234567890123456789012345678901234567890',
      ),
    )

    apiFetchMock.mockResolvedValueOnce(new Response('', { status: 200 }))
    mockAccountsResponse([SAMPLE_ACCOUNTS[1]])

    const confirmDialog = await screen.findByRole('alertdialog')
    await user.click(
      within(confirmDialog).getByRole('button', { name: 'Confirm' }),
    )

    await waitFor(() => {
      expect(apiFetchMock).toHaveBeenCalledWith(
        '/management/trusted-accounts?address=0x1234567890123456789012345678901234567890',
        { method: 'DELETE' },
      )
    })

    await waitFor(() => {
      expect(toastSuccessMock).toHaveBeenCalledWith(
        'Configuration saved successfully!',
      )
      expect(screen.queryByText('Alice LP')).not.toBeInTheDocument()
      expect(screen.getByText('Unknown')).toBeInTheDocument()
    })
  })

  it('toasts on remove failure', async () => {
    mockAccountsResponse(SAMPLE_ACCOUNTS)
    const user = userEvent.setup()
    render(<TrustedAccountsCard />)

    expect(await screen.findByText('Alice LP')).toBeInTheDocument()
    await user.click(
      screen.getByTestId(
        'remove-trusted-account-0x1234567890123456789012345678901234567890',
      ),
    )

    apiFetchMock.mockRejectedValueOnce(
      new ApiFetchError(500, 'Internal Server Error', {
        message: 'delete failed',
      }),
    )

    const confirmDialog = await screen.findByRole('alertdialog')
    await user.click(
      within(confirmDialog).getByRole('button', { name: 'Confirm' }),
    )

    await waitFor(() => {
      expect(toastErrorMock).toHaveBeenCalledWith(
        'Error removing trusted account: delete failed',
      )
    })
    expect(screen.getByText('Alice LP')).toBeInTheDocument()
  })

  it('disables confirm and blocks duplicate DELETE while pending', async () => {
    mockAccountsResponse(SAMPLE_ACCOUNTS)
    const user = userEvent.setup()
    render(<TrustedAccountsCard />)

    expect(await screen.findByText('Alice LP')).toBeInTheDocument()
    await user.click(
      screen.getByTestId(
        'remove-trusted-account-0x1234567890123456789012345678901234567890',
      ),
    )

    let resolveDelete: (value: Response) => void = () => undefined
    const deletePromise = new Promise<Response>((resolve) => {
      resolveDelete = resolve
    })
    apiFetchMock.mockReturnValueOnce(deletePromise)

    const confirmDialog = await screen.findByRole('alertdialog')
    const confirmButton = within(confirmDialog).getByRole('button', {
      name: 'Confirm',
    })
    await user.click(confirmButton)
    expect(confirmButton).toBeDisabled()

    await user.click(confirmButton)
    const deleteCalls = apiFetchMock.mock.calls.filter((call) => {
      const [url, init] = call as [string, { method?: string } | undefined]
      return (
        typeof url === 'string' &&
        url.includes('address=') &&
        init?.method === 'DELETE'
      )
    })
    expect(deleteCalls).toHaveLength(1)

    mockAccountsResponse([SAMPLE_ACCOUNTS[1]])
    resolveDelete(new Response('', { status: 200 }))

    await waitFor(() => {
      expect(toastSuccessMock).toHaveBeenCalledTimes(1)
    })
  })

  it('add success refreshes the list to include the new row', async () => {
    mockAccountsResponse([])
    const user = userEvent.setup()
    render(<TrustedAccountsCard />)

    expect(
      await screen.findByText('No trusted accounts found.'),
    ).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /add account/i }))
    await user.type(screen.getByLabelText('Account Name'), 'Bob')
    await user.type(
      screen.getByLabelText('Address'),
      '0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb',
    )
    await user.type(screen.getByLabelText('BTC Locking Cap'), '1')
    await user.type(screen.getByLabelText('rBTC Locking Cap'), '1')

    apiFetchMock.mockResolvedValueOnce(new Response('', { status: 200 }))
    mockAccountsResponse([
      {
        name: 'Bob',
        address: '0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb',
        btcLockingCap: '1000000000000000000',
        rbtcLockingCap: '1000000000000000000',
      },
    ])

    await user.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findByText('Bob')).toBeInTheDocument()
    expect(
      screen.queryByText('No trusted accounts found.'),
    ).not.toBeInTheDocument()
  })
})
