package http

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcutil"

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

const quoteCleaningInterval = 1 * time.Hour
const quoteExpTimeThreshold = 5 * time.Minute

type LiquidityProviderList struct {
	Endpoint                    string
	LBCAddr                     string
	BridgeAddr                  string
	RequiredBridgeConfirmations int64
	MaxQuoteValue               uint64
}

type ConfigData struct {
	MaxQuoteValue uint64
	RSK           LiquidityProviderList
}

type Server struct {
	srv             http.Server
	providers       []providers.LiquidityProvider
	rsk             connectors.RSKConnector
	btc             connectors.BTCConnector
	db              storage.DBConnector
	now             func() time.Time
	watchers        map[string]*BTCAddressWatcher
	addWatcherMu    sync.Mutex
	sharedWatcherMu sync.Mutex
	cfgData         ConfigData
}

type QuoteRequest struct {
	CallContractAddress   string     `json:"callContractAddress"`
	CallContractArguments string     `json:"callContractArguments"`
	ValueToTransfer       *types.Wei `json:"valueToTransfer"`
	RskRefundAddress      string     `json:"rskRefundAddress"`
	BitcoinRefundAddress  string     `json:"bitcoinRefundAddress"`
}

type QuoteReturn struct {
	Quote     *types.Quote `json:"quote"`
	QuoteHash string       `json:"quoteHash"`
}

type acceptReq struct {
	QuoteHash string
}

func New(rsk connectors.RSKConnector, btc connectors.BTCConnector, db storage.DBConnector, cfgData ConfigData) Server {
	return newServer(rsk, btc, db, time.Now, cfgData)
}

