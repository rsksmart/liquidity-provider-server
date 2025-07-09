package usecases_test

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	u "github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const id = "anyUseCase"

type rpcMock struct {
	mock.Mock
	blockchain.RootstockRpcServer
}

func (m *rpcMock) EstimateGas(ctx context.Context, addr string, value *entities.Wei, data []byte) (*entities.Wei, error) {
	args := m.Called(ctx, addr, value, data)
	return args.Get(0).(*entities.Wei), args.Error(1) // nolint:errcheck
}

type bridgeMock struct {
	mock.Mock
	blockchain.RootstockBridge
}

func (m *bridgeMock) GetMinimumLockTxValue() (*entities.Wei, error) {
	args := m.Called()
	return args.Get(0).(*entities.Wei), args.Error(1) // nolint:errcheck
}

func TestCalculateDaoAmounts(t *testing.T) {
	type testArgs struct {
		value      *entities.Wei
		percentage uint64
	}
	rpc := rpcMock{}
	rpc.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(entities.NewWei(500000000000000), nil)

	ctx := context.Background()
	feeCollectorAddress := "0x1234"

	cases := test.Table[testArgs, u.DaoAmounts]{
		{
			Value:  testArgs{entities.NewWei(1000000000000000000), 0},
			Result: u.DaoAmounts{DaoFeeAmount: entities.NewWei(0), DaoGasAmount: entities.NewWei(0)},
		},
		{
			Value: testArgs{entities.NewWei(500000000000000000), 50},
			Result: u.DaoAmounts{
				DaoGasAmount: entities.NewWei(500000000000000),
				DaoFeeAmount: entities.NewWei(250000000000000000),
			},
		},
		{
			Value: testArgs{entities.NewWei(6000000000000000000), 1},
			Result: u.DaoAmounts{
				DaoGasAmount: entities.NewWei(500000000000000),
				DaoFeeAmount: entities.NewWei(60000000000000000),
			},
		},
		{
			Value: testArgs{entities.NewWei(7700000000000000000), 17},
			Result: u.DaoAmounts{
				DaoGasAmount: entities.NewWei(500000000000000),
				DaoFeeAmount: entities.NewWei(1309000000000000000),
			},
		},
	}

	test.RunTable(t, cases, func(args testArgs) u.DaoAmounts {
		amounts, err := u.CalculateDaoAmounts(ctx, &rpc, args.value, args.percentage, feeCollectorAddress)
		require.NoError(t, err)
		return amounts
	})

}

func TestCalculateDaoAmounts_Fail(t *testing.T) {
	ctx := context.Background()
	rpc := rpcMock{}
	rpc.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(entities.NewWei(0), errors.New("some error"))
	result, err := u.CalculateDaoAmounts(ctx, &rpc, entities.NewUWei(500000000000000), 1, "0x1234")
	require.Equal(t, u.DaoAmounts{}, result)
	require.Error(t, err)
}

func TestValidateMinLockValue(t *testing.T) {
	var oneBtcInSatoshi uint64 = 1 * bitcoin.BtcToSatoshi
	var useCase u.UseCaseId = "anyUseCase"
	bridge := &bridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.SatoshiToWei(oneBtcInSatoshi), nil)

	err := u.ValidateMinLockValue(useCase, bridge, entities.SatoshiToWei(oneBtcInSatoshi))
	require.NoError(t, err)

	err = u.ValidateMinLockValue(useCase, bridge, entities.SatoshiToWei(oneBtcInSatoshi+1))
	require.NoError(t, err)

	value := new(entities.Wei).Sub(entities.SatoshiToWei(oneBtcInSatoshi), entities.NewWei(1))
	err = u.ValidateMinLockValue(useCase, bridge, value)
	require.Error(t, err)
	assert.Equal(t, "anyUseCase: requested amount below bridge's min transaction value. Args: {\"minimum\":\"1000000000000000000\",\"value\":\"999999999999999999\"}", err.Error())
}

