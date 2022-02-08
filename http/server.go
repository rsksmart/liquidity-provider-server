package http

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil"
	"io"
	"math/big"
	"net/http"
	"time"

	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/storage"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
)

const (
	svcStatusOk          = "ok"
	svcStatusDegraded    = "degraded"
	svcStatusUnreachable = "unreachable"
)

type Server struct {
	srv       http.Server
	providers []providers.LiquidityProvider
	rsk       connectors.RSKConnector
	btc       connectors.BTCConnector
	db        storage.DBConnector
}

type QuoteRequest struct {
	CallContractAddress   string `json:"callContractAddress"`
	CallContractArguments string `json:"callContractArguments"`
	ValueToTransfer       uint64 `json:"valueToTransfer"`
	GasLimit              uint32 `json:"gasLimit"`
	RskRefundAddress      string `json:"rskRefundAddress"`
	BitcoinRefundAddress  string `json:"bitcoinRefundAddress"`
}

type acceptReq struct {
	QuoteHash string
}

func New(rsk connectors.RSKConnector, btc connectors.BTCConnector, db storage.DBConnector) Server {
	return Server{
		rsk:       rsk,
		btc:       btc,
		db:        db,
		providers: make([]providers.LiquidityProvider, 0),
	}
}

func (s *Server) AddProvider(lp providers.LiquidityProvider) error {
	s.providers = append(s.providers, lp)
	addrStr := lp.Address()
	c, m, err := s.rsk.GetCollateral(addrStr)
	if err != nil {
		return err
	}
	addr := common.HexToAddress(addrStr)
	cmp := c.Cmp(big.NewInt(0))
	if cmp == 0 { // provider not registered
		opts := &bind.TransactOpts{
			Value:  m,
			From:   addr,
			Signer: lp.SignTx,
		}
		err := s.rsk.RegisterProvider(opts)
		if err != nil {
			return err
		}
	} else if cmp < 0 { // not enough collateral
		opts := &bind.TransactOpts{
			Value:  m.Sub(m, c),
			From:   addr,
			Signer: lp.SignTx,
		}
		err := s.rsk.AddCollateral(opts)
		if err != nil {
			return err
		}
	}
	liq, err := s.rsk.GetAvailableLiquidity(addrStr)
	if err != nil {
		return err
	}
	lp.SetLiquidity(liq)

	return nil
}

func (s *Server) Start(port uint) error {
	r := mux.NewRouter()
	r.Path("/health").Methods(http.MethodGet).HandlerFunc(s.checkHealthHandler)
	r.Path("/getQuote").Methods(http.MethodPost).HandlerFunc(s.getQuoteHandler)
	r.Path("/acceptQuote").Methods(http.MethodPost).HandlerFunc(s.acceptQuoteHandler)
	w := log.StandardLogger().WriterLevel(log.DebugLevel)
	h := handlers.LoggingHandler(w, r)
	defer func(w *io.PipeWriter) {
		_ = w.Close()
	}(w)

	err := s.initBtcWatchers()
	if err != nil {
		return err
	}

	s.srv = http.Server{
		Addr:    ":" + fmt.Sprint(port),
		Handler: h,
	}
	log.Info("server started at localhost:", s.srv.Addr)

	err = s.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) initBtcWatchers() error {
	retainedQuotes, err := s.db.GetRetainedQuotes()
	if err != nil {
		return err
	}

	for _, entry := range retainedQuotes {
		quote, err := s.db.GetQuote(entry.QuoteHash)
		if err != nil {
			return err
		}

		p := getProviderByAddress(s.providers, quote.LPRSKAddr)
		if p == nil {
			return errors.New(fmt.Sprintf("provider not found for LPRSKAddr: %s", quote.LPRSKAddr))
		}

		signB, err := hex.DecodeString(entry.Signature)
		if err != nil {
			return err
		}

		err = s.addAddressWatcher(quote, entry.QuoteHash, entry.DepositAddr, signB, p, entry.CalledForUser)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) addAddressWatcher(quote *types.Quote, hash string, depositAddr string, signB []byte, provider providers.LiquidityProvider, calledForUser bool) error {
	minAmount := btcutil.Amount(quote.Value + quote.CallFee)
	expTime := time.Unix(int64(quote.AgreementTimestamp + quote.TimeForDeposit), 0)
	watcher := NewBTCAddressWatcher(hash, s.btc, s.rsk, provider, s.db, quote, signB, calledForUser)
	err := s.btc.AddAddressWatcher(depositAddr, minAmount, time.Minute, expTime, watcher)
	return err
}

func (s *Server) Shutdown() {
	log.Info("stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown failed: ", err)
	}
	log.Info("server stopped")
}

func (s *Server) checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	type services struct {
		Db  string `json:"db"`
		Rsk string `json:"rsk"`
		Btc string `json:"btc"`
	}
	type healthRes struct {
		Status   string   `json:"status"`
		Services services `json:"services"`
	}

	lpsSvcStatus := svcStatusOk
	dbSvcStatus := svcStatusOk
	rskSvcStatus := svcStatusOk
	btcSvcStatus := svcStatusOk

	if err := s.db.CheckConnection(); err != nil {
		log.Error("error checking db connection status: ", err.Error())
		dbSvcStatus = svcStatusUnreachable
		lpsSvcStatus = svcStatusDegraded
	}

	if err := s.rsk.CheckConnection(); err != nil {
		log.Error("error checking rsk connection status: ", err.Error())
		rskSvcStatus = svcStatusUnreachable
		lpsSvcStatus = svcStatusDegraded
	}

	if err := s.btc.CheckConnection(); err != nil {
		log.Error("error checking btcd connection status: ", err.Error())
		btcSvcStatus = svcStatusUnreachable
		lpsSvcStatus = svcStatusDegraded
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	response := healthRes{
		Status: lpsSvcStatus,
		Services: services{
			Db:  dbSvcStatus,
			Rsk: rskSvcStatus,
			Btc: btcSvcStatus,
		},
	}
	err := enc.Encode(response)
	if err != nil {
		log.Error("error encoding response: ", err.Error())
		http.Error(w, "error processing request", http.StatusInternalServerError)
	}
}

