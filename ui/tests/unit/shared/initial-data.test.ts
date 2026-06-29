import {
  getInitialData,
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
