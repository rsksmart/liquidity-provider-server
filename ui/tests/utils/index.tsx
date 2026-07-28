import type { InitialDataPayload } from '@shared/types/initial-data'
import { replaceInitialDataPayload } from '@shared/utils/initial-data'
import { render } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'

import { App } from '@/App'

export const appBasename = '/management/next'

export function seedInitialData(
  payload: InitialDataPayload,
  options?: { csrfToken?: string },
): void {
  document.querySelector('meta[name="csrf-token"]')?.remove()
  if (options?.csrfToken) {
    document.head.insertAdjacentHTML(
      'beforeend',
      `<meta name="csrf-token" content="${options.csrfToken}" />`,
    )
  }

  if (!document.getElementById('initial-data')) {
    document.body.insertAdjacentHTML(
      'beforeend',
      '<script id="initial-data" type="application/json"></script>',
    )
  }
  replaceInitialDataPayload(payload)
}

/**
 * Production uses BrowserRouter in main.tsx with the real browser URL.
 * Tests use MemoryRouter with initialEntries so we can assert boot/redirect paths
 * without mutating window.location — same App route tree and basename only.
 */
export function renderAppAt(initialPath: string, payload: InitialDataPayload) {
  seedInitialData(payload)
  return render(
    <MemoryRouter basename={appBasename} initialEntries={[initialPath]}>
      <App />
    </MemoryRouter>,
  )
}
