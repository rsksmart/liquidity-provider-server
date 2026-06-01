import { screen } from '@testing-library/react'
import { loggedInFixture, loggedOutFixture } from '@tests/fixtures'
import { appBasename, renderAppAt } from '@tests/utils'
import { beforeEach, describe, expect, it } from 'vitest'

const rootPath = `${appBasename}/`
const errorPath = `${appBasename}/error`

describe('App router', () => {
  beforeEach(() => {
    document.body.innerHTML = ''
  })

  it('boots to login when loggedIn is false', () => {
    renderAppAt(rootPath, loggedOutFixture)
    expect(screen.getByRole('heading', { name: /login/i })).toBeInTheDocument()
  })

  it('boots to management when loggedIn is true', () => {
    renderAppAt(rootPath, loggedInFixture)
    expect(screen.getByRole('heading', { name: /management/i })).toBeInTheDocument()
  })

  it('renders the error placeholder route', () => {
    renderAppAt(errorPath, loggedOutFixture)
    expect(screen.getByRole('heading', { name: /error/i })).toBeInTheDocument()
  })
})