func TestSignConfiguration(t *testing.T) {
	var (
		signature = []byte{1, 2, 3}
		hash      = []byte{4, 5, 6}
	)
	configuration := liquidity_provider.DefaultPeginConfiguration()
	wallet := &mocks.RskWalletMock{}
	wallet.On("SignBytes", mock.Anything).Return(signature, nil)
	hashFunctionMock := &mocks.HashMock{}
	hashFunctionMock.On("Hash", mock.Anything).Return(hash)
	signed, err := u.SignConfiguration(id, wallet, hashFunctionMock.Hash, configuration)
	require.NoError(t, err)
	assert.Equal(t, entities.Signed[liquidity_provider.PeginConfiguration]{
		Value:     configuration,
		Hash:      "040506",
		Signature: "010203",
	}, signed)
}

func TestSignConfiguration_SignatureError(t *testing.T) {
	wallet := &mocks.RskWalletMock{}
	wallet.On("SignBytes", mock.Anything).Return(nil, assert.AnError)
	hashFunctionMock := &mocks.HashMock{}
	hashFunctionMock.On("Hash", mock.Anything).Return([]byte{1})
	configuration := liquidity_provider.DefaultPeginConfiguration()
	signed, err := u.SignConfiguration(id, wallet, hashFunctionMock.Hash, configuration)
	require.Equal(t, entities.Signed[liquidity_provider.PeginConfiguration]{}, signed)
	require.Error(t, err)
}

func TestRegisterCoinbaseTransaction(t *testing.T) {
	tx := blockchain.BitcoinTransactionInformation{
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{test.AnyAddress: {entities.NewWei(1)}},
		HasWitness:    true,
		Hash:          test.AnyHash,
	}
	coinbaseInfo := blockchain.BtcCoinbaseTransactionInformation{
		BtcTxSerialized:      utils.MustGetRandomBytes(32),
		BlockHash:            utils.To32Bytes(utils.MustGetRandomBytes(32)),
		BlockHeight:          big.NewInt(500),
		SerializedPmt:        utils.MustGetRandomBytes(64),
		WitnessMerkleRoot:    utils.To32Bytes(utils.MustGetRandomBytes(32)),
		WitnessReservedValue: utils.To32Bytes(utils.MustGetRandomBytes(32)),
	}
	t.Run("Should return if tx does not have witness data", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		rpc := &mocks.BtcRpcMock{}
		txWithoutWitness := tx
		txWithoutWitness.HasWitness = false
		err := u.RegisterCoinbaseTransaction(rpc, bridge, txWithoutWitness)
		require.NoError(t, err)
		bridge.AssertNotCalled(t, "RegisterCoinbaseTransaction")
		rpc.AssertNotCalled(t, "GetCoinbaseInformation")
	})
	t.Run("Should handle error fetching the coinbase information", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		rpc := &mocks.BtcRpcMock{}
		rpc.On("GetCoinbaseInformation", test.AnyHash).Return(blockchain.BtcCoinbaseTransactionInformation{}, assert.AnError)
		err := u.RegisterCoinbaseTransaction(rpc, bridge, tx)
		require.Error(t, err)
		bridge.AssertNotCalled(t, "RegisterCoinbaseTransaction")
		rpc.AssertExpectations(t)
	})
	t.Run("Should handle error registering the transaction", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		rpc := &mocks.BtcRpcMock{}
		rpc.On("GetCoinbaseInformation", test.AnyHash).Return(coinbaseInfo, nil)
		bridge.On("RegisterBtcCoinbaseTransaction", coinbaseInfo).Return("", assert.AnError)
		err := u.RegisterCoinbaseTransaction(rpc, bridge, tx)
		require.Error(t, err)
		bridge.AssertExpectations(t)
		rpc.AssertExpectations(t)
	})
	t.Run("Should register a coinbase tx successfully", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		rpc := &mocks.BtcRpcMock{}
		rpc.On("GetCoinbaseInformation", test.AnyHash).Return(coinbaseInfo, nil)
		bridge.On("RegisterBtcCoinbaseTransaction", coinbaseInfo).Return(test.AnyHash, nil)
		err := u.RegisterCoinbaseTransaction(rpc, bridge, tx)
		require.NoError(t, err)
		bridge.AssertExpectations(t)
		rpc.AssertExpectations(t)
	})
}

