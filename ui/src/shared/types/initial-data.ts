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
