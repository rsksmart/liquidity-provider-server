import { ApiFetchError } from '@api/management/types/errors'
import { apiFetch } from '@api/management/utils/api-fetch'
import {
  CheckboxInput,
  type ConfirmationRow,
  ConfirmationTiersEditor,
  DefaultInput,
  ExcessToleranceInput,
  FeeInput,
  ToggableFeeInput,
} from '@feature/management/components/config'
import type { SectionPrefix } from '@feature/management/components/config/types'
import { checkFeeWarnings } from '@feature/management/config/fee-warnings'
import { useConfigurationForm } from '@feature/management/hooks/use-configuration-form'
import {
  managementBootstrapButtonClass,
  managementCardClass,
  managementCardContentClass,
  managementCardHeaderClass,
  managementCardTitleClass,
  managementFieldTitleClass,
  managementTabTriggerClass,
} from '@feature/management/management-styles'
import type { ExcessTolerance } from '@shared/types/initial-data'
import { useCallback, useMemo, useState } from 'react'
import { toast } from 'sonner'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { cn } from '@/lib/utils'

/** Extracts a human-readable message from a section save failure. */
function apiErrorMessage(error: unknown): string {
  if (error instanceof ApiFetchError) {
    const body = error.body
    if (
      body !== null &&
      typeof body === 'object' &&
      'message' in body &&
      typeof (body).message === 'string'
    ) {
      return (body as { message: string }).message
    }
    return error.message
  }
  return error instanceof Error ? error.message : 'Unknown error'
}

