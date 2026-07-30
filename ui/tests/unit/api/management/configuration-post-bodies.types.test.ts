import type {
  SetGeneralConfigurationRequest,
  SetPeginConfigurationRequest,
  SetPegoutConfigurationRequest,
} from '@api/management/types/configuration'
import type { ManagementPostBodies } from '@api/management/types/post-bodies'
import { describe, expectTypeOf, it } from 'vitest'

describe('configuration POST body types', () => {
  it('registers each configuration endpoint with a { configuration } envelope', () => {
    expectTypeOf<
      ManagementPostBodies['/configuration']
    >().toEqualTypeOf<SetGeneralConfigurationRequest>()
    expectTypeOf<
      ManagementPostBodies['/pegin/configuration']
    >().toEqualTypeOf<SetPeginConfigurationRequest>()
    expectTypeOf<
      ManagementPostBodies['/pegout/configuration']
    >().toEqualTypeOf<SetPegoutConfigurationRequest>()
  })

  it('wraps the section config under the configuration key', () => {
    expectTypeOf<SetGeneralConfigurationRequest>().toHaveProperty('configuration')
    expectTypeOf<SetPeginConfigurationRequest>().toHaveProperty('configuration')
    expectTypeOf<SetPegoutConfigurationRequest>().toHaveProperty('configuration')
  })
})
