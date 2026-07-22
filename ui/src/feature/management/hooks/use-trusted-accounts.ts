import type {
  TrustedAccount,
  TrustedAccountsListResponse,
} from '@api/management/types/trusted-accounts'
import { apiFetch } from '@api/management/utils/api-fetch'
import { useCallback, useEffect, useRef, useState } from 'react'
import { toast } from 'sonner'

function errorMessage(err: unknown): string {
  if (err instanceof Error && err.message) {
    return err.message
  }
  return 'Unknown error'
}

async function fetchTrustedAccounts(): Promise<TrustedAccount[]> {
  const response = await apiFetch('/management/trusted-accounts')
  const data = (await response.json()) as TrustedAccountsListResponse
  return data.accounts
}

export function useTrustedAccounts() {
  const [accounts, setAccounts] = useState<TrustedAccount[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const requestIdRef = useRef(0)

  const applyResult = useCallback(
    (
      requestId: number,
      nextAccounts: TrustedAccount[] | null,
      nextError: string | null,
    ) => {
      if (requestId !== requestIdRef.current) {
        return
      }
      if (nextError) {
        setAccounts([])
        setError(nextError)
      } else {
        setAccounts(nextAccounts ?? [])
        setError(null)
      }
      setLoading(false)
    },
    [],
  )

  const runFetch = useCallback(
    (requestId: number) =>
      fetchTrustedAccounts()
        .then((nextAccounts) => {
          applyResult(requestId, nextAccounts, null)
        })
        .catch((err: unknown) => {
          const message = errorMessage(err)
          applyResult(requestId, null, message)
          if (requestId === requestIdRef.current) {
            toast.error(`Failed to load trusted accounts: ${message}`)
          }
        }),
    [applyResult],
  )

  const refresh = useCallback(() => {
    const requestId = ++requestIdRef.current
    setLoading(true)
    return runFetch(requestId)
  }, [runFetch])

  useEffect(() => {
    const requestId = ++requestIdRef.current
    void runFetch(requestId)
    return () => {
      requestIdRef.current += 1
    }
  }, [runFetch])

  return { accounts, loading, error, refresh }
}
