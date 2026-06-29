import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface LoginFieldProps {
  label: string
  name: string
  type: 'text' | 'password'
  testId?: string
}

export function LoginField({ label, name, type, testId }: LoginFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor={name}>{label}</Label>
      <Input id={name} name={name} type={type} required data-testid={testId} />
    </div>
  )
}
