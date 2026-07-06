import { expect, type Locator, type Page } from '@playwright/test'

import { getServerOrigin } from '../fixtures/session'

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
}

export function parsePx(value: string): number {
  return Number.parseFloat(value)
}

export function expectPxClose(actual: string, expected: number, tolerance = 2): void {
  expect(parsePx(actual)).toBeGreaterThanOrEqual(expected - tolerance)
  expect(parsePx(actual)).toBeLessThanOrEqual(expected + tolerance)
}

export function expectNextMatchesReference(
  nextValue: string,
  referenceValue: string,
  tolerance = 2,
): void {
  expectPxClose(nextValue, parsePx(referenceValue), tolerance)
}

export function expectNextAtLeastReference(nextValue: string, referenceValue: string, slack = 2): void {
  expect(parsePx(nextValue)).toBeGreaterThanOrEqual(parsePx(referenceValue) - slack)
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
    }
  })
}

export function effectiveHeight(styles: Pick<StyleMetrics, 'height' | 'minHeight'>): number {
  return Math.max(parsePx(styles.height), parsePx(styles.minHeight))
}

export interface ReferencePageHandle {
  referencePage: Page
  closeReference: () => Promise<void>
}

export async function openReferencePage(
  context: import('@playwright/test').BrowserContext,
): Promise<ReferencePageHandle> {
  const origin = getServerOrigin()
  const referencePage = await context.newPage()
  await referencePage.setViewportSize({ width: 1280, height: 900 })
  await referencePage.goto(`${origin}/management`, { waitUntil: 'networkidle' })
  return {
    referencePage,
    closeReference: async () => {
      await referencePage.close()
    },
  }
}

export function referenceProviderCard(referencePage: Page): Locator {
  return referencePage.locator('.compact-row .col-md-6 .card').filter({ hasText: 'Provider RSK Address' })
}

export function referenceCollateralCard(referencePage: Page): Locator {
  return referencePage.locator('.compact-row .col-md-6 .card').filter({ hasText: 'Pegin Collateral' })
}

export function nextProviderCard(page: Page): Locator {
  return page.locator('[data-slot="card"]').filter({ hasText: 'Provider RSK Address' })
}

export function nextCollateralCard(page: Page): Locator {
  return page.locator('[data-slot="card"]').filter({ hasText: 'Pegin Collateral' })
}

export type CollateralKind = 'pegin' | 'pegout'

export function referenceCollateralControls(referencePage: Page, kind: CollateralKind) {
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
    label: page.getByText(`Add ${kind === 'pegin' ? 'Pegin' : 'Pegout'} Collateral Amount`, {
      exact: true,
    }),
    input: page.getByTestId(`${prefix}-collateral-amount`),
    button: page.getByTestId(`${prefix}-add-collateral-button`),
    loadingBar: page.getByTestId(`${prefix}-loading-bar`),
    tab: page.getByRole('tab', { name: kind === 'pegin' ? 'Pegin' : 'Pegout' }),
  }
}
