import { useInitialData } from '@shared/utils/initial-data'
import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'

interface AuthGuardProps {
  requireAuth: boolean
  children: ReactNode
}

// Server session + csrf.Protect remain the authoritative access boundary.
// This guard only routes placeholder screens for operator UX.
export function AuthGuard({ requireAuth, children }: AuthGuardProps) {
  const { loggedIn } = useInitialData()

  if (requireAuth && !loggedIn) {
    return <Navigate to="/login" replace />
  }

  if (!requireAuth && loggedIn) {
    return <Navigate to="/management" replace />
  }

  return <>{children}</>
}
