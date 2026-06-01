import type { InitialDataPayload } from '@shared/types/initial-data'

export const loggedInFixture: InitialDataPayload = {
  loggedIn: true,
  data: {
    CredentialsSet: true,
    BaseUrl: 'http://localhost:8080',
    BtcAddress: 'tb1qloggedin',
    RskAddress: '0xloggedin',
    ProviderData: {
      id: 1,
      address: '0xprovider',
      name: 'Test LP',
      apiBaseUrl: 'http://localhost:8080',
      status: true,
      providerType: 'pegin',
    },
    ColdWallet: {
      BtcAddress: 'tb1qcold',
      RskAddress: '0xcold',
      Label: 'cold',
    },
    Configuration: {
      general: {},
      pegin: {},
      pegout: {},
    },
  },
}

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
