import { getDisplayLabel, getTooltipText } from '@feature/management/config/labels'
import { managementBootstrapLabelClass } from '@feature/management/management-styles'

import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

import { QuestionIcon } from './QuestionIcon'

interface ConfigFieldLabelProps {
  /** Associates the label with its control (`htmlFor`). */
  htmlFor: string
  /** Config key used to resolve the display label and tooltip copy. */
  fieldKey: string
  /** Overrides the resolved display text (defaults to `getDisplayLabel(fieldKey)`). */
  label?: string
  /** Renders the label bold. */
  bold?: boolean
  /** Appends the question-mark tooltip icon. */
  showTooltip?: boolean
  /** Config key used for the tooltip copy when it differs from `fieldKey`. */
  tooltipKey?: string
}

/** Form label for a configuration field, with optional question-mark tooltip. */
export function ConfigFieldLabel({
  htmlFor,
  fieldKey,
  label,
  bold = false,
  showTooltip = false,
  tooltipKey,
}: ConfigFieldLabelProps) {
  return (
    <Label htmlFor={htmlFor} className={cn(managementBootstrapLabelClass, bold && 'font-bold')}>
      <span>{label ?? getDisplayLabel(fieldKey)}</span>
      {showTooltip && <QuestionIcon text={getTooltipText(tooltipKey ?? fieldKey)} />}
    </Label>
  )
}
