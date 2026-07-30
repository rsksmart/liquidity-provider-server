import { ToggableFeeInput } from '@feature/management/components/config'
import { getTooltipText } from '@feature/management/config/labels'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it, vi } from 'vitest'

/** Test-only parent: keeps controlled props in sync the way a form would. */
function ToggableFeeInputHarness({
  initialValue,
  initialEnabled,
  onFeeToggleChange,
}: {
  initialValue: string
  initialEnabled: boolean
  onFeeToggleChange?: (enabled: boolean) => void
}) {
  const [value, setValue] = useState(initialValue)
  const [enabled, setEnabled] = useState(initialEnabled)
  return (
    <ToggableFeeInput
      sectionPrefix="pegin"
      fieldKey="fixedFee"
      value={value}
      onChange={setValue}
      enabled={enabled}
      onEnabledChange={setEnabled}
      onFeeToggleChange={onFeeToggleChange}
    />
  )
}

describe('ToggableFeeInput', () => {
  it('renders a checkbox + 40%-wide input with expected testids and tooltip', () => {
    render(<ToggableFeeInputHarness initialValue="0.5" initialEnabled />)

    const checkbox = screen.getByTestId('config-pegin-fixedFee-checkbox')
    const input = screen.getByTestId('config-pegin-fixedFee-input')
    expect(checkbox).toBeChecked()
    expect(input).toBeEnabled()
    expect(input).toHaveValue('0.5')
    expect(input).toHaveClass('w-[40%]')
    expect(
      screen.getByRole('button', { name: getTooltipText('fixedFee') }),
    ).toBeInTheDocument()
  })

  it('starts disabled showing 0 when value is 0', () => {
    render(<ToggableFeeInputHarness initialValue="0" initialEnabled={false} />)

    expect(screen.getByTestId('config-pegin-fixedFee-checkbox')).not.toBeChecked()
    const input = screen.getByTestId('config-pegin-fixedFee-input')
    expect(input).toBeDisabled()
    expect(input).toHaveValue('0')
  })

  it('unchecking forces the value to 0 and disables the input', async () => {
    const user = userEvent.setup()
    render(<ToggableFeeInputHarness initialValue="0.5" initialEnabled />)

    await user.click(screen.getByTestId('config-pegin-fixedFee-checkbox'))

    const input = screen.getByTestId('config-pegin-fixedFee-input')
    expect(input).toBeDisabled()
    expect(input).toHaveValue('0')
  })

  it('re-checking restores the original value and re-enables the input', async () => {
    const user = userEvent.setup()
    render(<ToggableFeeInputHarness initialValue="0.5" initialEnabled />)

    const checkbox = screen.getByTestId('config-pegin-fixedFee-checkbox')
    await user.click(checkbox) // off → value becomes 0
    await user.click(checkbox) // on → restore original

    const input = screen.getByTestId('config-pegin-fixedFee-input')
    expect(input).toBeEnabled()
    expect(input).toHaveValue('0.5')
  })

  it('re-checking a field whose original was 0 restores an empty input', async () => {
    const user = userEvent.setup()
    render(<ToggableFeeInputHarness initialValue="0" initialEnabled={false} />)

    await user.click(screen.getByTestId('config-pegin-fixedFee-checkbox'))

    const input = screen.getByTestId('config-pegin-fixedFee-input')
    expect(input).toBeEnabled()
    expect(input).toHaveValue('')
  })

  it('invokes onFeeToggleChange with the new enabled state', async () => {
    const user = userEvent.setup()
    const onFeeToggleChange = vi.fn()
    render(
      <ToggableFeeInputHarness
        initialValue="0.5"
        initialEnabled
        onFeeToggleChange={onFeeToggleChange}
      />,
    )

    await user.click(screen.getByTestId('config-pegin-fixedFee-checkbox'))
    expect(onFeeToggleChange).toHaveBeenLastCalledWith(false)
  })

  it('propagates typed input values through onChange', async () => {
    const user = userEvent.setup()
    const onDirty = vi.fn()
    function Harness() {
      const [value, setValue] = useState('0.5')
      const [enabled, setEnabled] = useState(true)
      return (
        <ToggableFeeInput
          sectionPrefix="pegin"
          fieldKey="fixedFee"
          value={value}
          onChange={setValue}
          enabled={enabled}
          onEnabledChange={setEnabled}
          onDirty={onDirty}
        />
      )
    }
    render(<Harness />)

    const input = screen.getByTestId('config-pegin-fixedFee-input')
    await user.clear(input)
    await user.type(input, '1.25')

    expect(input).toHaveValue('1.25')
    expect(onDirty).toHaveBeenCalled()
  })
})
