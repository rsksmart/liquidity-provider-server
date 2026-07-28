import {
  hasDuplicateConfirmationAmounts,
  inferType,
  isExcessToleranceFixedKey,
  isExcessToleranceKey,
  isExcessTolerancePercentageKey,
  isFeeKey,
  isfeePercentageKey,
  isMaxLiquidityKey,
  isToggableFeeKey,
  validateConfig,
  validateExcessToleranceFixed,
  validateExcessTolerancePercentage,
  validateMaxLiquidity,
} from '@feature/management/config/validation'
import { describe, expect, it } from 'vitest'

describe('validateMaxLiquidity', () => {
  it('requires a value', () => {
    expect(validateMaxLiquidity(null)).toEqual({ isValid: false, error: 'Max liquidity is required' })
    expect(validateMaxLiquidity(undefined)).toEqual({ isValid: false, error: 'Max liquidity is required' })
    expect(validateMaxLiquidity('')).toEqual({ isValid: false, error: 'Max liquidity is required' })
  })

  it('rejects non-numeric values', () => {
    expect(validateMaxLiquidity('abc')).toEqual({
      isValid: false,
      error: 'Max liquidity must be a valid number',
    })
  })

  it('rejects non-positive values', () => {
    expect(validateMaxLiquidity('0')).toEqual({
      isValid: false,
      error: 'Max liquidity must be a positive number',
    })
    expect(validateMaxLiquidity('-1')).toEqual({
      isValid: false,
      error: 'Max liquidity must be a positive number',
    })
  })

  it('rejects more than 18 decimal places', () => {
    expect(validateMaxLiquidity('1.0000000000000000001')).toEqual({
      isValid: false,
      error: 'Max liquidity cannot have more than 18 decimal places',
    })
  })

  it('accepts a valid positive value with up to 18 decimals', () => {
    expect(validateMaxLiquidity('1.5')).toEqual({ isValid: true, error: null })
    expect(validateMaxLiquidity('1.000000000000000001')).toEqual({ isValid: true, error: null })
    expect(validateMaxLiquidity(10)).toEqual({ isValid: true, error: null })
  })
})

describe('validateExcessToleranceFixed', () => {
  it('treats empty as valid (optional field)', () => {
    expect(validateExcessToleranceFixed('')).toEqual({ isValid: true, error: null })
    expect(validateExcessToleranceFixed(null)).toEqual({ isValid: true, error: null })
    expect(validateExcessToleranceFixed(undefined)).toEqual({ isValid: true, error: null })
  })

  it('rejects non-numeric values', () => {
    expect(validateExcessToleranceFixed('abc')).toEqual({
      isValid: false,
      error: 'Excess tolerance fixed must be a valid number',
    })
  })

  it('rejects negative values', () => {
    expect(validateExcessToleranceFixed('-1')).toEqual({
      isValid: false,
      error: 'Excess tolerance fixed must be a non-negative number',
    })
  })

  it('rejects more than 18 decimal places', () => {
    expect(validateExcessToleranceFixed('1.0000000000000000001')).toEqual({
      isValid: false,
      error: 'Excess tolerance fixed cannot have more than 18 decimal places',
    })
  })

  it('accepts zero and valid non-negative values', () => {
    expect(validateExcessToleranceFixed('0')).toEqual({ isValid: true, error: null })
    expect(validateExcessToleranceFixed('2.5')).toEqual({ isValid: true, error: null })
  })
})

describe('validateExcessTolerancePercentage', () => {
  it('treats empty as valid (optional field)', () => {
    expect(validateExcessTolerancePercentage('')).toEqual({ isValid: true, error: null })
    expect(validateExcessTolerancePercentage(null)).toEqual({ isValid: true, error: null })
    expect(validateExcessTolerancePercentage(undefined)).toEqual({ isValid: true, error: null })
  })

  it('rejects non-numeric values', () => {
    expect(validateExcessTolerancePercentage('abc')).toEqual({
      isValid: false,
      error: 'Excess tolerance percentage must be a valid number',
    })
  })

  it('rejects negative values', () => {
    expect(validateExcessTolerancePercentage('-1')).toEqual({
      isValid: false,
      error: 'Excess tolerance percentage must be non-negative',
    })
  })

  it('rejects values above 100', () => {
    expect(validateExcessTolerancePercentage('101')).toEqual({
      isValid: false,
      error: 'Excess tolerance percentage cannot exceed 100%',
    })
  })

  it('accepts values within 0-100', () => {
    expect(validateExcessTolerancePercentage('0')).toEqual({ isValid: true, error: null })
    expect(validateExcessTolerancePercentage('50')).toEqual({ isValid: true, error: null })
    expect(validateExcessTolerancePercentage('100')).toEqual({ isValid: true, error: null })
  })
})

