/**
 * Matches pkg.AddCollateralRequest. `amount` is wei, so it must reach the server
 * as an exact bare integer literal: Go's *big.Int rejects quoted strings, and a
 * JSON number cannot hold a full 18-decimal wei value. Hence bigint, which
 * stringifyJsonBody writes without precision loss.
 */
export interface AddCollateralRequest {
  amount: bigint
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
