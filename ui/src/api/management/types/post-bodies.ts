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
}

export type ManagementPostPath = keyof ManagementPostBodies

export type AddCollateralEndpoint =
  | '/pegin/addCollateral'
  | '/pegout/addCollateral'
