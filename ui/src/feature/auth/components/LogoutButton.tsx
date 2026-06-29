import { logout } from '@feature/auth/logout'

import { Button } from '@/components/ui/button'

function handleLogout(): void {
  void logout()
}

export function LogoutButton() {
  return (
    <Button type="button" onClick={handleLogout}>
      Logout
    </Button>
  )
}
