import type {
  ConfirmationRow,
} from '@feature/management/components/config'
import {
  confirmationRecordToRows,
  confirmationRowsToEntries,
  normalizePercentage,
  validateConfirmationRows,
} from '@feature/management/components/config'
import { formatGeneralConfig } from '@feature/management/config/format'
import {
  validateConfig,
  validateExcessToleranceFixed,
  validateExcessTolerancePercentage,
  validateMaxLiquidity,
} from '@feature/management/config/validation'
import type {
  ExcessTolerance,
  FullConfiguration,
  GeneralConfiguration,
  PeginConfiguration,
  PegoutConfiguration,
  QuoteConfigurationBase,
} from '@shared/types/initial-data'
import { useInitialData } from '@shared/utils/initial-data'
import { etherToWei, weiToEther } from '@shared/utils/wei'
import { useCallback, useMemo, useState } from 'react'

/** General section as edited in the form (display-domain strings). */
export interface GeneralFormState {
  publicLiquidityCheck: boolean
  maxLiquidity: string
  reimbursementWindowBlocks: string
  excessTolerance: ExcessTolerance
  rskConfirmations: ConfirmationRow[]
  btcConfirmations: ConfirmationRow[]
}

/** Shared pegin/pegout fields as edited in the form (display-domain strings). */
export interface QuoteFormStateBase {
  timeForDeposit: string
  penaltyFee: string
  fixedFee: string
  fixedFeeEnabled: boolean
  feePercentage: string
  feePercentageEnabled: boolean
  maxValue: string
  minValue: string
}

/** Pegin section as edited in the form (display-domain strings). */
export interface PeginFormState extends QuoteFormStateBase {
  callTime: string
}

/** Pegout section as edited in the form (display-domain strings). */
export interface PegoutFormState extends QuoteFormStateBase {
  expireTime: string
  expireBlocks: string
  bridgeTransactionMin: string
}

interface FormBaseline {
  general: GeneralFormState
  pegin: PeginFormState
  pegout: PegoutFormState
}

/** Validated payload for a single section, ready to POST when `config` is present. */
export interface SectionBuild<T> {
  dirty: boolean
  config: T | null
  errors: string[]
}

/** Result of validating and building all section payloads. */
export interface ConfigurationBuild {
  general: SectionBuild<GeneralConfiguration>
  pegin: SectionBuild<PeginConfiguration>
  pegout: SectionBuild<PegoutConfiguration>
}

export interface UseConfigurationForm {
  general: GeneralFormState
  pegin: PeginFormState
  pegout: PegoutFormState
  updateGeneral: (patch: Partial<GeneralFormState>) => void
  updatePegin: (patch: Partial<PeginFormState>) => void
  updatePegout: (patch: Partial<PegoutFormState>) => void
  dirty: {
    general: boolean
    pegin: boolean
    pegout: boolean
    any: boolean
  }
  build: () => ConfigurationBuild
  markSaved: () => void
}

function feePercentageEnabledFor(value: string): boolean {
  return value.trim() !== '' && value.trim() !== '0'
}

function feeEnabledFor(weiValue: string): boolean {
  return weiValue.trim() !== '' && weiValue.trim() !== '0'
}

function buildGeneralState(config: GeneralConfiguration): GeneralFormState {
  return {
    publicLiquidityCheck: config.publicLiquidityCheck,
    maxLiquidity: weiToEther(config.maxLiquidity),
    reimbursementWindowBlocks: String(config.reimbursementWindowBlocks),
    excessTolerance: {
      isFixed: config.excessTolerance.isFixed,
      fixedValue: config.excessTolerance.fixedValue,
      percentageValue: config.excessTolerance.percentageValue,
    },
    rskConfirmations: confirmationRecordToRows(config.rskConfirmations),
    btcConfirmations: confirmationRecordToRows(config.btcConfirmations),
  }
}

function buildQuoteFormBase(config: QuoteConfigurationBase): QuoteFormStateBase {
  return {
    timeForDeposit: String(config.timeForDeposit),
    penaltyFee: weiToEther(config.penaltyFee),
    fixedFee: weiToEther(config.fixedFee),
    fixedFeeEnabled: feeEnabledFor(config.fixedFee),
    feePercentage: config.feePercentage,
    feePercentageEnabled: feePercentageEnabledFor(config.feePercentage),
    maxValue: weiToEther(config.maxValue),
    minValue: weiToEther(config.minValue),
  }
}

function buildPeginState(config: PeginConfiguration): PeginFormState {
  return {
    ...buildQuoteFormBase(config),
    callTime: String(config.callTime),
  }
}

