import type {
  SetGeneralConfigurationRequest,
  SetPeginConfigurationRequest,
  SetPegoutConfigurationRequest,
} from '@api/management/types/configuration'
import type {
  AddCollateralRequest,
  LoginRequest,
  SetCredentialsRequest,
} from '@api/management/types/requests'
import type { AddTrustedAccountRequest } from '@api/management/types/trusted-accounts'

export interface ManagementPostBodies {
  '/pegin/addCollateral': AddCollateralRequest
  '/pegout/addCollateral': AddCollateralRequest
  '/management/login': LoginRequest
  '/management/credentials': SetCredentialsRequest
  '/management/trusted-accounts': AddTrustedAccountRequest
  '/configuration': SetGeneralConfigurationRequest
  '/pegin/configuration': SetPeginConfigurationRequest
  '/pegout/configuration': SetPegoutConfigurationRequest
}

export type ManagementPostPath = keyof ManagementPostBodies

export type AddCollateralEndpoint =
  | '/pegin/addCollateral'
  | '/pegout/addCollateral'
