import { getDisplayLabel, getTooltipText } from '@feature/management/config/labels'
import { describe, expect, it } from 'vitest'

describe('getDisplayLabel', () => {
  it('maps known keys to display labels', () => {
    expect(getDisplayLabel('maxLiquidity')).toBe('Maximum Liquidity')
    expect(getDisplayLabel('excessToleranceFixed')).toBe('Excess Tolerance')
    expect(getDisplayLabel('excessTolerance')).toBe('Excess Tolerance')
  })

  it('falls back to the key itself for unknown keys', () => {
    expect(getDisplayLabel('penaltyFee')).toBe('penaltyFee')
    expect(getDisplayLabel('someUnknownKey')).toBe('someUnknownKey')
  })
})

describe('getTooltipText', () => {
  it('maps known keys to tooltip copy', () => {
    expect(getTooltipText('timeForDeposit')).toBe(
      'The time (in seconds) for which a deposit is considered valid.',
    )
    expect(getTooltipText('expireTime')).toBe(
      'The time (in seconds) after which a quote is considered expired.',
    )
    expect(getTooltipText('penaltyFee')).toBe(
      'The penalty fee (in BTC) charged for invalid transactions.',
    )
    expect(getTooltipText('callFee')).toBe(
      'The fee (in BTC) charged by the LP for processing a transaction.',
    )
    expect(getTooltipText('maxValue')).toBe(
      'The maximum value (in BTC) allowed for a transaction.',
    )
    expect(getTooltipText('minValue')).toBe(
      'The minimum value (in BTC) allowed for a transaction.',
    )
    expect(getTooltipText('expireBlocks')).toBe(
      'The number of blocks after which a quote is considered expired.',
    )
    expect(getTooltipText('bridgeTransactionMin')).toBe(
      'The amount of rBTC that needs to be gathered in peg out refunds before executing a native peg out.',
    )
    expect(getTooltipText('fixedFee')).toBe('A fixed fee charged for transactions.')
    expect(getTooltipText('feePercentage')).toBe(
      'A percentage fee charged based on the transaction amount.',
    )
    expect(getTooltipText('maxLiquidity')).toBe(
      'The maximum liquidity (in rBTC) the provider is willing to offer. Must be a positive value with up to 18 decimal places.',
    )
    expect(getTooltipText('excessTolerance')).toBe(
      'The excess tolerance for transactions. Toggle ON for a fixed amount (in rBTC), OFF for a percentage (0-100%).',
    )
  })

  it('falls back to the default copy for unknown keys', () => {
    expect(getTooltipText('somethingElse')).toBe('No description available')
  })
})
