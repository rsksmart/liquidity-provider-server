import { managementConfigFieldInputClass } from '@feature/management/management-styles'
import { type ChangeEvent, useCallback, useRef } from 'react'

import { Input } from '@/components/ui/input'

import { FieldError } from '../FieldError'
import { ConfigCheckbox } from './ConfigCheckbox'
import { ConfigFieldLabel } from './ConfigFieldLabel'
import type { SectionPrefix } from './types'

interface ToggableFeeInputProps {
  sectionPrefix: SectionPrefix
  /** Toggable fee key: `fixedFee` or `feePercentage`. */
  fieldKey: string
  /** Display value (ether string for fees, percentage string for `feePercentage`). */
  value: string
  onChange: (next: string) => void
  /** Whether the fee is enabled (checkbox checked). */
  enabled: boolean
  onEnabledChange: (enabled: boolean) => void
  /** Called with the new enabled state on toggle (for zero-fee warnings). */
  onFeeToggleChange?: (enabled: boolean) => void
  onDirty?: () => void
  error?: string | null
}

/**
 * Checkbox-gated fee field.
 * Unchecking forces the value to `0` and disables the input; re-checking
 * restores the value captured when the component mounted (empty when that was `0`).
 */
export function ToggableFeeInput({
  sectionPrefix,
  fieldKey,
  value,
  onChange,
  enabled,
  onEnabledChange,
  onFeeToggleChange,
  onDirty,
  error,
}: ToggableFeeInputProps) {
  const inputTestId = `config-${sectionPrefix}-${fieldKey}-input`
  const checkboxTestId = `config-${sectionPrefix}-${fieldKey}-checkbox`

  // Value to restore when the fee is re-enabled.
  const originalValueRef = useRef(value)

  const handleToggle = useCallback(
    (checked: boolean) => {
      onEnabledChange(checked)
      if (checked) {
        const original = originalValueRef.current
        onChange(original === '0' || original === '' ? '' : original)
      } else {
        onChange('0')
      }
      onFeeToggleChange?.(checked)
      onDirty?.()
    },
    [onChange, onEnabledChange, onFeeToggleChange, onDirty],
  )

  const handleInputChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      onChange(event.target.value)
      onDirty?.()
    },
    [onChange, onDirty],
  )

  return (
    <div className="mb-3">
      <ConfigFieldLabel htmlFor={inputTestId} fieldKey={fieldKey} showTooltip />
      <div className="mt-1 flex items-center gap-2">
        <ConfigCheckbox
          testId={checkboxTestId}
          checked={enabled}
          onCheckedChange={handleToggle}
          className="mr-2.5"
          aria-label={`Enable ${fieldKey}`}
        />
        <Input
          id={inputTestId}
          data-testid={inputTestId}
          type="text"
          className={managementConfigFieldInputClass}
          value={value}
          disabled={!enabled}
          aria-invalid={error ? true : undefined}
          onChange={handleInputChange}
        />
      </div>
      <FieldError id={inputTestId} message={error ?? undefined} />
    </div>
  )
}
