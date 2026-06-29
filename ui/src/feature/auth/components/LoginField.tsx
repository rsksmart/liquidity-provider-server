import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

interface LoginFieldProps {
  label: string
  name: string
  type: 'text' | 'password'
  testId?: string
  inputClassName?: string
}

export function LoginField({ label, name, type, testId, inputClassName }: LoginFieldProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor={name}>{label}</Label>
      <Input id={name} name={name} type={type} required data-testid={testId} className={inputClassName} />
    </div>
  )
}
