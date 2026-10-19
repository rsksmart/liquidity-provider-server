import { ApiFetchError } from '@api/management/types/errors'
import { apiFetch } from '@api/management/utils/api-fetch'
import { FieldError } from '@feature/management/components/FieldError'
import {
  EMPTY_TRUSTED_ACCOUNT_FORM,
  getTrustedAccountApiErrorMessage,
  isValidationErrorBody,
  mapBackendValidationDetails,
  toAddTrustedAccountRequest,
  type TrustedAccountFormErrors,
  type TrustedAccountFormField,
  type TrustedAccountFormValues,
  validateTrustedAccountForm,
} from '@feature/management/trusted-account-form'
import {
  type ChangeEvent,
  type SubmitEvent,
  useCallback,
  useId,
  useState,
} from 'react'
import { toast } from 'sonner'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

import {
  managementBootstrapButtonClass,
  managementBootstrapInputClass,
  managementBootstrapLabelClass,
  managementBootstrapSmButtonClass,
  managementFormTextClass,
  managementModalFieldGroupClass,
  managementModalHeaderClass,
  managementModalTitleClass,
  managementSecondaryButtonClass,
} from '../management-styles'

interface AddTrustedAccountDialogProps {
  onSuccess: () => void | Promise<void>
}

/** testid kebab-case — matches other management UI ids (`pegin-collateral-amount`, …). */
const FIELD_TEST_IDS: Record<TrustedAccountFormField, string> = {
  name: 'account-name',
  address: 'account-address',
  btcLockingCap: 'btc-locking-cap',
  rbtcLockingCap: 'rbtc-locking-cap',
}

const FIELD_META: Array<{
  field: TrustedAccountFormField
  label: string
  placeholder: string
  helper: string
}> = [
  {
    field: 'name',
    label: 'Account Name',
    placeholder: 'Enter account name',
    helper: 'A friendly name to identify this account',
  },
  {
    field: 'address',
    label: 'Address',
    placeholder: 'Enter account address (0x...)',
    helper: 'The RSK address for this trusted account (starts with 0x)',
  },
  {
    field: 'btcLockingCap',
    label: 'BTC Locking Cap',
    placeholder: 'Enter BTC locking cap',
    helper: 'Maximum amount of BTC that can be locked by this account',
  },
  {
    field: 'rbtcLockingCap',
    label: 'rBTC Locking Cap',
    placeholder: 'Enter rBTC locking cap',
    helper: 'Maximum amount of rBTC that can be locked by this account',
  },
]