function buildPegoutState(config: PegoutConfiguration): PegoutFormState {
  return {
    ...buildQuoteFormBase(config),
    expireTime: String(config.expireTime),
    expireBlocks: String(config.expireBlocks),
    bridgeTransactionMin: weiToEther(config.bridgeTransactionMin),
  }
}

/** Parses a numeric field, keeping the raw string when it is not a valid number. */
function toNumericField(raw: string): number | string {
  const trimmed = raw.trim()
  if (trimmed !== '' && !Number.isNaN(Number(trimmed))) {
    return Number(trimmed)
  }
  return trimmed
}

/** Converts a display ether amount to wei, recording an error on failure. */
function feeToWei(display: string, key: string, errors: string[]): string {
  const trimmed = display.trim()
  try {
    return etherToWei(trimmed === '' ? '0' : trimmed)
  } catch {
    errors.push(
      `Invalid input "${display}" for field "${key}". Please enter a valid number.`,
    )
    return '0'
  }
}

function buildExcessTolerance(value: ExcessTolerance): ExcessTolerance {
  if (value.isFixed) {
    return {
      isFixed: true,
      fixedValue: value.fixedValue.trim() === '' ? '0' : value.fixedValue,
      percentageValue: '0',
    }
  }
  const percentage = value.percentageValue.trim()
  return {
    isFixed: false,
    fixedValue: '0',
    percentageValue: percentage === '' ? '0' : percentage,
  }
}

function buildGeneralConfig(
  state: GeneralFormState,
  original: GeneralConfiguration,
): SectionBuild<GeneralConfiguration> {
  const errors: string[] = []
  errors.push(...validateConfirmationRows(state.rskConfirmations, 'rskConfirmations'))
  errors.push(...validateConfirmationRows(state.btcConfirmations, 'btcConfirmations'))

  const maxLiquidityCheck = validateMaxLiquidity(state.maxLiquidity)
  if (!maxLiquidityCheck.isValid && maxLiquidityCheck.error) {
    errors.push(maxLiquidityCheck.error)
  }

  if (state.excessTolerance.isFixed) {
    const fixedCheck = validateExcessToleranceFixed(
      weiToEther(state.excessTolerance.fixedValue),
    )
    if (!fixedCheck.isValid && fixedCheck.error) {
      errors.push(fixedCheck.error)
    }
  } else {
    const percentageCheck = validateExcessTolerancePercentage(
      state.excessTolerance.percentageValue,
    )
    if (!percentageCheck.isValid && percentageCheck.error) {
      errors.push(percentageCheck.error)
    }
  }

  if (errors.length > 0) {
    return { dirty: true, config: null, errors }
  }

  const formatted = formatGeneralConfig({
    publicLiquidityCheck: state.publicLiquidityCheck,
    maxLiquidity: etherToWei(state.maxLiquidity),
    reimbursementWindowBlocks: toNumericField(state.reimbursementWindowBlocks),
    excessTolerance: buildExcessTolerance(state.excessTolerance),
    rskConfirmations: confirmationRowsToEntries(state.rskConfirmations),
    btcConfirmations: confirmationRowsToEntries(state.btcConfirmations),
  })

  const typeCheck = validateConfig(formatted, original as unknown as Record<string, unknown>)
  if (!typeCheck.isValid) {
    return { dirty: true, config: null, errors: typeCheck.errors }
  }

  return {
    dirty: true,
    config: formatted as unknown as GeneralConfiguration,
    errors: [],
  }
}

function buildQuoteFees(
  state: PeginFormState | PegoutFormState,
  errors: string[],
): {
  penaltyFee: string
  fixedFee: string
  feePercentage: string
  maxValue: string
  minValue: string
} {
  const penaltyFee = feeToWei(state.penaltyFee, 'penaltyFee', errors)
  const maxValue = feeToWei(state.maxValue, 'maxValue', errors)
  const minValue = feeToWei(state.minValue, 'minValue', errors)
  const fixedFee = state.fixedFeeEnabled
    ? feeToWei(state.fixedFee, 'fixedFee', errors)
    : '0'

  let feePercentage = '0'
  if (state.feePercentageEnabled) {
    const { value, error } = normalizePercentage(state.feePercentage)
    if (error) {
      errors.push(error)
    } else {
      feePercentage = value
    }
  }

  return { penaltyFee, fixedFee, feePercentage, maxValue, minValue }
}

