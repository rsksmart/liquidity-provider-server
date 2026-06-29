import { apiFetch, isSessionExpiredError } from '@api/management/utils/api-fetch'
import { toast } from 'sonner'

export async function logout(): Promise<void> {
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
}
