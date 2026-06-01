import { AuthGuard } from '@feature/auth/components/AuthGuard'
import { LoginPlaceholder } from '@feature/auth/components/LoginPlaceholder'
import { ErrorPlaceholder } from '@feature/error/components/ErrorPlaceholder'
import { ManagementPlaceholder } from '@feature/management/components/ManagementPlaceholder'
import { useInitialData } from '@shared/utils/initial-data'
import { Navigate, Route, Routes } from 'react-router-dom'

export function App() {
  const { loggedIn } = useInitialData()
  const defaultRoute = loggedIn ? '/management' : '/login'

  return (
    <Routes>
      <Route
        path="/login"
        element={
          <AuthGuard requireAuth={false}>
            <LoginPlaceholder />
          </AuthGuard>
        }
      />
      <Route
        path="/management"
        element={
          <AuthGuard requireAuth={true}>
            <ManagementPlaceholder />
          </AuthGuard>
        }
      />
      <Route path="/error" element={<ErrorPlaceholder />} />
      <Route path="/" element={<Navigate to={defaultRoute} replace />} />
      <Route path="*" element={<Navigate to={defaultRoute} replace />} />
    </Routes>
  )
}
