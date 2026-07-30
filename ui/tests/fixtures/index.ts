import type {
  FullConfiguration,
  InitialDataPayload,
  WireFullConfiguration,
  WireInitialDataPayload,
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
      providerType: 0,
    },
    ColdWallet: {
      BtcAddress: 'tb1qcold',
      RskAddress: '0xcold',
      Label: 'cold',
    },
    Configuration: emptyConfiguration,
  },
}

/**
 * Configuration exactly as a live LPS marshals it (regtest defaults captured
 * from `/management/next/management`): every wei and percentage field is a JSON
 * number, and large wei values exceed what `Number` prints without an exponent.
 */
export const wireConfigurationFixture: WireFullConfiguration = {
  general: {
    rskConfirmations: { '100000000000000000': 4, '2000000000000000000': 20 },
    btcConfirmations: { '100000000000000000': 2, '2000000000000000000': 10 },
    publicLiquidityCheck: true,
    maxLiquidity: 2e21,
    reimbursementWindowBlocks: 100,
    excessTolerance: {
      isFixed: false,
      percentageValue: 15,
      fixedValue: 0,
    },
  },
  pegin: {
    timeForDeposit: 3600,
    callTime: 7200,
    penaltyFee: 1000000000000000,
    fixedFee: 200000000000000,
    feePercentage: 0.33,
    maxValue: 10000000000000000000,
    minValue: 600000000000000000,
  },
  pegout: {
    timeForDeposit: 3600,
    expireTime: 28800,
    penaltyFee: 1000000000000000,
    fixedFee: 0,
    feePercentage: 0,
    maxValue: 10000000000000000000,
    minValue: 600000000000000000,
    expireBlocks: 1000,
    bridgeTransactionMin: 1500000000000000000,
  },
}

export const wireLoggedInFixture: WireInitialDataPayload = {
  ...loggedInFixture,
  data: { ...loggedInFixture.data, Configuration: wireConfigurationFixture },
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
      providerType: 0,
    },
    ColdWallet: {
      BtcAddress: '',
      RskAddress: '',
      Label: '',
    },
    Configuration: emptyConfiguration,
  },
}
