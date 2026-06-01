import type { InitialDataPayload } from '@shared/types/initial-data'
import { useMemo } from 'react'

let cachedInitialData: InitialDataPayload | null = null

function readInitialDataFromDom(): InitialDataPayload {
  const element = document.getElementById('initial-data')
  const content = element?.textContent
  if (!content?.trim()) {
    throw new Error('initial-data script element missing or empty')
  }

  return JSON.parse(content) as InitialDataPayload
}

export function getInitialData(): InitialDataPayload {
  if (!cachedInitialData) {
    cachedInitialData = readInitialDataFromDom()
  }
  return cachedInitialData
}

export function useInitialData(): InitialDataPayload {
  return useMemo(() => getInitialData(), [])
}

export function resetInitialDataCacheForTests(): void {
  cachedInitialData = null
}
