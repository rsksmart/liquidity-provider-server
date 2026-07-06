import {
  managementCardClass,
  managementCardContentClass,
  managementCardHeaderClass,
  managementCardTitleClass,
  managementFieldTextClass,
  managementFieldTitleClass,
} from '@feature/management/management-styles'
import { useInitialData } from '@shared/utils/initial-data'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'

export function ProviderCard() {
  const { data } = useInitialData()
  const operationalLabel = data.ProviderData.status ? 'Operational' : 'Not Operational'

  return (
    <Card className={managementCardClass}>
      <CardHeader className={managementCardHeaderClass}>
        <CardTitle className={managementCardTitleClass}>Provider</CardTitle>
      </CardHeader>
      <CardContent className={cn(managementCardContentClass, 'space-y-3')}>
        <div>
          <h2 className={managementFieldTitleClass}>Provider RSK Address</h2>
          <p className={managementFieldTextClass} data-testid="provider-rsk-address">
            {data.RskAddress}
          </p>
        </div>
        <div>
          <h2 className={managementFieldTitleClass}>Provider BTC Address</h2>
          <p className={managementFieldTextClass} data-testid="provider-btc-address">
            {data.BtcAddress}
          </p>
        </div>
        <div>
          <h2 className={managementFieldTitleClass}>Operational Status</h2>
          <p className={managementFieldTextClass} data-testid="provider-operational-status">
            {operationalLabel}
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
