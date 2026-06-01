import { ApiFetchError, CsrfTokenMissingError } from '@api/management/types/errors'
import { getInitialData } from '@shared/utils/initial-data'

const MUTATING_METHODS = new Set(['POST', 'PUT', 'PATCH', 'DELETE'])

function resolveRequestUrl(input: string): string {
  if (/^https?:\/\//i.test(input)) {
    return input
  }

  const baseUrl = getInitialData().data.BaseUrl.replace(/\/$/, '')
  const path = input.startsWith('/') ? input : `/${input}`
  return `${baseUrl}${path}`
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
  const contentType = response.headers.get('content-type') ?? ''
  if (contentType.includes('application/json')) {
    try {
      return await response.json()
    } catch {
      return await response.text()
    }
  }
  return await response.text()
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
