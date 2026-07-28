/**
 * Configuration card visual parity suite.
 *
 * Dual-gated golden reference:
 *   1. `test/e2e/fixtures/legacy-configuration.metrics.json` — sizes and colors
 *      captured from the Bootstrap Configuration `.card`.
 *   2. Live `/management` when the reference origin is up (`openReferencePage`).
 *
 * Run modes:
 *   - Authenticated (`LPS_E2E_USER` / `LPS_E2E_PASSWORD`): React card vs golden
 *     metrics and live reference.
 *   - Fixture (`LPS_E2E_FIXTURE_MODE=1`): Vite dev server with `/__dev/lps-shell`
 *     bootstrap; React card vs golden metrics only.
 *
 * Read-only: never POSTs configuration mutations.
 */
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

import { expect, test } from '../fixtures'
import { applyFreshManagementSession } from '../fixtures/session'
import {
  type ConfigTabKind,
  expectPxClose,
  expectRgbClose,
  nextConfigurationControls,
  openReferencePage,
  parsePx,
  readBox,
  readRgb,
  readStyles,
  referenceConfigurationControls,
  type RgbaColor,
} from './management-parity-helpers'

const currentDir = path.dirname(fileURLToPath(import.meta.url))

interface ElementMetrics {
  rect: { x: number; y: number; width: number; height: number }
  computed: Record<string, string>
}
interface GoldenMetrics {
  baseCommit: string
  viewport: { width: number; height: number }
  tabs: Record<ConfigTabKind, Record<string, ElementMetrics | null>>
}

const golden = JSON.parse(
  fs.readFileSync(
    path.join(currentDir, 'fixtures', 'legacy-configuration.metrics.json'),
    'utf8',
  ),
) as GoldenMetrics

function goldenEl(tab: ConfigTabKind, key: string): ElementMetrics {
  const el = golden.tabs[tab]?.[key]
  if (!el) {
    throw new Error(`golden metric missing: tabs.${tab}.${key}`)
  }
  return el
}

/** Parse a computed rgb()/rgba() string into sRGB channels for expectRgbClose. */
function goldenRgb(value: string): RgbaColor {
  const match = value.match(
    /rgba?\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(?:,\s*([\d.]+)\s*)?\)/,
  )
  if (!match) {
    throw new Error(`unparseable golden color: ${value}`)
  }
  return {
    r: Number(match[1]),
    g: Number(match[2]),
    b: Number(match[3]),
    a: match[4] === undefined ? 1 : Number(match[4]),
    raw: value,
  }
}

const hasE2ECredentials =
  Boolean(process.env.LPS_E2E_USER?.trim()) &&
  Boolean(process.env.LPS_E2E_PASSWORD?.trim())

const fixtureMode = Boolean(process.env.LPS_E2E_FIXTURE_MODE?.trim())

// Authenticated bootstrap injected in fixture mode so ManagementPage renders
// (loggedIn:true) with a populated Configuration mirroring the capture fixture.
const FIXTURE_CONFIGURATION = {
  general: {
    rskConfirmations: {
      '1000000000000000000': 10,
      '5000000000000000000': 40,
    },
    btcConfirmations: {
      '500000000000000000': 3,
      '2000000000000000000': 12,
    },
    publicLiquidityCheck: true,
    maxLiquidity: '10000000000000000000',
    reimbursementWindowBlocks: 120,
    excessTolerance: { isFixed: true, percentageValue: '0', fixedValue: '2000000000000000000' },
  },
  pegin: {
    timeForDeposit: 3600,
    penaltyFee: '1000000000000000',
    fixedFee: '5000000000000000',
    feePercentage: '1.5',
    maxValue: '1000000000000000000000',
    minValue: '10000000000000000',
    callTime: 7200,
  },
  pegout: {
    timeForDeposit: 3600,
    penaltyFee: '1000000000000000',
    fixedFee: '5000000000000000',
    feePercentage: '1.5',
    maxValue: '1000000000000000000000',
    minValue: '10000000000000000',
    expireTime: 7200,
    expireBlocks: 500,
    bridgeTransactionMin: '1500000000000000000',
  },
}

