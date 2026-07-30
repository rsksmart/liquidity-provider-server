import {
  confirmationRecordToRows,
  type ConfirmationRow,
  confirmationRowsToEntries,
  ConfirmationTiersEditor,
  validateConfirmationRows,
} from '@feature/management/components/config'
import { formatGeneralConfig } from '@feature/management/config/format'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it, vi } from 'vitest'

const ONE_RBTC_WEI = '1000000000000000000'
const TWO_RBTC_WEI = '2000000000000000000'
const FIVE_RBTC_WEI = '5000000000000000000'

/** Test-only parent: keeps controlled rows in sync the way a form would. */
function ConfirmationTiersHarness({
  initial,
  configKey = 'rskConfirmations',
  onChangeSpy,
  onDirty,
  errors,
}: {
  initial: ConfirmationRow[]
  configKey?: 'rskConfirmations' | 'btcConfirmations'
  onChangeSpy?: (rows: ConfirmationRow[]) => void
  onDirty?: () => void
  errors?: string[]
}) {
  const [rows, setRows] = useState(initial)
  return (
    <ConfirmationTiersEditor
      configKey={configKey}
      value={rows}
      onChange={(next) => {
        setRows(next)
        onChangeSpy?.(next)
      }}
      onDirty={onDirty}
      errors={errors}
    />
  )
}

describe('confirmationRecordToRows', () => {
  it('sorts entries ascending by ether amount and maps to display rows', () => {
    const rows = confirmationRecordToRows({
      [TWO_RBTC_WEI]: 4,
      [ONE_RBTC_WEI]: 2,
      [FIVE_RBTC_WEI]: 6,
    })

    expect(rows).toEqual([
      { amountEther: '1', confirmations: '2' },
      { amountEther: '2', confirmations: '4' },
      { amountEther: '5', confirmations: '6' },
    ])
  })

  it('returns an empty array for an empty record', () => {
    expect(confirmationRecordToRows({})).toEqual([])
  })
})

describe('confirmationRowsToEntries', () => {
  it('converts ether display rows to wei-keyed entries in order', () => {
    const entries = confirmationRowsToEntries([
      { amountEther: '2', confirmations: '4' },
      { amountEther: '1', confirmations: '2' },
    ])

    expect(entries).toEqual([
      { amount: TWO_RBTC_WEI, confirmation: 4 },
      { amount: ONE_RBTC_WEI, confirmation: 2 },
    ])
  })

  it('yields empty amount / undefined confirmation for incomplete or invalid rows', () => {
    const entries = confirmationRowsToEntries([
      { amountEther: '', confirmations: '3' },
      { amountEther: 'abc', confirmations: '3' },
      { amountEther: '1', confirmations: '' },
    ])

    expect(entries).toEqual([
      { amount: '', confirmation: 3 },
      { amount: '', confirmation: 3 },
      { amount: ONE_RBTC_WEI, confirmation: undefined },
    ])
  })
})

describe('validateConfirmationRows', () => {
  it('reports empty amount and confirmation entries with the expected copy', () => {
    const errors = validateConfirmationRows(
      [{ amountEther: '', confirmations: '' }],
      'rskConfirmations',
    )

    expect(errors).toContain('Please enter a non-empty value for "rBTC amount."')
    expect(errors).toContain('Please enter a non-empty value for "confirmations."')
  })

  it('uses the BTC amount label for btcConfirmations empties', () => {
    const errors = validateConfirmationRows(
      [{ amountEther: '', confirmations: '3' }],
      'btcConfirmations',
    )

    expect(errors).toContain('Please enter a non-empty value for "BTC amount."')
  })

  it('reports duplicate amounts with the expected copy', () => {
    const errors = validateConfirmationRows(
      [
        { amountEther: '1', confirmations: '2' },
        { amountEther: '1', confirmations: '4' },
      ],
      'rskConfirmations',
    )

    expect(errors).toContain(
      'Duplicate rBTC amounts found in rskConfirmations. Please remove duplicates before saving.',
    )
  })

  it('requires at least one complete entry', () => {
    const errors = validateConfirmationRows([], 'rskConfirmations')

    expect(errors).toContain(
      'Please provide at least one fully filled out entry for rskConfirmations.',
    )
  })

  it('reports invalid amount and non-integer confirmation values', () => {
    const errors = validateConfirmationRows(
      [
        { amountEther: 'not-a-number', confirmations: '2' },
        { amountEther: '1', confirmations: '1.5' },
      ],
      'rskConfirmations',
    )

    expect(errors).toContain(
      'Invalid input "not-a-number" for rBTC amount. Please enter a valid non-negative number.',
    )
    expect(errors).toContain(
      'Invalid input "1.5" for confirmations. Please enter a valid non-negative integer.',
    )
  })

  it('returns no errors for a valid, unique, complete set', () => {
    const errors = validateConfirmationRows(
      [
        { amountEther: '1', confirmations: '2' },
        { amountEther: '2', confirmations: '4' },
      ],
      'rskConfirmations',
    )

    expect(errors).toEqual([])
  })
})

