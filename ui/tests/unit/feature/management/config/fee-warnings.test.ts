import {
  checkFeeWarnings,
  ZERO_FEE_WARNING_MESSAGE,
} from '@feature/management/config/fee-warnings'
import { describe, expect, it } from 'vitest'

describe('checkFeeWarnings', () => {
  it('warns on the pegin tab when both fee toggles are off', () => {
    expect(
      checkFeeWarnings('pegin', { fixedFeeEnabled: false, feePercentageEnabled: false }),
    ).toEqual({ shouldWarn: true, message: ZERO_FEE_WARNING_MESSAGE })
  })

  it('warns on the pegout tab when both fee toggles are off', () => {
    expect(
      checkFeeWarnings('pegout', { fixedFeeEnabled: false, feePercentageEnabled: false }),
    ).toEqual({ shouldWarn: true, message: ZERO_FEE_WARNING_MESSAGE })
  })

  it('does not warn when at least one fee toggle is on', () => {
    expect(
      checkFeeWarnings('pegin', { fixedFeeEnabled: true, feePercentageEnabled: false }),
    ).toEqual({ shouldWarn: false, message: null })
    expect(
      checkFeeWarnings('pegout', { fixedFeeEnabled: false, feePercentageEnabled: true }),
    ).toEqual({ shouldWarn: false, message: null })
  })

  it('never warns on the general tab regardless of toggles', () => {
    expect(
      checkFeeWarnings('general', { fixedFeeEnabled: false, feePercentageEnabled: false }),
    ).toEqual({ shouldWarn: false, message: null })
  })

  it('uses the zero-fee warning copy', () => {
    expect(ZERO_FEE_WARNING_MESSAGE).toBe(
      "You have configured a zero-fee setting. This means you won't earn fees from bridging transactions.",
    )
  })
})
