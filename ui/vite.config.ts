/// <reference types="vitest/config" />
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import tsconfigPaths from 'vite-tsconfig-paths'

import { lpsDevBootstrapPlugin, lpsDevProxyTarget } from './vite.lps-dev-bootstrap.ts'

const devInitialData = {
  loggedIn: false,
  data: {
    CredentialsSet: true,
    BaseUrl: '',
    BtcAddress: 'tb1qexample',
    RskAddress: '0xabc',
    ProviderData: {
      id: 0,
      address: '',
      name: '',
      apiBaseUrl: '',
      status: false,
      providerType: 0,
    },
    ColdWallet: {
      BtcAddress: '',
      RskAddress: '',
      Label: '',
    },
    Configuration: {
      general: {
        rskConfirmations: {},
        btcConfirmations: {},
        publicLiquidityCheck: false,
        maxLiquidity: null,
        reimbursementWindowBlocks: 0,
        excessTolerance: {
          isFixed: false,
          percentageValue: '0',
          fixedValue: '0',
        },
      },
      pegin: {
        timeForDeposit: 0,
        callTime: 0,
        penaltyFee: '0',
        fixedFee: '0',
        feePercentage: '0',
        maxValue: '0',
        minValue: '0',
      },
      pegout: {
        timeForDeposit: 0,
        expireTime: 0,
        penaltyFee: '0',
        fixedFee: '0',
        feePercentage: '0',
        maxValue: '0',
        minValue: '0',
        expireBlocks: 0,
        bridgeTransactionMin: '0',
      },
    },
  },
}

export default defineConfig({
  base: '/management/next/',
  server: {
    proxy: {
      // Management API routes — never proxy the Vite app shell under /management/next/*
      '^/management/(?!next(?:/|$))': {
        target: lpsDevProxyTarget,
        changeOrigin: true,
      },
      '^/pegin': {
        target: lpsDevProxyTarget,
        changeOrigin: true,
      },
      '^/pegout': {
        target: lpsDevProxyTarget,
        changeOrigin: true,
      },
      '^/configuration': {
        target: lpsDevProxyTarget,
        changeOrigin: true,
      },
      '^/reports': {
        target: lpsDevProxyTarget,
        changeOrigin: true,
      },
      '^/providers': {
        target: lpsDevProxyTarget,
        changeOrigin: true,
      },
    },
  },
  plugins: [
    react(),
    tailwindcss(),
    tsconfigPaths(),
    lpsDevBootstrapPlugin(),
    {
      name: 'go-template-nonce',
      apply: 'build',
      transformIndexHtml: {
        order: 'post',
        handler(html) {
          return html
            .replace(/<script/g, '<script nonce="{{ .ScriptNonce }}"')
            .replace(/<style/g, '<style nonce="{{ .ScriptNonce }}"')
        },
      },
    },
    {
      name: 'go-template-dev-stubs',
      apply: 'serve',
      transformIndexHtml: {
        order: 'post',
        handler(html) {
          const initialDataJSON = JSON.stringify(devInitialData).replace(/</g, '\\u003c')
          return html
            .replace(/\{\{ \.CsrfToken \}\}/g, 'dev-csrf-token')
            .replace(/\{\{ \.InitialDataJSON \}\}/g, initialDataJSON)
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
        'src/shared/utils/wei.ts',
        'src/api/management/**',
        'src/feature/auth/**',
        'src/feature/management/**',
      ],
      thresholds: { lines: 90 },
    },
  },
})
