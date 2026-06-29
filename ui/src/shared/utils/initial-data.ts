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

export function replaceInitialDataPayload(payload: InitialDataPayload): void {
  const initialEl = document.getElementById('initial-data')
  if (!initialEl) {
    throw new Error('initial-data script element missing')
  }
  initialEl.textContent = JSON.stringify(payload)
  cachedInitialData = null
}

export function resetInitialDataCache(): void {
  cachedInitialData = null
}

interface LpsDevShellResponse {
  csrf: string
  initialData: InitialDataPayload
}

const LPS_DEV_BOOTSTRAP_SKIPPED_MESSAGE =
  'LPS dev bootstrap skipped — using Vite stubs. Start LPS on :8080 for login/logout flows.'

function warnDevBootstrapSkipped(reason: unknown): void {
  console.warn(LPS_DEV_BOOTSTRAP_SKIPPED_MESSAGE, reason)
}

export async function bootstrapDevEnvironment(): Promise<void> {
  try {
    const res = await fetch('/__dev/lps-shell', { credentials: 'include' })
    if (!res.ok) {
      const body = (await res.json().catch(() => null)) as { error?: string } | null
      warnDevBootstrapSkipped(body?.error ?? res.status)
      return
    }

    const { csrf, initialData } = (await res.json()) as LpsDevShellResponse
    document.querySelector('meta[name="csrf-token"]')?.setAttribute('content', csrf)
    replaceInitialDataPayload(initialData)
  } catch (err) {
    warnDevBootstrapSkipped(err)
  }
}
