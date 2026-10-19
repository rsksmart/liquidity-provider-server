/**
 * Trusted Accounts visual parity — live legacy `/management` is the golden reference.
 * Asserts sizes, colors, padding, and typography that Stage 06 reviews by eye.
 * Read-only: open dialog only; never POST/DELETE trusted-accounts.
 */
import { expect, test } from '../fixtures'
import { applyFreshManagementSession } from '../fixtures/session'
import {
  expectNextMatchesReference,
  expectPxClose,
  expectRgbClose,
  nextTrustedAccountsControls,
  openReferencePage,
  readBox,
  readRgb,
  readStyles,
  referenceTrustedAccountsControls,
} from './management-parity-helpers'

const hasE2ECredentials =
  Boolean(process.env.LPS_E2E_USER?.trim()) &&
  Boolean(process.env.LPS_E2E_PASSWORD?.trim())

test.describe('trusted accounts parity (P-10)', () => {
  test.beforeEach(async ({ request, context, page }) => {
    test.skip(!hasE2ECredentials, 'requires LPS_E2E_USER and LPS_E2E_PASSWORD')
    await page.setViewportSize({ width: 1280, height: 900 })
    await applyFreshManagementSession(request, context)
  })

  test('card chrome matches legacy Bootstrap card', async ({
    page,
    context,
  }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management', { waitUntil: 'domcontentloaded' })

    const reference = referenceTrustedAccountsControls(referencePage)
    const next = nextTrustedAccountsControls(page)

    await expect(next.card).toBeVisible()
    await expect(next.addButton).toHaveText('Add Account')

    const referenceCardStyles = await readStyles(reference.card)
    const nextCardStyles = await readStyles(next.card)
    const referenceHeaderStyles = await readStyles(reference.header)
    const nextHeaderStyles = await readStyles(next.header)

    expectNextMatchesReference(
      nextCardStyles.borderRadius,
      referenceCardStyles.borderRadius,
      1,
    )
    expectRgbClose(
      await readRgb(next.card, 'backgroundColor'),
      await readRgb(reference.card, 'backgroundColor'),
    )
    expectRgbClose(
      await readRgb(next.card, 'color'),
      await readRgb(reference.card, 'color'),
    )

    expectNextMatchesReference(
      nextHeaderStyles.paddingTop,
      referenceHeaderStyles.paddingTop,
      1,
    )
    expectNextMatchesReference(
      nextHeaderStyles.paddingBottom,
      referenceHeaderStyles.paddingBottom,
      1,
    )
    expectNextMatchesReference(
      nextHeaderStyles.paddingLeft,
      referenceHeaderStyles.paddingLeft,
      1,
    )
    expectNextMatchesReference(
      nextHeaderStyles.fontSize,
      referenceHeaderStyles.fontSize,
      1,
    )
    expectRgbClose(
      await readRgb(next.header, 'backgroundColor'),
      await readRgb(reference.header, 'backgroundColor'),
      8,
      0.02,
    )

    const referenceHeaderBox = await readBox(reference.header)
    const nextHeaderBox = await readBox(next.header)
    expectPxClose(nextHeaderBox.height, referenceHeaderBox.height, 2)

    await closeReference()
  })

  test('Add Account button matches legacy btn-primary btn-sm', async ({
    page,
    context,
  }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management', { waitUntil: 'domcontentloaded' })

    const reference = referenceTrustedAccountsControls(referencePage)
    const next = nextTrustedAccountsControls(page)

    const referenceBox = await readBox(reference.addButton)
    const nextBox = await readBox(next.addButton)
    const referenceStyles = await readStyles(reference.addButton)
    const nextStyles = await readStyles(next.addButton)

    expect(await next.addButton.innerText()).toBe('Add Account')
    expectPxClose(nextBox.width, referenceBox.width, 4)
    expectPxClose(nextBox.height, referenceBox.height, 2)
    expectNextMatchesReference(nextStyles.fontSize, referenceStyles.fontSize, 1)
    expectNextMatchesReference(
      nextStyles.paddingTop,
      referenceStyles.paddingTop,
      1,
    )
    expectNextMatchesReference(
      nextStyles.paddingLeft,
      referenceStyles.paddingLeft,
      1,
    )
    expectNextMatchesReference(
      nextStyles.borderRadius,
      referenceStyles.borderRadius,
      1,
    )
    expectRgbClose(
      await readRgb(next.addButton, 'backgroundColor'),
      await readRgb(reference.addButton, 'backgroundColor'),
    )
    expectRgbClose(
      await readRgb(next.addButton, 'color'),
      await readRgb(reference.addButton, 'color'),
    )

    await closeReference()
  })

  test('Remove button matches legacy btn-danger btn-sm', async ({
    page,
    context,
  }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management', { waitUntil: 'domcontentloaded' })

    const reference = referenceTrustedAccountsControls(referencePage)
    const next = nextTrustedAccountsControls(page)

    await expect(next.removeButton).toBeVisible()

    const referenceBox = await readBox(reference.removeButton)
    const nextBox = await readBox(next.removeButton)
    const referenceStyles = await readStyles(reference.removeButton)
    const nextStyles = await readStyles(next.removeButton)

    expectPxClose(nextBox.height, referenceBox.height, 2)
    expectPxClose(nextBox.width, referenceBox.width, 6)
    expectNextMatchesReference(nextStyles.fontSize, referenceStyles.fontSize, 1)
    expectNextMatchesReference(
      nextStyles.paddingTop,
      referenceStyles.paddingTop,
      1,
    )
    expectNextMatchesReference(
      nextStyles.paddingLeft,
      referenceStyles.paddingLeft,
      1,
    )
    expectRgbClose(
      await readRgb(next.removeButton, 'backgroundColor'),
      await readRgb(reference.removeButton, 'backgroundColor'),
      4,
      0.05,
    )
    expectRgbClose(
      await readRgb(next.removeButton, 'color'),
      await readRgb(reference.removeButton, 'color'),
    )

    await closeReference()
  })

  test('table container max-height and overflow match legacy', async ({
    page,
    context,
  }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management', { waitUntil: 'domcontentloaded' })

    const reference = referenceTrustedAccountsControls(referencePage)
    const next = nextTrustedAccountsControls(page)

    const referenceBox = await readBox(reference.tableContainer)
    const nextMaxHeight = await next.tableContainer.evaluate((el) => {
      const view = el.ownerDocument.defaultView
      if (!view) {
        throw new Error('defaultView missing')
      }
      return view.getComputedStyle(el).maxHeight
    })
    const nextOverflowX = await next.tableContainer.evaluate((el) => {
      const view = el.ownerDocument.defaultView
      if (!view) {
        throw new Error('defaultView missing')
      }
      return view.getComputedStyle(el).overflowX
    })
    const nextOverflowY = await next.tableContainer.evaluate((el) => {
      const view = el.ownerDocument.defaultView
      if (!view) {
        throw new Error('defaultView missing')
      }
      return view.getComputedStyle(el).overflowY
    })
    const referenceMaxHeight = await reference.tableContainer.evaluate((el) => {
      const view = el.ownerDocument.defaultView
      if (!view) {
        throw new Error('defaultView missing')
      }
      return view.getComputedStyle(el).maxHeight
    })
    const referenceOverflowX = await reference.tableContainer.evaluate((el) => {
      const view = el.ownerDocument.defaultView
      if (!view) {
        throw new Error('defaultView missing')
      }
      return view.getComputedStyle(el).overflowX
    })

    expectNextMatchesReference(nextMaxHeight, referenceMaxHeight, 1)
    expect(nextOverflowX).toBe(referenceOverflowX)
    expect(['auto', 'scroll']).toContain(nextOverflowY)
    expect(referenceBox.height).toBeLessThanOrEqual(205 + 2)

    await closeReference()
  })

  test('add dialog width, title, and close control match legacy modal-lg', async ({
    page,
    context,
  }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management', { waitUntil: 'domcontentloaded' })

    const reference = referenceTrustedAccountsControls(referencePage)
    const next = nextTrustedAccountsControls(page)

    await reference.addButton.click()
    await expect(reference.modal).toBeVisible()
    await expect(reference.modalTitle).toBeVisible()
    await next.addButton.click()
    await expect(next.dialog).toBeVisible()
    await expect(next.dialogTitle).toBeVisible()

    const referenceDialogBox = await readBox(reference.modalContent)
    const nextDialogBox = await readBox(next.dialog)
    const referenceTitleStyles = await readStyles(reference.modalTitle)
    const nextTitleStyles = await readStyles(next.dialogTitle)
    const referenceHeaderBox = await readBox(reference.modalHeader)
    const nextHeaderBox = await readBox(next.dialogHeader)
    const referenceCloseBox = await readBox(reference.modalClose)
    const nextCloseBox = await readBox(next.dialogClose)
    const referenceCloseStyles = await readStyles(reference.modalClose)
    const nextCloseStyles = await readStyles(next.dialogClose)
    const referenceContentStyles = await readStyles(reference.modalContent)
    const nextContentStyles = await readStyles(next.dialog)

    expectPxClose(nextDialogBox.width, referenceDialogBox.width, 4)
    expectPxClose(nextDialogBox.width, 800, 4)
    // Must stay fully inside the viewport (guards transform/animation centering regressions).
    expect(nextDialogBox.x).toBeGreaterThanOrEqual(0)
    expect(nextDialogBox.x + nextDialogBox.width).toBeLessThanOrEqual(1280 + 1)
    expect(nextDialogBox.y).toBeGreaterThanOrEqual(0)
    expectNextMatchesReference(
      nextContentStyles.borderRadius,
      referenceContentStyles.borderRadius,
      2,
    )

    expectNextMatchesReference(
      nextTitleStyles.fontSize,
      referenceTitleStyles.fontSize,
      1,
    )
    expectNextMatchesReference(
      nextTitleStyles.lineHeight,
      referenceTitleStyles.lineHeight,
      1,
    )
    expect(nextTitleStyles.fontWeight).toBe(referenceTitleStyles.fontWeight)
    expectRgbClose(
      await readRgb(next.dialogTitle, 'color'),
      await readRgb(reference.modalTitle, 'color'),
    )

    expectPxClose(nextHeaderBox.height, referenceHeaderBox.height, 2)
    expectPxClose(nextCloseBox.width, referenceCloseBox.width, 2)
    expectPxClose(nextCloseBox.height, referenceCloseBox.height, 2)
    expectPxClose(
      nextCloseStyles.opacity,
      parseFloat(referenceCloseStyles.opacity),
      0.05,
    )

    await closeReference()
  })

  test('add dialog fields match legacy form-label / form-control / form-text', async ({
    page,
    context,
  }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management', { waitUntil: 'domcontentloaded' })

    const reference = referenceTrustedAccountsControls(referencePage)
    const next = nextTrustedAccountsControls(page)

    await reference.addButton.click()
    await expect(reference.modal).toBeVisible()
    await expect(reference.accountNameInput).toBeVisible()
    await next.addButton.click()
    await expect(next.dialog).toBeVisible()
    await expect(next.accountNameInput).toBeVisible()

    const referenceLabelStyles = await readStyles(reference.accountNameLabel)
    const nextLabelStyles = await readStyles(next.accountNameLabel)
    const referenceInputStyles = await readStyles(reference.accountNameInput)
    const nextInputStyles = await readStyles(next.accountNameInput)
    const referenceHelperStyles = await readStyles(reference.accountNameHelper)
    const nextHelperStyles = await readStyles(next.accountNameHelper)
    const referenceInputBox = await readBox(reference.accountNameInput)
    const nextInputBox = await readBox(next.accountNameInput)

    expectNextMatchesReference(
      nextLabelStyles.fontSize,
      referenceLabelStyles.fontSize,
      1,
    )
    expectNextMatchesReference(
      nextLabelStyles.marginBottom,
      referenceLabelStyles.marginBottom,
      1,
    )
    expect(nextLabelStyles.fontWeight).toBe(referenceLabelStyles.fontWeight)

    expectPxClose(nextInputBox.height, referenceInputBox.height, 2)
    expectNextMatchesReference(
      nextInputStyles.fontSize,
      referenceInputStyles.fontSize,
      1,
    )
    expectNextMatchesReference(
      nextInputStyles.paddingTop,
      referenceInputStyles.paddingTop,
      1,
    )
    expectNextMatchesReference(
      nextInputStyles.paddingLeft,
      referenceInputStyles.paddingLeft,
      1,
    )
    expectNextMatchesReference(
      nextInputStyles.borderRadius,
      referenceInputStyles.borderRadius,
      1,
    )
    expectRgbClose(
      await readRgb(next.accountNameInput, 'backgroundColor'),
      await readRgb(reference.accountNameInput, 'backgroundColor'),
      4,
      0.05,
    )
    expectRgbClose(
      await readRgb(next.accountNameInput, 'borderTopColor'),
      await readRgb(reference.accountNameInput, 'borderTopColor'),
      8,
      0.1,
    )

    expectNextMatchesReference(
      nextHelperStyles.fontSize,
      referenceHelperStyles.fontSize,
      1,
    )
    expectNextMatchesReference(
      nextHelperStyles.marginTop,
      referenceHelperStyles.marginTop,
      1,
    )
    expectRgbClose(
      await readRgb(next.accountNameHelper, 'color'),
      await readRgb(reference.accountNameHelper, 'color'),
      6,
      0.1,
    )

    await closeReference()
  })

  test('add dialog footer buttons match legacy Cancel/Save', async ({
    page,
    context,
  }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management', { waitUntil: 'domcontentloaded' })

    const reference = referenceTrustedAccountsControls(referencePage)
    const next = nextTrustedAccountsControls(page)

    await reference.addButton.click()
    await expect(reference.modal).toBeVisible()
    await expect(reference.cancelButton).toBeVisible()
    await next.addButton.click()
    await expect(next.dialog).toBeVisible()
    await expect(next.cancelButton).toBeVisible()

    const referenceFooterStyles = await readStyles(reference.modalFooter)
    const nextFooterStyles = await readStyles(next.dialogFooter)
    const referenceCancelBox = await readBox(reference.cancelButton)
    const nextCancelBox = await readBox(next.cancelButton)
    const referenceSaveBox = await readBox(reference.saveButton)
    const nextSaveBox = await readBox(next.saveButton)

    expectNextMatchesReference(
      nextFooterStyles.paddingTop,
      referenceFooterStyles.paddingTop,
      2,
    )
    expectNextMatchesReference(
      nextFooterStyles.paddingLeft,
      referenceFooterStyles.paddingLeft,
      2,
    )

    expectPxClose(nextCancelBox.height, referenceCancelBox.height, 2)
    expectPxClose(nextSaveBox.height, referenceSaveBox.height, 2)
    expectRgbClose(
      await readRgb(next.cancelButton, 'backgroundColor'),
      await readRgb(reference.cancelButton, 'backgroundColor'),
    )
    expectRgbClose(
      await readRgb(next.saveButton, 'backgroundColor'),
      await readRgb(reference.saveButton, 'backgroundColor'),
    )
    expectRgbClose(
      await readRgb(next.cancelButton, 'color'),
      await readRgb(reference.cancelButton, 'color'),
    )
    expectRgbClose(
      await readRgb(next.saveButton, 'color'),
      await readRgb(reference.saveButton, 'color'),
    )

    await closeReference()
  })
})
