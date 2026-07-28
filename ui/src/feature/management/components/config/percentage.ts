/** Fee-percentage parsing: optional trailing `%`, then clamp to 0–100. */

const PERCENTAGE_PATTERN = /^\d+(\.\d+)?$/

export interface PercentageNormalizeResult {
  /** Normalised numeric string when valid, otherwise the raw input echoed back. */
  value: string
  /** Non-null when the raw input could not be parsed as a percentage. */
  error: string | null
}

export const PERCENTAGE_INVALID_MESSAGE =
  'Fee percentage must be a numeric value between 0% and 100%.'

/** Normalises a raw percentage string: strips a trailing `%` and clamps to 0–100. */
export function normalizePercentage(raw: string): PercentageNormalizeResult {
  const trimmed = raw.trim()
  if (trimmed === '') {
    return { value: '0', error: null }
  }

  const numericPart = trimmed.endsWith('%') ? trimmed.slice(0, -1).trim() : trimmed
  if (!PERCENTAGE_PATTERN.test(numericPart)) {
    return { value: raw, error: PERCENTAGE_INVALID_MESSAGE }
  }

  let num = parseFloat(numericPart)
  if (Number.isNaN(num)) {
    return { value: raw, error: PERCENTAGE_INVALID_MESSAGE }
  }

  if (num < 0) num = 0
  if (num > 100) num = 100

  return { value: String(num), error: null }
}