// nolint:funlen
func TestValidateBridgeUtxoMin(t *testing.T) {
	const (
		address       = "2N991MLUtYHfHzLQgtNfK9NtUVUSEe9Ncaf"
		txHash        = "9fa7040fbe1e442e970c231aaf830abf62c83b165fc436b1ce6e385b6ce40e59"
		confirmations = 10
	)
	t.Run("should fail if one of the UTXO is under the bridge min", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
		tx := blockchain.BitcoinTransactionInformation{
			Hash: txHash, Confirmations: confirmations,
			Outputs: map[string][]*entities.Wei{
				address: {entities.NewWei(1000), entities.NewWei(999), entities.NewWei(3000)},
			},
		}
		err := u.ValidateBridgeUtxoMin(bridge, tx, address)
		require.ErrorContains(t, err, "not all the UTXOs are above the min lock value")
		require.ErrorIs(t, err, u.TxBelowMinimumError)
		bridge.AssertExpectations(t)
	})
	t.Run("should fail if all of the UTXO is under the bridge min", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
		tx := blockchain.BitcoinTransactionInformation{
			Hash: txHash, Confirmations: confirmations,
			Outputs: map[string][]*entities.Wei{
				address: {entities.NewWei(1), entities.NewWei(999), entities.NewWei(100)},
			},
		}
		err := u.ValidateBridgeUtxoMin(bridge, tx, address)
		require.ErrorContains(t, err, "not all the UTXOs are above the min lock value")
		require.ErrorIs(t, err, u.TxBelowMinimumError)
		bridge.AssertExpectations(t)
	})
	t.Run("should not fail if all the UTXO is above the bridge min", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Twice()
		tx := blockchain.BitcoinTransactionInformation{
			Hash: txHash, Confirmations: confirmations,
			Outputs: map[string][]*entities.Wei{
				address: {entities.NewWei(1000), entities.NewWei(1000), entities.NewWei(1001)},
			},
		}
		err := u.ValidateBridgeUtxoMin(bridge, tx, address)
		require.NoError(t, err)

		tx = blockchain.BitcoinTransactionInformation{
			Hash: txHash, Confirmations: confirmations,
			Outputs: map[string][]*entities.Wei{address: {entities.NewWei(1000)}},
		}
		err = u.ValidateBridgeUtxoMin(bridge, tx, address)
		require.NoError(t, err)

		bridge.AssertExpectations(t)
	})
	t.Run("should return error if call to the bridge fails", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(nil, assert.AnError).Once()
		tx := blockchain.BitcoinTransactionInformation{
			Hash: txHash, Confirmations: confirmations,
			Outputs: map[string][]*entities.Wei{
				address: {entities.NewWei(1000), entities.NewWei(1000), entities.NewWei(1001)},
			},
		}
		err := u.ValidateBridgeUtxoMin(bridge, tx, address)
		require.Error(t, err)
		require.NotErrorIs(t, err, u.TxBelowMinimumError)
		bridge.AssertExpectations(t)
	})
	t.Run("should fail if there is no UTXO for a given address", func(t *testing.T) {
		bridge := &mocks.BridgeMock{}
		bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000), nil).Once()
		tx := blockchain.BitcoinTransactionInformation{
			Hash: txHash, Confirmations: confirmations,
			Outputs: map[string][]*entities.Wei{
				"other-address": {entities.NewWei(1000), entities.NewWei(1000), entities.NewWei(1001)},
			},
		}
		err := u.ValidateBridgeUtxoMin(bridge, tx, address)
		require.ErrorContains(t, err, "no UTXO directed to address 2N991MLUtYHfHzLQgtNfK9NtUVUSEe9Ncaf present in transaction")
		require.ErrorIs(t, err, u.TxBelowMinimumError)
		bridge.AssertExpectations(t)
	})
}

