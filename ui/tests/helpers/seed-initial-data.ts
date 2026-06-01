import type { InitialDataPayload } from '@shared/types/initial-data'
import { resetInitialDataCacheForTests } from '@shared/utils/initial-data'

export function seedInitialData(
  payload: InitialDataPayload,
  options?: { csrfToken?: string },
): void {
  if (options?.csrfToken) {
    document.head.innerHTML = `<meta name="csrf-token" content="${options.csrfToken}" />`
  }

  document.body.innerHTML = `<script id="initial-data" type="application/json">${JSON.stringify(payload)}</script>`
  resetInitialDataCacheForTests()
}
