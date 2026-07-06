/** Matches pkg.AddCollateralRequest — amount is a JSON number (Go *big.Int rejects quoted strings). */
export interface AddCollateralRequest {
  amount: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface SetCredentialsRequest {
  oldUsername: string
  oldPassword: string
  newUsername: string
  newPassword: string
}
