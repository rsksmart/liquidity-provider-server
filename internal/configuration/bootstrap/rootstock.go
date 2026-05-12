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
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

func Rootstock(ctx context.Context, env environment.Environment) (*rootstock.RskClient, error) {
	return createClient(ctx, env.Rsk.Endpoint, env.Rsk.ChainId)
}

func createClient(ctx context.Context, endpoint string, chainId uint64) (*rootstock.RskClient, error) {
	var err error
	var parsedUrl *url.URL
	var client *ethclient.Client
	var rpcClient *rpc.Client

	log.Info("Connecting to RSK node on ", endpoint)
	if parsedUrl, err = url.Parse(endpoint); err != nil {
		return nil, err
	}

	switch parsedUrl.Scheme {
	case "http", "https":
		defaultTransport, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			return nil, errors.New("failed to get default transport")
		}
		transport := defaultTransport.Clone()
		transport.DisableKeepAlives = true

		httpClient := new(http.Client)
		httpClient.Transport = transport

		if rpcClient, err = rpc.DialOptions(ctx, endpoint, rpc.WithHTTPClient(httpClient)); err != nil {
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
	if chainId != id.Uint64() {
		return nil, fmt.Errorf("chain id mismatch; expected chain id: %v, rsk node chain id: %v", chainId, id)
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

func ExternalRskSources(ctx context.Context, env environment.Environment) ([]blockchain.RootstockRpcServer, error) {
	sources := make([]blockchain.RootstockRpcServer, 0)
	for _, endpoint := range env.Rsk.RskExtraSources {
		client, err := createClient(ctx, endpoint, env.Rsk.ChainId)
		if err != nil {
			return nil, fmt.Errorf("failed to create RSK client for endpoint %s: %w", endpoint, err)
		}
		log.Info("Connected to external RSK node at ", endpoint)
		sources = append(sources, rootstock.NewRskjRpcServer(client, rootstock.DefaultRetryParams))
	}
	return sources, nil
}
