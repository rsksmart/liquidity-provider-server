import { ManagementPage } from '@feature/management/components'
import { render, screen, waitFor } from '@testing-library/react'
import { loggedInFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { apiFetchMock } = vi.hoisted(() => ({
  apiFetchMock: vi.fn(),
}))

vi.mock('@api/management/utils/api-fetch', () => ({
  apiFetch: apiFetchMock,
}))

vi.mock('@feature/auth/logout', () => ({
  logout: vi.fn(),
}))

vi.mock('sonner', () => ({
  toast: {
    error: vi.fn(),
    success: vi.fn(),
  },
}))

describe('ManagementPage', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
    apiFetchMock.mockImplementation((input: string) => {
      if (
        typeof input === 'string' &&
        input.includes('/management/trusted-accounts')
      ) {
        return Promise.resolve(
          new Response(JSON.stringify({ accounts: [] }), {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          }),
        )
      }
      return Promise.resolve(
        new Response(JSON.stringify({ collateral: '0' }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
    })
  })

  it('mounts provider, collateral, and trusted accounts cards with logout button', async () => {
    render(<ManagementPage />)

    expect(
      screen.getByRole('heading', { name: 'Management Dashboard' }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('provider-rsk-address')).toBeInTheDocument()
    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })
    expect(
      await screen.findByTestId('trusted-accounts-card'),
    ).toBeInTheDocument()
    expect(screen.getByText('Trusted Accounts')).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: /add account/i }),
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Logout' })).toBeInTheDocument()
  })
})
