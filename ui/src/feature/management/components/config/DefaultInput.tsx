import { managementConfigFieldInputClass } from '@feature/management/management-styles'
import { type ChangeEvent, useCallback } from 'react'

import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

import { FieldError } from '../FieldError'
import { ConfigFieldLabel } from './ConfigFieldLabel'
import type { SectionPrefix } from './types'

interface DefaultInputProps {
  sectionPrefix: SectionPrefix
  /** Any remaining key: timeForDeposit, expireTime, expireBlocks, callTime, reimbursementWindowBlocks, … */
  fieldKey: string
  value: string
  onChange: (next: string) => void
  onDirty?: () => void
  error?: string | null
}

/** Fallback text input for numeric/string config keys (40% width, no tooltip). */
export function DefaultInput({
  sectionPrefix,
  fieldKey,
  value,
  onChange,
  onDirty,
  error,
}: DefaultInputProps) {
  const testId = `config-${sectionPrefix}-${fieldKey}-input`

  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      onChange(event.target.value)
      onDirty?.()
    },
    [onChange, onDirty],
  )

  return (
    <div className="mb-3">
      <ConfigFieldLabel htmlFor={testId} fieldKey={fieldKey} />
      <Input
        id={testId}
        data-testid={testId}
        type="text"
        className={cn('mt-1', managementConfigFieldInputClass)}
        value={value}
        aria-invalid={error ? true : undefined}
        onChange={handleChange}
      />
      <FieldError id={testId} message={error ?? undefined} />
    </div>
  )
}
