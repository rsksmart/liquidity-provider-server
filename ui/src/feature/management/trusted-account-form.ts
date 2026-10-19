import type {
  AddTrustedAccountRequest,
  TrustedAccountValidationDetails,
} from '@api/management/types/trusted-accounts'
import { etherToWei, weiToApiAmount } from '@shared/utils/wei'

/** Form state keys — camelCase to match LPS JSON (`name`, `btcLockingCap`, …). */
export interface TrustedAccountFormValues {
  name: string
  address: string
  btcLockingCap: string
  rbtcLockingCap: string
}

export type TrustedAccountFormField = keyof TrustedAccountFormValues

export type TrustedAccountFormErrors = Partial<
  Record<TrustedAccountFormField, string>
>

/** Go validator detail keys → form field keys (API/JSON camelCase). */
const BACKEND_FIELD_MAP = {
  Name: 'name',
  Address: 'address',
  BtcLockingCap: 'btcLockingCap',
  RbtcLockingCap: 'rbtcLockingCap',
} as const satisfies Record<string, TrustedAccountFormField>

type BackendValidationField = keyof typeof BACKEND_FIELD_MAP

function isBackendValidationField(
  field: string,
): field is BackendValidationField {
  return Object.hasOwn(BACKEND_FIELD_MAP, field)
}

export const EMPTY_TRUSTED_ACCOUNT_FORM: TrustedAccountFormValues = {
  name: '',
  address: '',
  btcLockingCap: '',
  rbtcLockingCap: '',
}

function validatePositiveNumber(
  value: string,
  fieldName: string,
): string | null {
  if (!value || value.trim() === '') {
    return `${fieldName} is required`
  }

  const numValue = Number.parseFloat(value)
  if (Number.isNaN(numValue) || numValue <= 0) {
    return `${fieldName} must be a positive number`
  }

  return null
}

/** Legacy-equivalent client validation for the add-trusted-account form. */
export function validateTrustedAccountForm(
  values: TrustedAccountFormValues,
): TrustedAccountFormErrors {
  const errors: TrustedAccountFormErrors = {}

  if (!values.name.trim()) {
    errors.name = 'Account name is required'
  }

  if (!values.address.trim()) {
    errors.address = 'Account address is required'
  }

  const btcCapError = validatePositiveNumber(
    values.btcLockingCap,
    'BTC Locking Cap',
  )
  if (btcCapError) {
    errors.btcLockingCap = btcCapError
  }

  const rbtcCapError = validatePositiveNumber(
    values.rbtcLockingCap,
    'rBTC Locking Cap',
  )
  if (rbtcCapError) {
    errors.rbtcLockingCap = rbtcCapError
  }

  return errors
}

export function mapBackendValidationDetails(
  details: TrustedAccountValidationDetails,
): TrustedAccountFormErrors {
  const errors: TrustedAccountFormErrors = {}

  for (const [backendField, message] of Object.entries(details)) {
    if (!isBackendValidationField(backendField)) {
      continue
    }
    errors[BACKEND_FIELD_MAP[backendField]] = message
  }

  return errors
}

export function toAddTrustedAccountRequest(
  values: TrustedAccountFormValues,
): AddTrustedAccountRequest {
  return {
    name: values.name.trim(),
    address: values.address.trim(),
    btcLockingCap: weiToApiAmount(etherToWei(values.btcLockingCap.trim())),
    rbtcLockingCap: weiToApiAmount(etherToWei(values.rbtcLockingCap.trim())),
  }
}

export function isValidationErrorBody(body: unknown): body is {
  message: 'validation error'
  details: TrustedAccountValidationDetails
} {
  if (typeof body !== 'object' || body === null) {
    return false
  }
  if (!('message' in body) || !('details' in body)) {
    return false
  }
  return (
    body.message === 'validation error' &&
    typeof body.details === 'object' &&
    body.details !== null
  )
}

export function getTrustedAccountApiErrorMessage(
  body: unknown,
  fallback: string,
): string {
  if (typeof body === 'object' && body !== null) {
    if (
      'details' in body &&
      typeof body.details === 'object' &&
      body.details !== null &&
      'error' in body.details &&
      typeof body.details.error === 'string'
    ) {
      return body.details.error
    }
    if ('message' in body && typeof body.message === 'string') {
      return body.message
    }
  }
  return fallback
}
