package bootstrap

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
)

const (
	unknownBtcdVersion = -1
)

type CreatedClient struct {
	Client btcclient.ClientAdapter
	Params *chaincfg.Params
	Config rpcclient.ConnConfig
}

func BitcoinWallet(env environment.BtcEnv, walletId string) (*bitcoin.Connection, error) {
	if walletId == "" {
		return nil, errors.New("walletId cannot be empty")
	}
	endpoint := fmt.Sprintf("%s/wallet/%s", env.Endpoint, walletId)
	params, err := env.GetNetworkParams()
	if err != nil {
		return nil, err
	}
	createdClient, err := createBitcoinClient(params, env.Username, env.Password, endpoint)
	if err != nil {
		return nil, err
	}
	return bitcoin.NewWalletConnection(createdClient.Params, createdClient.Client, walletId), nil
}

func Bitcoin(env environment.BtcEnv) (*bitcoin.Connection, error) {
	params, err := env.GetNetworkParams()
	if err != nil {
		return nil, err
	}
	createdClient, err := createBitcoinClient(params, env.Username, env.Password, env.Endpoint)
	if err != nil {
		return nil, err
	}
	conn := bitcoin.NewConnection(createdClient.Params, createdClient.Client)
	return conn, nil
}

func ExternalBitcoinClients(env environment.Environment) ([]*bitcoin.Connection, error) {
	var createdClient CreatedClient
	clients := make([]*bitcoin.Connection, len(env.Btc.BtcExtraSources))
	params, err := env.Btc.GetNetworkParams()
	if err != nil {
		return nil, err
	}
	for i, source := range env.Btc.BtcExtraSources {
		createdClient, err = createBitcoinClient(params, "", "", source)
		if err != nil {
			return nil, fmt.Errorf("error creating external bitcoin client for %s: %w", source, err)
		}
		clients[i] = bitcoin.NewConnection(createdClient.Params, createdClient.Client)
	}
	return clients, nil
}

func createBitcoinClient(networkParams *chaincfg.Params, user, password, host string) (CreatedClient, error) {
	var params *chaincfg.Params
	log.Info("Connecting to BTC node at ", host, "...")

	config := rpcclient.ConnConfig{
		Host:   host,
		User:   user,
		Pass:   password,
		Params: networkParams.Name,
		// Rationale why this is disabled: https://en.bitcoin.it/wiki/Enabling_SSL_on_original_client_daemon
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	client, err := rpcclient.New(&config, nil)
	if err != nil {
		return CreatedClient{}, fmt.Errorf("RPC client error: %w", err)
	}

	version, err := checkBtcdVersion(client)
	if err != nil {
		return CreatedClient{}, err
	}

	if version == unknownBtcdVersion {
		log.Warn("unable to detect btcd version, but it is up and running")
	} else {
		log.Debugf("detected btcd version: %v\n", version)
	}
	return CreatedClient{
		Client: btcclient.NewBtcSuiteClientAdapter(config, client),
		Params: params,
		Config: config,
	}, nil
}

func checkBtcdVersion(c *rpcclient.Client) (int32, error) {
	var networkErr *btcjson.RPCError
	info, err := c.GetNetworkInfo()
	if err == nil {
		return info.Version, nil
	} else if errors.As(err, &networkErr) {
		if networkErr.Code != btcjson.ErrRPCMethodNotFound.Code {
			return 0, fmt.Errorf("unable to detect btcd version: %w", networkErr)
		}
		return unknownBtcdVersion, nil
	} else {
		return 0, fmt.Errorf("unable to detect btcd version: %w", err)
	}
}
