import type { InitialDataPayload } from '@shared/types/initial-data'
import { resetInitialDataCacheForTests } from '@shared/utils/initial-data'
import { render } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'

import { App } from '@/App'

export const appBasename = '/management/next'

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

/**
 * Production uses BrowserRouter in main.tsx with the real browser URL.
 * Tests use MemoryRouter with initialEntries so we can assert boot/redirect paths
 * without mutating window.location — same App route tree and basename only.
 */
export function renderAppAt(initialPath: string, payload: InitialDataPayload) {
  seedInitialData(payload)
  return render(
    // eslint-disable-next-line react-perf/jsx-no-new-array-as-prop -- RTL setup, not production re-renders
    <MemoryRouter basename={appBasename} initialEntries={[initialPath]}>
      <App />
    </MemoryRouter>,
  )
}
