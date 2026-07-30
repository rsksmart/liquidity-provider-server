import type { SectionPrefix } from '@feature/management/components/config/types'

export const ZERO_FEE_WARNING_MESSAGE =
  "You have configured a zero-fee setting. This means you won't earn fees from bridging transactions."

export interface FeeToggleState {
  fixedFeeEnabled: boolean
  feePercentageEnabled: boolean
}

export interface FeeWarningResult {
  shouldWarn: boolean
  message: string | null
}

/**
 * Determines whether the zero-fee warning should be shown for the active tab.
 * General never warns; Pegin/Pegout warn only when both fee toggles are off.
 */
export function checkFeeWarnings(
  activeTab: SectionPrefix,
  toggles: FeeToggleState,
): FeeWarningResult {
  if (activeTab !== 'pegin' && activeTab !== 'pegout') {
    return { shouldWarn: false, message: null }
  }

  const shouldWarn = !toggles.fixedFeeEnabled && !toggles.feePercentageEnabled
  return {
    shouldWarn,
    message: shouldWarn ? ZERO_FEE_WARNING_MESSAGE : null,
  }
}
