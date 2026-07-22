import { ApiFetchError } from '@api/management/types/errors'
import { AddTrustedAccountDialog } from '@feature/management/components'
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

async function openDialog(user: ReturnType<typeof userEvent.setup>) {
  await user.click(screen.getByRole('button', { name: /add account/i }))
  expect(await screen.findByRole('dialog')).toBeInTheDocument()
}

async function fillValidForm(user: ReturnType<typeof userEvent.setup>) {
  await user.type(screen.getByLabelText('Account Name'), 'Alice LP')
  await user.type(
    screen.getByLabelText('Address'),
    '0x1234567890123456789012345678901234567890',
  )
  await user.type(screen.getByLabelText('BTC Locking Cap'), '1.5')
  await user.type(screen.getByLabelText('rBTC Locking Cap'), '2.5')
}

describe('AddTrustedAccountDialog', () => {
  const onSuccess = vi.fn()

  beforeEach(() => {
    document.body.innerHTML = ''
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
    toastErrorMock.mockReset()
    toastSuccessMock.mockReset()
    onSuccess.mockReset()
  })

  it('focuses the account name field when opened', async () => {
    const user = userEvent.setup()
    render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

    await openDialog(user)

    await waitFor(() => {
      expect(screen.getByLabelText('Account Name')).toHaveFocus()
    })
  })

  it('returns focus to the trigger after close', async () => {
    const user = userEvent.setup()
    render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

    const trigger = screen.getByRole('button', { name: /add account/i })
    await openDialog(user)
    await user.click(screen.getByRole('button', { name: 'Cancel' }))

    await waitFor(() => {
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
    })
    expect(trigger).toHaveFocus()
  })

  it.each([
    {
      label: 'Account Name',
      clear: 'Account Name',
      message: 'Account name is required',
    },
    {
      label: 'Address',
      clear: 'Address',
      message: 'Account address is required',
    },
    {
      label: 'BTC Locking Cap',
      clear: 'BTC Locking Cap',
      message: 'BTC Locking Cap is required',
    },
    {
      label: 'rBTC Locking Cap',
      clear: 'rBTC Locking Cap',
      message: 'rBTC Locking Cap is required',
    },
  ])(
    'shows "$message" without POST when $clear is empty',
    async ({ clear, message }) => {
      const user = userEvent.setup()
      render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

      await openDialog(user)
      await fillValidForm(user)
      await user.clear(screen.getByLabelText(clear))
      await user.click(screen.getByRole('button', { name: 'Save' }))

      expect(await screen.findByText(message)).toBeInTheDocument()
      expect(apiFetchMock.post).not.toHaveBeenCalled()
      expect(toastSuccessMock).not.toHaveBeenCalled()
    },
  )

  it.each([
    { value: '0', message: 'BTC Locking Cap must be a positive number' },
    { value: '-1', message: 'BTC Locking Cap must be a positive number' },
    { value: 'abc', message: 'BTC Locking Cap must be a positive number' },
  ])(
    'rejects non-positive BTC cap "$value" without POST',
    async ({ value, message }) => {
      const user = userEvent.setup()
      render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

      await openDialog(user)
      await fillValidForm(user)
      await user.clear(screen.getByLabelText('BTC Locking Cap'))
      await user.type(screen.getByLabelText('BTC Locking Cap'), value)
      await user.click(screen.getByRole('button', { name: 'Save' }))

      expect(await screen.findByText(message)).toBeInTheDocument()
      expect(apiFetchMock.post).not.toHaveBeenCalled()
    },
  )

  it('rejects non-positive rBTC cap without POST', async () => {
    const user = userEvent.setup()
    render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

    await openDialog(user)
    await fillValidForm(user)
    await user.clear(screen.getByLabelText('rBTC Locking Cap'))
    await user.type(screen.getByLabelText('rBTC Locking Cap'), '0')
    await user.click(screen.getByRole('button', { name: 'Save' }))

    expect(
      await screen.findByText('rBTC Locking Cap must be a positive number'),
    ).toBeInTheDocument()
    expect(apiFetchMock.post).not.toHaveBeenCalled()
  })

  it('posts valid form, closes, refreshes, and toasts once', async () => {
    const user = userEvent.setup()
    apiFetchMock.mockResolvedValueOnce(new Response('', { status: 200 }))
    render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

    await openDialog(user)
    await fillValidForm(user)
    await user.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => {
      expect(apiFetchMock.post).toHaveBeenCalledTimes(1)
      expect(apiFetchMock.post).toHaveBeenCalledWith(
        '/management/trusted-accounts',
        {
          name: 'Alice LP',
          address: '0x1234567890123456789012345678901234567890',
          btcLockingCap: 1500000000000000000,
          rbtcLockingCap: 2500000000000000000,
        },
      )
    })

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalledTimes(1)
      expect(toastSuccessMock).toHaveBeenCalledWith(
        'Configuration saved successfully!',
      )
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
    })
  })

  it('maps backend validation details inline without toast', async () => {
    const user = userEvent.setup()
    apiFetchMock.mockRejectedValueOnce(
      new ApiFetchError(400, 'Bad Request', {
        message: 'validation error',
        details: {
          Address: 'invalid address format',
        },
      }),
    )
    render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

    await openDialog(user)
    await fillValidForm(user)
    await user.click(screen.getByRole('button', { name: 'Save' }))

    const dialog = await screen.findByRole('dialog')
    expect(
      await within(dialog).findByText('invalid address format'),
    ).toBeInTheDocument()
    expect(toastErrorMock).not.toHaveBeenCalled()
    expect(toastSuccessMock).not.toHaveBeenCalled()
    expect(onSuccess).not.toHaveBeenCalled()
    expect(screen.getByRole('dialog')).toBeInTheDocument()
  })

  it('toasts generic backend errors', async () => {
    const user = userEvent.setup()
    apiFetchMock.mockRejectedValueOnce(
      new ApiFetchError(500, 'Internal Server Error', {
        message: 'boom',
      }),
    )
    render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

    await openDialog(user)
    await fillValidForm(user)
    await user.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() => {
      expect(toastErrorMock).toHaveBeenCalledWith(
        'Error adding trusted account: boom',
      )
    })
    expect(onSuccess).not.toHaveBeenCalled()
  })

  it('disables save and blocks duplicate POST while pending', async () => {
    const user = userEvent.setup()
    let resolvePost: (value: Response) => void = () => undefined
    const postPromise = new Promise<Response>((resolve) => {
      resolvePost = resolve
    })
    apiFetchMock.mockReturnValueOnce(postPromise)

    render(<AddTrustedAccountDialog onSuccess={onSuccess} />)

    await openDialog(user)
    await fillValidForm(user)

    const saveButton = screen.getByRole('button', { name: 'Save' })
    await user.click(saveButton)
    expect(saveButton).toBeDisabled()

    await user.click(saveButton)
    expect(apiFetchMock.post).toHaveBeenCalledTimes(1)

    resolvePost(new Response('', { status: 200 }))

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalledTimes(1)
    })
  })
})
