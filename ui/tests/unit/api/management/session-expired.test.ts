import { ApiFetchError } from '@api/management/types/errors'
import { apiFetch } from '@api/management/utils/api-fetch'
import { resetInitialDataCache } from '@shared/utils/initial-data'
import { loggedOutFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

const { toastErrorMock } = vi.hoisted(() => ({
  toastErrorMock: vi.fn(),
}))

vi.mock('sonner', () => ({
  toast: {
    error: toastErrorMock,
    warning: vi.fn(),
  },
}))

describe('apiFetch session expiry', () => {
  const assignMock = vi.fn()

  beforeEach(() => {
    document.head.innerHTML = ''
    document.body.innerHTML = ''
    resetInitialDataCache()
    vi.stubGlobal('fetch', vi.fn())
    toastErrorMock.mockReset()
    assignMock.mockReset()
    vi.useFakeTimers()
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: { assign: assignMock },
    })
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('does not redirect on 401 for /management/login', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ error: 'unauthorized' }), {
        status: 401,
        statusText: 'Unauthorized',
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(
      apiFetch('/management/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: 'u', password: 'p' }),
      }),
    ).rejects.toBeInstanceOf(ApiFetchError)

    expect(toastErrorMock).not.toHaveBeenCalled()
    expect(assignMock).not.toHaveBeenCalled()
  })

  it('shows toast and redirects on 401 for authenticated endpoints', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response('unauthorized', {
        status: 401,
        statusText: 'Unauthorized',
      }),
    )

    await expect(apiFetch('/management/status', { method: 'GET' })).rejects.toBeInstanceOf(
      ApiFetchError,
    )

    expect(toastErrorMock).toHaveBeenCalledWith('Your session has expired. Please log in again.')
    expect(assignMock).not.toHaveBeenCalled()
    vi.runAllTimers()
    expect(assignMock).toHaveBeenCalledWith('/management/next/login')
  })

  it('shows toast and redirects on 403 session not recognized', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ message: 'session not recognized' }), {
        status: 403,
        statusText: 'Forbidden',
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(apiFetch('/management/status', { method: 'GET' })).rejects.toBeInstanceOf(
      ApiFetchError,
    )

    expect(toastErrorMock).toHaveBeenCalledWith('Your session has expired. Please log in again.')
    expect(assignMock).not.toHaveBeenCalled()
    vi.runAllTimers()
    expect(assignMock).toHaveBeenCalledWith('/management/next/login')
  })

  it('does not redirect on 403 CSRF errors', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ message: 'CRSF token validation error' }), {
        status: 403,
        statusText: 'Forbidden',
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await expect(apiFetch('/management/status', { method: 'GET' })).rejects.toBeInstanceOf(
      ApiFetchError,
    )

    expect(toastErrorMock).not.toHaveBeenCalled()
    expect(assignMock).not.toHaveBeenCalled()
  })
})
