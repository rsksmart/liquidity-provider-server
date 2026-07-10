import {
  ApiFetchError,
  CsrfTokenMissingError,
} from '@api/management/types/errors'
import { getInitialData } from '@shared/utils/initial-data'
import { toast } from 'sonner'

/**
 * Thin fetch wrapper for management API calls: resolves URLs against server BaseUrl,
 * attaches X-CSRF-Token on mutating methods (gorilla/csrf parity with legacy management.js),
 * maps non-2xx responses to ApiFetchError, and redirects on session expiry (401 or 403 session errors).
 */
const MUTATING_METHODS = new Set(['POST', 'PUT', 'PATCH', 'DELETE'])
const LOGIN_PATH = '/management/login'

function getManagementApiOrigin(): string {
  // Vite dev: API goes through the local proxy (same origin). Bootstrap loads real LPS BaseUrl.
  if (import.meta.env.MODE === 'development') return window.location.origin
  const configured = getInitialData().data.BaseUrl.trim()
  if (configured === '') return window.location.origin
  return configured
}

function resolveRequestUrl(input: string): string {
  if (input.startsWith('http://') || input.startsWith('https://')) {
    return new URL(input).href
  }
  const base = getManagementApiOrigin().replace(/\/?$/, '/')
  return new URL(input, base).href
}

function requestPath(input: string): string {
  if (input.startsWith('http://') || input.startsWith('https://')) {
    return new URL(input).pathname
  }
  return input.split('?')[0] ?? input
}

function isLoginPath(path: string): boolean {
  const normalized =
    path.endsWith('/') && path.length > 1 ? path.slice(0, -1) : path
  return normalized === LOGIN_PATH
}

function readCsrfToken(): string {
  const meta = document.querySelector('meta[name="csrf-token"]')
  const token = meta?.getAttribute('content')?.trim()
  if (!token) {
    throw new CsrfTokenMissingError()
  }
  return token
}

async function readErrorBody(response: Response): Promise<unknown> {
  const text = await response.text()
  if (!text) {
    return text
  }

  const contentType = response.headers.get('content-type') ?? ''
  if (!contentType.includes('application/json')) {
    return text
  }

  try {
    return JSON.parse(text) as unknown
  } catch {
    return text
  }
}

function isSessionExpiredBody(body: unknown): boolean {
  if (typeof body !== 'object' || body === null || !('message' in body)) {
    return false
  }

  const message = body.message
  return (
    message === 'session not recognized' ||
    message === 'session validation error'
  )
}

export function isSessionExpiredError(err: unknown): boolean {
  return (
    err instanceof ApiFetchError &&
    (err.status === 401 || isSessionExpiredBody(err.body))
  )
}

export const SESSION_EXPIRED_REDIRECT_DELAY_MS = 800

function handleSessionExpired(): void {
  toast.error('Your session has expired. Please log in again.')
  window.setTimeout(() => {
    window.location.assign('/management/next/login')
  }, SESSION_EXPIRED_REDIRECT_DELAY_MS)
}

export async function apiFetch(
  input: string,
  init: RequestInit = {},
): Promise<Response> {
  const method = (init.method ?? 'GET').toUpperCase()
  const headers = new Headers(init.headers)

  if (MUTATING_METHODS.has(method)) {
    headers.set('X-CSRF-Token', readCsrfToken())
  }

  const response = await fetch(resolveRequestUrl(input), {
    ...init,
    method,
    headers,
  })

  if (!response.ok) {
    const body = await readErrorBody(response)
    const path = requestPath(input)
    const sessionExpired =
      !isLoginPath(path) &&
      (response.status === 401 ||
        (response.status === 403 && isSessionExpiredBody(body)))

    if (sessionExpired) {
      handleSessionExpired()
    }

    throw new ApiFetchError(response.status, response.statusText, body)
  }

  return response
}
