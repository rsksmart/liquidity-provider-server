import { AppErrorBoundary } from '@feature/error/components/AppErrorBoundary'
import { render, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

function Boom(): never {
  throw new Error('config card exploded')
}

describe('AppErrorBoundary', () => {
  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders children while nothing throws', () => {
    render(
      <AppErrorBoundary>
        <p>dashboard</p>
      </AppErrorBoundary>,
    )

    expect(screen.getByText('dashboard')).toBeInTheDocument()
  })

  it('shows the render error instead of an empty page', () => {
    render(
      <AppErrorBoundary>
        <Boom />
      </AppErrorBoundary>,
    )

    expect(screen.getByTestId('app-error-boundary')).toBeInTheDocument()
    expect(screen.getByText('config card exploded')).toBeInTheDocument()
  })

  it('logs the render error for operators', () => {
    render(
      <AppErrorBoundary>
        <Boom />
      </AppErrorBoundary>,
    )

    expect(console.error).toHaveBeenCalledWith(
      'Management UI failed to render:',
      expect.objectContaining({ message: 'config card exploded' }),
      expect.anything(),
    )
  })
})
