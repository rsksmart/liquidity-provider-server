import { expect, test } from '../fixtures'
import { getShellUrl } from '../fixtures/session'

const SCRIPT_TAG_RE = /<script[^>]*>/gi
const STYLE_TAG_RE = /<style[^>]*>/gi
const CSRF_META_RE = /<meta name="csrf-token" content="([^"]*)"/
const NONCE_RE = /nonce="([a-f0-9]+)"/g

function validateCsrfMeta(html: string): void {
  const match = CSRF_META_RE.exec(html)
  if (!match?.[1]?.trim()) {
    throw new Error('csrf-token meta missing or empty')
  }
}

function validateEveryScriptAndStyleHasNonce(html: string): void {
  const nonceMatches = [...html.matchAll(NONCE_RE)]
  if (nonceMatches.length === 0) {
    throw new Error('no nonce attribute found in served HTML')
  }

  const expectedNonce = nonceMatches[0]?.[1]
  if (!expectedNonce) {
    throw new Error('no nonce attribute found in served HTML')
  }

  for (const tag of html.match(SCRIPT_TAG_RE) ?? []) {
    if (!tag.includes(`nonce="${expectedNonce}"`)) {
      throw new Error(`script tag missing expected nonce: ${tag}`)
    }
  }

  for (const tag of html.match(STYLE_TAG_RE) ?? []) {
    if (!tag.includes(`nonce="${expectedNonce}"`)) {
      throw new Error(`style tag missing expected nonce: ${tag}`)
    }
  }
}

test.describe('next UI shell security', () => {
  test('logged-out shell exposes CSRF meta and CSP nonces on served HTML', async ({ request }) => {
    const response = await request.get(getShellUrl())
    expect(response.ok()).toBeTruthy()

    const html = await response.text()
    validateCsrfMeta(html)
    validateEveryScriptAndStyleHasNonce(html)
  })
})
