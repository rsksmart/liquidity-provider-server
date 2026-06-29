import './index.css'

import { bootstrapDevEnvironment } from '@shared/utils/initial-data'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'

import { App } from '@/App'
import { Toaster } from '@/components/ui/sonner'

const rootElement = document.getElementById('root')
if (!rootElement) {
  throw new Error('Root element #root not found')
}

const appRoot: HTMLElement = rootElement

async function start() {
  if (import.meta.env.DEV) {
    await bootstrapDevEnvironment()
  }

  createRoot(appRoot).render(
    <StrictMode>
      <BrowserRouter basename="/management/next">
        <App />
        <Toaster />
      </BrowserRouter>
    </StrictMode>,
  )
}

void start()
