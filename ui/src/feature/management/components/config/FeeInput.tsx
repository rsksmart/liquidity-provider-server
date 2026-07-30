import { isMaxLiquidityKey, validateMaxLiquidity } from '@feature/management/config/validation'
import { managementConfigFieldInputClass } from '@feature/management/management-styles'
import { type ChangeEvent, useCallback, useState } from 'react'

import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

import { FieldError } from '../FieldError'
import { ConfigFieldLabel } from './ConfigFieldLabel'
import type { SectionPrefix } from './types'

interface FeeInputProps {
  sectionPrefix: SectionPrefix
  /** Non-toggable fee key: penaltyFee, callFee, maxValue, minValue, bridgeTransactionMin, maxLiquidity. */
  fieldKey: string
  /**
   * Display value (ether string). Callers convert stored wei via `weiToEther`
   * before passing it in and back via `etherToWei` on save.
   */
  value: string
  onChange: (next: string) => void
  onDirty?: () => void
  /** External error; takes precedence over the built-in maxLiquidity validation. */
  error?: string | null
}

/**
 * Text input for non-toggable fee fields.
 * Tooltip is skipped for `maxLiquidity`, which uses a bold label and inline validation.
 */
export function FeeInput({
  sectionPrefix,
  fieldKey,
  value,
  onChange,
  onDirty,
  error,
}: FeeInputProps) {
  const testId = `config-${sectionPrefix}-${fieldKey}-input`
  const isMaxLiquidity = isMaxLiquidityKey(fieldKey)
  const [touched, setTouched] = useState(false)

  const liquidityError =
    isMaxLiquidity && touched ? validateMaxLiquidity(value).error : null
  const shownError = error ?? liquidityError

  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      setTouched(true)
      onChange(event.target.value)
      onDirty?.()
    },
    [onChange, onDirty],
  )

  const handleBlur = useCallback(() => {
    setTouched(true)
  }, [])

  return (
    <div className="mb-3">
      <ConfigFieldLabel
        htmlFor={testId}
        fieldKey={fieldKey}
        bold={isMaxLiquidity}
        showTooltip={!isMaxLiquidity}
      />
      <Input
        id={testId}
        data-testid={testId}
        type="text"
        className={cn('mt-1', managementConfigFieldInputClass)}
        value={value}
        aria-invalid={shownError ? true : undefined}
        onChange={handleChange}
        onBlur={handleBlur}
      />
      <FieldError id={testId} message={shownError ?? undefined} />
    </div>
  )
}