const FIXTURE_INITIAL_DATA = {
  loggedIn: true,
  data: {
    CredentialsSet: true,
    BaseUrl: '',
    BtcAddress: 'bcrt1qparitycapture00000000000000000000000',
    RskAddress: '0x1111111111111111111111111111111111111111',
    ProviderData: {
      id: 0,
      address: '0x1111111111111111111111111111111111111111',
      name: 'Parity LP',
      apiBaseUrl: '',
      status: true,
      providerType: 0,
    },
    ColdWallet: {
      BtcAddress: 'bcrt1qcold000000000000000000000000000000',
      RskAddress: '0x2222222222222222222222222222222222222222',
      Label: 'Cold',
    },
    Configuration: FIXTURE_CONFIGURATION,
  },
}

const RATIO_STUB = {
  btcPercentage: 50,
  rbtcPercentage: 50,
  isPreview: false,
  cooldownActive: false,
  cooldownEndTimestamp: 0,
  btcTarget: '0',
  btcThreshold: '0',
  btcCurrentBalance: '0',
  rbtcTarget: '0',
  rbtcThreshold: '0',
  rbtcCurrentBalance: '0',
  btcImpact: { type: 'withinTolerance', amount: '0' },
  rbtcImpact: { type: 'withinTolerance', amount: '0' },
}

async function installFixtureRoutes(
  page: import('@playwright/test').Page,
): Promise<void> {
  await page.route('**/__dev/lps-shell', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        csrf: 'fixture-csrf-token',
        initialData: FIXTURE_INITIAL_DATA,
      }),
    })
  })
  const readOnly = async (
    route: import('@playwright/test').Route,
    body: unknown,
  ) => {
    if (route.request().method() !== 'GET') {
      await route.fulfill({
        status: 403,
        contentType: 'application/json',
        body: JSON.stringify({ message: 'parity e2e is read-only' }),
      })
      return
    }
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(body),
    })
  }
  await page.route(/\/(pegin|pegout)\/collateral(\?|$)/, (route) =>
    readOnly(route, { collateral: '0' }),
  )
  await page.route('**/management/trusted-accounts*', (route) =>
    readOnly(route, { accounts: [] }),
  )
  await page.route('**/management/liquidity-ratio*', (route) =>
    readOnly(route, RATIO_STUB),
  )
}

