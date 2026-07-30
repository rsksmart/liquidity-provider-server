import { expect, test } from '../fixtures'
import { applyFreshManagementSession } from '../fixtures/session'
import {
  effectiveHeight,
  expectNextAtLeastReference,
  expectNextMatchesReference,
  expectPxClose,
  horizontalGap,
  nextCollateralCard,
  nextCollateralControls,
  nextConfigurationCard,
  nextProviderCard,
  openReferencePage,
  parsePx,
  readBox,
  readStyles,
  referenceCollateralCard,
  referenceCollateralControls,
  referenceConfigurationCard,
  referenceProviderCard,
} from './management-parity-helpers'

const hasE2ECredentials =
  Boolean(process.env.LPS_E2E_USER?.trim()) && Boolean(process.env.LPS_E2E_PASSWORD?.trim())

test.describe('management cards parity (AC-5)', () => {
  test.beforeEach(async ({ request, context, page }) => {
    test.skip(!hasE2ECredentials, 'requires LPS_E2E_USER and LPS_E2E_PASSWORD')
    await page.setViewportSize({ width: 1280, height: 900 })
    await applyFreshManagementSession(request, context)
  })

  test('provider fields and collateral tabs render', async ({ page }) => {
    await page.goto('management')

    await expect(page.getByRole('heading', { level: 1, name: 'Management Dashboard' })).toBeVisible()
    await expect(page.getByTestId('provider-rsk-address')).toBeVisible()
    await expect(page.getByTestId('provider-btc-address')).toBeVisible()
    await expect(page.getByTestId('provider-operational-status')).toBeVisible()
    await expect(page.getByRole('tab', { name: 'Pegin' })).toBeVisible()
    await expect(page.getByRole('tab', { name: 'Pegout' })).toBeVisible()
    await expect(page.getByTestId('pegin-collateral-balance')).not.toHaveText('—')
  })

  test('tab switch preserves collateral amount input', async ({ page }) => {
    await page.goto('management')

    const peginAmount = page.getByTestId('pegin-collateral-amount')
    await peginAmount.fill('1.5')
    await page.getByRole('tab', { name: 'Pegout' }).click()
    await page.getByRole('tab', { name: 'Pegin' }).click()

    await expect(peginAmount).toHaveValue('1.5')
  })

  test('submit shows loading bar below button and disables submit', async ({ page }) => {
    await page.route('**/pegin/addCollateral', async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 400))
      await route.fulfill({ status: 200, body: '{}' })
    })

    await page.goto('management')
    await page.getByTestId('pegin-collateral-amount').fill('1')

    const submit = page.getByTestId('pegin-add-collateral-button')
    const loadingBar = page.getByTestId('pegin-loading-bar')

    await submit.click()

    await expect(submit).toBeDisabled()
    await expect(loadingBar).toHaveClass(/is-visible/)

    const submitBox = await readBox(submit)
    const barBox = await readBox(loadingBar)
    expect(barBox.y).toBeGreaterThanOrEqual(submitBox.y + submitBox.height - 2)

    const barStyles = await readStyles(loadingBar)
    expectPxClose(barStyles.height, 2, 1)
    expectPxClose(barStyles.marginTop, 10, 2)
    expect(barStyles.width).not.toBe('0px')
  })

  test('validation blocks non-positive amount without loading bar', async ({ page }) => {
    let postCount = 0
    await page.route('**/pegin/addCollateral', async (route) => {
      postCount += 1
      await route.fulfill({ status: 200, body: '{}' })
    })

    await page.goto('management')
    await page.getByTestId('pegin-collateral-amount').fill('0')
    await page.getByTestId('pegin-add-collateral-button').click()

    // Toasts are out of parity scope; assert behavioral guardrails only.
    expect(postCount).toBe(0)
    await expect(page.getByTestId('pegin-loading-bar')).not.toHaveClass(/is-visible/)
    await expect(page.getByTestId('pegin-add-collateral-button')).toBeEnabled()
  })

  test('provider card width matches reference col-md-6 cards', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    const referenceCardBox = await readBox(referenceProviderCard(referencePage))
    const nextCardBox = await readBox(nextProviderCard(page))

    expect(nextCardBox.width).toBeGreaterThanOrEqual(referenceCardBox.width - 12)
    expect(nextCardBox.width).toBeLessThanOrEqual(referenceCardBox.width + 12)

    await closeReference()
  })

  test('card chrome matches reference Bootstrap cards', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    const referenceHeader = referenceProviderCard(referencePage).locator('.card-header').first()
    const nextHeader = nextProviderCard(page).locator('[data-slot="card-header"]').first()
    const referenceBody = referenceProviderCard(referencePage).locator('.card-body').first()
    const nextBody = nextProviderCard(page).locator('[data-slot="card-content"]').first()

    const referenceHeaderStyles = await readStyles(referenceHeader)
    const nextHeaderStyles = await readStyles(nextHeader)
    const referenceBodyStyles = await readStyles(referenceBody)
    const nextBodyStyles = await readStyles(nextBody)

    expectNextMatchesReference(nextHeaderStyles.fontSize, referenceHeaderStyles.fontSize, 1)
    expectNextMatchesReference(nextHeaderStyles.paddingLeft, referenceHeaderStyles.paddingLeft, 2)
    expectNextMatchesReference(nextHeaderStyles.paddingRight, referenceHeaderStyles.paddingRight, 2)
    expectNextMatchesReference(nextBodyStyles.paddingTop, referenceBodyStyles.paddingTop, 2)
    expectNextMatchesReference(nextBodyStyles.paddingBottom, referenceBodyStyles.paddingBottom, 2)
    expectNextMatchesReference(nextBodyStyles.paddingLeft, referenceBodyStyles.paddingLeft, 2)
    expectNextMatchesReference(nextBodyStyles.paddingRight, referenceBodyStyles.paddingRight, 2)

    await closeReference()
  })

  test('provider field typography matches reference', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    const referenceFieldTitle = referencePage
      .locator('#providerRskAddress')
      .locator('xpath=preceding-sibling::h5[1]')
    const nextFieldTitle = page.getByText('Provider RSK Address', { exact: true })
    const referenceFieldText = referencePage.locator('#providerRskAddress')
    const nextFieldText = page.getByTestId('provider-rsk-address')

    const referenceTitleStyles = await readStyles(referenceFieldTitle)
    const nextTitleStyles = await readStyles(nextFieldTitle)
    const referenceTextStyles = await readStyles(referenceFieldText)
    const nextTextStyles = await readStyles(nextFieldText)

    expectPxClose(nextTitleStyles.fontSize, 16, 1)
    expectNextMatchesReference(nextTitleStyles.fontSize, referenceTitleStyles.fontSize, 1)
    expectNextMatchesReference(nextTextStyles.fontSize, referenceTextStyles.fontSize, 1)
    expectNextMatchesReference(nextTextStyles.marginTop, referenceTextStyles.marginTop, 2)

    await closeReference()
  })

  test('collateral tabs match reference nav-link sizing', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    const referenceTab = referencePage.locator('#pegin-tab')
    const nextTab = page.getByRole('tab', { name: 'Pegin' })

    const referenceTabStyles = await readStyles(referenceTab)
    const nextTabStyles = await readStyles(nextTab)

    expectNextMatchesReference(nextTabStyles.fontSize, referenceTabStyles.fontSize, 1)
    expectNextMatchesReference(nextTabStyles.paddingTop, referenceTabStyles.paddingTop, 2)
    expectNextMatchesReference(nextTabStyles.paddingBottom, referenceTabStyles.paddingBottom, 2)
    expectNextMatchesReference(nextTabStyles.paddingLeft, referenceTabStyles.paddingLeft, 2)
    expectNextMatchesReference(nextTabStyles.paddingRight, referenceTabStyles.paddingRight, 2)

    await closeReference()
  })

  test('pegin form controls match reference Bootstrap sizes', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    const reference = referenceCollateralControls(referencePage, 'pegin')
    const next = nextCollateralControls(page, 'pegin')

    await expect(next.input).toHaveAttribute('placeholder', 'Enter amount in rBTC')

    const pairs: Array<[LocatorKey, boolean]> = [
      ['fieldTitle', true],
      ['label', true],
      ['input', false],
      ['button', false],
    ]

    for (const [key, matchFont] of pairs) {
      const referenceStyles = await readStyles(reference[key])
      const nextStyles = await readStyles(next[key])
      if (matchFont) {
        expectNextMatchesReference(nextStyles.fontSize, referenceStyles.fontSize, 1)
      }
      if (key === 'input' || key === 'button') {
        expectNextAtLeastReference(
          String(effectiveHeight(nextStyles)),
          String(effectiveHeight(referenceStyles)),
        )
        expect(effectiveHeight(nextStyles)).toBeGreaterThanOrEqual(36)
      }
    }

    await closeReference()
  })

  test('pegout form controls match reference Bootstrap sizes', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    await page.getByRole('tab', { name: 'Pegout' }).click()

    const reference = referenceCollateralControls(referencePage, 'pegout')
    const next = nextCollateralControls(page, 'pegout')

    await referencePage.locator('#pegout-tab').click()

    const referenceButtonStyles = await readStyles(reference.button)
    const nextButtonStyles = await readStyles(next.button)
    const referenceInputStyles = await readStyles(reference.input)
    const nextInputStyles = await readStyles(next.input)

    expectNextAtLeastReference(
      String(effectiveHeight(nextButtonStyles)),
      String(effectiveHeight(referenceButtonStyles)),
    )
    expectNextAtLeastReference(
      String(effectiveHeight(nextInputStyles)),
      String(effectiveHeight(referenceInputStyles)),
    )
    expect(effectiveHeight(nextInputStyles)).toBeGreaterThanOrEqual(36)
    expect(effectiveHeight(nextButtonStyles)).toBeGreaterThanOrEqual(36)

    await closeReference()
  })

  test('page title and collateral field title sizes match reference', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    const referenceTitle = referencePage.locator('.main-content h1').first()
    const nextTitle = page.getByRole('heading', { level: 1, name: 'Management Dashboard' })
    const referenceFieldTitle = referencePage.locator('#pegin .card-title').first()
    const nextFieldTitle = page.getByRole('heading', { name: 'Pegin Collateral', level: 3 })

    const referenceTitleStyles = await readStyles(referenceTitle)
    const nextTitleStyles = await readStyles(nextTitle)
    const referenceFieldStyles = await readStyles(referenceFieldTitle)
    const nextFieldStyles = await readStyles(nextFieldTitle)

    expectPxClose(nextTitleStyles.fontSize, parsePx(referenceTitleStyles.fontSize), 4)
    expectPxClose(nextFieldStyles.fontSize, parsePx(referenceFieldStyles.fontSize), 1)
    expectPxClose(nextFieldStyles.fontSize, 16, 1)

    await closeReference()
  })

  test('column gutter matches reference row gutter', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    const referenceGutter = await horizontalGap(
      referenceProviderCard(referencePage),
      referenceConfigurationCard(referencePage),
    )
    const nextGutter = await horizontalGap(
      nextProviderCard(page),
      nextConfigurationCard(page),
    )

    // Guards against a false pass if the reference ever collapses to a single column.
    expect(referenceGutter).toBeGreaterThan(0)
    expectPxClose(nextGutter, referenceGutter, 2)

    await closeReference()
  })

  test('card stack spacing matches reference compact-row rhythm', async ({ page, context }) => {
    const { referencePage, closeReference } = await openReferencePage(context)
    await page.goto('management')

    const referenceProviderBox = await readBox(referenceProviderCard(referencePage))
    const referenceCollateralBox = await readBox(referenceCollateralCard(referencePage))
    const nextProviderBox = await readBox(nextProviderCard(page))
    const nextCollateralBox = await readBox(nextCollateralCard(page))

    const referenceGap =
      referenceCollateralBox.y - (referenceProviderBox.y + referenceProviderBox.height)
    const nextGap = nextCollateralBox.y - (nextProviderBox.y + nextProviderBox.height)

    expect(nextGap).toBeGreaterThanOrEqual(referenceGap - 4)
    expect(nextGap).toBeLessThanOrEqual(referenceGap + 8)

    await closeReference()
  })
})

type LocatorKey = 'fieldTitle' | 'label' | 'input' | 'button'
