import {
  validateExcessToleranceFixed,
  validateExcessTolerancePercentage,
} from '@feature/management/config/validation'
import { managementConfigFieldInputClass } from '@feature/management/management-styles'
import type { ExcessTolerance } from '@shared/types/initial-data'
import { etherToWeiOr, weiToEther } from '@shared/utils/wei'
import { type ChangeEvent, useCallback, useState } from 'react'

import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'

import { FieldError } from '../FieldError'
import { ConfigFieldLabel } from './ConfigFieldLabel'
import type { SectionPrefix } from './types'

const FIXED_PLACEHOLDER = 'Enter amount in rBTC'
const PERCENTAGE_PLACEHOLDER = 'Enter percentage (0-100)'

interface ExcessToleranceInputProps {
  sectionPrefix: SectionPrefix
  value: ExcessTolerance
  onChange: (next: ExcessTolerance) => void
  onDirty?: () => void
  error?: string | null
}

/**
 * Excess-tolerance field with a Fixed/Percentage switch
 * (ON = fixed rBTC amount, OFF = percentage 0–100).
 * Separate working values are kept per mode so toggling preserves edits.
 */
export function ExcessToleranceInput({
  sectionPrefix,
  value,
  onChange,
  onDirty,
  error,
}: ExcessToleranceInputProps) {
  const toggleTestId = `config-${sectionPrefix}-excessTolerance-toggle`
  const inputTestId = `config-${sectionPrefix}-excessTolerance-input`

  const [fixedText, setFixedText] = useState(() => weiToEther(value.fixedValue))
  const [percentageText, setPercentageText] = useState(() => value.percentageValue || '0')
  const [localError, setLocalError] = useState<string | null>(null)

  const isFixed = value.isFixed
  const currentText = isFixed ? fixedText : percentageText
  const shownError = error ?? localError

  const handleToggle = useCallback(
    (checked: boolean) => {
      setLocalError(null)
      onChange({
        isFixed: checked,
        fixedValue: etherToWeiOr(fixedText, value.fixedValue),
        percentageValue: percentageText,
      })
      onDirty?.()
    },
    [fixedText, percentageText, value.fixedValue, onChange, onDirty],
  )

  const handleInputChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      const text = event.target.value
      if (isFixed) {
        setFixedText(text)
        const validation = validateExcessToleranceFixed(text)
        setLocalError(validation.error)
        onChange({
          isFixed: true,
          fixedValue: validation.isValid ? etherToWeiOr(text, value.fixedValue) : value.fixedValue,
          percentageValue: value.percentageValue,
        })
      } else {
        setPercentageText(text)
        const validation = validateExcessTolerancePercentage(text)
        setLocalError(validation.error)
        onChange({
          isFixed: false,
          fixedValue: value.fixedValue,
          percentageValue: text,
        })
      }
      onDirty?.()
    },
    [isFixed, value.fixedValue, value.percentageValue, onChange, onDirty],
  )

  return (
    <div className="mb-3">
      <ConfigFieldLabel htmlFor={inputTestId} fieldKey="excessTolerance" showTooltip />
      <div className="mt-1 flex items-center gap-2">
        <span className="form-check form-switch mr-[15px] inline-flex items-center gap-1">
          <Switch
            id={toggleTestId}
            data-testid={toggleTestId}
            checked={isFixed}
            onCheckedChange={handleToggle}
            aria-label="Fixed"
          />
          <span className="text-[0.85em] text-[#666] select-none">Fixed</span>
        </span>
        <Input
          id={inputTestId}
          data-testid={inputTestId}
          type="text"
          className={managementConfigFieldInputClass}
          value={currentText}
          placeholder={isFixed ? FIXED_PLACEHOLDER : PERCENTAGE_PLACEHOLDER}
          aria-invalid={shownError ? true : undefined}
          onChange={handleInputChange}
        />
      </div>
      <FieldError id={inputTestId} message={shownError ?? undefined} />
    </div>
  )
}
