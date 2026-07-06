import type { ManagementPostBodies } from '@api/management/types/post-bodies'
import { describe, expectTypeOf, it } from 'vitest'

type AddCollateralBody = ManagementPostBodies['/pegin/addCollateral']

describe('apiFetch.post body types', () => {
  it('requires numeric add-collateral amount', () => {
    expectTypeOf<AddCollateralBody>().toEqualTypeOf<{ amount: number }>()
    expectTypeOf({ amount: '1000000000000000000' }).not.toEqualTypeOf<AddCollateralBody>()
    expectTypeOf({ amount: 1_000_000_000_000_000_000 }).toEqualTypeOf<AddCollateralBody>()
  })
})
