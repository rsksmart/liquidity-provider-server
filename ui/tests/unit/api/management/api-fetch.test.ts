import { ApiFetchError, CsrfTokenMissingError } from '@api/management/types/errors'
import { apiFetch } from '@api/management/utils/api-fetch'
import { resetInitialDataCacheForTests } from '@shared/utils/initial-data'
import { loggedOutFixture } from '@tests/fixtures/logged-out'
import { seedInitialData } from '@tests/helpers/seed-initial-data'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

describe('apiFetch', () => {
  beforeEach(() => {
    document.head.innerHTML = ''
    document.body.innerHTML = ''
    resetInitialDataCacheForTests()
    vi.stubGlobal('fetch', vi.fn())
  })

  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it.each(['POST', 'PUT', 'PATCH', 'DELETE'] as const)(
    'attaches X-CSRF-Token on %s requests',
    async (method) => {
      seedInitialData(loggedOutFixture, { csrfToken: 'csrf-mutate-token' })
      const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

      await apiFetch('/management/login', { method })

      const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
      expect(new Headers(init.headers).get('X-CSRF-Token')).toBe('csrf-mutate-token')
    },
  )

  it('omits X-CSRF-Token on GET requests', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-get-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch('/management/status', { method: 'GET' })

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(new Headers(init.headers).has('X-CSRF-Token')).toBe(false)
  })

  it('throws CsrfTokenMissingError before fetch when meta is missing on POST', async () => {
    seedInitialData(loggedOutFixture)
    const fetchMock = vi.mocked(fetch)

    await expect(apiFetch('/management/login', { method: 'POST' })).rejects.toBeInstanceOf(
      CsrfTokenMissingError,
    )
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('throws ApiFetchError with status and body on non-2xx responses', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ error: 'bad request' }), {
        status: 400,
        statusText: 'Bad Request',
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(apiFetch('/management/login', { method: 'POST' })).rejects.toMatchObject({
      name: 'ApiFetchError',
      status: 400,
      statusText: 'Bad Request',
      body: { error: 'bad request' },
    } satisfies Partial<ApiFetchError>)
  })

  it('resolves relative URLs against BaseUrl from initial data', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch('/management/pegin/configuration')

    expect(fetchMock.mock.calls[0]?.[0]).toBe('http://localhost:8080/management/pegin/configuration')
  })
})
