import { apiFetch } from '@api/management/utils/api-fetch'
import { useInitialData } from '@shared/utils/initial-data'
import { type SubmitEvent, useCallback, useState } from 'react'

import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'

import { LoginField } from './LoginField'

const LOGIN_ERROR = 'Invalid username or password.'

function fieldValue(formData: FormData, name: string): string {
  const value = formData.get(name)
  return typeof value === 'string' ? value : ''
}

export function LoginPage() {
  const { data } = useInitialData()
  const credentialsSet = data.CredentialsSet
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const submitLogin = useCallback(
    async (form: HTMLFormElement) => {
      setError(null)
      setSubmitting(true)

      const formData = new FormData(form)
      const username = fieldValue(formData, 'username')
      const password = fieldValue(formData, 'password')

      try {
        await apiFetch('/management/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, password }),
        })

        if (!credentialsSet) {
          const newUsername = fieldValue(formData, 'new-username')
          const newPassword = fieldValue(formData, 'new-password')
          await apiFetch('/management/credentials', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              oldUsername: username,
              oldPassword: password,
              newUsername,
              newPassword,
            }),
          })
        }

        window.location.assign('/management/next/management')
      } catch {
        setError(LOGIN_ERROR)
      } finally {
        setSubmitting(false)
      }
    },
    [credentialsSet],
  )

  const onFormSubmit = useCallback(
    (event: SubmitEvent<HTMLFormElement>) => {
      event.preventDefault()
      void submitLogin(event.currentTarget)
    },
    [submitLogin],
  )

  return (
    <main className="container mx-auto mt-5 px-4">
      <div className="mx-auto w-full md:w-1/2">
        <Card>
          <CardHeader>
            <h1 className="text-2xl font-semibold">Login</h1>
          </CardHeader>
          <CardContent>
            <form className="space-y-4" autoComplete="off" onSubmit={onFormSubmit}>
              <LoginField
                label="Username"
                name="username"
                type="text"
                testId="login-username-input"
                inputClassName="h-10"
              />
              <LoginField
                label="Password"
                name="password"
                type="password"
                testId="login-password-input"
                inputClassName="h-10"
              />
              {!credentialsSet ? (
                <>
                  <LoginField label="New Username" name="new-username" type="text" inputClassName="h-10" />
                  <LoginField label="New Password" name="new-password" type="password" inputClassName="h-10" />
                </>
              ) : null}
              {error ? (
                <Alert variant="destructive">
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              ) : null}
              <Button
                type="submit"
                size="lg"
                className="h-10 px-3"
                disabled={submitting}
                data-testid="login-submit-button"
              >
                Login
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