describe('fee-key helpers', () => {
  it('isFeeKey', () => {
    for (const key of [
      'penaltyFee',
      'callFee',
      'maxValue',
      'minValue',
      'bridgeTransactionMin',
      'fixedFee',
      'maxLiquidity',
    ]) {
      expect(isFeeKey(key)).toBe(true)
    }
    expect(isFeeKey('feePercentage')).toBe(false)
    expect(isFeeKey('timeForDeposit')).toBe(false)
  })

  it('isMaxLiquidityKey', () => {
    expect(isMaxLiquidityKey('maxLiquidity')).toBe(true)
    expect(isMaxLiquidityKey('fixedFee')).toBe(false)
  })

  it('isfeePercentageKey', () => {
    expect(isfeePercentageKey('feePercentage')).toBe(true)
    expect(isfeePercentageKey('fixedFee')).toBe(false)
  })

  it('isExcessToleranceKey', () => {
    expect(isExcessToleranceKey('excessToleranceFixed')).toBe(true)
    expect(isExcessToleranceKey('excessTolerancePercentage')).toBe(true)
    expect(isExcessToleranceKey('excessTolerance')).toBe(false)
  })

  it('isExcessToleranceFixedKey', () => {
    expect(isExcessToleranceFixedKey('excessToleranceFixed')).toBe(true)
    expect(isExcessToleranceFixedKey('excessTolerancePercentage')).toBe(false)
  })

  it('isExcessTolerancePercentageKey', () => {
    expect(isExcessTolerancePercentageKey('excessTolerancePercentage')).toBe(true)
    expect(isExcessTolerancePercentageKey('excessToleranceFixed')).toBe(false)
  })

  it('isToggableFeeKey', () => {
    expect(isToggableFeeKey('fixedFee')).toBe(true)
    expect(isToggableFeeKey('feePercentage')).toBe(true)
    expect(isToggableFeeKey('penaltyFee')).toBe(false)
  })
})

describe('inferType', () => {
  it('classifies null, arrays, and primitives', () => {
    expect(inferType(null)).toBe('undefined')
    expect(inferType(undefined)).toBe('undefined')
    expect(inferType([])).toBe('array')
    expect(inferType({})).toBe('object')
    expect(inferType('x')).toBe('string')
    expect(inferType(1)).toBe('number')
    expect(inferType(true)).toBe('boolean')
  })
})

describe('validateConfig', () => {
  it('passes when types match the original config', () => {
    const original = { penaltyFee: '0', timeForDeposit: 0 }
    const config = { penaltyFee: '5', timeForDeposit: 10 }
    expect(validateConfig(config, original)).toEqual({ isValid: true, errors: [] })
  })

  it('coerces string/number mismatches for fee keys', () => {
    const original = { penaltyFee: '0', feePercentage: '0', excessToleranceFixed: '0' }
    const config = { penaltyFee: 5, feePercentage: 1.5, excessToleranceFixed: 2 }
    expect(validateConfig(config, original)).toEqual({ isValid: true, errors: [] })
  })

  it('flags type mismatches on non-fee keys', () => {
    const original = { timeForDeposit: 0 }
    const config = { timeForDeposit: 'ten' }
    expect(validateConfig(config, original)).toEqual({
      isValid: false,
      errors: ['Invalid type for timeForDeposit: expected number, got string'],
    })
  })

  it('skips keys that are undefined in the original config', () => {
    const original = {}
    const config = { extra: 'value' }
    expect(validateConfig(config, original)).toEqual({ isValid: true, errors: [] })
  })

  it('validates confirmation maps as objects of numbers', () => {
    const original = { rskConfirmations: {} }
    const okConfig = { rskConfirmations: { '1000': 2 } }
    expect(validateConfig(okConfig, original)).toEqual({ isValid: true, errors: [] })

    const badTypeConfig = { rskConfirmations: [] }
    expect(validateConfig(badTypeConfig, original)).toEqual({
      isValid: false,
      errors: ['Invalid type for rskConfirmations: expected object, got array'],
    })

    const badValueConfig = { rskConfirmations: { '1000': 'two' } }
    expect(validateConfig(badValueConfig, original)).toEqual({
      isValid: false,
      errors: [
        'Invalid type for rskConfirmations confirmation value of amount 1000: expected number, got string',
      ],
    })
  })
})

describe('hasDuplicateConfirmationAmounts', () => {
  it('returns false when all amounts are unique', () => {
    expect(
      hasDuplicateConfirmationAmounts([
        { amount: '1000', confirmation: 2 },
        { amount: '2000', confirmation: 4 },
      ]),
    ).toBe(false)
  })

  it('returns true when an amount is duplicated', () => {
    expect(
      hasDuplicateConfirmationAmounts([
        { amount: '1000', confirmation: 2 },
        { amount: '1000', confirmation: 4 },
      ]),
    ).toBe(true)
  })
})
