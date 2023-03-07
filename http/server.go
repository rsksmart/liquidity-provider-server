package http

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	//"github.com/rsksmart/liquidity-provider-server/response"

	//"github.com/rsksmart/liquidity-provider-server/response"
	"io"
	"math"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	mongoDB "github.com/rsksmart/liquidity-provider-server/mongo"
	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider-server/pegout"
	"github.com/rsksmart/liquidity-provider-server/storage"

	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/connectors"

	// "github.com/rsksmart/liquidity-provider/providers"
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

const BadRequestError = "bad request"
const UnableToBuildResponse = "Unable to build response"
const UnableToDeserializePayloadError = "Unable to deserialize payload: %v"
const ErrorRetrievingFederationAddress = "error retrieving federation address: "
const ErrorRetrievingMinimumLockValue = "error retrieving minimum lock tx value: "
const ErrorRequestedAmountBelowBridgeMin = "requested amount below bridge's min pegin tx value"
const ErrorGetQuoteFailed = "error getting specified quote"
const ErrorEncodingQuotesList = "error encoding quote list for response"
const ErrorBadBodyRequest = "Body of the request is wrong: "
const ErrorEstimatingGas = "Error on RSK Network, couldnt estimate gas"
const ErrorValueTooHigh = "value to transfer too high"
const ErrorStoringProviderQuote = "Error storing the quote on server"
const ErrorFetchingMongoDBProviders = "Error Fetching Providers from MongoDB: "
const ErrorCreatingLocalProvider= "Error Creating New Local Provider"
const ErrorAddingProvider= "Error Adding New provider"
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
	providers       []pegin.LiquidityProvider
	pegoutProviders []pegout.LiquidityProvider
	rsk             connectors.RSKConnector
	btc             connectors.BTCConnector
	dbMongo         *mongoDB.DB
	now             func() time.Time
	watchers        map[string]*BTCAddressWatcher
	pegOutWatchers  map[string]*BTCAddressPegOutWatcher
	rskWatchers     map[string]*RegisterPegoutWatcher
	addWatcherMu    sync.Mutex
	sharedWatcherMu sync.Mutex
	cfgData         ConfigData
	ProviderRespository *storage.LPRepository
	ProviderConfig pegin.ProviderConfig
}

type QuoteRequest struct {
	CallContractAddress   string     `json:"callContractAddress"`
	CallContractArguments string     `json:"callContractArguments"`
	ValueToTransfer       *types.Wei `json:"valueToTransfer"`
	RskRefundAddress      string     `json:"rskRefundAddress"`
	LpAddress             string     `json:"lpAddress"`
	BitcoinRefundAddress  string     `json:"bitcoinRefundAddress"`
}

type QuoteReturn struct {
	Quote     *pegin.Quote `json:"quote"`
	QuoteHash string       `json:"quoteHash"`
}

type QuotePegOutRequest struct {
	From                 string `json:"from"`
	ValueToTransfer      uint64 `json:"valueToTransfer"`
	RskRefundAddress     string `json:"rskRefundAddress"`
	BitcoinRefundAddress string `json:"bitcoinRefundAddress"`
}

type QuotePegOutResponse struct {
	Quote             *pegout.Quote `json:"quote"`
	DerivationAddress string        `json:"derivationAddress"`
	QuoteHash         string        `json:"quoteHash"`
}

type acceptReq struct {
	QuoteHash string
}