export function AddTrustedAccountDialog({
  onSuccess,
}: AddTrustedAccountDialogProps) {
  const formId = useId()
  const [open, setOpen] = useState(false)
  const [values, setValues] = useState<TrustedAccountFormValues>(
    EMPTY_TRUSTED_ACCOUNT_FORM,
  )
  const [errors, setErrors] = useState<TrustedAccountFormErrors>({})
  const [submitting, setSubmitting] = useState(false)

  const resetForm = useCallback(() => {
    setValues(EMPTY_TRUSTED_ACCOUNT_FORM)
    setErrors({})
  }, [])

  const handleOpenChange = useCallback(
    (nextOpen: boolean) => {
      setOpen(nextOpen)
      if (!nextOpen) {
        resetForm()
        setSubmitting(false)
      } else {
        setErrors({})
      }
    },
    [resetForm],
  )

  const clearFieldError = useCallback((field: TrustedAccountFormField) => {
    setErrors((prev) => {
      if (!prev[field]) {
        return prev
      }
      const next: TrustedAccountFormErrors = {}
      for (const [key, message] of Object.entries(prev) as Array<
        [TrustedAccountFormField, string]
      >) {
        if (key !== field) {
          next[key] = message
        }
      }
      return next
    })
  }, [])

  const updateField = useCallback(
    (field: TrustedAccountFormField, next: string) => {
      setValues((prev) => ({ ...prev, [field]: next }))
      clearFieldError(field)
    },
    [clearFieldError],
  )

  const submit = useCallback(async () => {
    if (submitting) {
      return
    }

    const nextErrors = validateTrustedAccountForm(values)
    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors)
      return
    }

    setSubmitting(true)
    setErrors({})

    try {
      const body = toAddTrustedAccountRequest(values)
      await apiFetch.post('/management/trusted-accounts', body)
      resetForm()
      setOpen(false)
      toast.success('Configuration saved successfully!')
      await onSuccess()
    } catch (error) {
      if (error instanceof ApiFetchError && isValidationErrorBody(error.body)) {
        const fieldErrors = mapBackendValidationDetails(error.body.details)
        if (Object.keys(fieldErrors).length > 0) {
          setErrors(fieldErrors)
          return
        }
      }

      const message =
        error instanceof ApiFetchError
          ? getTrustedAccountApiErrorMessage(error.body, 'Unknown error')
          : error instanceof Error
            ? error.message
            : 'Unknown error'
      toast.error(`Error adding trusted account: ${message}`)
    } finally {
      setSubmitting(false)
    }
  }, [onSuccess, resetForm, submitting, values])

  const onSubmit = useCallback(
    (event: SubmitEvent<HTMLFormElement>) => {
      event.preventDefault()
      void submit()
    },
    [submit],
  )

  const onFieldChange = useCallback(
    (field: TrustedAccountFormField) =>
      (event: ChangeEvent<HTMLInputElement>) => {
        updateField(field, event.target.value)
      },
    [updateField],
  )

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger
        render={
          <Button
            type="button"
            variant="bootstrap"
            size="sm"
            className={managementBootstrapSmButtonClass}
            data-testid="add-trusted-account-button"
          />
        }
      >
        Add Account
      </DialogTrigger>
      <DialogContent
        className={cn(
          'w-full max-w-[min(800px,calc(100%-2rem))] sm:max-w-[min(800px,calc(100%-2rem))]',
          'gap-0 rounded-[8px] p-0 text-base text-[#212529]',
        )}
        showCloseButton
      >
        <DialogHeader className={managementModalHeaderClass}>
          <DialogTitle className={managementModalTitleClass}>
            Add Trusted Account
          </DialogTitle>
        </DialogHeader>
        <form
          id={`${formId}-form`}
          onSubmit={onSubmit}
          className="px-4 py-4"
          noValidate
        >
          {FIELD_META.map((meta, index) => {
            const field: TrustedAccountFormField = meta.field
            const fieldValue: string = values[field]
            const inputId = `${formId}-${field}`
            const errorId = `${inputId}-error`
            const errorMessage = errors[field]
            return (
              <div key={field} className={managementModalFieldGroupClass}>
                <Label
                  htmlFor={inputId}
                  className={cn(
                    managementBootstrapLabelClass,
                    'mb-2 leading-normal',
                  )}
                >
                  {meta.label}
                </Label>
                <Input
                  id={inputId}
                  name={field}
                  value={fieldValue}
                  placeholder={meta.placeholder}
                  className={managementBootstrapInputClass}
                  aria-invalid={errorMessage ? true : undefined}
                  aria-describedby={errorMessage ? errorId : undefined}
                  autoFocus={index === 0}
                  data-testid={FIELD_TEST_IDS[field]}
                  onChange={onFieldChange(field)}
                />
                <p className={managementFormTextClass}>{meta.helper}</p>
                <FieldError id={errorId} message={errorMessage} />
              </div>
            )
          })}
        </form>
        <DialogFooter className="m-0 rounded-b-[8px] border-[#dee2e6] bg-transparent px-3 py-3 sm:justify-end">
          <DialogClose
            render={
              <Button
                type="button"
                className={cn(
                  managementBootstrapButtonClass,
                  managementSecondaryButtonClass,
                )}
              />
            }
          >
            Cancel
          </DialogClose>
          <Button
            type="submit"
            form={`${formId}-form`}
            variant="bootstrap"
            disabled={submitting}
            className={managementBootstrapButtonClass}
            data-testid="save-trusted-account-button"
          >
            Save
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
