import { apiFetch, isSessionExpiredError } from '@api/management/utils/api-fetch'
import { toast } from 'sonner'

import { Button } from '@/components/ui/button'

function logout(): void {
  void (async () => {
    let sessionExpired = false
    try {
      await apiFetch('/management/logout', { method: 'POST' })
    } catch (err) {
      if (isSessionExpiredError(err)) {
        sessionExpired = true
      } else {
        console.error('Logout failed', err)
        toast.warning('Logout failed; local session cleared.')
      }
    } finally {
      if (!sessionExpired) {
        window.location.assign('/management/next/login')
      }
    }
  })()
}

export function LogoutButton() {
  return (
    <Button type="button" onClick={logout}>
      Logout
    </Button>
  )
}
