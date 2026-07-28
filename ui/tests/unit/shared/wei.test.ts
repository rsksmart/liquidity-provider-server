import { etherToWei, etherToWeiOr, weiToApiAmount, weiToEther } from '@shared/utils/wei'
import { describe, expect, it } from 'vitest'

describe('wei conversion', () => {
  it('converts wei to ether and back', () => {
    const wei = '1000000000000000000'
    expect(weiToEther(wei)).toBe('1')
    expect(etherToWei('1')).toBe(wei)
  })

  it('returns 0 for null and undefined wei', () => {
    expect(weiToEther(null)).toBe('0')
    expect(weiToEther(undefined)).toBe('0')
  })

  it('rejects negative ether values', () => {
    expect(() => etherToWei('-1')).toThrow(/not a valid number/)
  })

  it('rejects non-numeric ether input', () => {
    expect(() => etherToWei('abc')).toThrow(/Failed to convert ether to wei/)
  })

  it('converts wei string to API JSON number', () => {
    expect(weiToApiAmount('11000000000000000000')).toBe(11_000_000_000_000_000_000)
  })

  it('rejects invalid wei input', () => {
    expect(() => weiToEther('not-a-wei-value')).toThrow(/Failed to convert wei to ether/)
  })

  it('rejects non-string/non-number ether input', () => {
    expect(() => etherToWei(null as unknown as string)).toThrow(/Invalid input type for ether/)
  })

  it('keeps sub-wei fractional digits (caller must reject before BigInt)', () => {
    // Documented edge case: values smaller than 1 wei stay fractional after *1e18.
    expect(etherToWei('0.0000000000000000001')).toBe('0.1')
  })

  it('etherToWeiOr maps empty/zero to 0 and falls back on invalid input', () => {
    expect(etherToWeiOr('', 'fallback')).toBe('0')
    expect(etherToWeiOr('0', 'fallback')).toBe('0')
    expect(etherToWeiOr('1', 'fallback')).toBe('1000000000000000000')
    expect(etherToWeiOr('nope', 'fallback')).toBe('fallback')
  })
})