function buildPeginConfig(
  state: PeginFormState,
  original: PeginConfiguration,
): SectionBuild<PeginConfiguration> {
  const errors: string[] = []
  const fees = buildQuoteFees(state, errors)

  if (errors.length > 0) {
    return { dirty: true, config: null, errors }
  }

  const config = {
    timeForDeposit: toNumericField(state.timeForDeposit),
    callTime: toNumericField(state.callTime),
    ...fees,
  }

  const typeCheck = validateConfig(config, original as unknown as Record<string, unknown>)
  if (!typeCheck.isValid) {
    return { dirty: true, config: null, errors: typeCheck.errors }
  }

  return {
    dirty: true,
    config: config as unknown as PeginConfiguration,
    errors: [],
  }
}

function buildPegoutConfig(
  state: PegoutFormState,
  original: PegoutConfiguration,
): SectionBuild<PegoutConfiguration> {
  const errors: string[] = []
  const fees = buildQuoteFees(state, errors)
  const bridgeTransactionMin = feeToWei(
    state.bridgeTransactionMin,
    'bridgeTransactionMin',
    errors,
  )

  if (errors.length > 0) {
    return { dirty: true, config: null, errors }
  }

  const config = {
    timeForDeposit: toNumericField(state.timeForDeposit),
    expireTime: toNumericField(state.expireTime),
    expireBlocks: toNumericField(state.expireBlocks),
    bridgeTransactionMin,
    ...fees,
  }

  const typeCheck = validateConfig(config, original as unknown as Record<string, unknown>)
  if (!typeCheck.isValid) {
    return { dirty: true, config: null, errors: typeCheck.errors }
  }

  return {
    dirty: true,
    config: config as unknown as PegoutConfiguration,
    errors: [],
  }
}

function sectionEquals(a: unknown, b: unknown): boolean {
  return JSON.stringify(a) === JSON.stringify(b)
}

/**
 * Form state for the configuration card. Loads display-domain values from
 * `data.Configuration`, tracks per-section dirtiness against the loaded
 * baseline, and builds validated wei/percentage payloads for saving.
 */
export function useConfigurationForm(): UseConfigurationForm {
  const { data } = useInitialData()
  const configuration: FullConfiguration = data.Configuration

  const initial = useMemo<FormBaseline>(
    () => ({
      general: buildGeneralState(configuration.general),
      pegin: buildPeginState(configuration.pegin),
      pegout: buildPegoutState(configuration.pegout),
    }),
    [configuration],
  )

  const [general, setGeneral] = useState<GeneralFormState>(initial.general)
  const [pegin, setPegin] = useState<PeginFormState>(initial.pegin)
  const [pegout, setPegout] = useState<PegoutFormState>(initial.pegout)
  const [baseline, setBaseline] = useState<FormBaseline>(initial)

  const updateGeneral = useCallback((patch: Partial<GeneralFormState>) => {
    setGeneral((prev) => ({ ...prev, ...patch }))
  }, [])

  const updatePegin = useCallback((patch: Partial<PeginFormState>) => {
    setPegin((prev) => ({ ...prev, ...patch }))
  }, [])

  const updatePegout = useCallback((patch: Partial<PegoutFormState>) => {
    setPegout((prev) => ({ ...prev, ...patch }))
  }, [])

  const generalDirty = !sectionEquals(general, baseline.general)
  const peginDirty = !sectionEquals(pegin, baseline.pegin)
  const pegoutDirty = !sectionEquals(pegout, baseline.pegout)

  const build = useCallback((): ConfigurationBuild => {
    const generalBuild = buildGeneralConfig(general, configuration.general)
    const peginBuild = buildPeginConfig(pegin, configuration.pegin)
    const pegoutBuild = buildPegoutConfig(pegout, configuration.pegout)
    return {
      general: { ...generalBuild, dirty: generalDirty },
      pegin: { ...peginBuild, dirty: peginDirty },
      pegout: { ...pegoutBuild, dirty: pegoutDirty },
    }
  }, [
    general,
    pegin,
    pegout,
    configuration,
    generalDirty,
    peginDirty,
    pegoutDirty,
  ])

  const markSaved = useCallback(() => {
    setBaseline({ general, pegin, pegout })
  }, [general, pegin, pegout])

  return {
    general,
    pegin,
    pegout,
    updateGeneral,
    updatePegin,
    updatePegout,
    dirty: {
      general: generalDirty,
      pegin: peginDirty,
      pegout: pegoutDirty,
      any: generalDirty || peginDirty || pegoutDirty,
    },
    build,
    markSaved,
  }
}
