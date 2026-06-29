import { apiFetch } from '@api/management/utils/api-fetch'
import { toast } from 'sonner'

export async function logout(): Promise<void> {
  try {
    await apiFetch('/management/logout', { method: 'POST' })
  } catch (err) {
    console.error('Logout failed', err)
    toast.warning('Logout failed; local session cleared.')
  } finally {
    window.location.assign('/management/next/login')
  }
}
