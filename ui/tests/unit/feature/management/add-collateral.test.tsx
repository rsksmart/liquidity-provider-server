import { ApiFetchError } from '@api/management/types/errors'
import { AddCollateralForm } from '@feature/management/components'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { loggedInFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
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
  },
}))

function mockGetBalance(wei: string) {
  apiFetchMock.mockResolvedValueOnce(
    new Response(JSON.stringify({ collateral: wei }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    }),
  )
}

describe('AddCollateralForm', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
    toastErrorMock.mockReset()
    mockGetBalance('1000000000000000000')
  })

  it('submits valid collateral and refreshes balance', async () => {
    const user = userEvent.setup()
    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toHaveTextContent('1 rBTC')
    })

    apiFetchMock.mockResolvedValueOnce(new Response('', { status: 200 }))
    mockGetBalance('2000000000000000000')

    await user.type(screen.getByTestId('pegin-collateral-amount'), '1')
    await user.click(screen.getByTestId('pegin-add-collateral-button'))

    await waitFor(() => {
      expect(apiFetchMock.post).toHaveBeenCalledWith('/pegin/addCollateral', {
        amount: 1000000000000000000,
      })
    })

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toHaveTextContent('2 rBTC')
    })
    expect(toastErrorMock).not.toHaveBeenCalled()
  })

  it('disables submit while request is in flight', async () => {
    const user = userEvent.setup()
    let resolvePost: (value: Response) => void = () => undefined
    const postPromise = new Promise<Response>((resolve) => {
      resolvePost = resolve
    })

    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toHaveTextContent('1 rBTC')
    })

    apiFetchMock.mockReturnValueOnce(postPromise)
    await user.type(screen.getByTestId('pegin-collateral-amount'), '1')

    const submitButton = screen.getByTestId('pegin-add-collateral-button')
    await user.click(submitButton)

    expect(submitButton).toBeDisabled()
    expect(screen.getByTestId('pegin-loading-bar')).toHaveClass('is-visible')

    resolvePost(new Response('', { status: 200 }))
    mockGetBalance('1000000000000000000')

    await waitFor(() => {
      expect(submitButton).not.toBeDisabled()
    })
  })
})

describe('add-collateral validation', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
    toastErrorMock.mockReset()
    mockGetBalance('0')
  })

  it.each(['0'])('shows error toast for invalid amount %s without POST', async (value) => {
    const user = userEvent.setup()
    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })

    await user.type(screen.getByTestId('pegin-collateral-amount'), value)
    await user.click(screen.getByTestId('pegin-add-collateral-button'))

    expect(toastErrorMock).toHaveBeenCalledWith(
      `Invalid input "${value}" for collateral amount. Please enter a valid number.`,
    )
    expect(apiFetchMock).toHaveBeenCalledTimes(1)
  })

  it('shows error toast for sub-wei fractional amount without POST', async () => {
    const user = userEvent.setup()
    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })

    // Sub-wei ether (1e-19) yields a non-integral wei string; BigInt would throw
    // if left outside the validation try/catch (silent no-op regression).
    const amount = '0.0000000000000000001'
    const input = screen.getByTestId('pegin-collateral-amount')
    await user.clear(input)
    // HTML number inputs may normalize on type; set value directly.
    fireEvent.change(input, { target: { value: amount } })
    await user.click(screen.getByTestId('pegin-add-collateral-button'))

    expect(toastErrorMock).toHaveBeenCalledWith(
      `Invalid input "${amount}" for collateral amount. Please enter a valid number.`,
    )
    expect(apiFetchMock).toHaveBeenCalledTimes(1)
  })

  it('shows error toast for negative amount without POST', async () => {
    const user = userEvent.setup()
    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })

    await user.type(screen.getByTestId('pegin-collateral-amount'), '-1')
    await user.click(screen.getByTestId('pegin-add-collateral-button'))

    expect(toastErrorMock).toHaveBeenCalledWith(
      'Invalid input "-1" for collateral amount. Please enter a valid number.',
    )
    expect(apiFetchMock).toHaveBeenCalledTimes(1)
  })

  it('re-enables submit after validation error', async () => {
    const user = userEvent.setup()
    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })

    const submitButton = screen.getByTestId('pegin-add-collateral-button')
    await user.type(screen.getByTestId('pegin-collateral-amount'), '0')
    await user.click(submitButton)

    expect(toastErrorMock).toHaveBeenCalled()
    expect(submitButton).not.toBeDisabled()
    expect(screen.getByTestId('pegin-loading-bar')).not.toHaveClass('is-visible')
  })

  it('re-enables submit after API error', async () => {
    const user = userEvent.setup()
    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })

    apiFetchMock.mockRejectedValueOnce(
      new ApiFetchError(409, 'Conflict', { message: 'not enough for minimum collateral' }),
    )

    await user.type(screen.getByTestId('pegin-collateral-amount'), '1')
    await user.click(screen.getByTestId('pegin-add-collateral-button'))

    expect(toastErrorMock).toHaveBeenCalledWith(
      'Error adding collateral: not enough for minimum collateral',
    )
    expect(screen.getByTestId('pegin-add-collateral-button')).not.toBeDisabled()
  })

  it('shows unknown error when API body has no message', async () => {
    const user = userEvent.setup()
    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })

    apiFetchMock.mockRejectedValueOnce(new ApiFetchError(500, 'Server Error', { error: 'boom' }))

    await user.type(screen.getByTestId('pegin-collateral-amount'), '1')
    await user.click(screen.getByTestId('pegin-add-collateral-button'))

    expect(toastErrorMock).toHaveBeenCalledWith('Error adding collateral: Unknown error')
  })

  it('shows validation toast for non-API submit failures', async () => {
    const user = userEvent.setup()
    render(<AddCollateralForm kind="pegin" />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })

    apiFetchMock.mockRejectedValueOnce(new Error('network'))

    await user.type(screen.getByTestId('pegin-collateral-amount'), '1')
    await user.click(screen.getByTestId('pegin-add-collateral-button'))

    expect(toastErrorMock).toHaveBeenCalledWith(
      'Invalid input "1" for collateral amount. Please enter a valid number.',
    )
  })
})
