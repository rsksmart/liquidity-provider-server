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
	createdClient, err := createBitcoinClient(env, endpoint)
	if err != nil {
		return nil, err
	}
	return bitcoin.NewWalletConnection(createdClient.Params, createdClient.Client, walletId), nil
}

func Bitcoin(env environment.BtcEnv) (*bitcoin.Connection, error) {
	createdClient, err := createBitcoinClient(env, env.Endpoint)
	if err != nil {
		return nil, err
	}
	conn := bitcoin.NewConnection(createdClient.Params, createdClient.Client)
	return conn, nil
}

func createBitcoinClient(env environment.BtcEnv, host string) (CreatedClient, error) {
	var params *chaincfg.Params
	log.Info("Connecting to BTC node at ", host, "...")

	params, err := env.GetNetworkParams()
	if err != nil {
		return CreatedClient{}, err
	}

	config := rpcclient.ConnConfig{
		Host:   host,
		User:   env.Username,
		Pass:   env.Password,
		Params: params.Name,
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
		return 0, fmt.Errorf("unable to detect btcd version: %w", networkErr)
	}
}
