import { formatGeneralConfig } from '@feature/management/config/format'
import { describe, expect, it } from 'vitest'

describe('formatGeneralConfig', () => {
  it('converts confirmation arrays to amount-keyed records preserving insertion order', () => {
    // Wei-scale amounts (> 2^32) are non-array-index keys, so JS objects keep
    // forEach insertion order rather than numeric ordering.
    const result = formatGeneralConfig({
      rskConfirmations: [
        { amount: '2000000000000000000', confirmation: 4 },
        { amount: '1000000000000000000', confirmation: 2 },
      ],
      btcConfirmations: [{ amount: '5000000000000000000', confirmation: 6 }],
      maxLiquidity: '1000000000000000000',
      publicLiquidityCheck: true,
    })

    expect(result.rskConfirmations).toEqual({
      '2000000000000000000': 4,
      '1000000000000000000': 2,
    })
    expect(Object.keys(result.rskConfirmations)).toEqual([
      '2000000000000000000',
      '1000000000000000000',
    ])
    expect(result.btcConfirmations).toEqual({ '5000000000000000000': 6 })
    expect(result.maxLiquidity).toBe('1000000000000000000')
    expect(result.publicLiquidityCheck).toBe(true)
  })

  it('skips entries without an amount or confirmation', () => {
    const result = formatGeneralConfig({
      rskConfirmations: [
        { amount: '', confirmation: 4 },
        { amount: '1000', confirmation: undefined },
        { amount: '2000', confirmation: 3 },
      ],
      btcConfirmations: [],
    })

    expect(result.rskConfirmations).toEqual({ '2000': 3 })
    expect(result.btcConfirmations).toEqual({})
  })
})
