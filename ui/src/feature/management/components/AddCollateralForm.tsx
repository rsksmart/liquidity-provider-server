import { ApiFetchError } from '@api/management/types/errors'
import { apiFetch } from '@api/management/utils/api-fetch'
import { useCollateralBalance } from '@feature/management/hooks/use-collateral-balance'
import {
  managementBootstrapButtonClass,
  managementBootstrapInputClass,
  managementBootstrapLabelClass,
  managementCollateralButtonsClass,
  managementFieldTextClass,
  managementFieldTitleClass,
  managementLoadingBarClass,
} from '@feature/management/management-styles'
import { etherToWei, weiToApiAmount } from '@shared/utils/wei'
import { type ChangeEvent, type FormEvent, useCallback, useState } from 'react'
import { toast } from 'sonner'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

const COLLATERAL_CONFIG = {
  pegin: {
    title: 'Pegin Collateral',
    amountLabel: 'Add Pegin Collateral Amount',
    buttonLabel: 'Add Pegin Collateral',
    balanceEndpoint: '/pegin/collateral',
    addEndpoint: '/pegin/addCollateral',
    testIdPrefix: 'pegin',
  },
  pegout: {
    title: 'Pegout Collateral',
    amountLabel: 'Add Pegout Collateral Amount',
    buttonLabel: 'Add Pegout Collateral',
    balanceEndpoint: '/pegout/collateral',
    addEndpoint: '/pegout/addCollateral',
    testIdPrefix: 'pegout',
  },
} as const


function invalidAmountMessage(amount: string): string {
  return `Invalid input "${amount}" for collateral amount. Please enter a valid number.`
}

function getApiErrorMessage(body: unknown): string {
  if (typeof body === 'object' && body !== null && 'message' in body) {
    const message = body.message
    if (typeof message === 'string') {
      return message
    }
  }
  return 'Unknown error'
}

interface AddCollateralFormProps {
  kind: keyof typeof COLLATERAL_CONFIG
}

export function AddCollateralForm({ kind }: AddCollateralFormProps) {
  const config = COLLATERAL_CONFIG[kind]
  const { balance, loading: balanceLoading, refresh } = useCollateralBalance(config.balanceEndpoint)
  const [amount, setAmount] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const submit = useCallback(
    async (amountEther: string) => {
      let weiStr: string
      try {
        weiStr = etherToWei(amountEther)
      } catch {
        toast.error(invalidAmountMessage(amountEther))
        return
      }

      if (BigInt(weiStr) <= 0n) {
        toast.error(invalidAmountMessage(amountEther))
        return
      }

      setSubmitting(true)
      try {
        await apiFetch.post(config.addEndpoint, { amount: weiToApiAmount(weiStr) })

        setAmount('')
        await refresh()
      } catch (error) {
        if (error instanceof ApiFetchError) {
          toast.error(`Error adding collateral: ${getApiErrorMessage(error.body)}`)
          return
        }
        toast.error(invalidAmountMessage(amountEther))
      } finally {
        setSubmitting(false)
      }
    },
    [config.addEndpoint, refresh],
  )

  const onSubmit = useCallback(
    (event: FormEvent<HTMLFormElement>) => {
      event.preventDefault()
      void submit(amount)
    },
    [amount, submit],
  )

  const onAmountChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setAmount(event.target.value)
  }, [])

  return (
    <form className="space-y-3" onSubmit={onSubmit}>
      <div>
        <h3 className={managementFieldTitleClass}>{config.title}</h3>
        <p className={managementFieldTextClass} data-testid={`${config.testIdPrefix}-collateral-balance`}>
          {balanceLoading ? 'Loading…' : (balance ?? '—')}
        </p>
      </div>
      <div className="mt-[15px] space-y-2">
        <Label
          htmlFor={`${config.testIdPrefix}-collateral-amount`}
          className={managementBootstrapLabelClass}
        >
          {config.amountLabel}
        </Label>
        <Input
          id={`${config.testIdPrefix}-collateral-amount`}
          type="number"
          placeholder="Enter amount in rBTC"
          value={amount}
          onChange={onAmountChange}
          data-testid={`${config.testIdPrefix}-collateral-amount`}
          className={managementBootstrapInputClass}
        />
      </div>
      <div className={managementCollateralButtonsClass}>
        <Button
          type="submit"
          variant="bootstrap"
          disabled={submitting}
          className={cn(managementBootstrapButtonClass, 'mr-2.5')}
          data-testid={`${config.testIdPrefix}-add-collateral-button`}
        >
          {config.buttonLabel}
        </Button>
        <div
          className={cn(managementLoadingBarClass, submitting && 'is-visible')}
          data-testid={`${config.testIdPrefix}-loading-bar`}
          role="progressbar"
          aria-label="Adding collateral"
          aria-hidden={!submitting}
        />
      </div>
    </form>
  )
}
