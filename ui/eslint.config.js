import js from '@eslint/js'
import { defineConfig } from 'eslint/config'
import eslintConfigPrettier from 'eslint-config-prettier'
import react from 'eslint-plugin-react'
import reactHooks from 'eslint-plugin-react-hooks'
import reactPerf from 'eslint-plugin-react-perf'
import simpleImportSort from 'eslint-plugin-simple-import-sort'
import globals from 'globals'
import tseslint from 'typescript-eslint'

const toolingFiles = ['vite.config.ts', 'eslint.config.js', 'playwright.config.ts']

const sharedReactRules = {
  ...react.configs.recommended.rules,
  ...react.configs['jsx-runtime'].rules,
  ...reactHooks.configs.recommended.rules,
  '@typescript-eslint/consistent-type-definitions': ['error', 'interface'],
  // Required members before optional (interfaces, type literals, classes).
  '@typescript-eslint/member-ordering': [
    'error',
    {
      default: {
        memberTypes: 'never',
        optionalityOrder: 'required-first',
        order: 'as-written',
      },
    },
  ],
  'react-perf/jsx-no-new-object-as-prop': 'warn',
  'react-perf/jsx-no-new-array-as-prop': 'warn',
  'react-perf/jsx-no-new-function-as-prop': 'warn',
  'simple-import-sort/imports': 'error',
  'simple-import-sort/exports': 'error',
  'no-restricted-imports': [
    'error',
    {
      patterns: [
        {
          group: ['../../*', '../../../*', '../../../../*'],
          message: 'Use path aliases (@/, @shared/, @tests/, etc.) instead of deep relative imports.',
        },
      ],
    },
  ],
}

export default defineConfig(
  { ignores: ['dist', 'node_modules', 'coverage', 'test/.auth', 'playwright-report', 'test-results'] },
  js.configs.recommended,
  ...tseslint.configs.strictTypeChecked.map((config) => ({
    ...config,
    files: ['src/**/*.{ts,tsx}', 'tests/**/*.{ts,tsx}'],
  })),
  {
    files: ['src/**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2022,
      globals: globals.browser,
      parserOptions: {
        project: ['./tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
    },
    settings: {
      react: { version: 'detect' },
    },
    plugins: {
      react,
      'react-hooks': reactHooks,
      'react-perf': reactPerf,
      'simple-import-sort': simpleImportSort,
    },
    rules: sharedReactRules,
  },
  {
    files: ['tests/**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2022,
      globals: globals.browser,
      parserOptions: {
        project: ['./tsconfig.test.json'],
        tsconfigRootDir: import.meta.dirname,
      },
    },
    settings: {
      react: { version: 'detect' },
    },
    plugins: {
      react,
      'react-hooks': reactHooks,
      'react-perf': reactPerf,
      'simple-import-sort': simpleImportSort,
    },
    rules: sharedReactRules,
  },
  {
    files: ['test/**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2022,
      sourceType: 'module',
      globals: globals.node,
      parser: tseslint.parser,
      parserOptions: {
        ecmaVersion: 2022,
        sourceType: 'module',
      },
    },
    plugins: {
      '@typescript-eslint': tseslint.plugin,
      'simple-import-sort': simpleImportSort,
    },
    rules: {
      'simple-import-sort/imports': 'error',
      'simple-import-sort/exports': 'error',
      'no-restricted-imports': 'off',
    },
  },
  {
    ...tseslint.configs.disableTypeChecked,
    files: toolingFiles,
    languageOptions: {
      ...tseslint.configs.disableTypeChecked.languageOptions,
      globals: globals.node,
    },
  },
  eslintConfigPrettier,
)
