import { expect, type Locator, type Page } from '@playwright/test'

import { getLegacyOrigin } from '../fixtures/session'

export interface BoxMetrics {
  x: number
  y: number
  width: number
  height: number
}

export interface StyleMetrics {
  fontSize: string
  fontWeight: string
  minHeight: string
  height: string
  paddingTop: string
  paddingBottom: string
  paddingLeft: string
  paddingRight: string
  marginTop: string
  marginBottom: string
  backgroundColor: string
  borderTopWidth: string
  borderTopColor: string
  borderRadius: string
  color: string
  lineHeight: string
  maxWidth: string
  width: string
  opacity: string
}

export interface RgbaColor {
  r: number
  g: number
  b: number
  a: number
  raw: string
}

export function parsePx(value: string): number {
  return Number.parseFloat(value)
}

export function expectPxClose(
  actual: string | number,
  expected: number,
  tolerance = 2,
): void {
  const value = typeof actual === 'number' ? actual : parsePx(actual)
  expect(value).toBeGreaterThanOrEqual(expected - tolerance)
  expect(value).toBeLessThanOrEqual(expected + tolerance)
}

export function expectNextMatchesReference(
  nextValue: string,
  referenceValue: string,
  tolerance = 2,
): void {
  expectPxClose(nextValue, parsePx(referenceValue), tolerance)
}

export function expectNextAtLeastReference(
  nextValue: string,
  referenceValue: string,
  slack = 2,
): void {
  expect(parsePx(nextValue)).toBeGreaterThanOrEqual(
    parsePx(referenceValue) - slack,
  )
}

export async function readBox(locator: Locator): Promise<BoxMetrics> {
  return locator.evaluate((element) => {
    const rect = element.getBoundingClientRect()
    return { x: rect.x, y: rect.y, width: rect.width, height: rect.height }
  })
}

export async function readStyles(locator: Locator): Promise<StyleMetrics> {
  return locator.evaluate((element) => {
    const view = element.ownerDocument.defaultView
    if (!view) {
      throw new Error('defaultView missing')
    }
    const styles = view.getComputedStyle(element)
    return {
      fontSize: styles.fontSize,
      fontWeight: styles.fontWeight,
      minHeight: styles.minHeight,
      height: styles.height,
      paddingTop: styles.paddingTop,
      paddingBottom: styles.paddingBottom,
      paddingLeft: styles.paddingLeft,
      paddingRight: styles.paddingRight,
      marginTop: styles.marginTop,
      marginBottom: styles.marginBottom,
      backgroundColor: styles.backgroundColor,
      borderTopWidth: styles.borderTopWidth,
      borderTopColor: styles.borderTopColor,
      borderRadius: styles.borderRadius,
      color: styles.color,
      lineHeight: styles.lineHeight,
      maxWidth: styles.maxWidth,
      width: styles.width,
      opacity: styles.opacity,
    }
  })
}

/** Normalize any computed CSS color (rgb/oklab/oklch/…) to sRGB channels via canvas. */
export async function readRgb(
  locator: Locator,
  property: 'color' | 'backgroundColor' | 'borderTopColor',
): Promise<RgbaColor> {
  return locator.evaluate((element, prop) => {
    const view = element.ownerDocument.defaultView
    if (!view) {
      throw new Error('defaultView missing')
    }
    const raw = view.getComputedStyle(element)[prop]
    const canvas = element.ownerDocument.createElement('canvas')
    canvas.width = 1
    canvas.height = 1
    const ctx = canvas.getContext('2d')
    if (!ctx) {
      throw new Error('2d context missing')
    }
    ctx.clearRect(0, 0, 1, 1)
    ctx.fillStyle = raw
    ctx.fillRect(0, 0, 1, 1)
    const [r, g, b, a] = ctx.getImageData(0, 0, 1, 1).data
    return { r, g, b, a: a / 255, raw }
  }, property)
}

export function expectRgbClose(
  actual: RgbaColor,
  expected: RgbaColor,
  channelTolerance = 3,
  alphaTolerance = 0.08,
): void {
  expect(actual.a).toBeGreaterThanOrEqual(expected.a - alphaTolerance)
  expect(actual.a).toBeLessThanOrEqual(expected.a + alphaTolerance)
  // Near-transparent washes: RGB channels are unstable across color spaces — alpha is enough.
  if (expected.a <= 0.08 && actual.a <= 0.08) {
    return
  }
  expect(actual.r).toBeGreaterThanOrEqual(expected.r - channelTolerance)
  expect(actual.r).toBeLessThanOrEqual(expected.r + channelTolerance)
  expect(actual.g).toBeGreaterThanOrEqual(expected.g - channelTolerance)
  expect(actual.g).toBeLessThanOrEqual(expected.g + channelTolerance)
  expect(actual.b).toBeGreaterThanOrEqual(expected.b - channelTolerance)
  expect(actual.b).toBeLessThanOrEqual(expected.b + channelTolerance)
}

export function effectiveHeight(
  styles: Pick<StyleMetrics, 'height' | 'minHeight'>,
): number {
  return Math.max(parsePx(styles.height), parsePx(styles.minHeight) || 0)
}

export interface ReferencePageHandle {
  referencePage: Page
  closeReference: () => Promise<void>
}

export async function openReferencePage(
  context: import('@playwright/test').BrowserContext,
): Promise<ReferencePageHandle> {
  const legacyOrigin = getLegacyOrigin()
  const referencePage = await context.newPage()
  await referencePage.setViewportSize({ width: 1280, height: 900 })
  await referencePage.goto(`${legacyOrigin}/management`, {
    waitUntil: 'domcontentloaded',
  })
  return {
    referencePage,
    closeReference: async () => {
      await referencePage.close()
    },
  }
}

