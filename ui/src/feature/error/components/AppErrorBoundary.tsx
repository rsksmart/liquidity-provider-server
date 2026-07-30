import { managementShellClass } from '@feature/management/management-styles'
import { Component, type ErrorInfo, type ReactNode } from 'react'

interface AppErrorBoundaryProps {
  children: ReactNode
}

interface AppErrorBoundaryState {
  error: Error | null
}

/**
 * React error boundaries require the class lifecycle API. Keep this boundary
 * at the application root; feature UI remains function-component based.
 */
export class AppErrorBoundary extends Component<
  AppErrorBoundaryProps,
  AppErrorBoundaryState
> {
  state: AppErrorBoundaryState = { error: null }

  static getDerivedStateFromError(error: Error): AppErrorBoundaryState {
    return { error }
  }

  componentDidCatch(error: Error, info: ErrorInfo): void {
    console.error('Management UI failed to render:', error, info.componentStack)
  }

  render(): ReactNode {
    const { error } = this.state
    if (!error) {
      return this.props.children
    }

    return (
      <main className={managementShellClass} data-testid="app-error-boundary">
        <h1 className="text-2xl font-medium text-[#212529]">
          Management UI failed to load
        </h1>
        <p className="mt-3 text-base text-[#212529]">
          Reload the page. If it keeps failing, report the error below with the
          browser console output.
        </p>
        <pre className="mt-4 overflow-x-auto rounded-[6px] border border-[#dee2e6] bg-[rgba(33,37,41,0.03)] p-4 text-sm text-[#212529]">
          {error.message}
        </pre>
      </main>
    )
  }
}
