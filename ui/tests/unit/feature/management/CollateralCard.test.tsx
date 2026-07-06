import { CollateralCard } from '@feature/management/components'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { loggedInFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const { apiFetchMock } = vi.hoisted(() => ({
  apiFetchMock: vi.fn(),
}))

vi.mock('@api/management/utils/api-fetch', () => ({
  apiFetch: apiFetchMock,
}))

function mockCollateralBalance(wei = '1000000000000000000') {
  apiFetchMock.mockImplementation(() =>
    Promise.resolve(
      new Response(JSON.stringify({ collateral: wei }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    ),
  )
}

describe('CollateralCard', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedInitialData(loggedInFixture, { csrfToken: 'csrf-token' })
    apiFetchMock.mockReset()
    mockCollateralBalance()
  })

  it('renders pegin tab by default with balance', async () => {
    render(<CollateralCard />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toHaveTextContent('1 rBTC')
    })
    expect(screen.getByRole('tab', { name: 'Pegin' })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: 'Pegout' })).toBeInTheDocument()
  })

  it('preserves amount input when switching tabs', async () => {
    const user = userEvent.setup()
    render(<CollateralCard />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toHaveTextContent('1 rBTC')
    })

    await user.type(screen.getByTestId('pegin-collateral-amount'), '1.5')
    await user.click(screen.getByRole('tab', { name: 'Pegout' }))
    await user.click(screen.getByRole('tab', { name: 'Pegin' }))

    expect(screen.getByTestId('pegin-collateral-amount')).toHaveValue(1.5)
  })

  it('switches to pegout tab via keyboard', async () => {
    const user = userEvent.setup()
    render(<CollateralCard />)

    await waitFor(() => {
      expect(screen.getByTestId('pegin-collateral-balance')).toHaveTextContent('1 rBTC')
    })

    await user.tab()
    await user.keyboard('{ArrowRight}')
    await user.keyboard('{Enter}')

    await waitFor(() => {
      expect(screen.getByTestId('pegout-collateral-balance')).toHaveTextContent('1 rBTC')
    })
  })
})
