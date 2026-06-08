import type {
  FullConfiguration,
  InitialDataPayload,
  ProviderType,
} from '@shared/types/initial-data'

const emptyConfiguration: FullConfiguration = {
  general: {
    rskConfirmations: {},
    btcConfirmations: {},
    publicLiquidityCheck: false,
    maxLiquidity: null,
    reimbursementWindowBlocks: 0,
    excessTolerance: {
      isFixed: false,
      percentageValue: '0',
      fixedValue: '0',
    },
  },
  pegin: {
    timeForDeposit: 0,
    callTime: 0,
    penaltyFee: '0',
    fixedFee: '0',
    feePercentage: '0',
    maxValue: '0',
    minValue: '0',
  },
  pegout: {
    timeForDeposit: 0,
    expireTime: 0,
    penaltyFee: '0',
    fixedFee: '0',
    feePercentage: '0',
    maxValue: '0',
    minValue: '0',
    expireBlocks: 0,
    bridgeTransactionMin: '0',
  },
}

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
      providerType: 0 as ProviderType,
    },
    ColdWallet: {
      BtcAddress: 'tb1qcold',
      RskAddress: '0xcold',
      Label: 'cold',
    },
    Configuration: emptyConfiguration,
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
      providerType: 0 as ProviderType,
    },
    ColdWallet: {
      BtcAddress: '',
      RskAddress: '',
      Label: '',
    },
    Configuration: emptyConfiguration,
  },
}