func newServer(rsk connectors.RSKConnector, btc connectors.BTCConnector, db storage.DBConnector, now func() time.Time, cfgData ConfigData) Server {
	return Server{
		rsk:       rsk,
		btc:       btc,
		db:        db,
		providers: make([]providers.LiquidityProvider, 0),
		now:       now,
		watchers:  make(map[string]*BTCAddressWatcher),
		cfgData:   cfgData,
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

	return nil
}

func (s *Server) Start(port uint) error {
	r := mux.NewRouter()
	r.Path("/health").Methods(http.MethodGet).HandlerFunc(s.checkHealthHandler)
	r.Path("/getProviders").Methods(http.MethodGet).HandlerFunc(s.getProvidersHandler)
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

	s.initExpiredQuotesCleaner()

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
	quoteStatesToWatch := []types.RQState{types.RQStateWaitingForDeposit, types.RQStateCallForUserSucceeded}
	retainedQuotes, err := s.db.GetRetainedQuotes(quoteStatesToWatch)
	if err != nil {
		return err
	}

	for _, entry := range retainedQuotes {
		quote, err := s.db.GetQuote(entry.QuoteHash)
		if err != nil {
			return err
		}
		if quote == nil {
			return errors.New(fmt.Sprintf("initBtcWatchers: quote not found for hash: %s", entry.QuoteHash))
		}

		p := getProviderByAddress(s.providers, quote.LPRSKAddr)
		if p == nil {
			return errors.New(fmt.Sprintf("initBtcWatchers: provider not found for LPRSKAddr: %s", quote.LPRSKAddr))
		}

		signB, err := hex.DecodeString(entry.Signature)
		if err != nil {
			return err
		}

		err = s.addAddressWatcher(quote, entry.QuoteHash, entry.DepositAddr, signB, p, entry.State)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) addAddressWatcher(quote *types.Quote, hash string, depositAddr string, signB []byte, provider providers.LiquidityProvider, state types.RQState) error {
	s.addWatcherMu.Lock()
	defer s.addWatcherMu.Unlock()

	_, ok := s.watchers[hash]
	if ok {
		return nil
	}

	sat, _ := new(types.Wei).Add(quote.Value, quote.CallFee).ToSatoshi().Float64()
	minBtcAmount := btcutil.Amount(uint64(math.Ceil(sat)))
	expTime := getQuoteExpTime(quote)
	watcher := NewBTCAddressWatcher(hash, s.btc, s.rsk, provider, s.db, quote, signB, state, &s.sharedWatcherMu)
	err := s.btc.AddAddressWatcher(depositAddr, minBtcAmount, time.Minute, expTime, watcher, func(w connectors.AddressWatcher) {
		s.addWatcherMu.Lock()
		defer s.addWatcherMu.Unlock()
		delete(s.watchers, hash)
	})
	if err == nil {
		log.Info("added watcher for quote: : ", hash, "; deposit addr: ", depositAddr)
		s.watchers[hash] = watcher
	}
	return err
}

func (s *Server) initExpiredQuotesCleaner() {
	go func() {
		ticker := time.NewTicker(quoteCleaningInterval)
		quit := make(chan struct{})
		for {
			select {
			case <-ticker.C:
				err := s.db.DeleteExpiredQuotes(time.Now().Add(-1 * quoteExpTimeThreshold).Unix())
				if err != nil {
					log.Error("error deleting expired quites: ", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
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

func (s *Server) checkHealthHandler(w http.ResponseWriter, _ *http.Request) {
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
		log.Error("Heath Check - error encoding response: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) getProvidersHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rp, error := s.rsk.GetProviders()

	if error != nil {
		log.Error("GetProviders - error encoding response: ", error)
		http.Error(w, "internal server error "+error.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(&rp)
	if err != nil {
		log.Error("error encoding registered providers list: ", err.Error())
		http.Error(w, "internal server error "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	qr := QuoteRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&qr)

	if err != nil {
		log.Error("error decoding request: ", err.Error())
		http.Error(w, "bad request "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Debug("received quote request: ", fmt.Sprintf("%+v", qr))

	maxValueTotransfer := s.cfgData.MaxQuoteValue

	if maxValueTotransfer <= 0 {
		maxValueTotransfer = uint64(s.cfgData.RSK.MaxQuoteValue)
	}

	if qr.ValueToTransfer.Uint64() > maxValueTotransfer {
		log.Error("error on quote value, cannot be greater than: ", s.cfgData.MaxQuoteValue)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	gas, err := s.rsk.EstimateGas(qr.CallContractAddress, qr.ValueToTransfer.Copy().AsBigInt(), []byte(qr.CallContractArguments))
	if err != nil {
		log.Error("error estimating gas: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	price, err := s.rsk.GasPrice()
	if err != nil {
		log.Error("error estimating gas price: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var quotes []*QuoteReturn
	fedAddress, err := s.rsk.GetFedAddress()
	if err != nil {
		log.Error("error retrieving federation address: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	minLockTxValueInSatoshi, err := s.rsk.GetMinimumLockTxValue()
	if err != nil {
		log.Error("error retrieving minimum lock tx value: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	minLockTxValueInWei := types.SatoshiToWei(minLockTxValueInSatoshi.Uint64())

	getQuoteFailed := false
	amountBelowMinLockTxValue := false
	q := parseReqToQuote(qr, s.rsk.GetLBCAddress(), fedAddress, gas)
	for _, p := range s.providers {
		pq, err := p.GetQuote(q, gas, types.NewBigWei(price))
		if err != nil {
			log.Error("error getting quote: ", err)
			getQuoteFailed = true
			continue
		}
		if pq != nil {
			if new(types.Wei).Add(pq.Value, pq.CallFee).Cmp(minLockTxValueInWei) < 0 {
				log.Error("error getting quote; requested amount below bridge's min pegin tx value: ", qr.ValueToTransfer)
				amountBelowMinLockTxValue = true
				continue
			}

			hash, err := s.storeQuote(pq)

			if err != nil {
				log.Error(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			} else {
				quotes = append(quotes, &QuoteReturn{pq, hash})
			}
		}
	}

	if len(quotes) == 0 {
		if amountBelowMinLockTxValue {
			http.Error(w, "bad request; requested amount below bridge's min pegin tx value", http.StatusBadRequest)
			return
		}
		if getQuoteFailed {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(&quotes)
	if err != nil {
		log.Error("error encoding quote list: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) acceptQuoteHandler(w http.ResponseWriter, r *http.Request) {
	type acceptRes struct {
		Signature                 string `json:"signature"`
		BitcoinDepositAddressHash string `json:"bitcoinDepositAddressHash"`
	}
	returnQuoteSignFunc := func(w http.ResponseWriter, signature string, depositAddr string) {
		enc := json.NewEncoder(w)
		response := acceptRes{
			Signature:                 signature,
			BitcoinDepositAddressHash: depositAddr,
		}

		err := enc.Encode(response)
		if err != nil {
			log.Error("AcceptQuote - error encoding response: ", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}

	req := acceptReq{}
	w.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		log.Error("error decoding request: ", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	hashBytes, err := hex.DecodeString(req.QuoteHash)
	if err != nil {
		log.Error("error decoding quote hash: ", err.Error())
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	quote, err := s.db.GetQuote(req.QuoteHash)
	if err != nil {
		log.Error("error retrieving quote from db: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if quote == nil {
		log.Error("quote not found for hash: ", req.QuoteHash)
		http.Error(w, "quote not found", http.StatusNotFound)
		return
	}

	expTime := getQuoteExpTime(quote)
	if s.now().After(expTime) {
		log.Error("quote deposit time has elapsed; hash: ", req.QuoteHash)
		http.Error(w, "forbidden; quote deposit time has elapsed", http.StatusForbidden)
		return
	}

	rq, err := s.db.GetRetainedQuote(req.QuoteHash)
	if err != nil {
		log.Error("error fetching retained quote: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if rq != nil { // if the quote has already been accepted, just return signature and deposit addr
		returnQuoteSignFunc(w, rq.Signature, rq.DepositAddr)
		return
	}

	btcRefAddr, lpBTCAddr, lbcAddr, err := decodeAddresses(quote.BTCRefundAddr, quote.LPBTCAddr, quote.LBCAddr)
	if err != nil {
		log.Error("error decoding addresses: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	fedInfo, err := s.rsk.FetchFederationInfo()
	if err != nil {
		log.Error("error fetching fed info: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	depositAddress, err := s.btc.GetDerivedBitcoinAddress(fedInfo, btcRefAddr, lbcAddr, lpBTCAddr, hashBytes)
	if err != nil {
		log.Error("error getting derived bitcoin address: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	p := getProviderByAddress(s.providers, quote.LPRSKAddr)
	gasPrice, err := s.rsk.GasPrice()
	if err != nil {
		log.Error("error getting provider by address: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	adjustedGasLimit := types.NewUWei(uint64(CFUExtraGas) + uint64(quote.GasLimit))
	gasCost := new(types.Wei).Mul(adjustedGasLimit, types.NewBigWei(gasPrice))
	reqLiq := new(types.Wei).Add(gasCost, quote.Value)
	signB, err := p.SignQuote(hashBytes, depositAddress, reqLiq)
	if err != nil {
		log.Error("error signing quote: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.addAddressWatcher(quote, req.QuoteHash, depositAddress, signB, p, types.RQStateWaitingForDeposit)
	if err != nil {
		log.Error("error adding address watcher: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	signature := hex.EncodeToString(signB)
	returnQuoteSignFunc(w, signature, depositAddress)
}

func parseReqToQuote(qr QuoteRequest, lbcAddr string, fedAddr string, limitGas uint64) *types.Quote {
	return &types.Quote{
		LBCAddr:       lbcAddr,
		FedBTCAddr:    fedAddr,
		BTCRefundAddr: qr.BitcoinRefundAddress,
		RSKRefundAddr: qr.RskRefundAddress,
		ContractAddr:  qr.CallContractAddress,
		Data:          qr.CallContractArguments,
		Value:         qr.ValueToTransfer.Copy(),
		GasLimit:      uint32(limitGas),
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

func (s *Server) storeQuote(q *types.Quote) (string, error) {
	h, err := s.rsk.HashQuote(q)
	if err != nil {
		return "", err
	}

	err = s.db.InsertQuote(h, q)
	if err != nil {
		log.Fatalf("error inserting quote: %v", err)
	}
	return h, nil
}

func getQuoteExpTime(q *types.Quote) time.Time {
	return time.Unix(int64(q.AgreementTimestamp+q.TimeForDeposit), 0)
}
