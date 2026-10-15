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

const WEI_DECIMALS = 18

// Plain string math, not decimal.js: Decimal arithmetic rounds to 20 significant
// digits, which drops wei from any amount of 100 ether or more written to 18 decimals.
export function etherToWei(ether: string): string {
  if (typeof ether !== 'string') {
    throw new TypeError(`Invalid input type for ether: ${typeof ether}. Expected a string.`)
  }

  const match = /^(\d*)(?:\.(\d*))?$/.exec(ether)
  const whole = match?.[1] ?? ''
  const fraction = match?.[2] ?? ''
  if (whole === '' && fraction === '') {
    throw new Error(`Failed to convert ether to wei. Input: "${ether}". Error: not a valid number.`)
  }
  if (fraction.length > WEI_DECIMALS) {
    throw new Error(
      `Failed to convert ether to wei. Input: "${ether}". Error: more than 18 decimal places is a fraction of a wei.`,
    )
  }

  return BigInt(whole + fraction.padEnd(WEI_DECIMALS, '0')).toString()
}
