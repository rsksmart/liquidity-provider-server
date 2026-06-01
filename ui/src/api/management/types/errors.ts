export class ApiFetchError extends Error {
  readonly status: number
  readonly statusText: string
  readonly body: unknown

  constructor(status: number, statusText: string, body: unknown) {
    super(`API request failed: ${String(status)} ${statusText}`)
    this.name = 'ApiFetchError'
    this.status = status
    this.statusText = statusText
    this.body = body
  }
}

export class CsrfTokenMissingError extends Error {
  constructor() {
    super('CSRF token meta tag is missing or empty')
    this.name = 'CsrfTokenMissingError'
  }
}
