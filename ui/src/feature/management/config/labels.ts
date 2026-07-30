/** Display labels and tooltip copy for configuration fields. */

const DISPLAY_LABELS: Record<string, string> = {
  maxLiquidity: 'Maximum Liquidity',
  excessToleranceFixed: 'Excess Tolerance',
  excessTolerance: 'Excess Tolerance',
}

const TOOLTIP_TEXT: Record<string, string> = {
  timeForDeposit: 'The time (in seconds) for which a deposit is considered valid.',
  expireTime: 'The time (in seconds) after which a quote is considered expired.',
  penaltyFee: 'The penalty fee (in BTC) charged for invalid transactions.',
  callFee: 'The fee (in BTC) charged by the LP for processing a transaction.',
  maxValue: 'The maximum value (in BTC) allowed for a transaction.',
  minValue: 'The minimum value (in BTC) allowed for a transaction.',
  expireBlocks: 'The number of blocks after which a quote is considered expired.',
  bridgeTransactionMin:
    'The amount of rBTC that needs to be gathered in peg out refunds before executing a native peg out.',
  fixedFee: 'A fixed fee charged for transactions.',
  feePercentage: 'A percentage fee charged based on the transaction amount.',
  maxLiquidity:
    'The maximum liquidity (in rBTC) the provider is willing to offer. Must be a positive value with up to 18 decimal places.',
  excessTolerance:
    'The excess tolerance for transactions. Toggle ON for a fixed amount (in rBTC), OFF for a percentage (0-100%).',
}

/** Human-readable label for a config key; falls back to the key itself. */
export function getDisplayLabel(key: string): string {
  return DISPLAY_LABELS[key] ?? key
}

/** Tooltip copy for a config key; falls back to `'No description available'`. */
export function getTooltipText(key: string): string {
  return TOOLTIP_TEXT[key] ?? 'No description available'
}
