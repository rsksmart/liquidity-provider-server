import { logout } from '@feature/auth/logout'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

function handleLogout(): void {
  void logout()
}

interface LogoutButtonProps {
  className?: string
}

export function LogoutButton({ className }: LogoutButtonProps) {
  return (
    <Button
      type="button"
      variant="bootstrap"
      className={cn('w-auto', className)}
      onClick={handleLogout}
    >
      Logout
    </Button>
  )
}
