package scripts_test

import (
	"bytes"
	"context"
	"flag"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"path/filepath"
	"testing"
)

func TestGetWallet(t *testing.T) {
	t.Run("should return wallet", func(t *testing.T) {
		ctx := context.Background()
		keystorePath := filepath.Join("../../../", "docker-compose/localstack/local-key.json")
		env := environment.Environment{
			SecretSource:     "env",
			WalletManagement: "native",
			Rsk: environment.RskEnv{
				KeystoreFile:     keystorePath,
				KeystorePassword: "test",
			},
			Btc: environment.BtcEnv{Network: "regtest"},
		}
		rskClient := &rootstock.RskClient{}

		wallet, err := scripts.GetWallet(ctx, env, rskClient)
		require.NoError(t, err)
		require.NotNil(t, wallet)
	})

	t.Run("should return error", func(t *testing.T) {
		ctx := context.Background()
		rskClient := &rootstock.RskClient{}

		result, err := scripts.GetWallet(ctx, environment.Environment{}, rskClient)
		assert.Nil(t, result)
		require.Error(t, err)

		result, err = scripts.GetWallet(ctx, environment.Environment{SecretSource: "env"}, rskClient)
		assert.Nil(t, result)
		require.Error(t, err)
	})
}

func TestCreateLiquidityBridgeContract(t *testing.T) {
	t.Run("should return contract", func(t *testing.T) {
		keystorePath := filepath.Join("../../../", "docker-compose/localstack/local-key.json")
		env := environment.Environment{
			SecretSource:     "env",
			WalletManagement: "native",
			Rsk: environment.RskEnv{
				KeystoreFile:     keystorePath,
				KeystorePassword: "test",
			},
			Btc: environment.BtcEnv{Network: "regtest"},
		}
		factoryMock := func(ctx context.Context, env environment.RskEnv) (*rootstock.RskClient, error) {
			return &rootstock.RskClient{}, nil
		}
		contract, err := scripts.CreateLiquidityBridgeContract(context.Background(), factoryMock, env)
		require.NoError(t, err)
		require.NotNil(t, contract)
	})
}

func TestSetUsageMessage(t *testing.T) {
	msg := "Test usage message"
	scripts.SetUsageMessage(msg)

	buff := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buff)

	flag.Usage()

	readBytes, err := io.ReadAll(buff)
	require.NoError(t, err)
	assert.Contains(t, string(readBytes), msg)
}
