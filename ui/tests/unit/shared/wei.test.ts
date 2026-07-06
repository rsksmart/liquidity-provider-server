import { etherToWei, weiToApiAmount, weiToEther } from '@shared/utils/wei'
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
})