test.describe('configuration parity (P-9)', () => {
  test.beforeEach(async ({ request, context, page }) => {
    test.skip(
      !fixtureMode && !hasE2ECredentials,
      'requires LPS_E2E_USER/LPS_E2E_PASSWORD, or LPS_E2E_FIXTURE_MODE against a Vite dev server',
    )
    await page.setViewportSize({ width: 1280, height: 900 })
    if (fixtureMode) {
      await installFixtureRoutes(page)
    } else {
      await applyFreshManagementSession(request, context)
    }
  })

  test('card chrome matches golden Configuration card metrics', async ({
    page,
    context,
  }) => {
    await page.goto('management', { waitUntil: 'domcontentloaded' })
    const next = nextConfigurationControls(page)

    await expect(next.card).toBeVisible()
    await expect(next.header).toHaveText('Configuration')
    await expect(next.title).toBeVisible()
    await expect(next.saveButton).toHaveText('Save Configuration')

    const goldenCard = goldenEl('general', 'card')
    const goldenHeader = goldenEl('general', 'cardHeader')
    const nextCardStyles = await readStyles(next.card)
    const nextHeaderStyles = await readStyles(next.header)

    expectPxClose(nextCardStyles.borderRadius, parsePx(goldenCard.computed.borderRadius), 1)
    expectRgbClose(
      await readRgb(next.card, 'backgroundColor'),
      goldenRgb(goldenCard.computed.backgroundColor),
    )
    expectRgbClose(
      await readRgb(next.card, 'color'),
      goldenRgb(goldenCard.computed.color),
    )
    expectPxClose(nextHeaderStyles.paddingTop, parsePx(goldenHeader.computed.paddingTop), 1)
    expectPxClose(nextHeaderStyles.paddingLeft, parsePx(goldenHeader.computed.paddingLeft), 1)
    expectPxClose(nextHeaderStyles.fontSize, parsePx(goldenHeader.computed.fontSize), 1)

    if (!fixtureMode) {
      const { referencePage, closeReference } = await openReferencePage(context)
      const reference = referenceConfigurationControls(referencePage)
      await expect(reference.card).toBeVisible()
      expectRgbClose(
        await readRgb(next.card, 'backgroundColor'),
        await readRgb(reference.card, 'backgroundColor'),
      )
      expectPxClose(
        (await readBox(next.header)).height,
        (await readBox(reference.header)).height,
        2,
      )
      await closeReference()
    }
  })

  test('Save button matches golden Save Configuration metrics', async ({
    page,
    context,
  }) => {
    await page.goto('management', { waitUntil: 'domcontentloaded' })
    const next = nextConfigurationControls(page)

    await expect(next.saveButton).toBeVisible()

    const goldenSave = goldenEl('general', 'saveButton')
    const nextBox = await readBox(next.saveButton)
    const nextStyles = await readStyles(next.saveButton)

    expect(await next.saveButton.innerText()).toBe('Save Configuration')
    expectPxClose(nextBox.height, goldenSave.rect.height, 2)
    expectPxClose(nextStyles.paddingTop, parsePx(goldenSave.computed.paddingTop), 1)
    expectPxClose(nextStyles.paddingLeft, parsePx(goldenSave.computed.paddingLeft), 1)
    expectPxClose(nextStyles.borderRadius, parsePx(goldenSave.computed.borderRadius), 1)
    expectPxClose(nextStyles.fontSize, parsePx(goldenSave.computed.fontSize), 1)
    expectRgbClose(
      await readRgb(next.saveButton, 'backgroundColor'),
      goldenRgb(goldenSave.computed.backgroundColor),
    )
    expectRgbClose(
      await readRgb(next.saveButton, 'color'),
      goldenRgb(goldenSave.computed.color),
    )

    if (!fixtureMode) {
      const { referencePage, closeReference } = await openReferencePage(context)
      const reference = referenceConfigurationControls(referencePage)
      expectRgbClose(
        await readRgb(next.saveButton, 'backgroundColor'),
        await readRgb(reference.saveButton, 'backgroundColor'),
      )
      await closeReference()
    }
  })

  test('tabs match golden nav-tabs padding, size and colors', async ({
    page,
  }) => {
    await page.goto('management', { waitUntil: 'domcontentloaded' })
    const next = nextConfigurationControls(page)

    await expect(next.generalTab).toBeVisible()
    await expect(next.generalTab).toHaveText('General')
    await expect(next.peginTab).toHaveText('Pegin')
    await expect(next.pegoutTab).toHaveText('Pegout')

    const goldenGeneral = goldenEl('general', 'generalTab')
    const goldenPegin = goldenEl('general', 'peginTab')
    const nextGeneralStyles = await readStyles(next.generalTab)

    expectPxClose(nextGeneralStyles.paddingTop, parsePx(goldenGeneral.computed.paddingTop), 1)
    expectPxClose(nextGeneralStyles.paddingLeft, parsePx(goldenGeneral.computed.paddingLeft), 1)
    expectPxClose(nextGeneralStyles.fontSize, parsePx(goldenGeneral.computed.fontSize), 1)
    // Active tab dark, inactive tab primary-blue link.
    expectRgbClose(
      await readRgb(next.generalTab, 'color'),
      goldenRgb(goldenGeneral.computed.color),
    )
    expectRgbClose(
      await readRgb(next.peginTab, 'color'),
      goldenRgb(goldenPegin.computed.color),
    )
  })

  test('fee input ~40% width, confirmation controls ~180px, question icon 13x13', async ({
    page,
  }) => {
    await page.goto('management', { waitUntil: 'domcontentloaded' })
    const next = nextConfigurationControls(page)

    await expect(next.maxLiquidityInput).toBeVisible()

    // Fee-like input renders at ~40% of its container (viewport-independent ratio).
    const inputBox = await readBox(next.maxLiquidityInput)
    const containerWidth = await next.card.evaluate((el) => {
      const body = el.querySelector('[data-slot="card-content"], .card-body') ?? el
      return body.getBoundingClientRect().width
    })
    expect(inputBox.width / containerWidth).toBeGreaterThan(0.3)
    expect(inputBox.width / containerWidth).toBeLessThan(0.5)

    // Confirmation amount field max-width 180px (fixed px — viewport independent).
    const goldenAmount = goldenEl('general', 'rskAmountInput')
    const nextAmountStyles = await readStyles(next.confirmationAmountInput)
    expectPxClose(nextAmountStyles.maxWidth, parsePx(goldenAmount.computed.maxWidth), 2)

    // Question / tooltip icon 13x13.
    const goldenIcon = goldenEl('general', 'questionIcon')
    const iconBox = await readBox(next.questionIcon)
    expectPxClose(iconBox.width, goldenIcon.rect.width, 2)
    expectPxClose(iconBox.height, goldenIcon.rect.height, 2)
  })
})
