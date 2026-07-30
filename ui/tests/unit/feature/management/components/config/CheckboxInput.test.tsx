import { CheckboxInput } from '@feature/management/components/config'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it } from 'vitest'

/** Test-only parent: keeps the controlled `value` in sync the way a form would. */
function CheckboxInputHarness({ initial = false }: { initial?: boolean }) {
  const [value, setValue] = useState(initial)
  return (
    <CheckboxInput
      sectionPrefix="general"
      fieldKey="publicLiquidityCheck"
      value={value}
      onChange={setValue}
    />
  )
}

describe('CheckboxInput', () => {
  it('renders a form-check-input checkbox with the expected testid', () => {
    render(<CheckboxInputHarness initial />)

    const checkbox = screen.getByTestId('config-general-publicLiquidityCheck-checkbox')
    expect(checkbox).toHaveAttribute('type', 'checkbox')
    expect(checkbox).toHaveClass('form-check-input')
    expect(checkbox).toBeChecked()
  })

  it('toggles the boolean value on click', async () => {
    const user = userEvent.setup()
    render(<CheckboxInputHarness initial={false} />)

    const checkbox = screen.getByTestId('config-general-publicLiquidityCheck-checkbox')
    expect(checkbox).not.toBeChecked()

    await user.click(checkbox)
    expect(checkbox).toBeChecked()
  })
})
