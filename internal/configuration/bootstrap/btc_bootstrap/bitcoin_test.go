package btc_bootstrap_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/btc_bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/stretchr/testify/require"
	"testing"
)

// nolint:funlen
func TestExternalBitcoinSources(t *testing.T) {
	tests := []struct {
		name      string
		env       environment.Environment
		expected  int
		wantError bool
		errorMsg  string
	}{
		{
			name: "Returns N bitcoin external sources",
			env: environment.Environment{
				Btc: environment.BtcEnv{
					Network: "testnet",
					BtcExtraSources: []environment.BtcExtraSource{
						// we just test with mempool because we don't have a real RPC server to connect to
						{Format: "mempool", Url: "http://mempool-source/api"},
						{Format: "mempool", Url: "http://mempool-source/testnet/api"},
					},
				},
			},
			expected:  2,
			wantError: false,
		},
		{
			name: "Returns no sources for empty config",
			env: environment.Environment{
				Btc: environment.BtcEnv{
					Network:         "testnet",
					BtcExtraSources: []environment.BtcExtraSource{},
				},
			},
			expected:  0,
			wantError: false,
		},
		{
			name: "Returns error for invalid network params",
			env: environment.Environment{
				Btc: environment.BtcEnv{
					BtcExtraSources: []environment.BtcExtraSource{
						{Format: "rpc", Url: "http://rpc-source"},
					},
				},
			},
			wantError: true,
			errorMsg:  "invalid network name",
		},
		{
			name: "Returns error for invalid RPC client",
			env: environment.Environment{
				Btc: environment.BtcEnv{
					Network: "testnet",
					BtcExtraSources: []environment.BtcExtraSource{
						{Format: "rpc", Url: "invalid://url"},
					},
				},
			},
			wantError: true,
			errorMsg:  "error creating external btc_bootstrap client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sources, err := btc_bootstrap.ExternalBitcoinSources(tt.env)

			if tt.wantError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
				require.Nil(t, sources)
			} else {
				require.NoError(t, err)
				require.Len(t, sources, tt.expected)
			}
		})
	}
}
