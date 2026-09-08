import { stringifyJsonBody } from '@api/management/utils/json-body'
import { describe, expect, it } from 'vitest'

describe('stringifyJsonBody', () => {
  it('matches JSON.stringify when no bigint is present', () => {
    const payload = { username: 'u', password: 'p', nested: { list: [1, 'two', null, true] } }
    expect(stringifyJsonBody(payload)).toBe(JSON.stringify(payload))
  })

  it('writes a bigint as a bare integer literal', () => {
    expect(stringifyJsonBody({ amount: 42n })).toBe('{"amount":42}')
  })

  it('keeps every digit of a full 18-decimal wei amount', () => {
    expect(stringifyJsonBody({ amount: 1000000000000000001n })).toBe('{"amount":1000000000000000001}')
    expect(stringifyJsonBody({ amount: 123456789012345678n })).toBe('{"amount":123456789012345678}')
  })

  it('writes values of 1e21 and above in full, not in exponent notation', () => {
    // The float64 equivalent serializes as 1e+21, which Go's big.Int refuses.
    expect(JSON.stringify({ amount: 1e21 })).toBe('{"amount":1e+21}')
    expect(stringifyJsonBody({ amount: 10n ** 21n })).toBe('{"amount":1000000000000000000000}')
  })

  it('writes negative and nested bigint values', () => {
    expect(stringifyJsonBody({ outer: { inner: -7n }, list: [1n, 2n] })).toBe(
      '{"outer":{"inner":-7},"list":[1,2]}',
    )
  })

  it('leaves a string that merely looks numeric quoted', () => {
    expect(stringifyJsonBody({ amount: '1000000000000000001' })).toBe(
      '{"amount":"1000000000000000001"}',
    )
  })

  it('refuses a payload string carrying the internal marker', () => {
    expect(() => stringifyJsonBody({ note: '@@raw-bigint@@1,"admin":true' })).toThrow(
      /reserved marker/,
    )
  })

  it('returns undefined for values JSON.stringify drops', () => {
    expect(stringifyJsonBody(undefined)).toBeUndefined()
  })
})
