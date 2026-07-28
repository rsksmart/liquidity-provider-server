import Decimal from 'decimal.js'

const WEI_PER_ETHER = new Decimal('1000000000000000000')

export function weiToEther(wei: string | number | null | undefined): string {
  if (wei === null || wei === undefined) {
    return '0'
  }

  try {
    const decimalValue = new Decimal(wei)
    return decimalValue.dividedBy(WEI_PER_ETHER).toString()
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error)
    throw new Error(`Failed to convert wei to ether. Input: "${String(wei)}". Error: ${message}`)
  }
}

export function etherToWei(ether: string | number): string {
  if (typeof ether !== 'string' && typeof ether !== 'number') {
    throw new TypeError(`Invalid input type for ether: ${typeof ether}. Expected a number or string.`)
  }

  try {
    const num = new Decimal(ether)
    if (num.isNegative()) {
      throw new RangeError(`The input "${String(ether)}" is not a valid number.`)
    }
    return num.times(WEI_PER_ETHER).toFixed()
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error)
    throw new Error(`Failed to convert ether to wei. Input: "${String(ether)}". Error: ${message}`)
  }
}

/**
 * Converts ether to wei, returning `fallback` when the input is invalid.
 * Empty / `'0'` normalizes to `'0'` without throwing.
 */
export function etherToWeiOr(ether: string, fallback: string): string {
  const trimmed = ether.trim()
  if (trimmed === '' || trimmed === '0') {
    return '0'
  }
  try {
    return etherToWei(trimmed)
  } catch {
    return fallback
  }
}

/** Wei as JSON number for pkg.AddCollateralRequest (legacy management.js parity). */
export function weiToApiAmount(wei: string): number {
  return Number(wei)
}
