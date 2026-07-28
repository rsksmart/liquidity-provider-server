import { type ConfirmationEntry } from '@feature/management/config/format'
import { hasDuplicateConfirmationAmounts } from '@feature/management/config/validation'
import {
  managementBootstrapInputClass,
  managementBootstrapSmButtonClass,
  managementDangerButtonClass,
  managementModalTitleClass,
  managementSecondaryButtonClass,
} from '@feature/management/management-styles'
import type { ConfirmationsPerAmount } from '@shared/types/initial-data'
import { etherToWei, etherToWeiOr, weiToEther } from '@shared/utils/wei'
import { type ChangeEvent, type MouseEvent, useCallback } from 'react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

import { FieldError } from '../FieldError'

const FIELD_STYLE = { maxWidth: '180px' } as const

/** A confirmation tier as edited in the form: amount shown in ether. */
export interface ConfirmationRow {
  amountEther: string
  confirmations: string
}

type ConfigKey = 'rskConfirmations' | 'btcConfirmations'

function amountUnit(configKey: ConfigKey): string {
  return configKey === 'btcConfirmations' ? 'BTC' : 'rBTC'
}

function amountLabel(configKey: ConfigKey): string {
  return configKey === 'btcConfirmations' ? 'BTC amount' : 'rBTC amount'
}

function rowIndexFromTarget(target: EventTarget & HTMLElement): number {
  return Number(target.dataset.index)
}

/**
 * Seeds editor rows from a stored confirmations record, sorting ascending by
 * ether amount.
 */
export function confirmationRecordToRows(
  record: ConfirmationsPerAmount,
): ConfirmationRow[] {
  return Object.entries(record)
    .sort(([weiA], [weiB]) => parseFloat(weiToEther(weiA)) - parseFloat(weiToEther(weiB)))
    .map(([wei, confirmation]) => ({
      amountEther: weiToEther(wei),
      confirmations: String(confirmation),
    }))
}

/**
 * Converts display rows into the wei-keyed entries consumed by
 * `formatGeneralConfig`, preserving order. Empty/invalid amounts become `''`
 * and empty/invalid confirmations become `undefined` so `formatGeneralConfig`
 * drops incomplete rows.
 */
export function confirmationRowsToEntries(
  rows: ConfirmationRow[],
): ConfirmationEntry[] {
  return rows.map((row) => {
    const amountText = row.amountEther.trim()
    const amount =
      amountText === '' ? '' : etherToWeiOr(amountText, '')

    const confirmationText = row.confirmations.trim()
    let confirmation: number | undefined
    if (confirmationText !== '') {
      const parsed = Number(confirmationText)
      confirmation = Number.isNaN(parsed) ? undefined : parsed
    }

    return { amount, confirmation }
  })
}

/**
 * Validates confirmation rows: non-empty amount/confirmations, convertible
 * amounts, at least one complete row, and no duplicate amounts.
 */
export function validateConfirmationRows(
  rows: ConfirmationRow[],
  configKey: ConfigKey,
): string[] {
  const errors: string[] = []

  rows.forEach((row) => {
    const amountText = row.amountEther.trim()
    const confirmationText = row.confirmations.trim()

    if (amountText === '') {
      errors.push(`Please enter a non-empty value for "${amountLabel(configKey)}."`)
    } else {
      try {
        etherToWei(amountText)
      } catch {
        errors.push(
          `Invalid input "${row.amountEther}" for ${amountLabel(configKey)}. Please enter a valid non-negative number.`,
        )
      }
    }

    if (confirmationText === '') {
      errors.push('Please enter a non-empty value for "confirmations."')
    } else {
      const parsed = Number(confirmationText)
      if (Number.isNaN(parsed) || !Number.isInteger(parsed) || parsed < 0) {
        errors.push(
          `Invalid input "${row.confirmations}" for confirmations. Please enter a valid non-negative integer.`,
        )
      }
    }
  })

  const completeEntries = confirmationRowsToEntries(rows).filter(
    (entry) => entry.amount !== '' && entry.confirmation !== undefined,
  )

  if (completeEntries.length === 0) {
    errors.push(`Please provide at least one fully filled out entry for ${configKey}.`)
  }

  if (hasDuplicateConfirmationAmounts(completeEntries)) {
    errors.push(
      `Duplicate rBTC amounts found in ${configKey}. Please remove duplicates before saving.`,
    )
  }

  return errors
}

