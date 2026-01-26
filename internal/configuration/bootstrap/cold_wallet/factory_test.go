package cold_wallet_test

import (
	"encoding/json"
	cwFactory "github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	rpc := &mocks.BtcRpcMock{}
	validArgs := cold_wallet.StaticColdWalletArgs{BtcAddress: test.AnyBtcAddress, RskAddress: test.AnyRskAddress}

	validConfigBytes, err := json.Marshal(validArgs)
	require.NoError(t, err)

	tests := []struct {
		name        string
		config      secrets.ColdWalletConfiguration
		expectErr   bool
		errContains string
	}{
		{
			name: "static cold wallet - valid config",
			config: secrets.ColdWalletConfiguration{
				Type:          "static",
				Configuration: validConfigBytes,
			},
			expectErr: false,
		},
		{
			name: "static cold wallet - invalid json",
			config: secrets.ColdWalletConfiguration{
				Type:          "static",
				Configuration: []byte(`{invalid-json}`),
			},
			expectErr:   true,
			errContains: "invalid static cold wallet configuration",
		},
		{
			name: "unknown cold wallet type",
			config: secrets.ColdWalletConfiguration{
				Type:          "dynamic",
				Configuration: validConfigBytes,
			},
			expectErr:   true,
			errContains: "unknown cold wallet type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := cwFactory.Create(blockchain.Rpc{Btc: rpc}, tt.config)

			if tt.expectErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
				require.Nil(t, wallet)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, wallet)
		})
	}
}