export function ConfigurationCard() {
  const form = useConfigurationForm()
  const {
    general,
    pegin,
    pegout,
    updateGeneral,
    updatePegin,
    updatePegout,
    dirty,
  } = form

  const [activeTab, setActiveTab] = useState<SectionPrefix>('general')
  const [saving, setSaving] = useState(false)

  const handleTabChange = useCallback((value: string | number | null) => {
    if (value === 'general' || value === 'pegin' || value === 'pegout') {
      setActiveTab(value)
    }
  }, [])

  const generalHandlers = useMemo(
    () => ({
      rskConfirmations: (rskConfirmations: ConfirmationRow[]) => {
        updateGeneral({ rskConfirmations })
      },
      btcConfirmations: (btcConfirmations: ConfirmationRow[]) => {
        updateGeneral({ btcConfirmations })
      },
      publicLiquidityCheck: (publicLiquidityCheck: boolean) => {
        updateGeneral({ publicLiquidityCheck })
      },
      maxLiquidity: (maxLiquidity: string) => {
        updateGeneral({ maxLiquidity })
      },
      reimbursementWindowBlocks: (reimbursementWindowBlocks: string) => {
        updateGeneral({ reimbursementWindowBlocks })
      },
      excessTolerance: (excessTolerance: ExcessTolerance) => {
        updateGeneral({ excessTolerance })
      },
    }),
    [updateGeneral],
  )

  const peginHandlers = useMemo(
    () => ({
      timeForDeposit: (timeForDeposit: string) => {
        updatePegin({ timeForDeposit })
      },
      callTime: (callTime: string) => {
        updatePegin({ callTime })
      },
      penaltyFee: (penaltyFee: string) => {
        updatePegin({ penaltyFee })
      },
      fixedFee: (fixedFee: string) => {
        updatePegin({ fixedFee })
      },
      fixedFeeEnabled: (fixedFeeEnabled: boolean) => {
        updatePegin({ fixedFeeEnabled })
      },
      feePercentage: (feePercentage: string) => {
        updatePegin({ feePercentage })
      },
      feePercentageEnabled: (feePercentageEnabled: boolean) => {
        updatePegin({ feePercentageEnabled })
      },
      maxValue: (maxValue: string) => {
        updatePegin({ maxValue })
      },
      minValue: (minValue: string) => {
        updatePegin({ minValue })
      },
    }),
    [updatePegin],
  )

  const pegoutHandlers = useMemo(
    () => ({
      timeForDeposit: (timeForDeposit: string) => {
        updatePegout({ timeForDeposit })
      },
      expireTime: (expireTime: string) => {
        updatePegout({ expireTime })
      },
      penaltyFee: (penaltyFee: string) => {
        updatePegout({ penaltyFee })
      },
      fixedFee: (fixedFee: string) => {
        updatePegout({ fixedFee })
      },
      fixedFeeEnabled: (fixedFeeEnabled: boolean) => {
        updatePegout({ fixedFeeEnabled })
      },
      feePercentage: (feePercentage: string) => {
        updatePegout({ feePercentage })
      },
      feePercentageEnabled: (feePercentageEnabled: boolean) => {
        updatePegout({ feePercentageEnabled })
      },
      maxValue: (maxValue: string) => {
        updatePegout({ maxValue })
      },
      minValue: (minValue: string) => {
        updatePegout({ minValue })
      },
      expireBlocks: (expireBlocks: string) => {
        updatePegout({ expireBlocks })
      },
      bridgeTransactionMin: (bridgeTransactionMin: string) => {
        updatePegout({ bridgeTransactionMin })
      },
    }),
    [updatePegout],
  )

  const handleSave = useCallback(async () => {
    if (saving) {
      return
    }

    const toggles =
      activeTab === 'pegin'
        ? {
            fixedFeeEnabled: pegin.fixedFeeEnabled,
            feePercentageEnabled: pegin.feePercentageEnabled,
          }
        : activeTab === 'pegout'
          ? {
              fixedFeeEnabled: pegout.fixedFeeEnabled,
              feePercentageEnabled: pegout.feePercentageEnabled,
            }
          : { fixedFeeEnabled: true, feePercentageEnabled: true }

    const warning = checkFeeWarnings(activeTab, toggles)
    if (warning.shouldWarn && warning.message) {
      toast.warning(warning.message)
    }

    const result = form.build()
    setSaving(true)
    let success = true

    if (result.general.dirty) {
      if (result.general.errors.length > 0) {
        toast.error(result.general.errors.join('\n'))
        success = false
      } else if (result.general.config) {
        try {
          await apiFetch.post('/configuration', {
            configuration: result.general.config,
          })
        } catch (error) {
          toast.error(apiErrorMessage(error))
          success = false
        }
      }
    }

    if (result.pegin.dirty) {
      if (result.pegin.errors.length > 0) {
        toast.error(result.pegin.errors.join('\n'))
        success = false
      } else if (result.pegin.config) {
        try {
          await apiFetch.post('/pegin/configuration', {
            configuration: result.pegin.config,
          })
        } catch (error) {
          toast.error(apiErrorMessage(error))
          success = false
        }
      }
    }

    if (result.pegout.dirty) {
      if (result.pegout.errors.length > 0) {
        toast.error(result.pegout.errors.join('\n'))
        success = false
      } else if (result.pegout.config) {
        try {
          await apiFetch.post('/pegout/configuration', {
            configuration: result.pegout.config,
          })
        } catch (error) {
          toast.error(apiErrorMessage(error))
          success = false
        }
      }
    }

    if (success) {
      toast.success('Configuration saved successfully!')
      form.markSaved()
    }

    setSaving(false)
  }, [activeTab, form, pegout, pegin, saving])

  const handleSaveClick = useCallback(() => {
    void handleSave()
  }, [handleSave])

  return (
    <Card
      className={cn(managementCardClass, 'min-w-0')}
      data-testid="configuration-card"
    >
      <CardHeader className={managementCardHeaderClass}>
        <CardTitle className={managementCardTitleClass}>Configuration</CardTitle>
      </CardHeader>
      <CardContent className={managementCardContentClass}>
        <h2 className={cn(managementFieldTitleClass, 'mb-3')}>
          Current Configuration
        </h2>

        <Tabs value={activeTab} onValueChange={handleTabChange}>
          <TabsList
            variant="line"
            className="mb-0 h-auto w-full justify-start gap-0 rounded-none border-b border-[#dee2e6] bg-transparent p-0"
          >
            <TabsTrigger
              value="general"
              className={managementTabTriggerClass}
              data-testid="config-tab-general"
            >
              General
            </TabsTrigger>
            <TabsTrigger
              value="pegin"
              className={managementTabTriggerClass}
              data-testid="config-tab-pegin"
            >
              Pegin
            </TabsTrigger>
            <TabsTrigger
              value="pegout"
              className={managementTabTriggerClass}
              data-testid="config-tab-pegout"
            >
              Pegout
            </TabsTrigger>
          </TabsList>

          <TabsContent value="general" keepMounted className="pt-3 text-base">
            <ConfirmationTiersEditor
              configKey="rskConfirmations"
              value={general.rskConfirmations}
              onChange={generalHandlers.rskConfirmations}
            />
            <ConfirmationTiersEditor
              configKey="btcConfirmations"
              value={general.btcConfirmations}
              onChange={generalHandlers.btcConfirmations}
            />
            <CheckboxInput
              sectionPrefix="general"
              fieldKey="publicLiquidityCheck"
              value={general.publicLiquidityCheck}
              onChange={generalHandlers.publicLiquidityCheck}
            />
            <FeeInput
              sectionPrefix="general"
              fieldKey="maxLiquidity"
              value={general.maxLiquidity}
              onChange={generalHandlers.maxLiquidity}
            />
            <DefaultInput
              sectionPrefix="general"
              fieldKey="reimbursementWindowBlocks"
              value={general.reimbursementWindowBlocks}
              onChange={generalHandlers.reimbursementWindowBlocks}
            />
            <ExcessToleranceInput
              sectionPrefix="general"
              value={general.excessTolerance}
              onChange={generalHandlers.excessTolerance}
            />
          </TabsContent>

          <TabsContent value="pegin" keepMounted className="pt-3 text-base">
            <DefaultInput
              sectionPrefix="pegin"
              fieldKey="timeForDeposit"
              value={pegin.timeForDeposit}
              onChange={peginHandlers.timeForDeposit}
            />
            <DefaultInput
              sectionPrefix="pegin"
              fieldKey="callTime"
              value={pegin.callTime}
              onChange={peginHandlers.callTime}
            />
            <FeeInput
              sectionPrefix="pegin"
              fieldKey="penaltyFee"
              value={pegin.penaltyFee}
              onChange={peginHandlers.penaltyFee}
            />
            <ToggableFeeInput
              sectionPrefix="pegin"
              fieldKey="fixedFee"
              value={pegin.fixedFee}
              enabled={pegin.fixedFeeEnabled}
              onChange={peginHandlers.fixedFee}
              onEnabledChange={peginHandlers.fixedFeeEnabled}
            />
            <ToggableFeeInput
              sectionPrefix="pegin"
              fieldKey="feePercentage"
              value={pegin.feePercentage}
              enabled={pegin.feePercentageEnabled}
              onChange={peginHandlers.feePercentage}
              onEnabledChange={peginHandlers.feePercentageEnabled}
            />
            <FeeInput
              sectionPrefix="pegin"
              fieldKey="maxValue"
              value={pegin.maxValue}
              onChange={peginHandlers.maxValue}
            />
            <FeeInput
              sectionPrefix="pegin"
              fieldKey="minValue"
              value={pegin.minValue}
              onChange={peginHandlers.minValue}
            />
          </TabsContent>

          <TabsContent value="pegout" keepMounted className="pt-3 text-base">
            <DefaultInput
              sectionPrefix="pegout"
              fieldKey="timeForDeposit"
              value={pegout.timeForDeposit}
              onChange={pegoutHandlers.timeForDeposit}
            />
            <DefaultInput
              sectionPrefix="pegout"
              fieldKey="expireTime"
              value={pegout.expireTime}
              onChange={pegoutHandlers.expireTime}
            />
            <FeeInput
              sectionPrefix="pegout"
              fieldKey="penaltyFee"
              value={pegout.penaltyFee}
              onChange={pegoutHandlers.penaltyFee}
            />
            <ToggableFeeInput
              sectionPrefix="pegout"
              fieldKey="fixedFee"
              value={pegout.fixedFee}
              enabled={pegout.fixedFeeEnabled}
              onChange={pegoutHandlers.fixedFee}
              onEnabledChange={pegoutHandlers.fixedFeeEnabled}
            />
            <ToggableFeeInput
              sectionPrefix="pegout"
              fieldKey="feePercentage"
              value={pegout.feePercentage}
              enabled={pegout.feePercentageEnabled}
              onChange={pegoutHandlers.feePercentage}
              onEnabledChange={pegoutHandlers.feePercentageEnabled}
            />
            <FeeInput
              sectionPrefix="pegout"
              fieldKey="maxValue"
              value={pegout.maxValue}
              onChange={pegoutHandlers.maxValue}
            />
            <FeeInput
              sectionPrefix="pegout"
              fieldKey="minValue"
              value={pegout.minValue}
              onChange={pegoutHandlers.minValue}
            />
            <DefaultInput
              sectionPrefix="pegout"
              fieldKey="expireBlocks"
              value={pegout.expireBlocks}
              onChange={pegoutHandlers.expireBlocks}
            />
            <FeeInput
              sectionPrefix="pegout"
              fieldKey="bridgeTransactionMin"
              value={pegout.bridgeTransactionMin}
              onChange={pegoutHandlers.bridgeTransactionMin}
            />
          </TabsContent>
        </Tabs>

        <div className="mt-3">
          <Button
            type="button"
            variant="bootstrap"
            className={managementBootstrapButtonClass}
            data-testid="config-save-button"
            disabled={!dirty.any || saving}
            onClick={handleSaveClick}
          >
            Save Configuration
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
