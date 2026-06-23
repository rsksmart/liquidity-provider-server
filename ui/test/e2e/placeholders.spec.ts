import fs from 'node:fs'

import { expect, test } from '../fixtures'
import { MANAGEMENT_STORAGE_STATE_PATH } from '../fixtures/session'

test.describe('placeholder pages', () => {
  test('login route shows Login heading when logged out', async ({ page }) => {
    await page.goto('/login')
    await expect(page.getByRole('heading', { level: 1, name: 'Login' })).toBeVisible()
  })

  test('error route shows Error heading', async ({ page }) => {
    await page.goto('/error')
    await expect(page.getByRole('heading', { level: 1, name: 'Error' })).toBeVisible()
  })

  test.describe('authenticated', () => {
    test.use({ storageState: MANAGEMENT_STORAGE_STATE_PATH })

    test.beforeEach(() => {
      test.skip(
        !process.env.LPS_E2E_USER?.trim() ||
          !process.env.LPS_E2E_PASSWORD?.trim() ||
          !fs.existsSync(MANAGEMENT_STORAGE_STATE_PATH),
        'requires LPS_E2E_USER, LPS_E2E_PASSWORD, and session.setup',
      )
    })

    test('management route shows Management heading when logged in', async ({ page }) => {
      await page.goto('/management')
      await expect(page.getByRole('heading', { level: 1, name: 'Management' })).toBeVisible()
    })
  })
})
