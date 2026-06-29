import {
  bootstrapDevEnvironment,
  getInitialData,
  replaceInitialDataPayload,
  resetInitialDataCache,
  useInitialData,
} from '@shared/utils/initial-data'
import { renderHook } from '@testing-library/react'
import { loggedInFixture, loggedOutFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

describe('useInitialData', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    resetInitialDataCache()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('returns typed payload from the initial-data script', () => {
    seedInitialData(loggedOutFixture)
    const { result } = renderHook(() => useInitialData())

    expect(result.current.loggedIn).toBe(false)
    expect(result.current.data.BaseUrl).toBe('http://localhost:8080')
  })

  it('parses the DOM only once across repeated hook calls', () => {
    seedInitialData(loggedOutFixture)
    const parseSpy = vi.spyOn(JSON, 'parse')

    renderHook(() => useInitialData())
    renderHook(() => useInitialData())
    getInitialData()

    expect(parseSpy).toHaveBeenCalledTimes(1)
  })

  it('accepts Go fixture JSON for logged-in and logged-out payloads', () => {
    seedInitialData(loggedInFixture)
    const loggedIn = getInitialData()
    expect(loggedIn.loggedIn).toBe(true)
    expect(loggedIn.data.BtcAddress).toBe('tb1qloggedin')

    seedInitialData(loggedOutFixture)
    const loggedOut = getInitialData()
    expect(loggedOut.loggedIn).toBe(false)
    expect(loggedOut.data.CredentialsSet).toBe(true)
  })

  it('throws when the initial-data script is missing', () => {
    expect(() => getInitialData()).toThrow(/initial-data script element missing or empty/)
  })
})

describe('replaceInitialDataPayload', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    resetInitialDataCache()
  })

  it('throws when the initial-data script element is missing', () => {
    expect(() => {
      replaceInitialDataPayload(loggedOutFixture)
    }).toThrow(/initial-data script element missing/)
  })
})

describe('bootstrapDevEnvironment', () => {
  beforeEach(() => {
    document.head.innerHTML = ''
    document.body.innerHTML = ''
    resetInitialDataCache()
    vi.stubGlobal('fetch', vi.fn())
    vi.spyOn(console, 'warn').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('loads csrf and initial data from the LPS dev shell', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'old-csrf' })
    vi.mocked(fetch).mockResolvedValue(
      Response.json({ csrf: 'new-csrf', initialData: loggedInFixture }),
    )

    await bootstrapDevEnvironment()

    expect(document.querySelector('meta[name="csrf-token"]')?.getAttribute('content')).toBe(
      'new-csrf',
    )
    resetInitialDataCache()
    expect(getInitialData().loggedIn).toBe(true)
  })

  it('warns and keeps Vite stubs when the LPS dev shell is unavailable', async () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify({ error: 'LPS unavailable' }), {
        status: 503,
        headers: { 'Content-Type': 'application/json' },
      }),
    )

    await bootstrapDevEnvironment()

    expect(console.warn).toHaveBeenCalledWith(
      expect.stringContaining('LPS dev bootstrap skipped'),
      'LPS unavailable',
    )
  })

  it('warns when the LPS dev shell request fails', async () => {
    const networkError = new Error('network')
    vi.mocked(fetch).mockRejectedValue(networkError)

    await bootstrapDevEnvironment()

    expect(console.warn).toHaveBeenCalledWith(
      expect.stringContaining('LPS dev bootstrap skipped'),
      networkError,
    )
  })
})
