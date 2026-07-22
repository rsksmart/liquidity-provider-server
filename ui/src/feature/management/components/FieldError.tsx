import { cn } from '@/lib/utils'

import { managementFieldErrorClass } from '../management-styles'

interface FieldErrorProps {
  id?: string
  message?: string
  className?: string
}

/** Accessible inline field error — mirrors legacy Bootstrap invalid-feedback. */
export function FieldError({ id, message, className }: FieldErrorProps) {
  if (!message) {
    return null
  }

  return (
    <p
      id={id}
      role="alert"
      className={cn(managementFieldErrorClass, className)}
      data-testid={id ? `${id}-error` : undefined}
    >
      {message}
    </p>
  )
}
