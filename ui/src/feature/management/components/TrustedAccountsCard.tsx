import { ApiFetchError } from '@api/management/types/errors'
import { apiFetch } from '@api/management/utils/api-fetch'
import { AddTrustedAccountDialog } from '@feature/management/components/AddTrustedAccountDialog'
import { getTrustedAccountApiErrorMessage } from '@feature/management/components/trusted-account-form'
import { useTrustedAccounts } from '@feature/management/hooks/use-trusted-accounts'
import {
  managementBootstrapSmButtonClass,
  managementCardClass,
  managementCardContentClass,
  managementCardHeaderClass,
  managementCardTitleClass,
  managementDangerButtonClass,
  managementDangerTextClass,
  managementFieldTitleClass,
  managementLoadingBarClass,
} from '@feature/management/management-styles'
import { weiToEther } from '@shared/utils/wei'
import { type MouseEvent, useCallback, useState } from 'react'
import { toast } from 'sonner'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'

/** Legacy configUtils.formatCap — at most four decimals plus unit. */
export function formatCap(value: string, unit: string): string {
  try {
    const num = Number.parseFloat(value)
    return `${Number.parseFloat(num.toFixed(4)).toString()} ${unit}`
  } catch (error) {
    const message =
      error instanceof Error ? error.message : 'Failed to format value'
    return `Error: ${message}`
  }
}

export function TrustedAccountsCard() {
  const { accounts, loading, error, refresh } = useTrustedAccounts()
  const [pendingAddress, setPendingAddress] = useState<string | null>(null)
  const [removing, setRemoving] = useState(false)

  const handleConfirmRemove = useCallback(async () => {
    if (!pendingAddress || removing) {
      return
    }

    setRemoving(true)
    try {
      await apiFetch(
        `/management/trusted-accounts?address=${encodeURIComponent(pendingAddress)}`,
        { method: 'DELETE' },
      )
      setPendingAddress(null)
      toast.success('Configuration saved successfully!')
      await refresh()
    } catch (err) {
      const message =
        err instanceof ApiFetchError
          ? getTrustedAccountApiErrorMessage(err.body, 'Unknown error')
          : err instanceof Error
            ? err.message
            : 'Unknown error'
      toast.error(`Error removing trusted account: ${message}`)
    } finally {
      setRemoving(false)
    }
  }, [pendingAddress, refresh, removing])

  const handleRemoveClick = useCallback(
    (event: MouseEvent<HTMLButtonElement>) => {
      const address = event.currentTarget.dataset.address
      if (address) {
        setPendingAddress(address)
      }
    },
    [],
  )

  const handleAlertOpenChange = useCallback(
    (open: boolean) => {
      if (!open && !removing) {
        setPendingAddress(null)
      }
    },
    [removing],
  )

  const handleConfirmClick = useCallback(
    (event: MouseEvent<HTMLButtonElement>) => {
      event.preventDefault()
      void handleConfirmRemove()
    },
    [handleConfirmRemove],
  )

  return (
    <Card className={managementCardClass} data-testid="trusted-accounts-card">
      <CardHeader className={managementCardHeaderClass}>
        <CardTitle className={managementCardTitleClass}>
          Trusted Accounts
        </CardTitle>
      </CardHeader>
      <CardContent className={managementCardContentClass}>
        <div className="mb-3 flex items-center justify-between gap-3">
          <h2 className={managementFieldTitleClass}>Accounts List</h2>
          <AddTrustedAccountDialog onSuccess={refresh} />
        </div>

        <div
          className="max-h-[205px] w-full overflow-x-auto overflow-y-auto rounded"
          data-testid="trusted-accounts-container"
        >
          <table
            className="w-full caption-bottom text-base"
            data-testid="trusted-accounts-table"
          >
            <thead>
              <tr className="border-b border-[#dee2e6] text-left">
                <th className="px-2 py-2 font-medium whitespace-nowrap">
                  Name
                </th>
                <th className="px-2 py-2 font-medium whitespace-nowrap">
                  Address
                </th>
                <th className="px-2 py-2 text-right font-medium whitespace-nowrap">
                  BTC Cap
                </th>
                <th className="px-2 py-2 text-right font-medium whitespace-nowrap">
                  rBTC Cap
                </th>
                <th className="px-2 py-2 font-medium whitespace-nowrap">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {error ? (
                <tr>
                  <td
                    colSpan={5}
                    className={cn(
                      'px-2 py-3 text-center',
                      managementDangerTextClass,
                    )}
                  >
                    Error: {error}
                  </td>
                </tr>
              ) : accounts.length === 0 && !loading ? (
                <tr>
                  <td colSpan={5} className="px-2 py-3 text-center">
                    No trusted accounts found.
                  </td>
                </tr>
              ) : (
                accounts.map((account) => (
                  <tr
                    key={account.address}
                    className="border-b border-[#dee2e6] odd:bg-black/[0.02]"
                  >
                    <td className="px-2 py-2 whitespace-nowrap">
                      {account.name || 'Unknown'}
                    </td>
                    <td className="px-2 py-2 font-mono text-[0.85rem] whitespace-nowrap">
                      {account.address}
                    </td>
                    <td className="px-2 py-2 text-right font-mono whitespace-nowrap">
                      {formatCap(weiToEther(account.btcLockingCap), 'BTC')}
                    </td>
                    <td className="px-2 py-2 text-right font-mono whitespace-nowrap">
                      {formatCap(weiToEther(account.rbtcLockingCap), 'rBTC')}
                    </td>
                    <td className="px-2 py-2 whitespace-nowrap">
                      <Button
                        type="button"
                        variant="bootstrap"
                        size="sm"
                        className={cn(
                          managementBootstrapSmButtonClass,
                          managementDangerButtonClass,
                          'hover:bg-[#bb2d3b]',
                        )}
                        data-address={account.address}
                        data-testid={`remove-trusted-account-${account.address}`}
                        onClick={handleRemoveClick}
                      >
                        Remove
                      </Button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <div
          className={cn(managementLoadingBarClass, loading && 'is-visible')}
          data-testid="trusted-accounts-loading-bar"
          role="progressbar"
          aria-hidden={!loading}
        />
      </CardContent>

      <AlertDialog
        open={pendingAddress !== null}
        onOpenChange={handleAlertOpenChange}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Remove trusted account</AlertDialogTitle>
            <AlertDialogDescription>
              {pendingAddress
                ? `Are you sure you want to remove the trusted account with address ${pendingAddress}?`
                : null}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={removing}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              disabled={removing}
              onClick={handleConfirmClick}
              data-testid="confirm-remove-trusted-account"
            >
              Confirm
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </Card>
  )
}
