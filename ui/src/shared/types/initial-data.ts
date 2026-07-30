export type ProviderType = 0 | 1 | 2

export type ConfirmationsPerAmount = Record<string, number>

export type WeiValue = string

export type BigFloatValue = string

export interface RegisteredLiquidityProvider {
  id: number
  address: string
  name: string
  apiBaseUrl: string
  status: boolean
  providerType: ProviderType
}

export interface ManagementColdWallet {
  BtcAddress: string
  RskAddress: string
  Label: string
}

export interface ExcessTolerance {
  isFixed: boolean
  percentageValue: BigFloatValue
  fixedValue: WeiValue
}

export interface GeneralConfiguration {
  rskConfirmations: ConfirmationsPerAmount
  btcConfirmations: ConfirmationsPerAmount
  publicLiquidityCheck: boolean
  maxLiquidity: WeiValue | null
  reimbursementWindowBlocks: number
  excessTolerance: ExcessTolerance
}

export interface QuoteConfigurationBase {
  timeForDeposit: number
  penaltyFee: WeiValue
  fixedFee: WeiValue
  feePercentage: BigFloatValue
  maxValue: WeiValue
  minValue: WeiValue
}

export interface PeginConfiguration extends QuoteConfigurationBase {
  callTime: number
}

export interface PegoutConfiguration extends QuoteConfigurationBase {
  expireTime: number
  expireBlocks: number
  bridgeTransactionMin: WeiValue
}

export interface FullConfiguration {
  general: GeneralConfiguration
  pegin: PeginConfiguration
  pegout: PegoutConfiguration
}

export interface ManagementTemplateData {
  CredentialsSet: boolean
  BaseUrl: string
  BtcAddress: string
  RskAddress: string
  ProviderData: RegisteredLiquidityProvider
  ColdWallet: ManagementColdWallet
  Configuration: FullConfiguration
}

export interface InitialDataPayload {
  loggedIn: boolean
  data: ManagementTemplateData
}

/**
 * Wire shape of the bootstrap payload. `entities.Wei` and `utils.BigFloat`
 * marshal through `big.Int` / `big.Float`, so every wei and percentage field
 * reaches the browser as an unquoted JSON number. These are narrowed to the
 * display-domain strings above when initial-data is read.
 */
export type WireNumeric = string | number

export interface WireExcessTolerance {
  isFixed: boolean
  percentageValue: WireNumeric
  fixedValue: WireNumeric
}

export interface WireGeneralConfiguration
  extends Omit<GeneralConfiguration, 'maxLiquidity' | 'excessTolerance'> {
  maxLiquidity: WireNumeric | null
  excessTolerance: WireExcessTolerance
}

export interface WireQuoteConfigurationBase
  extends Omit<
    QuoteConfigurationBase,
    'penaltyFee' | 'fixedFee' | 'feePercentage' | 'maxValue' | 'minValue'
  > {
  penaltyFee: WireNumeric
  fixedFee: WireNumeric
  feePercentage: WireNumeric
  maxValue: WireNumeric
  minValue: WireNumeric
}

export interface WirePeginConfiguration extends WireQuoteConfigurationBase {
  callTime: number
}

export interface WirePegoutConfiguration extends WireQuoteConfigurationBase {
  expireTime: number
  expireBlocks: number
  bridgeTransactionMin: WireNumeric
}

export interface WireFullConfiguration {
  general: WireGeneralConfiguration
  pegin: WirePeginConfiguration
  pegout: WirePegoutConfiguration
}

export interface WireManagementTemplateData
  extends Omit<ManagementTemplateData, 'Configuration'> {
  Configuration: WireFullConfiguration
}

export interface WireInitialDataPayload {
  loggedIn: boolean
  data: WireManagementTemplateData
}