interface ConfirmationTiersEditorProps {
  configKey: ConfigKey
  /** Ordered display rows; use {@link confirmationRecordToRows} to seed. */
  value: ConfirmationRow[]
  onChange: (rows: ConfirmationRow[]) => void
  onDirty?: () => void
  /** Section-level validation errors (see {@link validateConfirmationRows}). */
  errors?: string[]
}

/** Editor for `rskConfirmations` / `btcConfirmations` amount→confirmations rows. */
export function ConfirmationTiersEditor({
  configKey,
  value,
  onChange,
  onDirty,
  errors,
}: ConfirmationTiersEditorProps) {
  const unit = amountUnit(configKey)

  const handleAmountChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      const index = rowIndexFromTarget(event.currentTarget)
      const next = value.map((row, i) =>
        i === index ? { ...row, amountEther: event.target.value } : row,
      )
      onChange(next)
      onDirty?.()
    },
    [value, onChange, onDirty],
  )

  const handleConfirmationChange = useCallback(
    (event: ChangeEvent<HTMLInputElement>) => {
      const index = rowIndexFromTarget(event.currentTarget)
      const next = value.map((row, i) =>
        i === index ? { ...row, confirmations: event.target.value } : row,
      )
      onChange(next)
      onDirty?.()
    },
    [value, onChange, onDirty],
  )

  const handleRemove = useCallback(
    (event: MouseEvent<HTMLButtonElement>) => {
      const index = rowIndexFromTarget(event.currentTarget)
      onChange(value.filter((_, i) => i !== index))
      onDirty?.()
    },
    [value, onChange, onDirty],
  )

  const handleAdd = useCallback(() => {
    onChange([...value, { amountEther: '', confirmations: '' }])
    onDirty?.()
  }, [value, onChange, onDirty])

  return (
    <div
      className="confirmation-config mb-3"
      data-config-key={configKey}
      data-testid={`config-${configKey}`}
    >
      <h5 className={cn(managementModalTitleClass, 'mb-2')}>{configKey}</h5>

      <div className="entries-container">
        {value.map((row, index) => {
          const confirmationTestId = `config-${configKey}-${String(index)}`
          const amountTestId = `config-${configKey}-amount-${String(index)}`
          return (
            <div
              key={index}
              className="mb-2 flex items-center gap-2"
              data-testid={`${confirmationTestId}-row`}
            >
              <div className="flex items-center gap-1">
                <Input
                  type="text"
                  data-testid={amountTestId}
                  data-index={String(index)}
                  data-config-key={configKey}
                  data-field="amount"
                  aria-label={`${amountLabel(configKey)} ${String(index + 1)}`}
                  className={managementBootstrapInputClass}
                  style={FIELD_STYLE}
                  placeholder="Amount"
                  value={row.amountEther}
                  onChange={handleAmountChange}
                />
                <span className="text-[#212529]">{unit}</span>
              </div>

              <div className="flex items-center gap-1">
                <Input
                  type="number"
                  data-testid={confirmationTestId}
                  data-index={String(index)}
                  data-config-key={configKey}
                  data-field="confirmation"
                  aria-label={`confirmations ${String(index + 1)}`}
                  className={managementBootstrapInputClass}
                  style={FIELD_STYLE}
                  placeholder="Confirmations"
                  value={row.confirmations}
                  onChange={handleConfirmationChange}
                />
                <span className="text-[#212529]">confirmations</span>
              </div>

              <Button
                type="button"
                variant="bootstrap"
                size="sm"
                className={cn(managementBootstrapSmButtonClass, managementDangerButtonClass)}
                data-testid={`${confirmationTestId}-remove`}
                data-index={String(index)}
                onClick={handleRemove}
              >
                Remove
              </Button>
            </div>
          )
        })}
      </div>

      <Button
        type="button"
        variant="bootstrap"
        size="sm"
        className={cn(managementBootstrapSmButtonClass, managementSecondaryButtonClass, 'mt-2')}
        data-testid={`config-${configKey}-add`}
        onClick={handleAdd}
      >
        Add Entry
      </Button>

      {errors?.map((message, index) => (
        <FieldError
          key={`${message}-${String(index)}`}
          id={`config-${configKey}-error-${String(index)}`}
          message={message}
        />
      ))}
    </div>
  )
}