func (s *Server) getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	qr := QuoteRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Debug("received quote request: ", fmt.Sprintf("%+v", qr))

	gas, err := s.rsk.EstimateGas(qr.CallContractAddress, qr.ValueToTransfer, []byte(qr.CallContractArguments))
	if err != nil {
		log.Error("error estimating gas: ", err.Error())
		http.Error(w, "error estimating gas", http.StatusInternalServerError)
		return
	}

	price, err := s.rsk.GasPrice()
	if err != nil {
		log.Error("error estimating gas price: ", err.Error())
		http.Error(w, "error estimating gas price", http.StatusInternalServerError)
		return
	}

	var quotes []*types.Quote
	fedAddress, err := s.rsk.GetFedAddress()
	if err != nil {
		log.Error("error retrieving federation address: ", err.Error())
		http.Error(w, "error retrieving federation address", http.StatusInternalServerError)
		return
	}

	q := parseReqToQuote(qr, s.rsk.GetLBCAddress(), fedAddress)
	for _, p := range s.providers {
		pq := p.GetQuote(q, gas, price)
		if pq != nil {
			err = s.storeQuote(pq)

			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				quotes = append(quotes, pq)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(&quotes)
	if err != nil {
		log.Error("error encoding quote list: ", err.Error())
		http.Error(w, "error processing quotes", http.StatusInternalServerError)
		return
	}
}

func (s *Server) acceptQuoteHandler(w http.ResponseWriter, r *http.Request) {
	type acceptRes struct {
		Signature                 string `json:"signature"`
		BitcoinDepositAddressHash string `json:"bitcoinDepositAddressHash"`
	}
	req := acceptReq{}
	w.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashBytes, err := hex.DecodeString(req.QuoteHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quote, err := s.db.GetQuote(req.QuoteHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	btcRefAddr, lpBTCAddr, lbcAddr, err := decodeAddresses(quote.BTCRefundAddr, quote.LPBTCAddr, quote.LBCAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	depositAddress, err := s.btc.GetDerivedBitcoinAddress(btcRefAddr, lbcAddr, lpBTCAddr, hashBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := getProviderByAddress(s.providers, quote.LPRSKAddr)
	gasPrice, err := s.rsk.GasPrice()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqLiq := (uint64(CFUExtraGas)+uint64(quote.GasLimit))*gasPrice + quote.Value
	signB, err := p.SignQuote(hashBytes, depositAddress, big.NewInt(int64(reqLiq)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.addAddressWatcher(quote, req.QuoteHash, depositAddress, signB, p, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	signature := hex.EncodeToString(signB)
	enc := json.NewEncoder(w)
	response := acceptRes{
		Signature:                 signature,
		BitcoinDepositAddressHash: depositAddress,
	}
	err = enc.Encode(response)

	// TODO: ensure that the quote is not processed if there is any kind of error in the communication with the client
	if err != nil {
		log.Error("error encoding response: ", err.Error())
		http.Error(w, "error processing request", http.StatusInternalServerError)
		return
	}
}

func parseReqToQuote(qr QuoteRequest, lbcAddr string, fedAddr string) types.Quote {
	return types.Quote{
		LBCAddr:       lbcAddr,
		FedBTCAddr:    fedAddr,
		BTCRefundAddr: qr.BitcoinRefundAddress,
		RSKRefundAddr: qr.RskRefundAddress,
		ContractAddr:  qr.CallContractAddress,
		Data:          qr.CallContractArguments,
		Value:         qr.ValueToTransfer,
		GasLimit:      qr.GasLimit,
	}
}

func decodeAddresses(btcRefundAddr string, lpBTCAddr string, lbcAddr string) ([]byte, []byte, []byte, error) {
	btcRefAddrB, err := connectors.DecodeBTCAddressWithVersion(btcRefundAddr)
	if err != nil {
		return nil, nil, nil, err
	}
	lpBTCAddrB, err := connectors.DecodeBTCAddressWithVersion(lpBTCAddr)
	if err != nil {
		return nil, nil, nil, err
	}
	lbcAddrB, err := connectors.DecodeRSKAddress(lbcAddr)
	if err != nil {
		return nil, nil, nil, err
	}
	return btcRefAddrB, lpBTCAddrB, lbcAddrB, nil
}

func getProviderByAddress(liquidityProviders []providers.LiquidityProvider, addr string) (ret providers.LiquidityProvider) {
	for _, p := range liquidityProviders {
		if p.Address() == addr {
			return p
		}
	}
	return nil
}

func (s *Server) storeQuote(q *types.Quote) error {
	h, err := s.rsk.HashQuote(q)
	if err != nil {
		return err
	}

	err = s.db.InsertQuote(h, q)
	if err != nil {
		log.Fatalf("error inserting quote: %v", err)
	}
	return nil
}
