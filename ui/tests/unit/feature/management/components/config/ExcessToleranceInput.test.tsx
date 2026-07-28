import { ExcessToleranceInput } from '@feature/management/components/config'
import { getTooltipText } from '@feature/management/config/labels'
import type { ExcessTolerance } from '@shared/types/initial-data'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it } from 'vitest'

const HALF_RBTC_WEI = '500000000000000000'

/** Test-only parent: keeps the controlled `value` in sync the way a form would. */
function ExcessToleranceInputHarness({ initial }: { initial: ExcessTolerance }) {
  const [value, setValue] = useState(initial)
  return (
    <ExcessToleranceInput sectionPrefix="general" value={value} onChange={setValue} />
  )
}

describe('ExcessToleranceInput', () => {
  it('renders the label, tooltip, switch and input with expected testids', () => {
    render(
      <ExcessToleranceInputHarness
        initial={{ isFixed: false, fixedValue: '0', percentageValue: '30' }}
      />,
    )

    expect(screen.getByText('Excess Tolerance')).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: getTooltipText('excessTolerance') }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('config-general-excessTolerance-toggle')).toBeInTheDocument()

    const input = screen.getByTestId('config-general-excessTolerance-input')
    expect(input).toHaveClass('w-[40%]')
    expect(input).toHaveValue('30')
    expect(input).toHaveAttribute('placeholder', 'Enter percentage (0-100)')
  })

  it('displays the fixed rBTC amount (ether) and placeholder in fixed mode', () => {
    render(
      <ExcessToleranceInputHarness
        initial={{ isFixed: true, fixedValue: HALF_RBTC_WEI, percentageValue: '0' }}
      />,
    )

    const input = screen.getByTestId('config-general-excessTolerance-input')
    expect(input).toHaveValue('0.5')
    expect(input).toHaveAttribute('placeholder', 'Enter amount in rBTC')
  })

  it('preserves separate fixed/percentage working values when toggling', async () => {
    const user = userEvent.setup()
    render(
      <ExcessToleranceInputHarness
        initial={{ isFixed: true, fixedValue: HALF_RBTC_WEI, percentageValue: '0' }}
      />,
    )

    const toggle = screen.getByTestId('config-general-excessTolerance-toggle')
    const input = screen.getByTestId('config-general-excessTolerance-input')

    // Switch to percentage mode and enter a value
    await user.click(toggle)
    expect(input).toHaveValue('0')
    await user.clear(input)
    await user.type(input, '25')
    expect(input).toHaveValue('25')

    // Back to fixed mode — original fixed display is preserved
    await user.click(toggle)
    expect(input).toHaveValue('0.5')
    expect(input).toHaveAttribute('placeholder', 'Enter amount in rBTC')

    // Back to percentage mode — typed percentage is preserved
    await user.click(toggle)
    expect(input).toHaveValue('25')
  })

  it('validates fixed-mode input as a number', async () => {
    const user = userEvent.setup()
    render(
      <ExcessToleranceInputHarness
        initial={{ isFixed: true, fixedValue: '0', percentageValue: '0' }}
      />,
    )

    const input = screen.getByTestId('config-general-excessTolerance-input')
    await user.clear(input)
    await user.type(input, 'abc')

    expect(screen.getByRole('alert')).toHaveTextContent(
      'Excess tolerance fixed must be a valid number',
    )
  })

  it('validates percentage-mode input is at most 100', async () => {
    const user = userEvent.setup()
    render(
      <ExcessToleranceInputHarness
        initial={{ isFixed: false, fixedValue: '0', percentageValue: '0' }}
      />,
    )

    const input = screen.getByTestId('config-general-excessTolerance-input')
    await user.clear(input)
    await user.type(input, '150')

    expect(screen.getByRole('alert')).toHaveTextContent(
      'Excess tolerance percentage cannot exceed 100%',
    )
  })
})
