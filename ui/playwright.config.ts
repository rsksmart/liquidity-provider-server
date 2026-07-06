import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './test/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: process.env.LPS_E2E_BASE_URL ?? 'http://localhost:8080/management/next/',
    trace: 'on-first-retry',
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
