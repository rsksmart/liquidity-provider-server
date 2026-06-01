import { AuthGuard } from '@feature/auth/components/AuthGuard'
import { resetInitialDataCacheForTests } from '@shared/utils/initial-data'
import { render, screen } from '@testing-library/react'
import { loggedInFixture } from '@tests/fixtures/logged-in'
import { loggedOutFixture } from '@tests/fixtures/logged-out'
import { seedInitialData } from '@tests/helpers/seed-initial-data'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { beforeEach, describe, expect, it } from 'vitest'

const loginRoute = '/login'
const managementRoute = '/management'

function renderGuardedRoute(path: string, requireAuth: boolean, label: string) {
  return render(
    // MemoryRouter needs a fresh initialEntries per call; no re-render loop in one-shot test renders.
    // eslint-disable-next-line react-perf/jsx-no-new-array-as-prop -- perf rule targets production re-renders, not RTL setup
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route
          path={path}
          element={
            <AuthGuard requireAuth={requireAuth}>
              <div>{label}</div>
            </AuthGuard>
          }
        />
        <Route path={loginRoute} element={<div>Login page</div>} />
        <Route path={managementRoute} element={<div>Management page</div>} />
      </Routes>
    </MemoryRouter>,
  )
}

describe('AuthGuard', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
    resetInitialDataCacheForTests()
  })

  it('redirects unauthenticated users away from /management', () => {
    seedInitialData(loggedOutFixture)
    renderGuardedRoute(managementRoute, true, 'Protected')

    expect(screen.getByText('Login page')).toBeInTheDocument()
    expect(screen.queryByText('Protected')).not.toBeInTheDocument()
  })

  it('redirects authenticated users away from /login', () => {
    seedInitialData(loggedInFixture)
    renderGuardedRoute(loginRoute, false, 'Login form')

    expect(screen.getByText('Management page')).toBeInTheDocument()
    expect(screen.queryByText('Login form')).not.toBeInTheDocument()
  })

  it('does not redirect when auth state matches the guarded route', () => {
    seedInitialData(loggedInFixture)
    renderGuardedRoute(managementRoute, true, 'Protected')

    expect(screen.getByText('Protected')).toBeInTheDocument()
  })

  it('performs a single redirect hop without looping', () => {
    seedInitialData(loggedOutFixture)
    renderGuardedRoute(managementRoute, true, 'Protected')

    expect(screen.getByText('Login page')).toBeInTheDocument()
    expect(screen.queryByText('Protected')).not.toBeInTheDocument()
    expect(screen.queryByText('Management page')).not.toBeInTheDocument()
  })
})
