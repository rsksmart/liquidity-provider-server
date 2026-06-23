import { test as setup } from '@playwright/test'

import { seedManagementSession } from '../fixtures/session'

setup('authenticate management session', async ({ request }) => {
  if (!process.env.LPS_E2E_USER?.trim() || !process.env.LPS_E2E_PASSWORD?.trim()) {
    setup.skip(true, 'LPS_E2E_USER and LPS_E2E_PASSWORD not set — authenticated specs will skip')
    return
  }

  await seedManagementSession(request)
})
