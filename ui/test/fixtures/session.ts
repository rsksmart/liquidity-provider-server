import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

import type { APIRequestContext, BrowserContext } from '@playwright/test'

const AUTH_DIR = path.join(path.dirname(fileURLToPath(import.meta.url)), '..', '.auth')
export const MANAGEMENT_STORAGE_STATE_PATH = path.join(AUTH_DIR, 'management.json')

export function getServerOrigin(): string {
  const baseUrl = process.env.LPS_E2E_BASE_URL ?? 'http://localhost:8080/management/next'
  return new URL(baseUrl).origin
}

export function getShellUrl(): string {
  const baseUrl = process.env.LPS_E2E_BASE_URL ?? 'http://localhost:8080/management/next'
  return baseUrl.endsWith('/') ? baseUrl : `${baseUrl}/`
}

export function extractCsrfFromResponseHtml(html: string): string {
  const match = html.match(/<meta name="csrf-token" content="([^"]*)"/)
  if (!match?.[1]?.trim()) {
    throw new Error('csrf-token meta missing or empty')
  }
  return match[1].replaceAll('&#43;', '+')
}

export async function seedManagementSession(request: APIRequestContext): Promise<void> {
  const user = process.env.LPS_E2E_USER
  const password = process.env.LPS_E2E_PASSWORD
  if (!user?.trim() || !password?.trim()) {
    throw new Error('LPS_E2E_USER and LPS_E2E_PASSWORD env vars are required for authenticated E2E')
  }

  const origin = getServerOrigin()
  const shellResponse = await request.get(getShellUrl())
  if (!shellResponse.ok()) {
    throw new Error(`Failed to GET next UI shell: ${shellResponse.status()}`)
  }

  const csrfToken = extractCsrfFromResponseHtml(await shellResponse.text())
  const loginResponse = await request.post(`${origin}/management/login`, {
    headers: {
      'Content-Type': 'application/json',
      'X-CSRF-Token': csrfToken,
    },
    data: { username: user, password },
  })

  if (!loginResponse.ok()) {
    throw new Error(
      `POST /management/login failed: ${loginResponse.status()} ${await loginResponse.text()}`,
    )
  }

  fs.mkdirSync(AUTH_DIR, { recursive: true })
  await request.storageState({ path: MANAGEMENT_STORAGE_STATE_PATH })
}

export async function applyFreshManagementSession(
  request: APIRequestContext,
  context: BrowserContext,
): Promise<void> {
  await seedManagementSession(request)
  const state = JSON.parse(fs.readFileSync(MANAGEMENT_STORAGE_STATE_PATH, 'utf8')) as {
    cookies: Parameters<BrowserContext['addCookies']>[0]
  }
  await context.clearCookies()
  await context.addCookies(state.cookies)
}

export async function clearSessionCookie(context: BrowserContext): Promise<void> {
  const cookies = await context.cookies()
  await context.clearCookies()
  const withoutSession = cookies.filter((cookie) => cookie.name !== 'lp-session')
  if (withoutSession.length > 0) {
    await context.addCookies(withoutSession)
  }
}
