import type { InitialDataPayload } from '@shared/types/initial-data'

export const loggedOutFixture: InitialDataPayload = {
  loggedIn: false,
  data: {
    CredentialsSet: true,
    BaseUrl: 'http://localhost:8080',
    BtcAddress: 'tb1qexample',
    RskAddress: '0xabc',
    ProviderData: {
      id: 0,
      address: '',
      name: '',
      apiBaseUrl: '',
      status: false,
      providerType: '',
    },
    ColdWallet: {
      BtcAddress: '',
      RskAddress: '',
      Label: '',
    },
    Configuration: {
      general: {},
      pegin: {},
      pegout: {},
    },
  },
}