// nolint:funlen
func TestRecoverSignerAddress(t *testing.T) {
	testCases := []struct {
		name            string
		quoteHash       string
		signature       string
		expectedAddress string
		expectError     bool
	}{
		{
			name:            "valid signature and hash",
			quoteHash:       "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf",
			signature:       "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1c",
			expectedAddress: "0x233845a26a4dA08E16218e7B401501D048670674",
			expectError:     false,
		},
		{
			name:            "valid signature and hash mainnet",
			quoteHash:       "249e59bf1a92a867629f111222fa63946370640030faaa5fdb79744fd0539f81",
			signature:       "229ac43839c3a66520f3f96b1b8f124608dd652c0c9071629a9d8457c292b0257eb3300d2c3b7f5be30c35c9bd15b05582e00070052a4006b8768afab853872a1c",
			expectedAddress: "0x82a06ebdb97776a2da4041df8f2b2ea8d3257852",
			expectError:     false,
		},
		{
			name:            "valid signature and hash testnet",
			quoteHash:       "3431bf06ab6ebde3e1297d5eaa5563edf500fea7dcdc466be1a8492e78dd6d80",
			signature:       "a316d2dd76e2a325efdfee6b6686fa2c73aedad0090a1cbf286a5799f03548892943d436af4543c7b33548a6e24c1546ed51b844e0b9ab9b8da37aa30dd8f7031b",
			expectedAddress: "0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b",
			expectError:     false,
		},
		{
			name:            "valid signature and hash dev",
			quoteHash:       "e795e119e597f411d379197f11680840f995128a1502024dee5d8a1480658ed5",
			signature:       "72a8cfd0da1b4c008b209b404574469a1c935c5e6252ef3d53688320d2f0060449c0033709f3460ab34bd17aa6141d8e5485476793382391ca33d8cc5c6895831c",
			expectedAddress: "0xdfcf32644e6cc5badd1188cddf66f66e21b24375",
			expectError:     false,
		},
		{
			name:            "invalid hash format",
			quoteHash:       "invalid-hex",
			signature:       "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1c",
			expectedAddress: "",
			expectError:     true,
		},
		{
			name:            "invalid signature format",
			quoteHash:       "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf",
			signature:       "invalid-signature",
			expectedAddress: "",
			expectError:     true,
		},
		{
			name:            "signature too short",
			quoteHash:       "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf",
			signature:       "5f1a", // Too short
			expectedAddress: "",
			expectError:     true,
		},
		{
			name:            "hash too short",
			quoteHash:       "c8d4", // Too short
			signature:       "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1c",
			expectedAddress: "",
			expectError:     true,
		},
		{
			name:            "empty hash",
			quoteHash:       "",
			signature:       "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1c",
			expectedAddress: "",
			expectError:     true,
		},
		{
			name:            "empty signature",
			quoteHash:       "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf",
			signature:       "",
			expectedAddress: "",
			expectError:     true,
		},
		{
			name:      "signature with wrong recovery ID",
			quoteHash: "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf",
			// This is the valid signature with the last byte changed from 1c to 1d (invalid recovery ID)
			signature:       "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1d",
			expectedAddress: "",
			expectError:     true,
		},
		{
			name:      "signature with recovery ID = 27 (0x1b)",
			quoteHash: "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf",
			// This is a signature with recovery ID 27 (0x1b) which should adjust to 0
			signature:       "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1b",
			expectedAddress: "0xf719893448f705385d9f7192d2f72c894767d9dc", // Actual recovered address for this signature
			expectError:     false,
		},
		{
			name:      "signature with recovery ID = 28 (0x1c)",
			quoteHash: "c8d4ad8d5d717371b92950cbe43a6a4e891cf27bcd7603c988595866944bd9cf",
			// Valid signature with recovery ID 28 (0x1c) which should adjust to 1
			signature:       "5f1a75f55f92c23be729adfb9eff21a00feb1ba99c5e7c2ea9c98a6430e3958f2db856b6260730b6aeeab83571bbafb77730ef1a9cb3a09ce3fa07065c8b200d1c",
			expectedAddress: "0x233845a26a4dA08E16218e7B401501D048670674",
			expectError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recoveredAddress, err := u.RecoverSignerAddress(tc.quoteHash, tc.signature)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t,
					strings.ToLower(tc.expectedAddress),
					strings.ToLower(recoveredAddress),
					"Address recovery should match the expected signer address")
			}
		})
	}
}
