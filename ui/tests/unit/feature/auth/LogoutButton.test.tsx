import { apiFetch } from '@api/management/utils/api-fetch'
import { LogoutButton } from '@feature/auth/components/LogoutButton'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

const { toastWarningMock } = vi.hoisted(() => ({
  toastWarningMock: vi.fn(),
}))

vi.mock('@api/management/utils/api-fetch', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@api/management/utils/api-fetch')>()
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

describe('LogoutButton', () => {
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

  // success path covered by Playwright logout E2E
  it('shows warning toast and still redirects when logout fails', async () => {
    vi.mocked(apiFetch).mockRejectedValue(new Error('network'))
    const user = userEvent.setup()

    render(<LogoutButton />)
    await user.click(screen.getByRole('button', { name: 'Logout' }))

    await waitFor(() => {
      expect(toastWarningMock).toHaveBeenCalledWith('Logout failed; local session cleared.')
      expect(assignMock).toHaveBeenCalledWith('/management/next/login')
    })
    expect(console.error).toHaveBeenCalled()
  })
})