func enableCors(res *http.ResponseWriter) {
	headers := (*res).Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

type acceptRes struct {
	Signature                 string `json:"signature"`
	BitcoinDepositAddressHash string `json:"bitcoinDepositAddressHash"`
}

type AcceptResPegOut struct {
	Signature string `json:"signature"`
}

type acceptReqPegout struct {
	QuoteHash         string `json:"quoteHash"`
	DerivationAddress string `json:"derivationAddress"`
}

type pegOutQuoteReq struct {
	Quote *pegout.Quote `json:"quote"`
}

type pegOutQuoteResponse struct {
	QuoteHash string `json:"quoteHash"`
}

func New(rsk connectors.RSKConnector, btc connectors.BTCConnector, dbMongo *mongoDB.DB, cfgData ConfigData,LPRep *storage.LPRepository,ProviderConfig pegin.ProviderConfig) Server {
	return newServer(rsk, btc, dbMongo, time.Now, cfgData,LPRep,ProviderConfig)
}

func newServer(rsk connectors.RSKConnector, btc connectors.BTCConnector, dbMongo *mongoDB.DB, now func() time.Time, cfgData ConfigData,LPRep *storage.LPRepository, ProviderConfig pegin.ProviderConfig) Server {
	return Server{
		rsk:             rsk,
		btc:             btc,
		dbMongo:         dbMongo,
		providers:       make([]pegin.LiquidityProvider, 0),
		pegoutProviders: make([]pegout.LiquidityProvider, 0),
		now:             now,
		watchers:        make(map[string]*BTCAddressWatcher),
		pegOutWatchers:  make(map[string]*BTCAddressPegOutWatcher),
		rskWatchers:     make(map[string]*RegisterPegoutWatcher),
		cfgData:         cfgData,
		ProviderRespository: LPRep,
		ProviderConfig: ProviderConfig,
	}
}

func (s *Server) AddProvider(lp pegin.LiquidityProvider,ProviderDetails types.ProviderRegisterRequest) error {
	s.providers = append(s.providers, lp)
	addrStr := lp.Address()
	c, m, err := s.rsk.GetCollateral(addrStr)
	if err != nil {
		return err
	}
	addr := common.HexToAddress(addrStr)
	cmp := c.Cmp(big.NewInt(0))
	if cmp >= 0 { 
		opts := &bind.TransactOpts{
			Value:  m,
			From:   addr,
			Signer: lp.SignTx,
		}
		providerID,err := s.rsk.RegisterProvider(opts,ProviderDetails.Name,big.NewInt(int64(ProviderDetails.Fee)),big.NewInt(int64(ProviderDetails.QuoteExpiration)),big.NewInt(int64(ProviderDetails.AcceptedQuoteExpiration)),big.NewInt(int64(ProviderDetails.MinTransactionValue)),big.NewInt(int64(ProviderDetails.MaxTransactionValue)),ProviderDetails.ApiBaseUrl,ProviderDetails.Status)
		if err != nil {
			return err
		}
		err2 := s.dbMongo.InsertProvider(providerID)
		if(err2 != nil){
			return err2
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

func (s *Server) AddPegOutProvider(lp pegout.LiquidityProvider,ProviderDetails types.ProviderRegisterRequest) error {
	s.pegoutProviders = append(s.pegoutProviders, lp)
	addrStr := lp.Address()
	c, m, err := s.rsk.GetCollateral(addrStr)
	if err != nil {
		return err
	}
	addr := common.HexToAddress(addrStr)
	cmp := c.Cmp(big.NewInt(0))
	if cmp >= 0 { 
		opts := &bind.TransactOpts{
			Value:  m,
			From:   addr,
			Signer: lp.SignTx,
		}
		providerID,err := s.rsk.RegisterProvider(opts,ProviderDetails.Name,big.NewInt(int64(ProviderDetails.Fee)),big.NewInt(int64(ProviderDetails.QuoteExpiration)),big.NewInt(int64(ProviderDetails.AcceptedQuoteExpiration)),big.NewInt(int64(ProviderDetails.MinTransactionValue)),big.NewInt(int64(ProviderDetails.MaxTransactionValue)),ProviderDetails.ApiBaseUrl,ProviderDetails.Status)
		if err != nil {
			return err
		}
		err2 := s.dbMongo.InsertProvider(providerID)
		if(err2 != nil){
			return err2
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

func (s *Server) registerProviderHandler(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)
	enableCors(&w)
	payload := types.ProviderRegisterRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Errorf(UnableToDeserializePayloadError, err)
		http.Error(w, UnableToDeserializePayloadError, http.StatusBadRequest)
		return
	}
	lp, err := pegin.NewLocalProvider(s.ProviderConfig, s.ProviderRespository)
	if err != nil {
		log.Errorf(ErrorCreatingLocalProvider, err)
		http.Error(w, ErrorCreatingLocalProvider, http.StatusBadRequest)
		return
	}
	err = s.AddProvider(lp,payload)
	if err != nil {
		log.Errorf(ErrorAddingProvider, err)
		http.Error(w, ErrorAddingProvider, http.StatusBadRequest)
		return
	}
	response := "Provider Created Successfully";
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&response)
	if err != nil {
		http.Error(w, UnableToBuildResponse, http.StatusInternalServerError)
		return
	}
}
func (s *Server) Start(port uint) error {
	r := mux.NewRouter()
	r.Path("/health").Methods(http.MethodGet).HandlerFunc(s.checkHealthHandler)
	r.Path("/getProviders").Methods(http.MethodGet).HandlerFunc(s.getProvidersHandler)
	r.Path("/pegin/getQuote").Methods(http.MethodPost).HandlerFunc(s.getQuoteHandler)
	r.Path("/pegin/acceptQuote").Methods(http.MethodPost).HandlerFunc(s.acceptQuoteHandler)
	r.Path("/pegout/getQuotes").Methods(http.MethodPost).HandlerFunc(s.getQuotesPegOutHandler)
	r.Path("/pegout/acceptQuote").Methods(http.MethodPost).HandlerFunc(s.acceptQuotePegOutHandler)
	r.Path("/pegout/refundPegOut").Methods(http.MethodPost).HandlerFunc(s.refundPegOutHandler)
	r.Path("/pegout/sendBTC").Methods(http.MethodPost).HandlerFunc(s.sendBTC)
	r.Path("/addCollateral").Methods(http.MethodPost).HandlerFunc(s.addCollateral)
	r.Path("/provider/register").Methods(http.MethodPost).HandlerFunc(s.registerProviderHandler)
	r.Methods("OPTIONS").HandlerFunc(s.handleOptions)
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

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) initBtcWatchers() error {
	quoteStatesToWatch := []types.RQState{types.RQStateWaitingForDeposit, types.RQStateCallForUserSucceeded}
	retainedQuotes, err := s.dbMongo.GetRetainedQuotes(quoteStatesToWatch)
	if err != nil {
		return err
	}

	for _, entry := range retainedQuotes {
		quote, err := s.dbMongo.GetQuote(entry.QuoteHash)
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

func (s *Server) addAddressWatcher(quote *pegin.Quote, hash string, depositAddr string, signB []byte, provider pegin.LiquidityProvider, state types.RQState) error {
	s.addWatcherMu.Lock()
	defer s.addWatcherMu.Unlock()

	_, ok := s.watchers[hash]
	if ok {
		return nil
	}

	sat, _ := new(types.Wei).Add(quote.Value, quote.CallFee).ToSatoshi().Float64()
	minBtcAmount := btcutil.Amount(uint64(math.Ceil(sat)))
	expTime := getQuoteExpTime(quote)
	watcher := NewBTCAddressWatcher(hash, s.btc, s.rsk, provider, *s.dbMongo, quote, signB, state, &s.sharedWatcherMu)
	err := s.btc.AddAddressWatcher(depositAddr, minBtcAmount, time.Minute, expTime, watcher, func(w connectors.AddressWatcher) {
		s.addWatcherMu.Lock()
		defer s.addWatcherMu.Unlock()
		delete(s.watchers, hash)
	})
	if err == nil {
		escapedDepositAddr := strings.Replace(depositAddr, "\n", "", -1)
		escapedDepositAddr = strings.Replace(escapedDepositAddr, "\r", "", -1)
		s.watchers[hash] = watcher
	}
	return err
}

func (s *Server) addAddressPegOutWatcher(quote *pegout.Quote, hash string, depositAddr string, signB []byte, provider pegout.LiquidityProvider, state types.RQState) error {
	_, ok := s.pegOutWatchers[hash]

	if ok {
		return nil
	}

	minBtcAmount := btcutil.Amount(quote.Value)
	expTime := getPegOutQuoteExpTime(quote)
	watcher := &BTCAddressPegOutWatcher{
		hash:         hash,
		btc:          s.btc,
		rsk:          s.rsk,
		lp:           provider,
		quote:        quote,
		state:        state,
		signature:    signB,
		done:         make(chan struct{}),
		sharedLocker: &s.sharedWatcherMu,
	}
	err := s.btc.AddAddressPegOutWatcher(depositAddr, minBtcAmount, time.Minute, expTime, watcher, func(w connectors.AddressWatcher) {
		log.Debugln("Done: addAddressPegOutWatcher")
		s.addWatcherMu.Lock()
		defer s.addWatcherMu.Unlock()
		delete(s.pegOutWatchers, hash)
	})
	if err == nil {
		escapedDepositAddr := strings.Replace(depositAddr, "\n", "", -1)
		escapedDepositAddr = strings.Replace(escapedDepositAddr, "\r", "", -1)
		s.pegOutWatchers[hash] = watcher
	}
	return err
}

func (s *Server) addAddressWatcherToVerifyRegisterPegOut(quote *pegout.Quote, hash string, derivationAddress string, signB []byte, provider pegout.LiquidityProvider, state types.RQState) error {
	s.addWatcherMu.Lock()
	defer s.addWatcherMu.Unlock()

	_, ok := s.watchers[hash]
	if ok {
		return nil
	}

	expTime := getPegOutQuoteExpTime(quote)
	watcher := &RegisterPegoutWatcher{
		hash:              hash,
		btc:               s.btc,
		rsk:               s.rsk,
		lp:                provider,
		quote:             quote,
		state:             state,
		signature:         signB,
		done:              make(chan struct{}),
		sharedLocker:      &s.sharedWatcherMu,
		derivationAddress: derivationAddress,
	}
	err := s.rsk.AddQuoteToWatch(hash, time.Minute, expTime, watcher, func(w connectors.QuotePegOutWatcher) {
		s.addWatcherMu.Lock()
		defer s.addWatcherMu.Unlock()
		delete(s.rskWatchers, hash)
	})
	if err == nil {
		s.rskWatchers[hash] = watcher
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
				err := s.dbMongo.DeleteExpiredQuotes(time.Now().Add(-1 * quoteExpTimeThreshold).Unix())
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
	enableCors(&w)
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

	if err := s.dbMongo.CheckConnection(); err != nil {
		log.Error("error checking mongo DB connection status: ", err.Error())
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

	toRestAPI(w)
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

func toRestAPI(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func (a *QuoteRequest) validateQuoteRequest() string {
	err := ""

	if len(a.RskRefundAddress) == 0 {
		err += "RskRefundAddress is empty; "
	}
	if len(a.BitcoinRefundAddress) == 0 {
		err += "BitcoinRefundAddress is empty; "
	}
	if len(a.CallContractAddress) == 0 {
		err += "CallContractAddress is empty; "
	}

	return err
}

func (a *QuotePegOutRequest) validateQuoteRequest() string {
	err := ""

	if a.ValueToTransfer == 0 {
		err += "Value to Transfer cannot be empty or zero!"
	}

	return err
}

func (s *Server) getProvidersHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	
	providerList, error := s.dbMongo.GetProviders()
	if(error != nil){
		log.Error("Error fetching providers. Error: ", error)
		customError := NewServerError(ErrorFetchingMongoDBProviders + error.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}
	rp, error := s.rsk.GetProviders(providerList)

	if error != nil {
		log.Error("GetProviders - error encoding response: ", error)
		customError := NewServerError("GetProviders - error encoding response: "+error.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(&rp)
	if err != nil {
		log.Error("error encoding registered providers list: ", err.Error())
		customError := NewServerError("error encoding registered providers list: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}
}

func (s *Server) getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	qr := QuoteRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&qr)

	if err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}
	log.Debug("received quote request: ", fmt.Sprintf("%+v", qr))

	maxValueTotransfer := s.cfgData.MaxQuoteValue

	if maxValueTotransfer <= 0 {
		maxValueTotransfer = uint64(s.cfgData.RSK.MaxQuoteValue)
	}

	if qr.LpAddress == "" || !common.IsHexAddress(qr.LpAddress) {
		log.Error("Liquidity Provider Address lpAddress not sent")
		customError := NewServerError("Liquidity Provider Address lpAddress not sent", createEmptyInterfaceMap(), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	if qr.ValueToTransfer.Uint64() > maxValueTotransfer {
		log.Error(ErrorRetrievingFederationAddress, err.Error())
		customError := NewServerError(ErrorValueTooHigh, make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	if errval := qr.validateQuoteRequest(); len(errval) > 0 {
		customError := NewServerError(ErrorBadBodyRequest+errval, make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	gas, err := s.rsk.EstimateGas(qr.CallContractAddress, qr.ValueToTransfer.Copy().AsBigInt(), []byte(qr.CallContractArguments))
	if err != nil {
		log.Error(ErrorEstimatingGas, err.Error())
		customError := NewServerError(ErrorEstimatingGas, make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}

	price, err := s.rsk.GasPrice()
	if err != nil {
		log.Error(ErrorEstimatingGas+" price", err.Error())
		customError := NewServerError(ErrorEstimatingGas+" price", make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}

	var quotes []*QuoteReturn
	fedAddress, err := s.rsk.GetFedAddress()
	if err != nil {
		log.Error(ErrorRetrievingFederationAddress, err.Error())
		customError := NewServerError(ErrorRetrievingFederationAddress, make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}

	minLockTxValueInSatoshi, err := s.rsk.GetMinimumLockTxValue()
	if err != nil {
		log.Error(ErrorRetrievingMinimumLockValue, err.Error())
		customError := NewServerError(ErrorRetrievingMinimumLockValue, make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusInternalServerError)
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
				customError := NewServerError(ErrorStoringProviderQuote, make(map[string]interface{}), false)
				ResponseError(w, customError, http.StatusInternalServerError)
				return
			} else {
				quotes = append(quotes, &QuoteReturn{pq, hash})
			}
		}
	}

	if len(quotes) == 0 {
		if amountBelowMinLockTxValue {
			details := map[string]interface{}{
				"value":               q.Value,
				"callFee":             q.CallFee,
				"minLockTxValueInWei": minLockTxValueInWei,
			}

			customError := NewServerError(ErrorRequestedAmountBelowBridgeMin, details, true)
			ResponseError(w, customError, http.StatusBadRequest)
			return
		}
		if getQuoteFailed {
			details := map[string]interface{}{
				"quote": q,
				"gas":   gas,
			}
			customError := NewServerError(ErrorGetQuoteFailed, details, true)
			ResponseError(w, customError, http.StatusNotFound) // StatusBadRequest or StatusInternalServerError?
			return
		}
	}

	toRestAPI(w)
	enc := json.NewEncoder(w)
	err = enc.Encode(&quotes)
	if err != nil {
		log.Error("error encoding quote list: ", err.Error())
		details := map[string]interface{}{
			"quotes": quotes,
			"check":  true,
		}

		customError := NewServerError(ErrorEncodingQuotesList, details, true)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}
}

func (s *Server) acceptQuoteHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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
			customError := NewServerError("AcceptQuote - error encoding response: "+err.Error(), make(map[string]interface{}), true)
			ResponseError(w, customError, http.StatusBadRequest)
			return
		}
	}

	req := acceptReq{}
	toRestAPI(w)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}

	hashBytes, err := hex.DecodeString(req.QuoteHash)
	if err != nil {
		log.Error("error decoding quote hash: ", err.Error())
		customError := NewServerError("error decoding quote hash: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	quote, err := s.dbMongo.GetQuote(req.QuoteHash)
	if err != nil {
		log.Error("error retrieving quote from db: ", err.Error())
		customError := NewServerError("error retrieving quote from db: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}
	if quote == nil {
		log.Error("quote not found for hash: ", req.QuoteHash)
		customError := NewServerError("quote not found for hash: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	expTime := getQuoteExpTime(quote)
	if s.now().After(expTime) {
		log.Error("quote deposit time has elapsed; hash: ", req.QuoteHash)
		customError := NewServerError("quote deposit time has elapsed; hash: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	rq, err := s.dbMongo.GetRetainedQuote(req.QuoteHash)
	if err != nil {
		log.Error("error fetching retained quote: ", err.Error())
		customError := NewServerError("error fetching retained quote: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}
	if rq != nil { // if the quote has already been accepted, just return signature and deposit addr
		returnQuoteSignFunc(w, rq.Signature, rq.DepositAddr)
		return
	}

	btcRefAddr, lpBTCAddr, lbcAddr, err := decodeAddresses(quote.BTCRefundAddr, quote.LPBTCAddr, quote.LBCAddr)
	if err != nil {
		log.Error("error decoding addresses: ", err.Error())
		customError := NewServerError("error decoding addresses: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	fedInfo, err := s.rsk.FetchFederationInfo()
	if err != nil {
		log.Error("error fetching fed info: ", err.Error())
		customError := NewServerError("error fetching fed info: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	depositAddress, err := s.rsk.GetDerivedBitcoinAddress(fedInfo, s.btc.GetParams(), btcRefAddr, lbcAddr, lpBTCAddr, hashBytes)
	if err != nil {
		log.Error("error getting derived bitcoin address: ", err.Error())
		customError := NewServerError("error getting derived bitcoin address: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	p := getProviderByAddress(s.providers, quote.LPRSKAddr)
	gasPrice, err := s.rsk.GasPrice()
	if err != nil {
		log.Error("error getting provider by address: ", err.Error())
		customError := NewServerError("error getting provider by address: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	adjustedGasLimit := types.NewUWei(uint64(CFUExtraGas) + uint64(quote.GasLimit))
	gasCost := new(types.Wei).Mul(adjustedGasLimit, types.NewBigWei(gasPrice))
	reqLiq := new(types.Wei).Add(gasCost, quote.Value)
	signB, err := p.SignQuote(hashBytes, depositAddress, reqLiq)
	if err != nil {
		log.Error("error signing quote: ", err.Error())
		customError := NewServerError("error signing quote: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	err = s.addAddressWatcher(quote, req.QuoteHash, depositAddress, signB, p, types.RQStateWaitingForDeposit)
	if err != nil {
		log.Error("error adding address watcher: ", err.Error())
		customError := NewServerError("error adding address watcher: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	signature := hex.EncodeToString(signB)
	returnQuoteSignFunc(w, signature, depositAddress)
}

func parseReqToQuote(qr QuoteRequest, lbcAddr string, fedAddr string, limitGas uint64) *pegin.Quote {
	return &pegin.Quote{
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

func parseQuotePegOutRequestToQuote(qr QuotePegOutRequest) *pegout.Quote {
	return &pegout.Quote{
		RSKRefundAddr: qr.RskRefundAddress,
		Value:         qr.ValueToTransfer,
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

func getProviderByAddress(liquidityProviders []pegin.LiquidityProvider, addr string) (ret pegin.LiquidityProvider) {
	for _, p := range liquidityProviders {
		if p.Address() == addr {
			return p
		}
	}
	return nil
}

func getPegOutProviderByAddress(liquidityProviders []pegout.LiquidityProvider, addr string) (ret pegout.LiquidityProvider) {
	for _, p := range liquidityProviders {
		if p.Address() == addr {
			return p
		}
	}
	return nil
}

func (s *Server) storeQuote(q *pegin.Quote) (string, error) {
	h, err := s.rsk.HashQuote(q)
	if err != nil {
		return "", err
	}

	err = s.dbMongo.InsertQuote(h, q)
	if err != nil {
		log.Fatalf("error inserting quote: %v", err)
	}

	return h, nil
}

func (s *Server) storePegoutQuote(q *pegout.Quote, derivationAddress string) (string, error) {
	h, err := s.rsk.HashPegOutQuote(q)
	if err != nil {
		return "", err
	}

	err = s.dbMongo.InsertPegOutQuote(h, q, derivationAddress)
	if err != nil {
		log.Fatalf("error inserting quote: %v", err)
		return "", err
	}
	return h, nil
}

func getQuoteExpTime(q *pegin.Quote) time.Time {
	return time.Unix(int64(q.AgreementTimestamp+q.TimeForDeposit), 0)
}

func getPegOutQuoteExpTime(q *pegout.Quote) time.Time {
	return time.Unix(int64(q.AgreementTimestamp+q.DepositDateLimit), 0)
}

func (s *Server) getQuotesPegOutHandler(w http.ResponseWriter, r *http.Request) {
	qr := QuotePegOutRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&qr)
	if err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}
	log.Debug("received quote request: ", fmt.Sprintf("%+v", qr))

	q := parseQuotePegOutRequestToQuote(qr)
	quotes := make([]QuotePegOutResponse, 0)

	rskBlockNumber, err := s.rsk.GetRskHeight()

	if errval := qr.validateQuoteRequest(); len(errval) > 0 {
		log.Error("[pegout] [getquote] - error validating body params: ", errval)
		toRestAPI(w)
		http.Error(w, "bad request body: "+errval, http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Error(ErrorRetrievingFederationAddress, err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	quotes, ok := s.generateQuotesByProviders(q, rskBlockNumber, qr, quotes)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	buildResponseGetQuotePegOut(w, quotes)
}

func (s *Server) generateQuotesByProviders(q *pegout.Quote, rskBlockNumber uint64, qr QuotePegOutRequest, quotes []QuotePegOutResponse) ([]QuotePegOutResponse, bool) {
	for _, p := range s.pegoutProviders {

		pq, err := p.GetQuote(q, rskBlockNumber)

		if err != nil {
			log.Error("error getting quote: ", err)
			return nil, false
		}

		if pq != nil {

			pq.LBCAddr = s.rsk.GetLBCAddress()

			h, err := s.rsk.HashPegOutQuote(pq)

			if err != nil {
				log.Error("error getting quote: unable to hash quote", err)
				return nil, false
			}

			derivationAddress, ok := s.buildDerivationAddress(qr, h)

			if !ok {
				return nil, false
			}

			quoteHash, err := s.storePegoutQuote(pq, derivationAddress)

			if err != nil {
				log.Error(err)
				return nil, false
			}

			quote := &QuotePegOutResponse{
				Quote:             pq,
				DerivationAddress: derivationAddress,
				QuoteHash:         quoteHash,
			}
			quotes = append(quotes, *quote)

		}
	}
	return quotes, true
}

func (s *Server) buildDerivationAddress(qr QuotePegOutRequest, h string) (string, bool) {
	pubKey, err := hex.DecodeString(qr.From)

	if err != nil {
		log.Error("Unable to decode bitocin user public key")
		log.Error(err)
		return "", false
	}

	decodedQuoteHash, err := hex.DecodeString(h)

	if err != nil {
		log.Error("Unable to decode quote hash")
		log.Error(err)
		return "", false
	}

	derivationAddress, err := s.btc.ComputeDerivationAddresss(pubKey, decodedQuoteHash)
	return derivationAddress, true
}

func buildResponseGetQuotePegOut(w http.ResponseWriter, quotes []QuotePegOutResponse) {
	toRestAPI(w)
	enc := json.NewEncoder(w)
	err := enc.Encode(&quotes)
	if err != nil {
		log.Error("error encoding quote list: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func buildErrorDecodingRequest(w http.ResponseWriter, err error) {
	log.Error("error decoding request: ", err.Error())
	http.Error(w, BadRequestError, http.StatusBadRequest)
	return
}

func returnQuoteSignFunc(w http.ResponseWriter, signature string, depositAddr string) {
	enc := json.NewEncoder(w)
	response := acceptRes{
		Signature:                 signature,
		BitcoinDepositAddressHash: depositAddr,
	}

	err := enc.Encode(response)
	if err != nil {
		log.Error("error encoding response: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func returnQuotePegOutSignFunc(w http.ResponseWriter, signature string) {
	enc := json.NewEncoder(w)
	response := AcceptResPegOut{
		Signature: signature,
	}

	err := enc.Encode(response)
	if err != nil {
		log.Error("error encoding response: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) acceptQuotePegOutHandler(w http.ResponseWriter, r *http.Request) {
	req := acceptReqPegout{}
	toRestAPI(w)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)

	if err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}

	quote, err := s.dbMongo.GetPegOutQuote(req.QuoteHash)

	if err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}

	if quote == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	p := getPegOutProviderByAddress(s.pegoutProviders, quote.LPRSKAddr)

	quoteHashInBytes, err := hex.DecodeString(req.QuoteHash)

	if err != nil {
		log.Error("error decoding string: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	signB, err := p.SignQuote(quoteHashInBytes, req.DerivationAddress, quote.Value)

	if err != nil {
		log.Error("error signing quote: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = s.addAddressWatcherToVerifyRegisterPegOut(quote, req.QuoteHash, req.DerivationAddress, signB, p, types.RQStateWaitingForDeposit)
	if err != nil {
		log.Error("error adding address watcher: ", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	signature := hex.EncodeToString(signB)
	returnQuotePegOutSignFunc(w, signature)
}

type SendBTCReq struct {
	Amount uint64 `json:"amount"`
	To     string `json:"to"`
}

type RegisterPegOutReg struct {
	quote     *pegout.Quote
	signature string
}

type BuildRefundPegOutPayloadRequest struct {
	QuoteHash         string `json:"quoteHash"`
	BtcTxHash         string `json:"btcTxHash"`
	DerivationAddress string `json:"derivationAddress"`
}

type BuildRefundPegOutPayloadResponse struct {
	Quote              *pegout.Quote `json:"quote"`
	MerkleBranchPath   int           `json:"merkleBranchPath"`
	MerkleBranchHashes []string      `json:"merkleBranchHashes"`
}

func (s *Server) refundPegOutHandler(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)
	payload := BuildRefundPegOutPayloadRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		log.Errorf(UnableToDeserializePayloadError, err)
		http.Error(w, UnableToDeserializePayloadError, http.StatusBadRequest)
		return
	}

	log.Printf("payload ::: %v", payload)

	quote, err := s.dbMongo.GetPegOutQuote(payload.QuoteHash)

	if err != nil {
		log.Errorf("Quote not found: %v", err)
		http.Error(w, "Quote not found", http.StatusBadRequest)
		return
	}

	branch, err := s.btc.BuildMerkleBranchByEndpoint(payload.BtcTxHash, payload.DerivationAddress)

	if err != nil {
		log.Errorf("Unable to create merkle branch: %v", err)
		http.Error(w, "Unable to create merkle branch", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	var hashes = make([]string, len(branch.Hashes))
	for i, hash := range branch.Hashes {
		hashes[i] = hash.String()
	}

	response := &BuildRefundPegOutPayloadResponse{
		Quote:              quote,
		MerkleBranchPath:   branch.Path,
		MerkleBranchHashes: hashes,
	}

	encoder := json.NewEncoder(w)

	err = encoder.Encode(&response)

	if err != nil {
		http.Error(w, UnableToBuildResponse, http.StatusInternalServerError)
		return
	}
}

type SenBTCRequest struct {
	Address string `json:"address"`
	Amount  uint   `json:"amount"`
}

type SenBTCResponse struct {
	TxHash string `json:"txHash"`
}

func (s *Server) sendBTC(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)
	enableCors(&w)
	payload := SenBTCRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		log.Errorf(UnableToDeserializePayloadError, err)
		http.Error(w, UnableToDeserializePayloadError, http.StatusBadRequest)
		return
	}

	txHash, err := s.btc.SendBTC(payload.Address, payload.Amount)

	if err != nil {
		log.Errorf(UnableToDeserializePayloadError, err)
		http.Error(w, "Unable to sendAddress", http.StatusBadRequest)
		return
	}

	response := &SenBTCResponse{
		TxHash: txHash,
	}

	encoder := json.NewEncoder(w)

	err = encoder.Encode(&response)

	if err != nil {
		http.Error(w, UnableToBuildResponse, http.StatusInternalServerError)
		return
	}
}

type AddCollateralRequest struct {
	Amount       uint64 `json:"amount" validate:"required"`
	LpRskAddress string `json:"lpRskAddress" validate:"required"`
}

type AddCollateralResponse struct {
	NewCollateralBalance uint64 `json:"newCollateralBalance"`
}

func (s *Server) addCollateral(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)
	enableCors(&w)
	payload := AddCollateralRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		log.Errorf(UnableToDeserializePayloadError, err)
		http.Error(w, UnableToDeserializePayloadError, http.StatusBadRequest)
		return
	}

	if isValid := Validate(payload)(w); !isValid {
		return
	}

	var lp pegout.LiquidityProvider
	for _, provider := range s.pegoutProviders {
		if provider.Address() == payload.LpRskAddress {
			lp = provider
		}
	}

	addrStr := lp.Address()

	c, min, err := s.rsk.GetCollateral(addrStr)

	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to get collateral", http.StatusInternalServerError)
		return
	}

	if min.Uint64()+payload.Amount < min.Uint64() {
		http.Error(w, "Amount is lower than min collateral", http.StatusBadRequest)
		return
	}

	addr := common.HexToAddress(addrStr)

	cmp := c.Cmp(big.NewInt(0))

	if cmp == 0 {
		http.Error(w, "LP not registered", http.StatusBadRequest)
		return
	}

	opts := &bind.TransactOpts{
		Value:  big.NewInt(int64(payload.Amount)),
		From:   addr,
		Signer: lp.SignTx,
	}

	err = s.rsk.AddCollateral(opts)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to add collateral", http.StatusInternalServerError)
		return
	}

	collateral, _, err := s.rsk.GetCollateral(addrStr)
	if err != nil {
		log.Error(err)
		http.Error(w, "Unable to get collateral", http.StatusInternalServerError)
		return
	}

	response := &AddCollateralResponse{
		NewCollateralBalance: collateral.Uint64(),
	}

	encoder := json.NewEncoder(w)

	err = encoder.Encode(&response)

	if err != nil {
		http.Error(w, UnableToBuildResponse, http.StatusInternalServerError)
		return
	}
}
