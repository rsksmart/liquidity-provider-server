/** Trusted-account DTOs for management UI — mirrors existing LPS JSON shapes. */

export interface TrustedAccount {
  name?: string
  address: string
  btcLockingCap: string | number
  rbtcLockingCap: string | number
}

export interface TrustedAccountsListResponse {
  accounts: TrustedAccount[]
}

export interface AddTrustedAccountRequest {
  name: string
  address: string
  btcLockingCap: number
  rbtcLockingCap: number
}

/** Backend validation-error `details` map (field → message). */
export type TrustedAccountValidationDetails = Record<string, string>
