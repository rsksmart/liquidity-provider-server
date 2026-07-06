import {
  apiFetch,
  isSessionExpiredError,
  SESSION_EXPIRED_REDIRECT_DELAY_MS,
} from '@api/management/utils/api-fetch'
import { toast } from 'sonner'

export async function logout(): Promise<void> {
  let sessionExpired = false
  let logoutFailed = false

  try {
    await apiFetch.post('/management/logout')
  } catch (err) {
    if (isSessionExpiredError(err)) {
      sessionExpired = true
    } else {
      console.error('Logout failed', err)
      toast.warning('Logout failed. Redirecting to login.')
      logoutFailed = true
    }
  }

  if (!sessionExpired) {
    const redirect = () => {
      window.location.assign('/management/next/login')
    }

    if (logoutFailed) {
      window.setTimeout(redirect, SESSION_EXPIRED_REDIRECT_DELAY_MS)
    } else {
      redirect()
    }
  }
}
