import Decimal from 'decimal.js'

export function weiToEther(wei: string | number | null | undefined): string {
  if (wei === null || wei === undefined) {
    return '0'
  }

  try {
    const decimalValue = new Decimal(wei)
    return decimalValue.dividedBy(new Decimal(1e18)).toString()
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
    return num.times(new Decimal(1e18)).toFixed()
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error)
    throw new Error(`Failed to convert ether to wei. Input: "${String(ether)}". Error: ${message}`)
  }
}

/** Wei as JSON number for pkg.AddCollateralRequest (legacy management.js parity). */
export function weiToApiAmount(wei: string): number {
  return Number(wei)
}