export function referenceProviderCard(referencePage: Page): Locator {
  return referencePage
    .locator('.compact-row .col-md-6 .card')
    .filter({ hasText: 'Provider RSK Address' })
}

export function referenceCollateralCard(referencePage: Page): Locator {
  return referencePage
    .locator('.compact-row .col-md-6 .card')
    .filter({ hasText: 'Pegin Collateral' })
}

export function nextProviderCard(page: Page): Locator {
  return page
    .locator('[data-slot="card"]')
    .filter({ hasText: 'Provider RSK Address' })
}

export function nextCollateralCard(page: Page): Locator {
  return page
    .locator('[data-slot="card"]')
    .filter({ hasText: 'Pegin Collateral' })
}

export type CollateralKind = 'pegin' | 'pegout'

export function referenceCollateralControls(
  referencePage: Page,
  kind: CollateralKind,
) {
  const prefix = kind === 'pegin' ? 'Pegin' : 'Pegout'
  return {
    fieldTitle: referencePage.locator(`#${kind} .card-title`).first(),
    balance: referencePage.locator(`#${kind}Collateral`),
    label: referencePage.locator(`label[for="add${prefix}CollateralAmount"]`),
    input: referencePage.locator(`#add${prefix}CollateralAmount`),
    button: referencePage.locator(`#add${prefix}CollateralButton`),
    loadingBar: referencePage.locator(`#${kind}LoadingBar`),
    tab: referencePage.locator(`#${kind}-tab`),
  }
}

export function nextCollateralControls(page: Page, kind: CollateralKind) {
  const prefix = kind
  const title = kind === 'pegin' ? 'Pegin Collateral' : 'Pegout Collateral'
  return {
    fieldTitle: page.getByRole('heading', { name: title, level: 3 }),
    balance: page.getByTestId(`${prefix}-collateral-balance`),
    label: page.getByText(
      `Add ${kind === 'pegin' ? 'Pegin' : 'Pegout'} Collateral Amount`,
      {
        exact: true,
      },
    ),
    input: page.getByTestId(`${prefix}-collateral-amount`),
    button: page.getByTestId(`${prefix}-add-collateral-button`),
    loadingBar: page.getByTestId(`${prefix}-loading-bar`),
    tab: page.getByRole('tab', { name: kind === 'pegin' ? 'Pegin' : 'Pegout' }),
  }
}

export function referenceTrustedAccountsCard(referencePage: Page): Locator {
  return referencePage
    .locator('.card')
    .filter({ has: referencePage.locator('#addTrustedAccountButton') })
}

export function nextTrustedAccountsCard(page: Page): Locator {
  return page.getByTestId('trusted-accounts-card')
}

export function referenceTrustedAccountsControls(referencePage: Page) {
  const card = referenceTrustedAccountsCard(referencePage)
  return {
    card,
    header: card.locator('.card-header').first(),
    listTitle: card.locator('h5.card-title').first(),
    addButton: referencePage.locator('#addTrustedAccountButton'),
    removeButton: card.locator('button.btn-danger').first(),
    tableContainer: referencePage.locator('#trustedAccountsContainer'),
    modal: referencePage.locator('#addTrustedAccountModal'),
    modalDialog: referencePage.locator('#addTrustedAccountModal .modal-dialog'),
    modalContent: referencePage.locator(
      '#addTrustedAccountModal .modal-content',
    ),
    modalHeader: referencePage.locator('#addTrustedAccountModal .modal-header'),
    modalTitle: referencePage.locator('#addTrustedAccountModalLabel'),
    modalClose: referencePage.locator('#addTrustedAccountModal .btn-close'),
    accountNameLabel: referencePage.locator('label[for="accountName"]'),
    accountNameInput: referencePage.locator('#accountName'),
    accountNameHelper: referencePage
      .locator('#accountName')
      .locator('xpath=following-sibling::div[contains(@class,"form-text")][1]'),
    modalFooter: referencePage.locator('#addTrustedAccountModal .modal-footer'),
    cancelButton: referencePage.locator(
      '#addTrustedAccountModal .modal-footer .btn-secondary',
    ),
    saveButton: referencePage.locator('#saveAccountButton'),
  }
}

export function nextTrustedAccountsControls(page: Page) {
  const card = nextTrustedAccountsCard(page)
  return {
    card,
    header: card.locator('[data-slot="card-header"]').first(),
    listTitle: card.getByRole('heading', { name: 'Accounts List' }),
    addButton: page.getByTestId('add-trusted-account-button'),
    removeButton: card
      .locator('[data-testid^="remove-trusted-account-"]')
      .first(),
    tableContainer: page.getByTestId('trusted-accounts-container'),
    dialog: page.locator('[data-slot="dialog-content"]'),
    dialogHeader: page.locator('[data-slot="dialog-header"]'),
    dialogTitle: page.locator('[data-slot="dialog-title"]'),
    dialogClose: page.locator(
      '[data-slot="dialog-content"] > [data-slot="dialog-close"]',
    ),
    accountNameLabel: page
      .locator('label')
      .filter({ hasText: /^Account Name$/ }),
    accountNameInput: page.getByTestId('account-name'),
    accountNameHelper: page.getByText(
      'A friendly name to identify this account',
      { exact: true },
    ),
    dialogFooter: page.locator('[data-slot="dialog-footer"]'),
    cancelButton: page.getByRole('button', { name: 'Cancel' }),
    saveButton: page.getByTestId('save-trusted-account-button'),
  }
}
