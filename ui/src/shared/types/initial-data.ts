export type ProviderType = string

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

export interface FullConfiguration {
  general: Record<string, unknown>
  pegin: Record<string, unknown>
  pegout: Record<string, unknown>
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
