import { expect, test } from '../fixtures'
import {
  applyFreshManagementSession,
  clearSessionCookie,
  extractCsrfFromResponseHtml,
  getServerOrigin,
  getShellUrl,
} from '../fixtures/session'

const hasE2ECredentials =
  Boolean(process.env.LPS_E2E_USER?.trim()) && Boolean(process.env.LPS_E2E_PASSWORD?.trim())

async function getCsrfToken(
  request: import('@playwright/test').APIRequestContext,
): Promise<string> {
  const shellResponse = await request.get(getShellUrl())
  expect(shellResponse.ok()).toBeTruthy()
  return extractCsrfFromResponseHtml(await shellResponse.text())
}

test.describe('auth login flows', () => {
  test('happy path — login redirects to management', async ({ page }) => {
    test.skip(!hasE2ECredentials, 'requires LPS_E2E_USER and LPS_E2E_PASSWORD')

    await page.goto('login')
    await page.getByTestId('login-username-input').fill(process.env.LPS_E2E_USER!)
    await page.getByTestId('login-password-input').fill(process.env.LPS_E2E_PASSWORD!)
    await page.getByTestId('login-submit-button').click()

    await expect(page).toHaveURL(/\/management\/next\/management$/)
    await expect(page.getByRole('heading', { level: 1, name: 'Management' })).toBeVisible()
  })

  test('bad credentials show legacy error banner', async ({ page }) => {
    await page.goto('login')
    await page.getByTestId('login-username-input').fill('not-a-real-user')
    await page.getByTestId('login-password-input').fill('not-a-real-password')
    await page.getByTestId('login-submit-button').click()

    await expect(page.getByRole('alert')).toHaveText('Invalid username or password.')
    await expect(page).toHaveURL(/\/management\/next\/login$/)
  })

  test('network error on login shows legacy error banner', async ({ page }) => {
    await page.route('**/management/login', (route) => route.abort())

    await page.goto('login')
    await page.getByTestId('login-username-input').fill('user')
    await page.getByTestId('login-password-input').fill('pass')
    await page.getByTestId('login-submit-button').click()

    await expect(page.getByRole('alert')).toHaveText('Invalid username or password.')
  })

  test('CSRF rejection surfaces an error', async ({ page, request }) => {
    const origin = getServerOrigin()
    const csrfToken = await getCsrfToken(request)

    await page.route('**/management/login', async (route) => {
      const postData = route.request().postData()
      await route.continue({
        headers: {
          ...route.request().headers(),
          'X-CSRF-Token': 'invalid-token',
        },
        postData,
      })
    })

    await page.goto('login')
    await page.getByTestId('login-username-input').fill('user')
    await page.getByTestId('login-password-input').fill('pass')
    await page.getByTestId('login-submit-button').click()

    await expect(page.getByRole('alert')).toBeVisible()

    const directResponse = await request.post(`${origin}/management/login`, {
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': 'invalid-token',
      },
      data: { username: 'user', password: 'pass' },
    })
    expect(directResponse.status()).toBe(403)
    expect(csrfToken.length).toBeGreaterThan(0)
  })

  test('successful login sets lp-session cookie flags', async ({ request }) => {
    test.skip(!hasE2ECredentials, 'requires LPS_E2E_USER and LPS_E2E_PASSWORD')

    const origin = getServerOrigin()
    const csrfToken = await getCsrfToken(request)
    const loginResponse = await request.post(`${origin}/management/login`, {
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': csrfToken,
      },
      data: {
        username: process.env.LPS_E2E_USER,
        password: process.env.LPS_E2E_PASSWORD,
      },
    })

    expect(loginResponse.ok()).toBeTruthy()
    const setCookie = loginResponse.headers()['set-cookie'] ?? ''
    expect(setCookie).toMatch(/lp-session/i)
    expect(setCookie).toMatch(/HttpOnly/i)
    expect(setCookie).toMatch(/SameSite=Strict/i)

    const baseUrl = process.env.LPS_E2E_BASE_URL ?? 'http://localhost:8080/management/next'
    if (baseUrl.startsWith('https://')) {
      expect(setCookie).toMatch(/Secure/i)
    }
  })

  test.describe('authenticated', () => {
    test.describe.configure({ mode: 'serial' })

    test.beforeEach(async ({ request, context }) => {
      test.skip(
        !hasE2ECredentials,
        'requires LPS_E2E_USER and LPS_E2E_PASSWORD',
      )
      await applyFreshManagementSession(request, context)
    })

    test('session expired — toast and redirect to login', async ({ page, context }) => {
      await page.goto('management')
      await expect(page.getByRole('heading', { level: 1, name: 'Management' })).toBeVisible()

      await clearSessionCookie(context)
      await page.getByRole('button', { name: 'Logout' }).click()

      await expect(page.getByText('Your session has expired. Please log in again.')).toBeVisible()
      await expect(page).toHaveURL(/\/management\/next\/login$/)
    })

    test('logout clears session and redirects to login', async ({ page }) => {
      await page.goto('management')
      await page.getByRole('button', { name: 'Logout' }).click()
      await expect(page).toHaveURL(/\/management\/next\/login$/)

      const origin = getServerOrigin()
      const statusResponse = await page.request.get(`${origin}/management/trusted-accounts`)
      expect(statusResponse.status()).toBe(403)
    })
  })
})
