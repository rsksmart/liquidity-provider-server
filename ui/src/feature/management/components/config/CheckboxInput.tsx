import { getDisplayLabel } from '@feature/management/config/labels'
import { managementBootstrapLabelClass } from '@feature/management/management-styles'
import { useCallback } from 'react'

import { Label } from '@/components/ui/label'

import { ConfigCheckbox } from './ConfigCheckbox'
import type { SectionPrefix } from './types'

interface CheckboxInputProps {
  sectionPrefix: SectionPrefix
  /** Boolean config key (e.g. `publicLiquidityCheck`). */
  fieldKey: string
  value: boolean
  onChange: (next: boolean) => void
  onDirty?: () => void
}

/** Boolean checkbox field (`form-check-input` styling, no tooltip). */
export function CheckboxInput({
  sectionPrefix,
  fieldKey,
  value,
  onChange,
  onDirty,
}: CheckboxInputProps) {
  const testId = `config-${sectionPrefix}-${fieldKey}-checkbox`

  const handleCheckedChange = useCallback(
    (checked: boolean) => {
      onChange(checked)
      onDirty?.()
    },
    [onChange, onDirty],
  )

  return (
    <div className="mb-3">
      <Label htmlFor={testId} className={managementBootstrapLabelClass}>
        {getDisplayLabel(fieldKey)}
      </Label>
      <ConfigCheckbox
        id={testId}
        testId={testId}
        checked={value}
        onCheckedChange={handleCheckedChange}
        className="mt-1 mr-2.5"
      />
    </div>
  )
}
