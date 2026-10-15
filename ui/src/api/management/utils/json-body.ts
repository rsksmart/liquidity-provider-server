/**
 * JSON serializer for management request bodies that must carry exact integers.
 *
 * `JSON.stringify` has no way to emit an arbitrary-precision integer: a JSON
 * number is a float64 and loses digits past 2^53, and above 1e21 it serializes
 * as exponent notation. Go's `big.Int` unmarshal rejects both exponent notation
 * and quoted strings, so a wei amount has to reach the server as a bare integer
 * literal. `bigint` values are written that way here; every other value goes
 * through `JSON.stringify` unchanged.
 */
const BIGINT_SENTINEL = '@@raw-bigint@@'
const SENTINEL_PATTERN = new RegExp(`"${BIGINT_SENTINEL}(-?\\d+)"`, 'g')

export function stringifyJsonBody(value: unknown): string | undefined {
  // JSON.stringify is typed as returning string, but it answers undefined for
  // undefined, functions and symbols.
  const json = JSON.stringify(value, (_key, entry: unknown) => {
    if (typeof entry === 'bigint') {
      return `${BIGINT_SENTINEL}${entry.toString()}`
    }
    // A payload string carrying the sentinel would be unquoted below, which
    // would let caller data forge JSON structure. Refuse instead.
    if (typeof entry === 'string' && entry.includes(BIGINT_SENTINEL)) {
      throw new TypeError(`Request body string contains the reserved marker "${BIGINT_SENTINEL}".`)
    }
    return entry
  }) as string | undefined

  return json?.replace(SENTINEL_PATTERN, '$1')
}
