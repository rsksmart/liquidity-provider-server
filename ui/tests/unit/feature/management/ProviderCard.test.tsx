import { ProviderCard } from '@feature/management/components'
import { render, screen } from '@testing-library/react'
import { loggedInFixture } from '@tests/fixtures'
import { seedInitialData } from '@tests/utils'
import { beforeEach, describe, expect, it } from 'vitest'

describe('ProviderCard', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    seedInitialData(loggedInFixture)
  })

  it('renders provider addresses and operational status from initial data', () => {
    render(<ProviderCard />)

    expect(screen.getByTestId('provider-rsk-address')).toHaveTextContent('0xloggedin')
    expect(screen.getByTestId('provider-btc-address')).toHaveTextContent('tb1qloggedin')
    expect(screen.getByTestId('provider-operational-status')).toHaveTextContent('Operational')
  })

  it('shows not operational when provider status is false', () => {
    seedInitialData({
      ...loggedInFixture,
      data: {
        ...loggedInFixture.data,
        ProviderData: {
          ...loggedInFixture.data.ProviderData,
          status: false,
        },
      },
    })

    render(<ProviderCard />)
    expect(screen.getByTestId('provider-operational-status')).toHaveTextContent('Not Operational')
  })
})
