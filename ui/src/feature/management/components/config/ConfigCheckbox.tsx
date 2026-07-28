import { type ChangeEvent, useCallback } from 'react'

import { cn } from '@/lib/utils'

interface ConfigCheckboxProps {
  testId: string
  checked: boolean
  onCheckedChange: (checked: boolean) => void
  id?: string
  className?: string
  'aria-label'?: string
}

/** Bootstrap `form-check-input` checkbox control. */
export function ConfigCheckbox({
  id,
  testId,
  checked,
  onCheckedChange,
  className,
  'aria-label': ariaLabel,
}: ConfigCheckboxProps) {
  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      onCheckedChange(event.target.checked)
    },
    [onCheckedChange],
  )

  return (
    <input
      id={id}
      type="checkbox"
      className={cn('form-check-input', className)}
      data-testid={testId}
      checked={checked}
      onChange={handleChange}
      aria-label={ariaLabel}
    />
  )
}
