import { ConfigurationCard } from '@feature/management/components'
import { ZERO_FEE_WARNING_MESSAGE } from '@feature/management/config/fee-warnings'
import type {
  FullConfiguration,
  InitialDataPayload,
} from '@shared/types/initial-data'
import { etherToWei } from '@shared/utils/wei'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { loggedInFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { apiFetchMock, apiFetchPostMock, toastErrorMock, toastSuccessMock, toastWarningMock } =
  vi.hoisted(() => {
    const base = vi.fn()
    const post = vi.fn()
    return {
      apiFetchMock: Object.assign(base, { get: base, post }),
      apiFetchPostMock: post,
      toastErrorMock: vi.fn(),
      toastSuccessMock: vi.fn(),
      toastWarningMock: vi.fn(),
    }
  })

vi.mock('@api/management/utils/api-fetch', () => ({
  apiFetch: apiFetchMock,
}))

vi.mock('sonner', () => ({
  toast: {
    error: toastErrorMock,
    success: toastSuccessMock,
    warning: toastWarningMock,
  },
}))

const sampleConfiguration: FullConfiguration = {
  general: {
    rskConfirmations: { '1000000000000000000': 2 },
    btcConfirmations: { '2000000000000000000': 6 },
    publicLiquidityCheck: true,
    maxLiquidity: '5000000000000000000',
    reimbursementWindowBlocks: 100,
    excessTolerance: {
      isFixed: false,
      percentageValue: '10',
      fixedValue: '0',
    },
  },
  pegin: {
    timeForDeposit: 3600,
    callTime: 7200,
    penaltyFee: '1000000000000000',
    fixedFee: '2000000000000000',
    feePercentage: '1.5',
    maxValue: '10000000000000000000',
    minValue: '1000000000000000000',
  },
  pegout: {
    timeForDeposit: 3600,
    expireTime: 7200,
    expireBlocks: 500,
    penaltyFee: '1000000000000000',
    fixedFee: '0',
    feePercentage: '0',
    maxValue: '10000000000000000000',
    minValue: '1000000000000000000',
    bridgeTransactionMin: '3000000000000000000',
  },
}

function seedConfiguration(): void {
  const payload: InitialDataPayload = {
    ...loggedInFixture,
    data: { ...loggedInFixture.data, Configuration: sampleConfiguration },
  }
  seedInitialData(payload, { csrfToken: 'csrf-token' })
}

function okResponse(): Response {
  return new Response('', { status: 200 })
}

describe('ConfigurationCard', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedConfiguration()
    apiFetchMock.mockReset()
    apiFetchPostMock.mockReset()
    toastErrorMock.mockReset()
    toastSuccessMock.mockReset()
    toastWarningMock.mockReset()
  })

  it('renders card chrome with tabs and disabled save', () => {
    render(<ConfigurationCard />)

    expect(screen.getByTestId('configuration-card')).toBeInTheDocument()
    expect(screen.getByText('Configuration')).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: 'Current Configuration' }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('config-tab-general')).toBeInTheDocument()
    expect(screen.getByTestId('config-tab-pegin')).toBeInTheDocument()
    expect(screen.getByTestId('config-tab-pegout')).toBeInTheDocument()

    const save = screen.getByTestId('config-save-button')
    expect(save).toHaveTextContent('Save Configuration')
    expect(save).toBeDisabled()
  })

  it('enables save on edit, disables on round-trip, and disables after a successful save', async () => {
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    const save = screen.getByTestId('config-save-button')
    const input = screen.getByTestId('config-general-maxLiquidity-input')

    expect(save).toBeDisabled()

    await user.clear(input)
    await user.type(input, '7')
    expect(save).toBeEnabled()

    await user.clear(input)
    await user.type(input, '5')
    expect(save).toBeDisabled()

    await user.clear(input)
    await user.type(input, '7')
    expect(save).toBeEnabled()

    apiFetchPostMock.mockResolvedValueOnce(okResponse())
    await user.click(save)

    await waitFor(() => {
      expect(toastSuccessMock).toHaveBeenCalledWith(
        'Configuration saved successfully!',
      )
    })
    expect(save).toBeDisabled()
  })

  it('posts only the dirty section with the { configuration } body shape', async () => {
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    const input = screen.getByTestId('config-general-maxLiquidity-input')
    await user.clear(input)
    await user.type(input, '7')

    apiFetchPostMock.mockResolvedValueOnce(okResponse())
    await user.click(screen.getByTestId('config-save-button'))

    await waitFor(() => {
      expect(apiFetchPostMock).toHaveBeenCalledTimes(1)
    })
    expect(apiFetchPostMock).toHaveBeenCalledWith('/configuration', {
      configuration: {
        publicLiquidityCheck: true,
        maxLiquidity: etherToWei('7'),
        reimbursementWindowBlocks: 100,
        excessTolerance: {
          isFixed: false,
          fixedValue: '0',
          percentageValue: '10',
        },
        rskConfirmations: { '1000000000000000000': 2 },
        btcConfirmations: { '2000000000000000000': 6 },
      },
    })
  })

  it('sends the save through apiFetch.post (CSRF handled by the wrapper)', async () => {
    const fetchSpy = vi.spyOn(globalThis, 'fetch')
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    const input = screen.getByTestId('config-pegout-maxValue-input')
    await user.click(screen.getByTestId('config-tab-pegout'))
    await user.clear(input)
    await user.type(input, '20')

    apiFetchPostMock.mockResolvedValueOnce(okResponse())
    await user.click(screen.getByTestId('config-save-button'))

    await waitFor(() => {
      expect(apiFetchPostMock).toHaveBeenCalledWith('/pegout/configuration', {
        configuration: {
          timeForDeposit: 3600,
          expireTime: 7200,
          expireBlocks: 500,
          bridgeTransactionMin: etherToWei('3'),
          penaltyFee: etherToWei('0.001'),
          fixedFee: '0',
          feePercentage: '0',
          maxValue: etherToWei('20'),
          minValue: etherToWei('1'),
        },
      })
    })
    expect(fetchSpy).not.toHaveBeenCalled()
    fetchSpy.mockRestore()
  })

  it('warns about zero fees before saving pegin but still proceeds', async () => {
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    await user.click(screen.getByTestId('config-tab-pegin'))
    await user.click(screen.getByTestId('config-pegin-fixedFee-checkbox'))
    await user.click(screen.getByTestId('config-pegin-feePercentage-checkbox'))

    apiFetchPostMock.mockResolvedValueOnce(okResponse())
    await user.click(screen.getByTestId('config-save-button'))

    await waitFor(() => {
      expect(toastWarningMock).toHaveBeenCalledWith(ZERO_FEE_WARNING_MESSAGE)
    })
    await waitFor(() => {
      expect(apiFetchPostMock).toHaveBeenCalledWith('/pegin/configuration', {
        configuration: {
          timeForDeposit: 3600,
          callTime: 7200,
          penaltyFee: etherToWei('0.001'),
          fixedFee: '0',
          feePercentage: '0',
          maxValue: etherToWei('10'),
          minValue: etherToWei('1'),
        },
      })
    })
    await waitFor(() => {
      expect(toastSuccessMock).toHaveBeenCalledWith(
        'Configuration saved successfully!',
      )
    })
  })

  it('shows a section error toast and does not report success on save failure', async () => {
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    const input = screen.getByTestId('config-general-maxLiquidity-input')
    await user.clear(input)
    await user.type(input, '7')

    apiFetchPostMock.mockRejectedValueOnce(new Error('server down'))
    await user.click(screen.getByTestId('config-save-button'))

    await waitFor(() => {
      expect(toastErrorMock).toHaveBeenCalledWith('server down')
    })
    expect(toastSuccessMock).not.toHaveBeenCalled()
    expect(screen.getByTestId('config-save-button')).toBeEnabled()
  })
})
