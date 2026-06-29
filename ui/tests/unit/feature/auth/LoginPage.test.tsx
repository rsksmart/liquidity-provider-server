import { ApiFetchError, CsrfTokenMissingError } from '@api/management/types/errors'
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
  const assignMock = vi.fn()

  beforeEach(() => {
    document.head.innerHTML = ''
    document.body.innerHTML = ''
    resetInitialDataCache()
    vi.mocked(apiFetch).mockReset()
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

  it('redirects to management after a successful login', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(apiFetch).mockResolvedValue(new Response('ok', { status: 200 }))
    const user = userEvent.setup()

    render(<LoginPage />)
    await user.type(screen.getByTestId('login-username-input'), 'user')
    await user.type(screen.getByTestId('login-password-input'), 'pass')
    await user.click(screen.getByTestId('login-submit-button'))

    await waitFor(() => {
      expect(assignMock).toHaveBeenCalledWith('/management/next/management')
    })
  })

  it.each([
    {
      name: 'CSRF token missing',
      error: new CsrfTokenMissingError(),
      expectedMessage: 'Invalid username or password.',
    },
    {
      name: 'network failure',
      error: new TypeError('Failed to fetch'),
      expectedMessage: 'Invalid username or password.',
    },
    {
      name: '401 response',
      error: new ApiFetchError(401, 'Unauthorized', { error: 'unauthorized' }),
      expectedMessage: 'Invalid username or password.',
    },
    {
      name: '403 without a server message',
      error: new ApiFetchError(403, 'Forbidden', { error: 'forbidden' }),
      expectedMessage: 'Invalid username or password.',
    },
    {
      name: 'unexpected error',
      error: new Error('boom'),
      expectedMessage: 'Invalid username or password.',
    },
  ])('shows "$expectedMessage" for $name', async ({ error, expectedMessage }) => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(apiFetch).mockRejectedValue(error)
    const user = userEvent.setup()

    render(<LoginPage />)
    await user.type(screen.getByTestId('login-username-input'), 'user')
    await user.type(screen.getByTestId('login-password-input'), 'pass')
    await user.click(screen.getByTestId('login-submit-button'))

    expect(await screen.findByRole('alert')).toHaveTextContent(expectedMessage)
  })

  it('shows the server message for a 403 CSRF rejection', async () => {
    seedInitialData(loggedOutFixture, { csrfToken: 'csrf-token' })
    vi.mocked(apiFetch).mockRejectedValue(
      new ApiFetchError(403, 'Forbidden', { message: 'CSRF token validation error' }),
    )
    const user = userEvent.setup()

    render(<LoginPage />)
    await user.type(screen.getByTestId('login-username-input'), 'user')
    await user.type(screen.getByTestId('login-password-input'), 'pass')
    await user.click(screen.getByTestId('login-submit-button'))

    expect(await screen.findByRole('alert')).toHaveTextContent('CSRF token validation error')
  })

  // login happy path covered above; credentials-set chain below
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
