import { ApiFetchError, CsrfTokenMissingError } from '@api/management/types/errors'
import {
  apiFetch,
  resetSessionExpiredHandling,
} from '@api/management/utils/api-fetch'
import { resetInitialDataCache } from '@shared/utils/initial-data'
import { loggedOutFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

describe('apiFetch', () => {
  beforeEach(() => {
    document.head.innerHTML = ''
    document.body.innerHTML = ''
    resetInitialDataCache()
    resetSessionExpiredHandling()
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

  it('sets Content-Type application/json on mutating requests with a body', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch('/management/login', {
      json: { username: 'u', password: 'p' },
    })

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(new Headers(init.headers).get('Content-Type')).toBe('application/json')
  })

  it('defaults method to POST when json is provided', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch('/management/login', { json: { username: 'u', password: 'p' } })

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(init.method).toBe('POST')
    expect(init.body).toBe(JSON.stringify({ username: 'u', password: 'p' }))
  })

  it('does not override an explicit Content-Type header', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch('/management/login', {
      method: 'POST',
      headers: { 'Content-Type': 'text/plain' },
      json: { username: 'u' },
    })

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(new Headers(init.headers).get('Content-Type')).toBe('text/plain')
    expect(init.body).toBe(JSON.stringify({ username: 'u' }))
  })

  it('omits Content-Type on body-less POST requests', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch.post('/management/logout')

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(new Headers(init.headers).has('Content-Type')).toBe(false)
  })

  it('treats the second logout argument as RequestInit, not JSON body', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))
    const controller = new AbortController()

    await apiFetch.post('/management/logout', { signal: controller.signal })

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(init.signal).toBe(controller.signal)
    expect(init.body).toBeUndefined()
  })

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

  it('falls back to raw text when JSON content-type body is not valid JSON', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response('not-json', {
        status: 500,
        statusText: 'Internal Server Error',
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(apiFetch('/management/login', { method: 'POST' })).rejects.toMatchObject({
      body: 'not-json',
    } satisfies Partial<ApiFetchError>)
  })

  it('returns an empty body as-is for non-JSON error responses', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response('', {
        status: 500,
        statusText: 'Internal Server Error',
      }),
    )

    await expect(apiFetch('/management/login', { method: 'POST' })).rejects.toMatchObject({
      body: '',
    } satisfies Partial<ApiFetchError>)
  })

  it('does not treat a 403 without a session message as session expiry', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ error: 'forbidden' }), {
        status: 403,
        statusText: 'Forbidden',
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(apiFetch('/management/status', { method: 'GET' })).rejects.toMatchObject({
      status: 403,
    } satisfies Partial<ApiFetchError>)
  })

  it('treats 403 session validation errors as session expiry', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ message: 'session validation error' }), {
        status: 403,
        statusText: 'Forbidden',
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(apiFetch('/management/status', { method: 'GET' })).rejects.toBeInstanceOf(
      ApiFetchError,
    )
  })

  it('resolves relative URLs against BaseUrl from initial data', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch('/management/pegin/configuration')

    expect(fetchMock.mock.calls[0]?.[0]).toBe('http://localhost:8080/management/pegin/configuration')
  })

  it('leaves absolute http(s) URLs unchanged', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch('https://example.com/management/status')

    expect(fetchMock.mock.calls[0]?.[0]).toBe('https://example.com/management/status')
  })

  it('preserves the response body on successful requests', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ collateral: '1000000000000000000' }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    const response = await apiFetch('/pegin/collateral')
    await expect(response.json()).resolves.toEqual({ collateral: '1000000000000000000' })
  })

  it('posts JSON via apiFetch.post', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch.post('/pegin/addCollateral', { amount: 1_000_000_000_000_000_000 })

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(init.method).toBe('POST')
    expect(init.body).toBe(JSON.stringify({ amount: 1_000_000_000_000_000_000 }))
  })

  it('fetches via apiFetch.get', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch.get('/pegin/collateral')

    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(init.method).toBe('GET')
  })

  it('resolves relative URLs against window origin when BaseUrl is empty', async () => {
    seedInitialData(
      { ...loggedOutFixture, data: { ...loggedOutFixture.data, BaseUrl: '' } },
      { csrfToken: 'csrf-token' },
    )
    const fetchMock = vi.mocked(fetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await apiFetch('/pegin/collateral')

    expect(fetchMock.mock.calls[0]?.[0]).toBe(`${window.location.origin}/pegin/collateral`)
  })

  it('classifies session expiry for absolute URLs on error', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ message: 'session not recognized' }), {
        status: 403,
        statusText: 'Forbidden',
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(
      apiFetch('https://example.com/management/status', { method: 'GET' }),
    ).rejects.toBeInstanceOf(ApiFetchError)
  })
})
