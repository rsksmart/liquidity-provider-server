package btcclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type BtcSuiteClientAdapter struct {
	RpcClient
	config     rpcclient.ConnConfig
	httpClient utils.HttpClient
}

func NewBtcSuiteClientAdapter(config rpcclient.ConnConfig, client RpcClient) *BtcSuiteClientAdapter {
	return &BtcSuiteClientAdapter{config: config, RpcClient: client}
}

// SetClient sets the http client to be used by the adapter, only for testing purposes
func (c *BtcSuiteClientAdapter) SetClient(httpClient utils.HttpClient) {
	c.httpClient = httpClient
}

func (c *BtcSuiteClientAdapter) signRawTransactionWithKeyAsync(tx *wire.MsgTx, privateKeysWIFs []string) FutureSignRawTransactionWithKeyResult {
	cmd := &SignRawTransactionWithKeyCmd{RawTx: "", WifKeys: privateKeysWIFs}
	if tx == nil {
		return c.SendCmd(cmd)
	}

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		if err = tx.SerializeNoWitness(buf); err != nil {
			log.Errorf("Error serializing transaction to sign: %v", err)
			responseChan := make(chan *rpcclient.Response, 1)
			responseChan <- &rpcclient.Response{}
			return responseChan
		}
	}
	cmd.RawTx = hex.EncodeToString(buf.Bytes())
	return c.SendCmd(cmd)
}

func (c *BtcSuiteClientAdapter) SignRawTransactionWithKey(tx *wire.MsgTx, privateKeysWIFs []string) (*wire.MsgTx, bool, error) {
	return c.signRawTransactionWithKeyAsync(tx, privateKeysWIFs).Receive()
}

func (c *BtcSuiteClientAdapter) CreateReadonlyWallet(bodyParams ReadonlyWalletRequest) error {
	var err error
	var bodyBytes []byte
	var response btcjson.Response
	body := RpcRequestParamsObject[ReadonlyWalletRequest]{
		Jsonrpc: btcjson.RpcVersion1,
		Method:  "createwallet",
		Params:  bodyParams,
		ID:      c.NextID(),
	}

	bodyBytes, err = json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getUrl(), bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.config.User, c.config.Pass)

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	// we're alreay closing the body in utils.CloseBodyIfExists
	// nolint:bodyclose
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer utils.CloseBodyIfExists(res)

	if res == nil || res.Body == nil {
		return errors.New("received emtpy response from RPC server on createwallet")
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	} else if response.Error != nil {
		return fmt.Errorf("error creating wallet: %w", response.Error)
	}
	return nil
}

func (c *BtcSuiteClientAdapter) getUrl() string {
	if c.config.DisableTLS {
		return "http://" + c.config.Host
	} else {
		return "https://" + c.config.Host
	}
}
