export {
  applyFreshManagementSession,
  clearSessionCookie,
  getServerOrigin,
  getShellUrl,
  MANAGEMENT_STORAGE_STATE_PATH,
  seedManagementSession,
} from './session'
import { expect, test as base } from '@playwright/test'

export const test = base.extend({
  page: async ({ page }, use) => {
    const goto = page.goto.bind(page)
    page.goto = (url, options) => goto(url, { waitUntil: 'networkidle', ...options })
    await use(page)
  },
})

export { expect }
