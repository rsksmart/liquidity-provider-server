import { ApiFetchError } from '@api/management/types/errors'
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

  /**
   * The legacy card renders `Object.entries(config)`, so its field order is the
   * JSON key order of the Go structs — `callTime` sits directly under
   * `timeForDeposit`, and `expireTime` likewise for pegout.
   */
  it.each([
    {
      section: 'pegin',
      expected: [
        'timeForDeposit',
        'callTime',
        'penaltyFee',
        'fixedFee',
        'feePercentage',
        'maxValue',
        'minValue',
      ],
    },
    {
      section: 'pegout',
      expected: [
        'timeForDeposit',
        'expireTime',
        'penaltyFee',
        'fixedFee',
        'feePercentage',
        'maxValue',
        'minValue',
        'expireBlocks',
        'bridgeTransactionMin',
      ],
    },
  ])('orders $section fields as the legacy card does', ({ section, expected }) => {
    render(<ConfigurationCard />)

    const inputs = Array.from(
      document.querySelectorAll<HTMLElement>(
        `[data-testid^="config-${section}-"][data-testid$="-input"]`,
      ),
    ).map((input) =>
      input.dataset.testid?.replace(`config-${section}-`, '').replace('-input', ''),
    )

    expect(inputs).toEqual(expected)
  })

  it('orders general fields as the legacy card does', () => {
    render(<ConfigurationCard />)

    const markers = [
      'config-rskConfirmations',
      'config-btcConfirmations',
      'config-general-publicLiquidityCheck-checkbox',
      'config-general-maxLiquidity-input',
      'config-general-reimbursementWindowBlocks-input',
      'config-general-excessTolerance-input',
    ]
    const rendered = Array.from(
      document.querySelectorAll<HTMLElement>(
        markers.map((testId) => `[data-testid="${testId}"]`).join(','),
      ),
    ).map((element) => element.dataset.testid)

    expect(rendered).toEqual(markers)
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

  it('ignores overlapping save clicks while a request is in flight', async () => {
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    await user.clear(screen.getByTestId('config-general-maxLiquidity-input'))
    await user.type(screen.getByTestId('config-general-maxLiquidity-input'), '7')

    let release!: (value: Response) => void
    apiFetchPostMock.mockImplementationOnce(
      () =>
        new Promise<Response>((resolve) => {
          release = resolve
        }),
    )

    const save = screen.getByTestId('config-save-button')
    await user.click(save)
    await user.click(save)

    expect(apiFetchPostMock).toHaveBeenCalledTimes(1)
    release(okResponse())

    await waitFor(() => {
      expect(toastSuccessMock).toHaveBeenCalledWith(
        'Configuration saved successfully!',
      )
    })
  })

  it('surfaces ApiFetchError body.message for general/pegin/pegout save failures', async () => {
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    await user.clear(screen.getByTestId('config-general-maxLiquidity-input'))
    await user.type(screen.getByTestId('config-general-maxLiquidity-input'), '7')
    await user.clear(screen.getByTestId('config-pegin-maxValue-input'))
    await user.type(screen.getByTestId('config-pegin-maxValue-input'), '11')
    await user.clear(screen.getByTestId('config-pegout-maxValue-input'))
    await user.type(screen.getByTestId('config-pegout-maxValue-input'), '12')

    apiFetchPostMock
      .mockRejectedValueOnce(
        new ApiFetchError(400, 'Bad Request', { message: 'general rejected' }),
      )
      .mockRejectedValueOnce(
        new ApiFetchError(400, 'Bad Request', { message: 'pegin rejected' }),
      )
      .mockRejectedValueOnce(new ApiFetchError(500, 'Error', null))

    await user.click(screen.getByTestId('config-save-button'))

    await waitFor(() => {
      expect(toastErrorMock).toHaveBeenCalledWith('general rejected')
    })
    expect(toastErrorMock).toHaveBeenCalledWith('pegin rejected')
    expect(toastErrorMock).toHaveBeenCalledWith('API request failed: 500 Error')
    expect(toastSuccessMock).not.toHaveBeenCalled()
  })

  it('toasts validation errors instead of posting when a dirty section is invalid', async () => {
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    await user.clear(screen.getByTestId('config-general-maxLiquidity-input'))
    await user.type(screen.getByTestId('config-general-maxLiquidity-input'), '-1')
    await user.clear(screen.getByTestId('config-pegin-penaltyFee-input'))
    await user.type(screen.getByTestId('config-pegin-penaltyFee-input'), 'not-a-number')
    await user.clear(screen.getByTestId('config-pegout-bridgeTransactionMin-input'))
    await user.type(
      screen.getByTestId('config-pegout-bridgeTransactionMin-input'),
      'also-bad',
    )

    await user.click(screen.getByTestId('config-save-button'))

    await waitFor(() => {
      expect(toastErrorMock).toHaveBeenCalled()
    })
    expect(apiFetchPostMock).not.toHaveBeenCalled()
    expect(toastSuccessMock).not.toHaveBeenCalled()
  })

  it('saves dirty general/pegin/pegout sections with edited values and preserved siblings', async () => {
    const user = userEvent.setup()
    render(<ConfigurationCard />)

    // General — checkbox, numeric, percentage, confirmation tiers
    await user.click(screen.getByTestId('config-general-publicLiquidityCheck-checkbox'))
    await user.clear(screen.getByTestId('config-general-reimbursementWindowBlocks-input'))
    await user.type(
      screen.getByTestId('config-general-reimbursementWindowBlocks-input'),
      '200',
    )
    await user.clear(screen.getByTestId('config-general-excessTolerance-input'))
    await user.type(screen.getByTestId('config-general-excessTolerance-input'), '12')
    await user.clear(screen.getByTestId('config-rskConfirmations-amount-0'))
    await user.type(screen.getByTestId('config-rskConfirmations-amount-0'), '1.5')
    await user.clear(screen.getByTestId('config-btcConfirmations-0'))
    await user.type(screen.getByTestId('config-btcConfirmations-0'), '7')

    // Pegin — timing + fee fields (fixed/% already enabled in fixture)
    await user.clear(screen.getByTestId('config-pegin-timeForDeposit-input'))
    await user.type(screen.getByTestId('config-pegin-timeForDeposit-input'), '4000')
    await user.clear(screen.getByTestId('config-pegin-fixedFee-input'))
    await user.type(screen.getByTestId('config-pegin-fixedFee-input'), '0.003')
    await user.clear(screen.getByTestId('config-pegin-feePercentage-input'))
    await user.type(screen.getByTestId('config-pegin-feePercentage-input'), '2')

    // Pegout — enable toggled fees, timing, bridge min
    await user.clear(screen.getByTestId('config-pegout-expireBlocks-input'))
    await user.type(screen.getByTestId('config-pegout-expireBlocks-input'), '600')
    await user.click(screen.getByTestId('config-pegout-fixedFee-checkbox'))
    await user.clear(screen.getByTestId('config-pegout-fixedFee-input'))
    await user.type(screen.getByTestId('config-pegout-fixedFee-input'), '0.005')
    await user.click(screen.getByTestId('config-pegout-feePercentage-checkbox'))
    await user.clear(screen.getByTestId('config-pegout-feePercentage-input'))
    await user.type(screen.getByTestId('config-pegout-feePercentage-input'), '3')
    await user.clear(screen.getByTestId('config-pegout-bridgeTransactionMin-input'))
    await user.type(
      screen.getByTestId('config-pegout-bridgeTransactionMin-input'),
      '3.5',
    )

    apiFetchPostMock
      .mockResolvedValueOnce(okResponse())
      .mockResolvedValueOnce(okResponse())
      .mockResolvedValueOnce(okResponse())

    await user.click(screen.getByTestId('config-save-button'))

    await waitFor(() => {
      expect(apiFetchPostMock).toHaveBeenCalledTimes(3)
    })

    expect(apiFetchPostMock).toHaveBeenCalledWith('/configuration', {
      configuration: {
        publicLiquidityCheck: false,
        maxLiquidity: etherToWei('5'),
        reimbursementWindowBlocks: 200,
        excessTolerance: {
          isFixed: false,
          fixedValue: '0',
          percentageValue: '12',
        },
        rskConfirmations: { [etherToWei('1.5')]: 2 },
        btcConfirmations: { [etherToWei('2')]: 7 },
      },
    })
    expect(apiFetchPostMock).toHaveBeenCalledWith('/pegin/configuration', {
      configuration: {
        timeForDeposit: 4000,
        callTime: 7200,
        penaltyFee: etherToWei('0.001'),
        fixedFee: etherToWei('0.003'),
        feePercentage: '2',
        maxValue: etherToWei('10'),
        minValue: etherToWei('1'),
      },
    })
    expect(apiFetchPostMock).toHaveBeenCalledWith('/pegout/configuration', {
      configuration: {
        timeForDeposit: 3600,
        expireTime: 7200,
        expireBlocks: 600,
        bridgeTransactionMin: etherToWei('3.5'),
        penaltyFee: etherToWei('0.001'),
        fixedFee: etherToWei('0.005'),
        feePercentage: '3',
        maxValue: etherToWei('10'),
        minValue: etherToWei('1'),
      },
    })
    await waitFor(() => {
      expect(toastSuccessMock).toHaveBeenCalledWith(
        'Configuration saved successfully!',
      )
    })
  })
})
