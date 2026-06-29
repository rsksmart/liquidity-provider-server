import { apiFetch } from '@api/management/utils/api-fetch'
import { LoginPage } from '@feature/auth/components/LoginPage'
import { resetInitialDataCache } from '@shared/utils/initial-data'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { loggedOutFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@api/management/utils/api-fetch', () => ({
  apiFetch: vi.fn(),
}))

const credentialsNotSetFixture = {
  ...loggedOutFixture,
  data: {
    ...loggedOutFixture.data,
    CredentialsSet: false,
  },
}

describe('LoginPage', () => {
  beforeEach(() => {
    document.head.innerHTML = ''
    document.body.innerHTML = ''
    resetInitialDataCache()
    vi.mocked(apiFetch).mockReset()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  // login errors + happy path covered by Playwright; add Vitest rows if E2E gaps appear
  it('chains credentials POST when CredentialsSet is false', async () => {
    seedInitialData(credentialsNotSetFixture, { csrfToken: 'csrf-token' })
    vi.mocked(apiFetch)
      .mockResolvedValueOnce(new Response('ok', { status: 200 }))
      .mockResolvedValueOnce(new Response('ok', { status: 200 }))
    const user = userEvent.setup()

    render(<LoginPage />)

    expect(screen.getByLabelText('New Username')).toBeInTheDocument()

    await user.type(screen.getByTestId('login-username-input'), 'old-user')
    await user.type(screen.getByTestId('login-password-input'), 'old-pass')
    await user.type(screen.getByLabelText('New Username'), 'new-user')
    await user.type(screen.getByLabelText('New Password'), 'new-pass')
    await user.click(screen.getByTestId('login-submit-button'))

    await waitFor(() => {
      expect(apiFetch).toHaveBeenNthCalledWith(
        2,
        '/management/credentials',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            oldUsername: 'old-user',
            oldPassword: 'old-pass',
            newUsername: 'new-user',
            newPassword: 'new-pass',
          }),
        }),
      )
    })
  })
})
