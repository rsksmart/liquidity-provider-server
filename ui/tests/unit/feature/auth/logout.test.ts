import { ApiFetchError } from '@api/management/types/errors'
import { apiFetch } from '@api/management/utils/api-fetch'
import { logout } from '@feature/auth/logout'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

const { toastWarningMock } = vi.hoisted(() => ({
  toastWarningMock: vi.fn(),
}))

vi.mock('@api/management/utils/api-fetch', async (importOriginal) => {
  const actual =
    await importOriginal<typeof import('@api/management/utils/api-fetch')>()
  return {
    ...actual,
    apiFetch: vi.fn(),
  }
})

vi.mock('sonner', () => ({
  toast: {
    error: vi.fn(),
    warning: toastWarningMock,
  },
}))

describe('logout', () => {
  const assignMock = vi.fn()

  beforeEach(() => {
    vi.mocked(apiFetch).mockReset()
    toastWarningMock.mockReset()
    vi.spyOn(console, 'error').mockImplementation(() => {})
    assignMock.mockReset()
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: { assign: assignMock },
    })
  })

  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('redirects to login after a successful logout', async () => {
    vi.mocked(apiFetch).mockResolvedValue(new Response('ok', { status: 200 }))

    await logout()

    expect(assignMock).toHaveBeenCalledWith('/management/next/login')
  })

  it('skips redirect when the session is already expired (403)', async () => {
    vi.mocked(apiFetch).mockRejectedValue(
      new ApiFetchError(403, 'Forbidden', {
        message: 'session not recognized',
      }),
    )

    await logout()

    expect(toastWarningMock).not.toHaveBeenCalled()
    expect(assignMock).not.toHaveBeenCalled()
  })

  it('skips redirect when the session is already expired (401)', async () => {
    vi.mocked(apiFetch).mockRejectedValue(
      new ApiFetchError(401, 'Unauthorized', ''),
    )

    await logout()

    expect(toastWarningMock).not.toHaveBeenCalled()
    expect(assignMock).not.toHaveBeenCalled()
  })

  it('shows warning toast and still redirects when logout fails', async () => {
    vi.mocked(apiFetch).mockRejectedValue(new Error('network'))

    await logout()

    expect(toastWarningMock).toHaveBeenCalledWith(
      'Logout failed; local session cleared.',
    )
    expect(assignMock).toHaveBeenCalledWith('/management/next/login')
    expect(console.error).toHaveBeenCalled()
  })
})
