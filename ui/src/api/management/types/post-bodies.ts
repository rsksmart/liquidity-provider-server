import type {
  AddCollateralRequest,
  LoginRequest,
  SetCredentialsRequest,
} from '@api/management/types/requests'

export interface ManagementPostBodies {
  '/pegin/addCollateral': AddCollateralRequest
  '/pegout/addCollateral': AddCollateralRequest
  '/management/login': LoginRequest
  '/management/credentials': SetCredentialsRequest
}

export type ManagementPostPath = keyof ManagementPostBodies

export type AddCollateralEndpoint = '/pegin/addCollateral' | '/pegout/addCollateral'