describe('ConfirmationTiersEditor', () => {
  it('renders the config key header, unit suffix and expected testids, sorted on load', () => {
    render(
      <ConfirmationTiersHarness
        initial={confirmationRecordToRows({
          [TWO_RBTC_WEI]: 4,
          [ONE_RBTC_WEI]: 2,
          [FIVE_RBTC_WEI]: 6,
        })}
      />,
    )

    expect(
      screen.getByRole('heading', { name: 'rskConfirmations' }),
    ).toBeInTheDocument()

    const amountInputs = screen.getAllByPlaceholderText('Amount')
    expect(amountInputs.map((input) => (input as HTMLInputElement).value)).toEqual([
      '1',
      '2',
      '5',
    ])

    // Legacy confirmation testid: config-${configKey}-${index}
    expect(screen.getByTestId('config-rskConfirmations-0')).toHaveValue(2)
    expect(screen.getByTestId('config-rskConfirmations-2')).toHaveValue(6)

    // Amount fields carry the rBTC unit suffix; confirmations carry the label
    expect(screen.getAllByText('rBTC').length).toBe(3)
    expect(screen.getAllByText('confirmations').length).toBe(3)
  })

  it('uses the BTC unit suffix for btcConfirmations', () => {
    render(
      <ConfirmationTiersHarness
        configKey="btcConfirmations"
        initial={[{ amountEther: '1', confirmations: '2' }]}
      />,
    )

    expect(screen.getByText('BTC')).toBeInTheDocument()
    expect(screen.queryByText('rBTC')).not.toBeInTheDocument()
  })

  it('adds an empty entry and marks dirty when Add Entry is clicked', async () => {
    const user = userEvent.setup()
    const onChangeSpy = vi.fn()
    const onDirty = vi.fn()
    render(
      <ConfirmationTiersHarness
        initial={[{ amountEther: '1', confirmations: '2' }]}
        onChangeSpy={onChangeSpy}
        onDirty={onDirty}
      />,
    )

    await user.click(screen.getByRole('button', { name: 'Add Entry' }))

    expect(onChangeSpy).toHaveBeenLastCalledWith([
      { amountEther: '1', confirmations: '2' },
      { amountEther: '', confirmations: '' },
    ])
    expect(onDirty).toHaveBeenCalled()
    expect(screen.getAllByPlaceholderText('Amount')).toHaveLength(2)
  })

  it('removes an entry and marks dirty when Remove is clicked', async () => {
    const user = userEvent.setup()
    const onChangeSpy = vi.fn()
    const onDirty = vi.fn()
    render(
      <ConfirmationTiersHarness
        initial={[
          { amountEther: '1', confirmations: '2' },
          { amountEther: '2', confirmations: '4' },
        ]}
        onChangeSpy={onChangeSpy}
        onDirty={onDirty}
      />,
    )

    const removeButtons = screen.getAllByRole('button', { name: 'Remove' })
    await user.click(removeButtons[0])

    expect(onChangeSpy).toHaveBeenLastCalledWith([
      { amountEther: '2', confirmations: '4' },
    ])
    expect(onDirty).toHaveBeenCalled()
    expect(screen.getAllByPlaceholderText('Amount')).toHaveLength(1)
  })

  it('emits edited amount and confirmation values and marks dirty', async () => {
    const user = userEvent.setup()
    const onChangeSpy = vi.fn()
    const onDirty = vi.fn()
    render(
      <ConfirmationTiersHarness
        initial={[{ amountEther: '', confirmations: '' }]}
        onChangeSpy={onChangeSpy}
        onDirty={onDirty}
      />,
    )

    await user.type(screen.getByPlaceholderText('Amount'), '3')
    expect(onChangeSpy).toHaveBeenLastCalledWith([
      { amountEther: '3', confirmations: '' },
    ])

    await user.type(screen.getByTestId('config-rskConfirmations-0'), '7')
    expect(onChangeSpy).toHaveBeenLastCalledWith([
      { amountEther: '3', confirmations: '7' },
    ])
    expect(onDirty).toHaveBeenCalled()
  })

  it('renders externally-supplied validation errors', () => {
    render(
      <ConfirmationTiersHarness
        initial={[{ amountEther: '1', confirmations: '2' }]}
        errors={[
          'Duplicate rBTC amounts found in rskConfirmations. Please remove duplicates before saving.',
        ]}
      />,
    )

    expect(screen.getByRole('alert')).toHaveTextContent(
      'Duplicate rBTC amounts found in rskConfirmations. Please remove duplicates before saving.',
    )
  })

  it('produces a save map whose keys match displayed order via formatGeneralConfig', async () => {
    const user = userEvent.setup()
    let latestRows: ConfirmationRow[] = []
    const onChangeSpy = vi.fn((rows: ConfirmationRow[]) => {
      latestRows = rows
    })
    render(
      <ConfirmationTiersHarness
        initial={confirmationRecordToRows({
          [ONE_RBTC_WEI]: 2,
          [TWO_RBTC_WEI]: 4,
        })}
        onChangeSpy={onChangeSpy}
      />,
    )

    // Append a lower amount at the end; display order is by position, not value.
    await user.click(screen.getByRole('button', { name: 'Add Entry' }))
    const amountInputs = screen.getAllByPlaceholderText('Amount')
    await user.type(amountInputs[2], '0.5')
    await user.type(screen.getByTestId('config-rskConfirmations-2'), '1')

    const formatted = formatGeneralConfig({
      rskConfirmations: confirmationRowsToEntries(latestRows),
      btcConfirmations: [],
    })

    // Keys follow the displayed row order (1, 2, then the appended 0.5), not
    // numeric sorting of the wei keys.
    expect(Object.keys(formatted.rskConfirmations)).toEqual([
      ONE_RBTC_WEI,
      TWO_RBTC_WEI,
      '500000000000000000',
    ])
  })
})
