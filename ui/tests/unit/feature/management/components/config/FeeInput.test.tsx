import { FeeInput } from '@feature/management/components/config'
import { getDisplayLabel, getTooltipText } from '@feature/management/config/labels'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it, vi } from 'vitest'

/** Test-only parent: keeps the controlled `value` in sync the way a form would. */
function FeeInputHarness({
  fieldKey,
  initialValue = '',
  onChangeSpy,
}: {
  fieldKey: string
  initialValue?: string
  onChangeSpy?: (next: string) => void
}) {
  const [value, setValue] = useState(initialValue)
  return (
    <FeeInput
      sectionPrefix="pegin"
      fieldKey={fieldKey}
      value={value}
      onChange={(next) => {
        setValue(next)
        onChangeSpy?.(next)
      }}
    />
  )
}

describe('FeeInput', () => {
  it('renders a 40%-wide text input with expected testid', () => {
    render(<FeeInputHarness fieldKey="penaltyFee" initialValue="0.5" />)

    const input = screen.getByTestId('config-pegin-penaltyFee-input')
    expect(input).toHaveValue('0.5')
    expect(input).toHaveAttribute('type', 'text')
    expect(input).toHaveClass('w-[40%]')
  })

  it('shows the display label and tooltip copy for fee keys', () => {
    render(<FeeInputHarness fieldKey="callFee" />)

    // callFee has no friendly label → falls back to the key
    expect(screen.getByText(getDisplayLabel('callFee'))).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: getTooltipText('callFee') }),
    ).toBeInTheDocument()
  })

  it('renders maxLiquidity bold and without a tooltip icon', () => {
    render(<FeeInputHarness fieldKey="maxLiquidity" initialValue="1" />)

    expect(screen.getByText('Maximum Liquidity')).toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: getTooltipText('maxLiquidity') }),
    ).not.toBeInTheDocument()
  })

  it('surfaces maxLiquidity validation errors on change', async () => {
    const user = userEvent.setup()
    render(<FeeInputHarness fieldKey="maxLiquidity" initialValue="" />)

    const input = screen.getByTestId('config-pegin-maxLiquidity-input')

    await user.type(input, 'abc')
    expect(screen.getByRole('alert')).toHaveTextContent(
      'Max liquidity must be a valid number',
    )

    await user.clear(input)
    await user.type(input, '-5')
    expect(screen.getByRole('alert')).toHaveTextContent(
      'Max liquidity must be a positive number',
    )

    await user.clear(input)
    await user.type(input, '1.5')
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })

  it('does not validate non-maxLiquidity fee keys internally', async () => {
    const user = userEvent.setup()
    render(<FeeInputHarness fieldKey="penaltyFee" initialValue="" />)

    await user.type(screen.getByTestId('config-pegin-penaltyFee-input'), 'abc')
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })

  it('emits typed values through onChange', async () => {
    const user = userEvent.setup()
    const onChangeSpy = vi.fn()
    render(<FeeInputHarness fieldKey="minValue" onChangeSpy={onChangeSpy} />)

    await user.type(screen.getByTestId('config-pegin-minValue-input'), '7')
    expect(onChangeSpy).toHaveBeenLastCalledWith('7')
  })
})
