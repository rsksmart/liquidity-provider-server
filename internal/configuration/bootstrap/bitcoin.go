package bootstrap

import (
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
)

const (
	unknownBtcdVersion = -1
)

func Bitcoin(env environment.BtcEnv) (*bitcoin.Connection, error) {
	var params chaincfg.Params
	log.Info("Connecting to BTC node")

	switch env.Network {
	case "mainnet":
		params = chaincfg.MainNetParams
	case "testnet":
		params = chaincfg.TestNet3Params
	case "regtest":
		params = chaincfg.RegressionNetParams
	default:
		return nil, fmt.Errorf("invalid network name: %v", env.Network)
	}

	config := rpcclient.ConnConfig{
		Host:         env.Endpoint,
		User:         env.Username,
		Pass:         env.Password,
		Params:       params.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	client, err := rpcclient.New(&config, nil)
	if err != nil {
		return nil, fmt.Errorf("RPC client error: %w", err)
	}

	version, err := checkBtcdVersion(client)
	if err != nil {
		return nil, err
	}

	if version == unknownBtcdVersion {
		log.Warn("unable to detect btcd version, but it is up and running")
	} else {
		log.Debugf("detected btcd version: %v\n", version)
	}
	conn := bitcoin.NewConnection(&params, client)
	return conn, nil
}

func checkBtcdVersion(c *rpcclient.Client) (int32, error) {
	info, err := c.GetNetworkInfo()
	switch networkErr := err.(type) {
	case nil:
		return info.Version, nil
	case *btcjson.RPCError:
		if networkErr.Code != btcjson.ErrRPCMethodNotFound.Code {
			return 0, fmt.Errorf("unable to detect btcd version: %w", networkErr)
		}
		return unknownBtcdVersion, nil
	default:
		return 0, fmt.Errorf("unable to detect btcd version: %w", networkErr)
	}
}
