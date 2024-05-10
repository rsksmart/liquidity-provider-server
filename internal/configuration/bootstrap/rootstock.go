package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

func Rootstock(ctx context.Context, env environment.RskEnv) (*rootstock.RskClient, error) {
	var err error
	var parsedUrl *url.URL
	var client *ethclient.Client
	var rpcClient *rpc.Client

	log.Info("Connecting to RSK node on ", env.Endpoint)
	if parsedUrl, err = url.Parse(env.Endpoint); err != nil {
		return nil, err
	}

	switch parsedUrl.Scheme {
	case "http", "https":
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.DisableKeepAlives = true

		httpClient := new(http.Client)
		httpClient.Transport = transport

		if rpcClient, err = rpc.DialOptions(ctx, env.Endpoint, rpc.WithHTTPClient(httpClient)); err != nil {
			return nil, err
		}

		client = ethclient.NewClient(rpcClient)
	default:
		return nil, errors.New("unknown scheme for rsk connection string")
	}

	log.Debug("Verifying connection to RSK node")
	id, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	log.Debug("Connection verified")
	if env.ChainId != id.Uint64() {
		return nil, fmt.Errorf("chain id mismatch; expected chain id: %v, rsk node chain id: %v", env.ChainId, id)
	}
	return rootstock.NewRskClient(client), nil
}

func RootstockAccount(
	rskEnv environment.RskEnv,
	btcEnv environment.BtcEnv,
	secrets secrets.DerivativeWalletSecrets) (*account.RskAccount, error) {
	networkParams, err := btcEnv.GetNetworkParams()
	if err != nil {
		return nil, err
	}
	return account.GetRskAccountWithDerivation(account.CreationWithDerivationArgs{
		CreationArgs: account.CreationArgs{
			KeyDir:        "geth_keystore",
			AccountNum:    rskEnv.AccountNumber,
			EncryptedJson: secrets.EncryptedJson,
			Password:      secrets.EncryptedJsonPassword,
		},
		BtcParams: networkParams,
	})
}
