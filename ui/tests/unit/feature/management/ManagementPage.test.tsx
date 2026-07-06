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

describe('ManagementPage', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
    apiFetchMock.mockResolvedValue(
      new Response(JSON.stringify({ collateral: '0' }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    )
  })

  it('mounts provider and collateral cards with logout button', async () => {
    render(<ManagementPage />)

    expect(screen.getByRole('heading', { name: 'Management Dashboard' })).toBeInTheDocument()
    expect(screen.getByTestId('provider-rsk-address')).toBeInTheDocument()
    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toBeInTheDocument()
    })
    expect(screen.getByRole('button', { name: 'Logout' })).toBeInTheDocument()
  })
})
