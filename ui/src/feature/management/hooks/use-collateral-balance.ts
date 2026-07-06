import { apiFetch } from '@api/management/utils/api-fetch'
import { weiToEther } from '@shared/utils/wei'
import { useCallback, useEffect, useRef, useState } from 'react'

interface CollateralResponse {
  collateral: string | number
}

function formatCollateralBalance(wei: string | number): string {
  return `${weiToEther(wei)} rBTC`
}

async function readCollateralBalance(endpoint: string): Promise<string | null> {
  try {
    const response = await apiFetch(endpoint)
    const data = (await response.json()) as CollateralResponse
    return formatCollateralBalance(data.collateral)
  } catch {
    return null
  }
}

export function useCollateralBalance(endpoint: string) {
  const [balance, setBalance] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const requestIdRef = useRef(0)

  const applyBalance = useCallback((requestId: number, nextBalance: string | null) => {
    if (requestId !== requestIdRef.current) {
      return
    }
    setBalance(nextBalance)
    setLoading(false)
  }, [])

  const refresh = useCallback(() => {
    const requestId = ++requestIdRef.current
    setLoading(true)
    return readCollateralBalance(endpoint).then((nextBalance) => {
      applyBalance(requestId, nextBalance)
    })
  }, [applyBalance, endpoint])

  useEffect(() => {
    const requestId = ++requestIdRef.current
    void readCollateralBalance(endpoint).then((nextBalance) => {
      applyBalance(requestId, nextBalance)
    })
    return () => {
      requestIdRef.current += 1
    }
  }, [applyBalance, endpoint])

  return { balance, loading, refresh }
}
