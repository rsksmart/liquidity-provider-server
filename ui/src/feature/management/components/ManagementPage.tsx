import { LogoutButton } from '@feature/auth/components/LogoutButton'
import {
  CollateralCard,
  ConfigurationCard,
  ProviderCard,
  TrustedAccountsCard,
} from '@feature/management/components'
import {
  managementCardsColumnClass,
  managementCardsGridClass,
  managementPageTitleClass,
  managementShellClass,
} from '@feature/management/management-styles'

export function ManagementPage() {
  return (
    <main className={managementShellClass}>
      <div className="relative mb-2 pr-28">
        <h1 className={managementPageTitleClass}>Management Dashboard</h1>
        <LogoutButton className="absolute top-0 right-0" />
      </div>
      <hr className="mb-4 border-[#dee2e6]" />
      <div className={managementCardsGridClass}>
        <div
          className={managementCardsColumnClass}
          data-testid="management-cards-column"
        >
          <ProviderCard />
          <CollateralCard />
          <TrustedAccountsCard />
        </div>
        <ConfigurationCard />
      </div>
    </main>
  )
}
