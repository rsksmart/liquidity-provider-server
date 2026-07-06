import { AddCollateralForm } from '@feature/management/components/AddCollateralForm'
import {
  managementCardClass,
  managementCardContentClass,
  managementCardHeaderClass,
  managementCardTitleClass,
  managementTabTriggerClass,
} from '@feature/management/management-styles'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

export function CollateralCard() {
  return (
    <Card className={managementCardClass}>
      <CardHeader className={managementCardHeaderClass}>
        <CardTitle className={managementCardTitleClass}>Collateral</CardTitle>
      </CardHeader>
      <CardContent className={managementCardContentClass}>
        <Tabs defaultValue="pegin">
          <TabsList
            variant="line"
            className="mb-0 h-auto w-full justify-start gap-0 rounded-none border-b border-[#dee2e6] bg-transparent p-0"
          >
            <TabsTrigger value="pegin" className={managementTabTriggerClass}>
              Pegin
            </TabsTrigger>
            <TabsTrigger value="pegout" className={managementTabTriggerClass}>
              Pegout
            </TabsTrigger>
          </TabsList>
          <TabsContent value="pegin" keepMounted className="pt-3 text-base">
            <AddCollateralForm kind="pegin" />
          </TabsContent>
          <TabsContent value="pegout" keepMounted className="pt-3 text-base">
            <AddCollateralForm kind="pegout" />
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  )
}
