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
    setupFiles: ['./tests/setup/vitest-setup.ts'],
    include: ['tests/unit/**/*.test.{ts,tsx}'],
    globals: true,
    coverage: {
      provider: 'v8',
      reportsDirectory: './coverage',
      reporter: ['text', 'json', 'lcov'],
      include: [
        'src/shared/utils/initial-data.ts',
        'src/api/management/**',
        'src/feature/auth/**',
      ],
      thresholds: { lines: 90 },
    },
  },
})
