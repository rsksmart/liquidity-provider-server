import { DefaultInput } from '@feature/management/components/config'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it } from 'vitest'

/** Test-only parent: keeps the controlled `value` in sync the way a form would. */
function DefaultInputHarness({ initialValue = '' }: { initialValue?: string }) {
  const [value, setValue] = useState(initialValue)
  return (
    <DefaultInput
      sectionPrefix="pegout"
      fieldKey="expireBlocks"
      value={value}
      onChange={setValue}
    />
  )
}

describe('DefaultInput', () => {
  it('renders a text input with the expected testid and no tooltip icon', () => {
    render(<DefaultInputHarness initialValue="500" />)

    const input = screen.getByTestId('config-pegout-expireBlocks-input')
    expect(input).toHaveAttribute('type', 'text')
    expect(input).toHaveValue('500')
    expect(screen.queryByRole('button')).not.toBeInTheDocument()
  })

  it('emits typed values through onChange', async () => {
    const user = userEvent.setup()
    render(<DefaultInputHarness />)

    const input = screen.getByTestId('config-pegout-expireBlocks-input')
    await user.type(input, '42')
    expect(input).toHaveValue('42')
  })
})
