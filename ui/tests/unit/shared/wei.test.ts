import { etherToWei, weiToEther } from '@shared/utils/wei'
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

  it('keeps every digit of a full 18-decimal ether amount', () => {
    expect(etherToWei('1.000000000000000001')).toBe('1000000000000000001')
    expect(etherToWei('0.123456789012345678')).toBe('123456789012345678')
  })

  // decimal.js arithmetic keeps 20 significant digits; these need 21 and 27.
  it('keeps every digit past 20 significant digits', () => {
    expect(etherToWei('100.000000000000000001')).toBe('100000000000000000001')
    expect(etherToWei('123456789.123456789123456789')).toBe('123456789123456789123456789')
  })

  it('accepts a bare leading or trailing decimal point', () => {
    expect(etherToWei('.5')).toBe('500000000000000000')
    expect(etherToWei('1.')).toBe('1000000000000000000')
  })

  it('rejects input that is not plain decimal notation', () => {
    for (const input of ['', '.', '1e3', '0x10', ' 1', '1.2.3', '+1']) {
      expect(() => etherToWei(input)).toThrow(/not a valid number/)
    }
  })

  it('converts amounts of 1000 ether and above without exponent notation', () => {
    expect(etherToWei('1000')).toBe('1000000000000000000000')
  })

  it('rejects invalid wei input', () => {
    expect(() => weiToEther('not-a-wei-value')).toThrow(/Failed to convert wei to ether/)
  })

  it('rejects non-string/non-number ether input', () => {
    expect(() => etherToWei(null as unknown as string)).toThrow(/Invalid input type for ether/)
  })

  it('rejects more than 18 decimal places', () => {
    expect(() => etherToWei('0.0000000000000000001')).toThrow(/fraction of a wei/)
    expect(() => etherToWei('1.0000000000000000000')).toThrow(/fraction of a wei/)
  })
})
