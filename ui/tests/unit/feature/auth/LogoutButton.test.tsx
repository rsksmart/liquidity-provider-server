import { LogoutButton } from '@feature/auth/components/LogoutButton'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

const { logoutMock } = vi.hoisted(() => ({
  logoutMock: vi.fn(),
}))

vi.mock('@feature/auth/logout', () => ({
  logout: logoutMock,
}))

describe('LogoutButton', () => {
  it('calls logout when clicked', async () => {
    logoutMock.mockResolvedValue(undefined)
    const user = userEvent.setup()

    render(<LogoutButton />)
    await user.click(screen.getByRole('button', { name: 'Logout' }))

    expect(logoutMock).toHaveBeenCalledTimes(1)
  })
})
