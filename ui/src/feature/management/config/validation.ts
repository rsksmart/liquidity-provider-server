/** Configuration field validation and key classification helpers. */

import type { ConfirmationEntry } from './format'
export type { ConfirmationEntry } from './format'

export interface ValidationResult {
  isValid: boolean
  error: string | null
}

export interface ConfigValidationResult {
  isValid: boolean
  errors: string[]
}

type NumericInput = string | number | null | undefined

function isAbsent(value: NumericInput): boolean {
  return value === null || value === undefined || value === ''
}

const FEE_KEYS = [
  'penaltyFee',
  'callFee',
  'maxValue',
  'minValue',
  'bridgeTransactionMin',
  'fixedFee',
  'maxLiquidity',
]

export function isFeeKey(key: string): boolean {
  return FEE_KEYS.includes(key)
}

export function isMaxLiquidityKey(key: string): boolean {
  return key === 'maxLiquidity'
}

export function isfeePercentageKey(key: string): boolean {
  return key === 'feePercentage'
}

export function isExcessToleranceKey(key: string): boolean {
  return key === 'excessToleranceFixed' || key === 'excessTolerancePercentage'
}

export function isExcessToleranceFixedKey(key: string): boolean {
  return key === 'excessToleranceFixed'
}

export function isExcessTolerancePercentageKey(key: string): boolean {
  return key === 'excessTolerancePercentage'
}

export function isToggableFeeKey(key: string): boolean {
  return key === 'fixedFee' || key === 'feePercentage'
}

/** Validates a maxLiquidity value (in RBTC). */
export function validateMaxLiquidity(value: NumericInput): ValidationResult {
  if (isAbsent(value)) {
    return { isValid: false, error: 'Max liquidity is required' }
  }

  const strValue = String(value).trim()

  const num = parseFloat(strValue)
  if (isNaN(num)) {
    return { isValid: false, error: 'Max liquidity must be a valid number' }
  }

  if (num <= 0) {
    return { isValid: false, error: 'Max liquidity must be a positive number' }
  }

  const decimalPart = strValue.split('.')[1]
  if (decimalPart && decimalPart.length > 18) {
    return { isValid: false, error: 'Max liquidity cannot have more than 18 decimal places' }
  }

  return { isValid: true, error: null }
}

/** Validates an excessToleranceFixed value (optional, non-negative RBTC). */
export function validateExcessToleranceFixed(value: NumericInput): ValidationResult {
  if (isAbsent(value)) {
    return { isValid: true, error: null }
  }

  const strValue = String(value).trim()

  const num = parseFloat(strValue)
  if (isNaN(num)) {
    return { isValid: false, error: 'Excess tolerance fixed must be a valid number' }
  }

  if (num < 0) {
    return { isValid: false, error: 'Excess tolerance fixed must be a non-negative number' }
  }

  const decimalPart = strValue.split('.')[1]
  if (decimalPart && decimalPart.length > 18) {
    return { isValid: false, error: 'Excess tolerance fixed cannot have more than 18 decimal places' }
  }

  return { isValid: true, error: null }
}

/** Validates an excessTolerancePercentage value (optional, 0-100). */
export function validateExcessTolerancePercentage(value: NumericInput): ValidationResult {
  if (isAbsent(value)) {
    return { isValid: true, error: null }
  }

  const strValue = String(value).trim()
  const num = parseFloat(strValue)

  if (isNaN(num)) {
    return { isValid: false, error: 'Excess tolerance percentage must be a valid number' }
  }

  if (num < 0) {
    return { isValid: false, error: 'Excess tolerance percentage must be non-negative' }
  }

  if (num > 100) {
    return { isValid: false, error: 'Excess tolerance percentage cannot exceed 100%' }
  }

  return { isValid: true, error: null }
}

/** Runtime type label: null/undefined → `'undefined'`, arrays → `'array'`, else `typeof`. */
export function inferType(value: unknown): string {
  if (value === null || value === undefined) return 'undefined'
  if (Array.isArray(value)) return 'array'
  return typeof value
}

/**
 * Validates that a submitted config matches the types of the original config.
 * Fee/percentage/excess-tolerance keys tolerate string↔number differences.
 */
export function validateConfig(
  config: Record<string, unknown>,
  originalConfig: Record<string, unknown>,
): ConfigValidationResult {
  const errors: string[] = []
  const confirmationKeys = ['rskConfirmations', 'btcConfirmations']

  for (const [key, value] of Object.entries(config)) {
    const expectedValue = originalConfig[key]
    const expectedType = inferType(expectedValue)
    let actualType = inferType(value)
    if (
      (isFeeKey(key) || isfeePercentageKey(key) || isExcessToleranceKey(key)) &&
      ((expectedType === 'number' && actualType === 'string') ||
        (expectedType === 'string' && actualType === 'number'))
    ) {
      actualType = expectedType
    }
    if (expectedType === 'undefined') continue
    if (confirmationKeys.includes(key)) {
      if (actualType !== 'object') {
        errors.push(`Invalid type for ${key}: expected object, got ${actualType}`)
        continue
      }
      for (const [subKey, subValue] of Object.entries(value as Record<string, unknown>)) {
        if (inferType(subValue) !== 'number') {
          errors.push(
            `Invalid type for ${key} confirmation value of amount ${subKey}: expected number, got ${inferType(subValue)}`,
          )
        }
      }
    } else if (actualType !== expectedType) {
      errors.push(`Invalid type for ${key}: expected ${expectedType}, got ${actualType}`)
    }
  }
  return { isValid: errors.length === 0, errors }
}

/** Returns true when a confirmation array contains duplicate rBTC amounts. */
export function hasDuplicateConfirmationAmounts(confirmationArray: ConfirmationEntry[]): boolean {
  const amounts = confirmationArray.map((entry) => entry.amount)
  const uniqueAmounts = new Set(amounts)
  return uniqueAmounts.size < amounts.length
}
