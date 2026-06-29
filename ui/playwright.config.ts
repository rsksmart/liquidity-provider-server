import { defineConfig, devices } from '@playwright/test'

import { getShellUrl } from './test/fixtures/session'

const hasE2ECredentials =
  Boolean(process.env.LPS_E2E_USER?.trim()) && Boolean(process.env.LPS_E2E_PASSWORD?.trim())

export default defineConfig({
  testDir: './test/e2e',
  fullyParallel: !hasE2ECredentials,
  workers: hasE2ECredentials ? 1 : undefined,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: getShellUrl(),
    trace: 'on-first-retry',
    launchOptions: {
      slowMo: Number(process.env.PLAYWRIGHT_SLOW_MO ?? 0),
    },
  },
  projects: [
    { name: 'setup', testMatch: /session\.setup\.ts/ },
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
      dependencies: ['setup'],
      testIgnore: /session\.setup\.ts/,
    },
  ],
})
