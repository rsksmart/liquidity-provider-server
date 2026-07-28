import { FeePercentageInput, normalizePercentage } from '@feature/management/components/config'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it } from 'vitest'

/** Test-only parent: keeps the controlled `value` in sync the way a form would. */
function FeePercentageInputHarness({ initialValue = '' }: { initialValue?: string }) {
  const [value, setValue] = useState(initialValue)
  return (
    <FeePercentageInput
      sectionPrefix="pegout"
      fieldKey="feePercentage"
      value={value}
      onChange={setValue}
    />
  )
}

describe('normalizePercentage', () => {
  it('strips a trailing percent sign', () => {
    expect(normalizePercentage('50%')).toEqual({ value: '50', error: null })
  })

  it('clamps values into the 0–100 range', () => {
    expect(normalizePercentage('150')).toEqual({ value: '100', error: null })
    expect(normalizePercentage('3.5')).toEqual({ value: '3.5', error: null })
  })

  it('treats empty input as zero', () => {
    expect(normalizePercentage('')).toEqual({ value: '0', error: null })
  })

  it('reports an error for non-numeric input', () => {
    const result = normalizePercentage('abc')
    expect(result.error).toMatch(/between 0% and 100%/)
    expect(result.value).toBe('abc')
  })
})

describe('FeePercentageInput', () => {
  it('renders a 40%-wide text input with expected testid', () => {
    render(<FeePercentageInputHarness initialValue="5" />)

    const input = screen.getByTestId('config-pegout-feePercentage-input')
    expect(input).toHaveValue('5')
    expect(input).toHaveClass('w-[40%]')
  })

  it('normalizes a trailing percent on blur', async () => {
    const user = userEvent.setup()
    render(<FeePercentageInputHarness />)

    const input = screen.getByTestId('config-pegout-feePercentage-input')
    await user.type(input, '50%')
    await user.tab()

    expect(input).toHaveValue('50')
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  })

  it('clamps out-of-range values to 100 on blur', async () => {
    const user = userEvent.setup()
    render(<FeePercentageInputHarness />)

    const input = screen.getByTestId('config-pegout-feePercentage-input')
    await user.type(input, '150')
    await user.tab()

    expect(input).toHaveValue('100')
  })

  it('shows an error for invalid percentage input', async () => {
    const user = userEvent.setup()
    render(<FeePercentageInputHarness />)

    const input = screen.getByTestId('config-pegout-feePercentage-input')
    await user.type(input, 'abc')
    await user.tab()

    expect(screen.getByRole('alert')).toHaveTextContent('between 0% and 100%')
  })
})
