import { expect, test } from '../fixtures'
import { applyFreshManagementSession } from '../fixtures/session'

test.describe('placeholder pages', () => {
  test('login route shows Login heading when logged out', async ({ page }) => {
    await page.goto('login')
    await expect(page.getByRole('heading', { level: 1, name: 'Login' })).toBeVisible()
  })

  test('error route shows Error heading', async ({ page }) => {
    await page.goto('error')
    await expect(page.getByRole('heading', { level: 1, name: 'Error' })).toBeVisible()
  })

  test.describe('authenticated', () => {
    test.beforeEach(async ({ request, context }) => {
      test.skip(
        !process.env.LPS_E2E_USER?.trim() || !process.env.LPS_E2E_PASSWORD?.trim(),
        'requires LPS_E2E_USER and LPS_E2E_PASSWORD',
      )
      await applyFreshManagementSession(request, context)
    })

    test('management route shows dashboard when logged in', async ({ page }) => {
      await page.goto('management')
      await expect(page.getByRole('heading', { level: 1, name: 'Management Dashboard' })).toBeVisible()
      await expect(page.getByText('Provider RSK Address')).toBeVisible()
      await expect(page.getByRole('tab', { name: 'Pegin' })).toBeVisible()
    })
  })
})
