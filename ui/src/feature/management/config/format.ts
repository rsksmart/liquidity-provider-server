/** Confirmation row before formatting for the API (confirmation may be empty). */
export interface ConfirmationEntry {
  amount: string
  confirmation: number | undefined
}

/** General config as edited in the form, before formatting for the API. */
export interface GeneralConfigFormInput {
  rskConfirmations: ConfirmationEntry[]
  btcConfirmations: ConfirmationEntry[]
  [key: string]: unknown
}

/** General config with confirmation arrays collapsed to amount→confirmation maps. */
export interface FormattedGeneralConfig {
  rskConfirmations: Record<string, number>
  btcConfirmations: Record<string, number>
  [key: string]: unknown
}

/**
 * Converts `rskConfirmations`/`btcConfirmations` from an array of
 * `{ amount, confirmation }` rows into a `Record<amountWei, confirmations>`,
 * preserving insertion order. All other keys pass through unchanged.
 */
export function formatGeneralConfig(config: GeneralConfigFormInput): FormattedGeneralConfig {
  const formattedConfig: FormattedGeneralConfig = {
    rskConfirmations: {},
    btcConfirmations: {},
  }
  Object.keys(config).forEach((key) => {
    if (key === 'rskConfirmations' || key === 'btcConfirmations') {
      const record: Record<string, number> = {}
      config[key].forEach((entry) => {
        if (entry.amount && entry.confirmation !== undefined) {
          record[entry.amount] = entry.confirmation
        }
      })
      formattedConfig[key] = record
    } else {
      formattedConfig[key] = config[key]
    }
  })
  return formattedConfig
}
