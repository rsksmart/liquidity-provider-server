import type {
  GeneralConfiguration,
  InitialDataPayload,
  PeginConfiguration,
  PegoutConfiguration,
  QuoteConfigurationBase,
  WireGeneralConfiguration,
  WireInitialDataPayload,
  WireNumeric,
  WirePeginConfiguration,
  WirePegoutConfiguration,
  WireQuoteConfigurationBase,
} from '@shared/types/initial-data'
import Decimal from 'decimal.js'
import { useMemo } from 'react'

let cachedInitialData: InitialDataPayload | null = null

/**
 * Wei and percentage fields arrive as JSON numbers (see `WireNumeric`). Decimal
 * is used instead of `String` so large wei values keep their full digits rather
 * than becoming exponential (`2e+21`), which downstream parsers reject.
 *
 * Values above `Number.MAX_SAFE_INTEGER` have already lost precision by the time
 * `JSON.parse` hands them over; the legacy card had the same limit.
 */
function numericToString(value: WireNumeric): string {
  return typeof value === 'number' ? new Decimal(value).toFixed() : value
}

function normalizeGeneral(general: WireGeneralConfiguration): GeneralConfiguration {
  return {
    rskConfirmations: general.rskConfirmations,
    btcConfirmations: general.btcConfirmations,
    publicLiquidityCheck: general.publicLiquidityCheck,
    maxLiquidity:
      general.maxLiquidity === null ? null : numericToString(general.maxLiquidity),
    reimbursementWindowBlocks: general.reimbursementWindowBlocks,
    excessTolerance: {
      isFixed: general.excessTolerance.isFixed,
      percentageValue: numericToString(general.excessTolerance.percentageValue),
      fixedValue: numericToString(general.excessTolerance.fixedValue),
    },
  }
}

function normalizeQuoteBase(config: WireQuoteConfigurationBase): QuoteConfigurationBase {
  return {
    timeForDeposit: config.timeForDeposit,
    penaltyFee: numericToString(config.penaltyFee),
    fixedFee: numericToString(config.fixedFee),
    feePercentage: numericToString(config.feePercentage),
    maxValue: numericToString(config.maxValue),
    minValue: numericToString(config.minValue),
  }
}

function normalizePegin(config: WirePeginConfiguration): PeginConfiguration {
  return { ...normalizeQuoteBase(config), callTime: config.callTime }
}

function normalizePegout(config: WirePegoutConfiguration): PegoutConfiguration {
  return {
    ...normalizeQuoteBase(config),
    expireTime: config.expireTime,
    expireBlocks: config.expireBlocks,
    bridgeTransactionMin: numericToString(config.bridgeTransactionMin),
  }
}

/**
 * Single boundary where the server payload becomes the display domain the
 * management feature is typed against.
 */
export function normalizeInitialData(payload: WireInitialDataPayload): InitialDataPayload {
  const { Configuration } = payload.data
  return {
    loggedIn: payload.loggedIn,
    data: {
      ...payload.data,
      Configuration: {
        general: normalizeGeneral(Configuration.general),
        pegin: normalizePegin(Configuration.pegin),
        pegout: normalizePegout(Configuration.pegout),
      },
    },
  }
}

function readInitialDataFromDom(): InitialDataPayload {
  const element = document.getElementById('initial-data')
  const content = element?.textContent
  if (!content?.trim()) {
    throw new Error('initial-data script element missing or empty')
  }

  return normalizeInitialData(JSON.parse(content) as WireInitialDataPayload)
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

export function replaceInitialDataPayload(payload: WireInitialDataPayload): void {
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
  initialData: WireInitialDataPayload
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
