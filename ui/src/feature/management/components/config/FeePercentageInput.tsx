import { managementConfigFieldInputClass } from '@feature/management/management-styles'
import { type ChangeEvent, useCallback, useState } from 'react'

import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

import { FieldError } from '../FieldError'
import { ConfigFieldLabel } from './ConfigFieldLabel'
import { normalizePercentage } from './percentage'
import type { SectionPrefix } from './types'

interface FeePercentageInputProps {
  sectionPrefix: SectionPrefix
  /** Percentage key (e.g. `feePercentage`). */
  fieldKey: string
  /** Percentage display string (0–100, optionally with a trailing `%` while editing). */
  value: string
  onChange: (next: string) => void
  onDirty?: () => void
  error?: string | null
}

/**
 * Percentage text field. Accepts an optional trailing `%` and, on blur,
 * strips it and clamps the value to 0–100.
 */
export function FeePercentageInput({
  sectionPrefix,
  fieldKey,
  value,
  onChange,
  onDirty,
  error,
}: FeePercentageInputProps) {
  const testId = `config-${sectionPrefix}-${fieldKey}-input`
  const [localError, setLocalError] = useState<string | null>(null)

  const shownError = error ?? localError

  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      setLocalError(null)
      onChange(event.target.value)
      onDirty?.()
    },
    [onChange, onDirty],
  )

  const handleBlur = useCallback(() => {
    const { value: normalized, error: normalizeError } = normalizePercentage(value)
    setLocalError(normalizeError)
    if (!normalizeError && normalized !== value) {
      onChange(normalized)
    }
  }, [onChange, value])

  return (
    <div className="mb-3">
      <ConfigFieldLabel htmlFor={testId} fieldKey={fieldKey} />
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
