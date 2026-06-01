import type { InitialDataPayload } from '@shared/types/initial-data'
import { resetInitialDataCacheForTests } from '@shared/utils/initial-data'
import { render, screen } from '@testing-library/react'
import { loggedInFixture } from '@tests/fixtures/logged-in'
import { loggedOutFixture } from '@tests/fixtures/logged-out'
import { seedInitialData } from '@tests/helpers/seed-initial-data'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it } from 'vitest'

import { App } from '@/App'

const appBasename = '/management/next'
const rootPath = `${appBasename}/`
const errorPath = `${appBasename}/error`

function renderApp(initialPath: string, payload: InitialDataPayload) {
  seedInitialData(payload)
  return render(
    // MemoryRouter needs a fresh initialEntries per call; no re-render loop in one-shot test renders.
    // eslint-disable-next-line react-perf/jsx-no-new-array-as-prop -- perf rule targets production re-renders, not RTL setup
    <MemoryRouter basename={appBasename} initialEntries={[initialPath]}>
      <App />
    </MemoryRouter>,
  )
}

describe('App router', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    resetInitialDataCacheForTests()
  })

  it('boots to login when loggedIn is false', () => {
    renderApp(rootPath, loggedOutFixture)
    expect(screen.getByRole('heading', { name: /login/i })).toBeInTheDocument()
  })

  it('boots to management when loggedIn is true', () => {
    renderApp(rootPath, loggedInFixture)
    expect(screen.getByRole('heading', { name: /management/i })).toBeInTheDocument()
  })

  it('renders the error placeholder route', () => {
    renderApp(errorPath, loggedOutFixture)
    expect(screen.getByRole('heading', { name: /error/i })).toBeInTheDocument()
  })
})
