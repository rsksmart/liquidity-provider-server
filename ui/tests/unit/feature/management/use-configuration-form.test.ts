import { useConfigurationForm } from '@feature/management/hooks/use-configuration-form'
import type {
  FullConfiguration,
  WireFullConfiguration,
  WireInitialDataPayload,
} from '@shared/types/initial-data'
import { etherToWei, weiToEther } from '@shared/utils/wei'
import { act, renderHook } from '@testing-library/react'
import { loggedInFixture, wireConfigurationFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { beforeEach, describe, expect, it } from 'vitest'

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

function seedConfiguration(
  configuration: FullConfiguration | WireFullConfiguration,
): void {
  const payload: WireInitialDataPayload = {
    ...loggedInFixture,
    data: { ...loggedInFixture.data, Configuration: configuration },
  }
  seedInitialData(payload)
}

describe('useConfigurationForm', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedConfiguration(sampleConfiguration)
  })

  it('loads display-domain values from the configuration', () => {
    const { result } = renderHook(() => useConfigurationForm())

    expect(result.current.general.maxLiquidity).toBe('5')
    expect(result.current.general.publicLiquidityCheck).toBe(true)
    expect(result.current.general.reimbursementWindowBlocks).toBe('100')
    expect(result.current.general.rskConfirmations).toEqual([
      { amountEther: '1', confirmations: '2' },
    ])
    expect(result.current.pegin.penaltyFee).toBe('0.001')
    expect(result.current.pegin.fixedFeeEnabled).toBe(true)
    expect(result.current.pegout.fixedFeeEnabled).toBe(false)
    expect(result.current.pegout.feePercentageEnabled).toBe(false)
  })

  it('starts clean and dirties per section on edit', () => {
    const { result } = renderHook(() => useConfigurationForm())

    expect(result.current.dirty.any).toBe(false)

    act(() => {
      result.current.updatePegin({ penaltyFee: '0.5' })
    })

    expect(result.current.dirty.pegin).toBe(true)
    expect(result.current.dirty.general).toBe(false)
    expect(result.current.dirty.any).toBe(true)
  })

  it('is clean again after a round-trip edit restoring the baseline', () => {
    const { result } = renderHook(() => useConfigurationForm())

    act(() => {
      result.current.updatePegin({ penaltyFee: '0.5' })
    })
    expect(result.current.dirty.pegin).toBe(true)

    act(() => {
      result.current.updatePegin({
        penaltyFee: weiToEther(sampleConfiguration.pegin.penaltyFee),
      })
    })
    expect(result.current.dirty.pegin).toBe(false)
    expect(result.current.dirty.any).toBe(false)
  })

  it('builds a wei-encoded pegin payload', () => {
    const { result } = renderHook(() => useConfigurationForm())

    act(() => {
      result.current.updatePegin({ penaltyFee: '0.5' })
    })

    const built = result.current.build()
    expect(built.pegin.config).toEqual({
      timeForDeposit: 3600,
      callTime: 7200,
      penaltyFee: etherToWei('0.5'),
      fixedFee: etherToWei('0.002'),
      feePercentage: '1.5',
      maxValue: etherToWei('10'),
      minValue: etherToWei('1'),
    })
    expect(built.pegin.errors).toEqual([])
  })

  it('builds a formatted general payload with confirmation records', () => {
    const { result } = renderHook(() => useConfigurationForm())

    act(() => {
      result.current.updateGeneral({ maxLiquidity: '7' })
    })

    const built = result.current.build()
    expect(built.general.config).toEqual({
      publicLiquidityCheck: true,
      maxLiquidity: etherToWei('7'),
      reimbursementWindowBlocks: 100,
      excessTolerance: { isFixed: false, fixedValue: '0', percentageValue: '10' },
      rskConfirmations: { '1000000000000000000': 2 },
      btcConfirmations: { '2000000000000000000': 6 },
    })
    expect(built.general.errors).toEqual([])
  })

  it('reports errors for an invalid maxLiquidity', () => {
    const { result } = renderHook(() => useConfigurationForm())

    act(() => {
      result.current.updateGeneral({ maxLiquidity: '-1' })
    })

    const built = result.current.build()
    expect(built.general.config).toBeNull()
    expect(built.general.errors.length).toBeGreaterThan(0)
  })

  it('derives fee toggles from the numeric configuration a live LPS sends', () => {
    seedConfiguration(wireConfigurationFixture)
    const { result } = renderHook(() => useConfigurationForm())

    expect(result.current.pegin.fixedFee).toBe('0.0002')
    expect(result.current.pegin.fixedFeeEnabled).toBe(true)
    expect(result.current.pegin.feePercentage).toBe('0.33')
    expect(result.current.pegin.feePercentageEnabled).toBe(true)
    expect(result.current.pegout.fixedFeeEnabled).toBe(false)
    expect(result.current.pegout.feePercentageEnabled).toBe(false)
    expect(result.current.general.maxLiquidity).toBe('2000')
    expect(result.current.dirty.any).toBe(false)
  })

  it('clears dirtiness after markSaved', () => {
    const { result } = renderHook(() => useConfigurationForm())

    act(() => {
      result.current.updatePegout({ maxValue: '99' })
    })
    expect(result.current.dirty.pegout).toBe(true)

    act(() => {
      result.current.markSaved()
    })
    expect(result.current.dirty.any).toBe(false)
  })
})
