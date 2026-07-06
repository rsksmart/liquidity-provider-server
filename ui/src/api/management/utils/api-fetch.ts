import { ApiFetchError, CsrfTokenMissingError } from '@api/management/types/errors'
import type { ManagementPostBodies } from '@api/management/types/post-bodies'
import { getInitialData } from '@shared/utils/initial-data'
import { toast } from 'sonner'

/**
 * Thin fetch wrapper for management API calls: resolves URLs against server BaseUrl,
 * attaches X-CSRF-Token on mutating methods (gorilla/csrf parity with legacy management.js),
 * maps non-2xx responses to ApiFetchError, and redirects on session expiry (401 or 403 session errors).
 */
const MUTATING_METHODS = new Set(['POST', 'PUT', 'PATCH', 'DELETE'])
const LOGIN_PATH = '/management/login'
const LOGOUT_PATH = '/management/logout'

export type ApiFetchInit = Omit<RequestInit, 'body'> & {
  /** JSON payload — stringified automatically; defaults method to POST when method is omitted. */
  json?: unknown
}

type PostInit = Omit<ApiFetchInit, 'method' | 'json'>

interface ApiFetch {
  (input: string, init?: ApiFetchInit): Promise<Response>
  get: (input: string, init?: Omit<ApiFetchInit, 'method' | 'json'>) => Promise<Response>
  post(input: '/management/logout', init?: PostInit): Promise<Response>
  post<Path extends keyof ManagementPostBodies>(
    input: Path,
    json: ManagementPostBodies[Path],
    init?: PostInit,
  ): Promise<Response>
}

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
  const normalized = path.endsWith('/') && path.length > 1 ? path.slice(0, -1) : path
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
  return message === 'session not recognized' || message === 'session validation error'
}

export function isSessionExpiredError(err: unknown): boolean {
  return (
    err instanceof ApiFetchError &&
    (err.status === 401 || isSessionExpiredBody(err.body))
  )
}

export const SESSION_EXPIRED_REDIRECT_DELAY_MS = 800

let sessionExpiredHandled = false

/** Test-only: clear the once-per-page-load session-expired guard. */
export function resetSessionExpiredHandling(): void {
  sessionExpiredHandled = false
}

function handleSessionExpired(): void {
  if (sessionExpiredHandled) return
  sessionExpiredHandled = true
  toast.error('Your session has expired. Please log in again.')
  window.setTimeout(() => {
    window.location.assign('/management/next/login')
  }, SESSION_EXPIRED_REDIRECT_DELAY_MS)
}

async function apiFetchImpl(input: string, init: ApiFetchInit = {}): Promise<Response> {
  const { json, ...requestInit } = init
  const method = (requestInit.method ?? (json !== undefined ? 'POST' : 'GET')).toUpperCase()
  const body = json !== undefined ? JSON.stringify(json) : undefined
  const headers = new Headers(requestInit.headers)

  if (MUTATING_METHODS.has(method)) {
    headers.set('X-CSRF-Token', readCsrfToken())
    if (body != null && !headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }
  }

  const response = await fetch(resolveRequestUrl(input), {
    ...requestInit,
    method,
    body,
    headers,
  })

  if (!response.ok) {
    const responseBody = await readErrorBody(response)
    const path = requestPath(input)
    const sessionExpired =
      !isLoginPath(path) &&
      (response.status === 401 ||
        (response.status === 403 && isSessionExpiredBody(responseBody)))

    if (sessionExpired) {
      handleSessionExpired()
    }

    throw new ApiFetchError(response.status, response.statusText, responseBody)
  }

  return response
}

function apiFetchPost(input: '/management/logout', init?: PostInit): Promise<Response>
function apiFetchPost<Path extends keyof ManagementPostBodies>(
  input: Path,
  json: ManagementPostBodies[Path],
  init?: PostInit,
): Promise<Response>
function apiFetchPost(input: string, jsonOrInit?: unknown, init: PostInit = {}): Promise<Response> {
  if (input === LOGOUT_PATH) {
    const logoutInit = jsonOrInit !== undefined ? (jsonOrInit as PostInit) : init
    return apiFetchImpl(input, { ...logoutInit, method: 'POST' })
  }

  return apiFetchImpl(input, {
    ...init,
    method: 'POST',
    ...(jsonOrInit !== undefined ? { json: jsonOrInit } : {}),
  })
}

export const apiFetch: ApiFetch = Object.assign(apiFetchImpl, {
  get: (input: string, init: Omit<ApiFetchInit, 'method' | 'json'> = {}) =>
    apiFetchImpl(input, { ...init, method: 'GET' }),
  post: apiFetchPost,
})
