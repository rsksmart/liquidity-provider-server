import { ApiFetchError, CsrfTokenMissingError } from '@api/management/types/errors'
import { getInitialData } from '@shared/utils/initial-data'

/**
 * Thin fetch wrapper for management API calls: resolves URLs against server BaseUrl,
 * attaches X-CSRF-Token on mutating methods (gorilla/csrf parity with legacy management.js),
 * and maps non-2xx responses to ApiFetchError.
 */
const MUTATING_METHODS = new Set(['POST', 'PUT', 'PATCH', 'DELETE'])

function resolveRequestUrl(input: string): string {
  const baseUrl = getInitialData().data.BaseUrl
  const base = baseUrl.endsWith('/') ? baseUrl : `${baseUrl}/`
  return new URL(input, base).href
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

export async function apiFetch(input: string, init: RequestInit = {}): Promise<Response> {
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
    throw new ApiFetchError(response.status, response.statusText, body)
  }

  return response
}
