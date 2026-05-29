/// <reference types="vitest/config" />
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  base: '/management/next/',
  plugins: [
    react(),
    tsconfigPaths(),
    {
      name: 'go-template-nonce',
      transformIndexHtml: {
        order: 'post',
        handler(html) {
          return html
            .replace(/<script/g, '<script nonce="{{ .ScriptNonce }}"')
            .replace(/<style/g, '<style nonce="{{ .ScriptNonce }}"')
        },
      },
    },
  ],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    globals: true,
  },
})
