/**
 * Configuration save request bodies.
 * Each section POST wraps its payload in `{ configuration: … }`.
 */

import type {
  GeneralConfiguration,
  PeginConfiguration,
  PegoutConfiguration,
} from '@shared/types/initial-data'

export interface SetGeneralConfigurationRequest {
  configuration: GeneralConfiguration
}

export interface SetPeginConfigurationRequest {
  configuration: PeginConfiguration
}

export interface SetPegoutConfigurationRequest {
  configuration: PegoutConfiguration
}

/** Union of all section save bodies (all share the `{ configuration }` shape). */
export type SetConfigurationRequest =
  | SetGeneralConfigurationRequest
  | SetPeginConfigurationRequest
  | SetPegoutConfigurationRequest
