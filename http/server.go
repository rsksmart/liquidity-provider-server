package http

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/account"

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
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
const ErrorEstimatingGas = "Error on RSK Network, couldnt estimate gas price"
const ErrorValueHigherThanMaxAllowed = "value to transfer is higher than max allowed"
const ErrorStoringProviderQuote = "Error storing the quote on server"
const ErrorFetchingMongoDBProviders = "Error Fetching Providers from MongoDB: "
const ErrorSigningQuote = "error signing quote: "
const ErrorAddingAddressWatcher = "error signing quote: "
const ErrorBech32AddressNotSupported = "BECH32 address type is not supported yet"
const ErrorCreatingLocalProvider = "Error Creating New Local Provider"
const GetCollateralError = "Unable to get collateral"
const ErrorAddingProvider = "Error Adding New provider: %v"
const ErrorRetrivingProviderAddress = "Error Retrieving Provider Address from MongoDB"

type LiquidityProviderList struct {
	Endpoint                    string `env:"RSK_ENDPOINT"`
	LBCAddr                     string `env:"LBC_ADDR"`
	BridgeAddr                  string `env:"RSK_BRIDGE_ADDR"`
	RequiredBridgeConfirmations int64  `env:"RSK_REQUIRED_BRIDGE_CONFIRMATONS"`
	MaxQuoteValue               uint64 `env:"RSK_MAX_QUOTE_VALUE"`
}

type ConfigData struct {
	MaxQuoteValue uint64
	EncryptKey    string
	RSK           LiquidityProviderList
}

type Server struct {
	srv                  http.Server
	providers            []pegin.LiquidityProvider
	pegoutProviders      []pegout.LiquidityProvider
	rsk                  connectors.RSKConnector
	btc                  connectors.BTCConnector
	dbMongo              mongoDB.DBConnector
	now                  func() time.Time
	watchers             map[string]*BTCAddressWatcher
	pegOutWatchers       map[string]*BTCAddressPegOutWatcher
	pegOutDepositWatcher DepositEventWatcher
	addWatcherMu         sync.Mutex
	sharedPeginMutex     sync.Mutex
	sharedPegoutMutex    sync.Mutex
	cfgData              ConfigData
	ProviderRespository  *storage.LPRepository
	ProviderConfig       pegin.ProviderConfig
	AccountProvider      account.AccountProvider
}

type QuoteRequest struct {
	CallEoaOrContractAddress string `json:"callEoaOrContractAddress" required:"" validate:"required" example:"0x0" description:"Contract address or EOA address"`
	CallContractArguments    string `json:"callContractArguments" required:"" example:"0x0" description:"Contract data"`
	ValueToTransfer          uint64 `json:"valueToTransfer" required:"" example:"0x0" description:"Value to send in the call"`
	RskRefundAddress         string `json:"rskRefundAddress" required:"" validate:"required" example:"0x0" description:"User RSK refund address"`
	BitcoinRefundAddress     string `json:"bitcoinRefundAddress" required:"" validate:"required" example:"0x0" description:"User Bitcoin refund address. Note: Must be a legacy address, segwit addresses are not accepted"`
}

type QuoteReturn struct {
	Quote     *PeginQuoteDTO `json:"quote" required:"" description:"Detail of the quote"`
	QuoteHash string         `json:"quoteHash" required:"" description:"This is a 64 digit number that derives from a quote object"`
}

type QuotePegOutRequest struct {
	To                   string `json:"to" required:"" description:"Bitcoin address that will receive the BTC amount"`
	ValueToTransfer      uint64 `json:"valueToTransfer" required:"" example:"10000000000000" description:"ValueToTransfer"`
	RskRefundAddress     string `json:"rskRefundAddress" required:"" example:"0x0" description:"RskRefundAddress"`
	BitcoinRefundAddress string `json:"bitcoinRefundAddress" required:"" example:"0x0" description:"BitcoinRefundAddress"`
}

type QuotePegOutResponse struct {
	Quote     *PegoutQuoteDTO `json:"quote" required:"" description:"Quote"`
	QuoteHash string          `json:"quoteHash" required:"" example:"0x0" description:"QuoteHash"`
}

type acceptReq struct {
	QuoteHash string `json:"quoteHash" required:"" example:"0x0" description:"QuoteHash"`
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
	Signature                 string `json:"signature" required:"" example:"0x0" description:"Signature of the quote"`
	BitcoinDepositAddressHash string `json:"bitcoinDepositAddressHash" required:"" example:"0x0" description:"Hash of the deposit BTC address"`
}
type acceptResPegOut struct {
	Signature         string `json:"signature" required:"" example:"0x0" description:"Signature of the quote"`
	RskDepositAddress string `json:"rskDepositAddress" required:"" example:"0x0" description:"Hash of the deposit RSK address"`
}

type AcceptResPegOut struct {
	Signature string `json:"signature" required:"" example:"0x0" description:"Signature"`
}

func New(rsk connectors.RSKConnector, btc connectors.BTCConnector, dbMongo mongoDB.DBConnector, cfgData ConfigData,
	LPRep *storage.LPRepository, ProviderConfig pegin.ProviderConfig, accountProvider account.AccountProvider) Server {
	return newServer(rsk, btc, dbMongo, time.Now, cfgData, LPRep, ProviderConfig, accountProvider)
}

func newServer(rsk connectors.RSKConnector, btc connectors.BTCConnector, dbMongo mongoDB.DBConnector, now func() time.Time,
	cfgData ConfigData, LPRep *storage.LPRepository, ProviderConfig pegin.ProviderConfig, accountProvider account.AccountProvider) Server {
	return Server{
		rsk:                 rsk,
		btc:                 btc,
		dbMongo:             dbMongo,
		providers:           make([]pegin.LiquidityProvider, 0),
		pegoutProviders:     make([]pegout.LiquidityProvider, 0),
		now:                 now,
		watchers:            make(map[string]*BTCAddressWatcher),
		pegOutWatchers:      make(map[string]*BTCAddressPegOutWatcher),
		cfgData:             cfgData,
		ProviderRespository: LPRep,
		ProviderConfig:      ProviderConfig,
		AccountProvider:     accountProvider,
	}
}

func (s *Server) AddProvider(lp pegin.LiquidityProvider, ProviderDetails types.ProviderRegisterRequest) error {
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
			Value:  new(big.Int).Mul(m, big.NewInt(2)),
			From:   addr,
			Signer: lp.SignTx,
		}
		providerID, err := s.rsk.RegisterProvider(opts, ProviderDetails.Name, big.NewInt(int64(ProviderDetails.Fee)), big.NewInt(int64(ProviderDetails.QuoteExpiration)), big.NewInt(int64(ProviderDetails.AcceptedQuoteExpiration)), big.NewInt(int64(ProviderDetails.MinTransactionValue)), big.NewInt(int64(ProviderDetails.MaxTransactionValue)), ProviderDetails.ApiBaseUrl, ProviderDetails.Status)
		if err != nil {
			return err
		}
		err2 := s.dbMongo.InsertProvider(providerID, lp.Address())
		if err2 != nil {
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

func (s *Server) AddPegOutProvider(lp pegout.LiquidityProvider, ProviderDetails types.ProviderRegisterRequest) error {
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
			Value:  new(big.Int).Mul(m, big.NewInt(2)),
			From:   addr,
			Signer: lp.SignTx,
		}
		providerID, err := s.rsk.RegisterProvider(opts, ProviderDetails.Name, big.NewInt(int64(ProviderDetails.Fee)), big.NewInt(int64(ProviderDetails.QuoteExpiration)), big.NewInt(int64(ProviderDetails.AcceptedQuoteExpiration)), big.NewInt(int64(ProviderDetails.MinTransactionValue)), big.NewInt(int64(ProviderDetails.MaxTransactionValue)), ProviderDetails.ApiBaseUrl, ProviderDetails.Status)
		if err != nil {
			return err
		}
		err2 := s.dbMongo.InsertProvider(providerID, lp.Address())
		if err2 != nil {
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

type RegistrationStatus struct {
	Status string `json:"Status" example:"Provider Created Successfully" description:"Returned Status"`
}

// @Title Register Provider
// @Description Registers New Provider
// @Param  RegisterRequest  body types.ProviderRegisterRequest true "Provider Register Request"
// @Success  200 object RegistrationStatus
// @Route /provider/register [post]
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
	lp, err := pegin.NewLocalProvider(s.ProviderConfig, s.ProviderRespository, s.AccountProvider)
	if err != nil {
		log.Error(ErrorCreatingLocalProvider, err)
		http.Error(w, ErrorCreatingLocalProvider, http.StatusBadRequest)
		return
	}
	err = s.AddProvider(lp, payload)
	if err != nil {
		log.Errorf(ErrorAddingProvider, err)
		http.Error(w, ErrorAddingProvider, http.StatusBadRequest)
		return
	}
	response := RegistrationStatus{Status: "Provider Created Successfully"}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&response)
	if err != nil {
		http.Error(w, UnableToBuildResponse, http.StatusInternalServerError)
		return
	}
}

type ChangeStatusRequest struct {
	ProviderId uint64 `json:"providerId"`
	Status     bool   `json:"status"`
}
type ProviderStatusChangeStatus struct {
	Status string `json:"Status" example:"Provider Updated Successfully" description:"Returned Status"`
}

// @Title Change Provider Status
// @Description Changes the status of the provider
// @Param  ChangeStatusRequest  body ChangeStatusRequest true "Change Provider Status Request"
// @Success  200 object ProviderStatusChangeStatus
// @Route /provider/changeStatus [post]
func (s *Server) changeStatusHandler(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)
	enableCors(&w)
	payload := ChangeStatusRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Errorf(UnableToDeserializePayloadError, err)
		http.Error(w, UnableToDeserializePayloadError, http.StatusBadRequest)
		return
	}
	providerAddress, err := s.dbMongo.GetProvider(payload.ProviderId)
	if err != nil {
		log.Errorf(ErrorAddingProvider, err)
		http.Error(w, ErrorAddingProvider, http.StatusBadRequest)
		return
	}
	var lp pegin.LiquidityProvider
	for _, provider := range s.providers {
		if provider.Address() == providerAddress {
			lp = provider
		}
	}
	addrStr := lp.Address()
	addr := common.HexToAddress(addrStr)
	opts := &bind.TransactOpts{
		From:   addr,
		Signer: lp.SignTx,
	}
	err = s.rsk.ChangeStatus(opts, new(big.Int).SetUint64(payload.ProviderId), payload.Status)
	log.Debug(err)
	if err != nil {
		log.Errorf(ErrorAddingProvider, err)
		http.Error(w, ErrorAddingProvider, http.StatusBadRequest)
		return
	}
	response := "Provider Updated Successfully"
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
	r.Path("/pegout/getQuotes").Methods(http.MethodPost).HandlerFunc(s.getPegoutQuoteHandler)
	r.Path("/pegout/acceptQuote").Methods(http.MethodPost).HandlerFunc(s.acceptQuotePegOutHandler)
	r.Path("/pegout/refundPegOut").Methods(http.MethodPost).HandlerFunc(s.refundPegOutHandler)
	r.Path("/pegout/sendBTC").Methods(http.MethodPost).HandlerFunc(s.sendBTC)
	r.Path("/collateral").Methods(http.MethodGet).HandlerFunc(s.getCollateralHandler)
	r.Path("/addCollateral").Methods(http.MethodPost).HandlerFunc(s.addCollateral)
	r.Path("/withdrawCollateral").Methods(http.MethodPost).HandlerFunc(s.withdrawCollateral)
	r.Path("/provider/register").Methods(http.MethodPost).HandlerFunc(s.registerProviderHandler)
	r.Path("/provider/changeStatus").Methods(http.MethodPost).HandlerFunc(s.changeStatusHandler)
	r.Path("/provider/resignation").Methods(http.MethodPost).HandlerFunc(s.providerResignHandler)
	r.Methods("OPTIONS").HandlerFunc(s.handleOptions)
	w := log.StandardLogger().WriterLevel(log.DebugLevel)
	h := handlers.LoggingHandler(w, r)
	defer func(w *io.PipeWriter) {
		_ = w.Close()
	}(w)

	err := s.initPeginWatchers()
	if err != nil {
		return err
	}

	provider := s.pegoutProviders[0] // TODO convert providers array into normal variable since its going to be only 1 LP per LPS instance
	s.pegOutDepositWatcher = NewDepositEventWatcher(time.Minute*2, provider, &s.addWatcherMu, &s.sharedPegoutMutex, make(chan bool), s.rsk, s.btc, s.dbMongo,
		func(hash string, quote *WatchedQuote, endState types.RQState) {
			if endState != types.RQStateCallForUserSucceeded {
				return
			}
			signB, err := hex.DecodeString(quote.Signature)
			if err != nil {
				log.Error("Error decoding pegout quote signature: ", err)
			}
			err = s.addAddressPegOutWatcher(quote.Data, hash, quote.Data.DepositAddr, signB, provider, types.RQStateCallForUserSucceeded)
			if err != nil {
				log.Error("Error starting BTC pegout watcher: ", err)
			}
		})

	err = s.initPegoutWatchers()
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
		s.pegOutDepositWatcher.EndChannel() <- true
		return err
	}
	s.pegOutDepositWatcher.EndChannel() <- true
	return nil
}

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) initPegoutWatchers() error {
	quoteStatesToWatch := []types.RQState{types.RQStateCallForUserSucceeded, types.RQStateWaitingForDeposit}
	quotes, err := s.dbMongo.GetRetainedPegOutQuoteByState(quoteStatesToWatch)
	if err != nil {
		return err
	}
	var waitingForDepositQuotes, waitingForConfirmationQuotes map[string]*WatchedQuote
	for _, entry := range quotes {
		quote, err := s.dbMongo.GetPegOutQuote(entry.QuoteHash)
		if err != nil || quote == nil {
			log.Errorf("initPegoutWatchers: quote not found for hash: %s. Watcher not initialized for address %s", entry.QuoteHash, entry.DepositAddr)
			continue
		}

		p := pegout.GetPegoutProviderByAddress(s.pegoutProviders, quote.LPRSKAddr)
		if p == nil {
			log.Errorf("initPegoutWatchers: provider not found for LPRSKAddr: %s. Watcher not initialized for address %s", quote.LPRSKAddr, entry.DepositAddr)
			continue
		}

		signB, err := hex.DecodeString(entry.Signature)
		if err != nil {
			log.Errorf("initPeginBtcWatchers: couldn't decode signature %s for quote %s. Watcher not initialized for address %s", entry.Signature, entry.QuoteHash, entry.DepositAddr)
			continue
		}

		if entry.State == types.RQStateCallForUserSucceeded {
			err = s.addAddressPegOutWatcher(quote, entry.QuoteHash, quote.DepositAddr, signB, p, entry.State)
		} else if entry.State == types.RQStateWaitingForDepositConfirmations {
			waitingForConfirmationQuotes[entry.QuoteHash] = &WatchedQuote{Signature: entry.Signature, Data: quote, DepositBlock: entry.DepositBlockNumber}
		} else {
			waitingForDepositQuotes[entry.QuoteHash] = &WatchedQuote{Signature: entry.Signature, Data: quote}
		}

		if err != nil {
			log.Errorf("initPegoutWatchers: error initializing watcher for quote hash %s: %v", entry.QuoteHash, err)
		}
	}
	go s.pegOutDepositWatcher.Init(waitingForDepositQuotes, waitingForConfirmationQuotes)
	return nil
}

func (s *Server) initPeginWatchers() error {
	quoteStatesToWatch := []types.RQState{types.RQStateWaitingForDeposit, types.RQStateCallForUserSucceeded}
	retainedQuotes, err := s.dbMongo.GetRetainedQuotes(quoteStatesToWatch)
	if err != nil {
		return err
	}

	for _, entry := range retainedQuotes {
		quote, err := s.dbMongo.GetQuote(entry.QuoteHash)
		if err != nil || quote == nil {
			log.Errorf("initPeginBtcWatchers: quote not found for hash: %s. Watcher not initialized for address %s", entry.QuoteHash, entry.DepositAddr)
			continue
		}

		p := pegin.GetPeginProviderByAddress(s.providers, quote.LPRSKAddr)
		if p == nil {
			log.Errorf("initPeginBtcWatchers: provider not found for LPRSKAddr: %s. Watcher not initialized for address %s", quote.LPRSKAddr, entry.DepositAddr)
			continue
		}

		signB, err := hex.DecodeString(entry.Signature)
		if err != nil {
			log.Errorf("initPeginBtcWatchers: couldn't decode signature %s for quote %s. Watcher not initialized for address %s", entry.Signature, entry.QuoteHash, entry.DepositAddr)
			continue
		}

		err = s.addAddressWatcher(quote, entry.QuoteHash, entry.DepositAddr, signB, p, entry.State)
		if err != nil {
			log.Errorf("initPeginBtcWatchers: error initializing watcher for quote hash %s: %v", entry.QuoteHash, err)
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
	watcher := NewBTCAddressWatcher(hash, s.btc, s.rsk, provider, s.dbMongo, quote, signB, state, &s.sharedPeginMutex)
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

	satoshis, _ := quote.Value.ToSatoshi().Float64()
	minBtcAmount := btcutil.Amount(uint64(math.Ceil(satoshis)))
	expTime := quote.GetExpirationTime()
	watcher := &BTCAddressPegOutWatcher{
		hash:                 hash,
		btc:                  s.btc,
		addressDecryptionKey: s.cfgData.EncryptKey,
		rsk:                  s.rsk,
		lp:                   provider,
		dbMongo:              s.dbMongo,
		quote:                quote,
		state:                state,
		signature:            signB,
		done:                 make(chan struct{}),
		sharedLocker:         &s.sharedPegoutMutex,
	}
	err := s.btc.AddAddressPegOutWatcher(depositAddr, minBtcAmount, time.Minute, expTime, watcher, func(w connectors.AddressWatcher) {
		s.addWatcherMu.Lock()
		defer s.addWatcherMu.Unlock()
		delete(s.pegOutWatchers, hash)
	})
	if err == nil {
		s.addWatcherMu.Lock()
		s.pegOutWatchers[hash] = watcher
		s.addWatcherMu.Unlock()
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

type services struct {
	Db  string `json:"db"`
	Rsk string `json:"rsk"`
	Btc string `json:"btc"`
}
type healthRes struct {
	Status   string   `json:"status" example:"ok" description:"Overall LPS Health Status"`
	Services services `json:"services" example:"{\"db\":\"ok\",\"rsk\":\"ok\",\"btc\":\"ok\"}" description:"LPS Services Status"`
}

// @Title Health
// @Description Returns server health.
// @Success  200  object healthRes
// @Route /health [get]
func (s *Server) checkHealthHandler(w http.ResponseWriter, _ *http.Request) {
	enableCors(&w)
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

func (a *QuotePegOutRequest) validateQuoteRequest() string {
	err := ""

	if a.ValueToTransfer == 0 {
		err += "Value to Transfer cannot be empty or zero!"
	}

	return err
}

// @Title Get Providers
// @Description Returns a list of providers.
// @Success  200  array ProviderDTO
// @Route /getProviders [get]
func (s *Server) getProvidersHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	providerList, error := s.dbMongo.GetProviders()
	if error != nil {
		log.Error("Error fetching providers. Error: ", error)
		customError := NewServerError(ErrorFetchingMongoDBProviders+error.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}
	providers, error := s.rsk.GetProviders(providerList)

	if error != nil {
		log.Error("GetProviders - error encoding response: ", error)
		customError := NewServerError("GetProviders - error encoding response: "+error.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	response := make([]*ProviderDTO, 0)
	for _, provider := range providers {
		response = append(response, toProviderDTO(&provider))
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(&response)
	if err != nil {
		log.Error("error encoding registered providers list: ", err.Error())
		customError := NewServerError("error encoding registered providers list: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}
}

// @Title Pegin GetQuote
// @Description Gets Pegin Quote
// @Param  PeginQuoteRequest  body QuoteRequest true "Interface with parameters for computing possible quotes for the service"
// @Success  200  array QuoteReturn The quote structure defines the conditions of a service, and acts as a contract between users and LPs
// @Route /pegin/getQuote [post]
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
	if isValid := Validate(qr)(w); !isValid {
		return
	}

	maxValueTotransfer := s.cfgData.MaxQuoteValue

	if maxValueTotransfer <= 0 {
		maxValueTotransfer = uint64(s.cfgData.RSK.MaxQuoteValue)
	}

	if qr.ValueToTransfer > maxValueTotransfer {
		log.Error(ErrorValueHigherThanMaxAllowed)
		details := map[string]interface{}{
			"maxValueTotransfer": maxValueTotransfer,
			"valueToTransfer":    qr.ValueToTransfer,
		}
		customError := NewServerError(ErrorValueHigherThanMaxAllowed, details, true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	var gas uint64
	gas, err = s.rsk.EstimateGas(qr.CallEoaOrContractAddress, big.NewInt(int64(qr.ValueToTransfer)), []byte(qr.CallContractArguments))

	if err != nil {
		log.Error(ErrorEstimatingGas, err.Error())
		customError := NewServerError(ErrorEstimatingGas, make(Details), true)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}

	price, err := s.rsk.GasPrice()
	if err != nil {
		log.Error(ErrorEstimatingGas, err.Error())
		customError := NewServerError(ErrorEstimatingGas, make(map[string]interface{}), true)
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
				errmsg := ErrorStoringProviderQuote + ": " + err.Error()
				status := http.StatusInternalServerError
				if strings.HasPrefix(err.Error(), "VM Exception") {
					_, vmString, _ := strings.Cut(err.Error(), "VM Exception while processing transaction: revert ")
					status = http.StatusBadRequest
					errmsg = "LBC error: " + vmString
				}
				customError := NewServerError(errmsg, make(map[string]interface{}), false)
				ResponseError(w, customError, status)
				return
			} else {
				quotes = append(quotes, &QuoteReturn{toPeginQuote(pq), hash})
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

// @Title Pegout GetQuote
// @Description Gets Pegout Quote
// @Param PegoutQuoteRequest body QuotePegOutRequest true "Interface with parameters for computing possible quotes for the service"
// @Success 200 array QuotePegOutResponse The quote structure defines the conditions of a service, and acts as a contract between users and LPs
// @Route /pegout/getQuotes [post]
func (s *Server) getPegoutQuoteHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	qr := QuotePegOutRequest{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&qr)

	if err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}
	log.Debug("received peg out quote request: ", fmt.Sprintf("%+v", qr))
	if isValid := Validate(qr)(w); !isValid {
		return
	}

	maxValueTotransfer := s.cfgData.MaxQuoteValue

	if maxValueTotransfer <= 0 {
		maxValueTotransfer = uint64(s.cfgData.RSK.MaxQuoteValue)
	}

	if qr.ValueToTransfer > maxValueTotransfer {
		log.Error(ErrorValueHigherThanMaxAllowed)
		details := map[string]interface{}{
			"maxValueTotransfer": maxValueTotransfer,
			"valueToTransfer":    qr.ValueToTransfer,
		}
		customError := NewServerError(ErrorValueHigherThanMaxAllowed, details, true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	var gas uint64
	gas, err = s.rsk.EstimateGas(s.rsk.GetBridgeAddress().Hex(), big.NewInt(int64(qr.ValueToTransfer)), []byte(nil))

	if err != nil {
		log.Error(ErrorEstimatingGas, err.Error())
		customError := NewServerError(ErrorEstimatingGas, make(Details), true)
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

	var quotes []*QuotePegOutResponse
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
	q := parseReqToPegOutQuote(qr, s.rsk.GetLBCAddress(), gas)
	rskBlockNumber, err := s.rsk.GetRskHeight()
	if err != nil {
		log.Error("Error getting last block", err.Error())
		customError := NewServerError("Error getting last block", make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}
	for _, p := range s.pegoutProviders {
		pq, err := p.GetQuote(q, rskBlockNumber, gas, types.NewBigWei(price))
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

			hash, err := s.storePegoutQuote(pq)

			if err != nil {
				log.Error(err)
				errmsg := ErrorStoringProviderQuote + ": " + err.Error()
				status := http.StatusInternalServerError
				if strings.HasPrefix(err.Error(), "VM Exception") {
					_, vmString, _ := strings.Cut(err.Error(), "VM Exception while processing transaction: revert ")
					status = http.StatusBadRequest
					errmsg = "LBC error: " + vmString
				}
				customError := NewServerError(errmsg, make(map[string]interface{}), false)
				ResponseError(w, customError, status)
				return
			} else {
				quotes = append(quotes, &QuotePegOutResponse{toPegoutQuote(pq), hash})
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
			details := Details{
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

// @Title Accept Quote
// @Description Accepts Quote
// @Param  QuoteHash  body acceptReq true "Quote Hash"
// @Success  200  object acceptRes Interface that represents that the quote has been successfully accepted
// @Route /pegin/acceptQuote [post]
func (s *Server) acceptQuoteHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	returnQuoteSignFunc := func(w http.ResponseWriter, signature string, depositAddr string) {
		enc := json.NewEncoder(w)
		response := acceptRes{
			Signature:                 signature,
			BitcoinDepositAddressHash: depositAddr,
		}

		err := enc.Encode(response)
		if err != nil {
			const errorMsg = "AcceptQuote - error encoding response: "
			log.Error(errorMsg, err.Error())
			customError := NewServerError(errorMsg+err.Error(), make(map[string]interface{}), true)
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

	p := pegin.GetPeginProviderByAddress(s.providers, quote.LPRSKAddr)
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
		log.Error(ErrorSigningQuote, err.Error())
		customError := NewServerError(ErrorSigningQuote+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	err = s.addAddressWatcher(quote, req.QuoteHash, depositAddress, signB, p, types.RQStateWaitingForDeposit)
	if err != nil {
		log.Error(ErrorAddingAddressWatcher, err.Error())
		customError := NewServerError(ErrorAddingAddressWatcher+err.Error(), make(map[string]interface{}), true)
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
		ContractAddr:  qr.CallEoaOrContractAddress,
		Data:          qr.CallContractArguments,
		Value:         types.NewWei(int64(qr.ValueToTransfer)),
		GasLimit:      uint32(limitGas),
	}
}

func parseReqToPegOutQuote(qr QuotePegOutRequest, lbcAddr string, limitGas uint64) *pegout.Quote {
	return &pegout.Quote{
		LBCAddr:       lbcAddr,
		BtcRefundAddr: qr.BitcoinRefundAddress,
		RSKRefundAddr: qr.RskRefundAddress,
		DepositAddr:   qr.To,
		Value:         types.NewWei(int64(qr.ValueToTransfer)),
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

func decodeAddressesPegOut(rskRefundAddr string, lpRSKAddr string, lbcAddr string) ([]byte, []byte, []byte, error) {
	rskRefAddrB, err := connectors.DecodeRSKAddress(rskRefundAddr)
	if err != nil {
		return nil, nil, nil, err
	}
	lpRSKAddrB, err := connectors.DecodeRSKAddress(lpRSKAddr)
	if err != nil {
		return nil, nil, nil, err
	}
	lbcAddrB, err := connectors.DecodeRSKAddress(lbcAddr)
	if err != nil {
		return nil, nil, nil, err
	}
	return rskRefAddrB, lpRSKAddrB, lbcAddrB, nil
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

func (s *Server) storePegoutQuote(q *pegout.Quote) (string, error) {
	h, err := s.rsk.HashPegOutQuote(q)
	if err != nil {
		return "", err
	}

	err = s.dbMongo.InsertPegOutQuote(h, q)
	if err != nil {
		log.Fatalf("error inserting quote: %v", err)
		return "", err
	}
	return h, nil
}

func getQuoteExpTime(q *pegin.Quote) time.Time {
	return time.Unix(int64(q.AgreementTimestamp+q.TimeForDeposit), 0)
}

func getQuoteExpTimePegOut(q *pegout.Quote) time.Time {
	return time.Unix(int64(q.AgreementTimestamp+q.DepositDateLimit), 0)
}

func buildErrorDecodingRequest(w http.ResponseWriter, err error) {
	log.Error("Error decoding request: ", err.Error())
	customError := NewServerError(fmt.Sprintf("Error decoding request: %s", err.Error()), make(Details), true)
	ResponseError(w, customError, http.StatusBadRequest)
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

func generateRskEthereumAddress() ([]byte, []byte, common.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("SAVE BUT DO NOT SHARE THIS (Private Key):", hexutil.Encode(privateKeyBytes))

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("Public Key:", hexutil.Encode(publicKeyBytes))

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("Address:", address)

	return privateKeyBytes, publicKeyBytes, address, nil
}

// @Title Accept Quote Pegout
// @Description Accepts Quote Pegout
// @Param  QuoteHash  body acceptReq true "Quote Hash"
// @Success 200 object acceptResPegOut
// @Route /pegout/acceptQuote [post]
func (s *Server) acceptQuotePegOutHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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

	quote, err := s.dbMongo.GetPegOutQuote(req.QuoteHash)
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

	expTime := getQuoteExpTimePegOut(quote)
	if s.now().After(expTime) {
		log.Error("quote deposit time has elapsed; hash: ", req.QuoteHash)
		customError := NewServerError("quote deposit time has elapsed; hash: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	rq, err := s.dbMongo.GetRetainedPegOutQuote(req.QuoteHash)
	if err != nil {
		log.Error("error fetching retained quote: ", err.Error())
		customError := NewServerError("error fetching retained quote: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	if rq != nil { // if the quote has already been accepted, just return signature and deposit addr
		signAndReturnPegoutQuote(w, rq.Signature, rq.DepositAddr)
		return
	}

	p := pegout.GetPegoutProviderByAddress(s.pegoutProviders, quote.LPRSKAddr)
	gasPrice, err := s.rsk.GasPrice()
	if err != nil {
		log.Error("error getting provider by address: ", err.Error())
		customError := NewServerError("error getting provider by address: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	adjustedGasLimit := types.NewUWei(uint64(CFUExtraGas) + uint64(quote.GasLimit))
	gasCost := new(types.Wei).Mul(adjustedGasLimit, types.NewBigWei(gasPrice))
	reqLiq := gasCost.Uint64() + quote.Value.Uint64()
	signB, err := p.SignQuote(hashBytes, s.rsk.GetLBCAddress(), reqLiq)
	if err != nil {
		log.Error(ErrorSigningQuote, err.Error())
		customError := NewServerError(ErrorSigningQuote+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	signature := hex.EncodeToString(signB)

	err = s.pegOutDepositWatcher.WatchNewQuote(req.QuoteHash, signature, quote)
	if err != nil {
		log.Error(ErrorAddingAddressWatcher, err.Error())
		customError := NewServerError(ErrorAddingAddressWatcher+err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusConflict)
		return
	}
	signAndReturnPegoutQuote(w, signature, s.rsk.GetLBCAddress())
}

func signAndReturnPegoutQuote(w http.ResponseWriter, signature string, depositAddr string) {
	enc := json.NewEncoder(w)
	response := acceptResPegOut{
		Signature:         signature,
		RskDepositAddress: depositAddr,
	}

	err := enc.Encode(response)
	if err != nil {
		const errorMsg = "AcceptQuotePegout - error encoding response: "
		log.Error(errorMsg, err.Error())
		customError := NewServerError(errorMsg+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}
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
	QuoteHash         string `json:"quoteHash" example:"0x0" description:"QuoteHash"`
	BtcTxHash         string `json:"btcTxHash" example:"0x0" description:"BtcTxHash"`
	DerivationAddress string `json:"derivationAddress" example:"0x0" description:"DerivationAddress"`
}

type BuildRefundPegOutPayloadResponse struct {
	Quote              *pegout.Quote `json:"quote" example:"0x0" description:"Quote"`
	MerkleBranchPath   int           `json:"merkleBranchPath" example:"0x0" description:"MerkleBranchPath"`
	MerkleBranchHashes []string      `json:"merkleBranchHashes" example:"0x0" description:"MerkleBranchHashes"`
}

// @Title Refund Pegout
// @Description Refunds Pegout
// @Param  RefundPegout  body BuildRefundPegOutPayloadRequest true "Pegout Refund Details"
// @Success  200  object BuildRefundPegOutPayloadResponse
// @Route /pegout/refundPegOut [post]
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
	Address string `json:"address" example:"0x0" description:"Address to send BTC to"`
	Amount  uint   `json:"amount" example:"100000000000" description:"Amount to send BTC to address"`
}

type SenBTCResponse struct {
	TxHash string `json:"txHash" example:"0x0" description:"TxHash of the BTC transaction sent to the address"`
}

// @Title Send BTC
// @Description Sends BTC
// @Param  SendBTCRequest  body SenBTCRequest true "Send BTC Request"
// @Success  200  object SenBTCResponse
// @Route /pegout/sendBTC [post]
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

	txHash, err := s.btc.SendBtc(payload.Address, uint64(payload.Amount))

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
	Amount       uint64 `json:"amount" validate:"required" example:"100000000000" description:"Amount to add to the collateral"`
	LpRskAddress string `json:"lpRskAddress" validate:"required,eth_addr" example:"0x0" description:"Liquidity Provider RSK Address"`
}

type AddCollateralResponse struct {
	NewCollateralBalance uint64 `json:"newCollateralBalance" example:"100000000000" description:"New Collateral Balance`
}

// @Title Add Collateral
// @Description Adds Collateral
// @Param  AddCollateralRequest  body AddCollateralRequest true "Add Collateral Request"
// @Success  200  object SenBTCResponse
// @Route /addCollateral [post]
func (s *Server) addCollateral(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)
	enableCors(&w)
	payload := AddCollateralRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		customError := NewServerError(fmt.Sprintf(UnableToDeserializePayloadError, err.Error()), make(Details), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	if isValid := Validate(payload)(w); !isValid {
		return
	}

	lp := pegin.GetPeginProviderByAddress(s.providers, payload.LpRskAddress)
	if lp == nil {
		customError := NewServerError("missing liquidity provider", make(Details), true)
		ResponseError(w, customError, http.StatusNotFound)
		return
	}

	addrStr := lp.Address()

	collateral, min, err := s.rsk.GetCollateral(addrStr)

	if err != nil {
		log.Error(err)
		customError := NewServerError(GetCollateralError, *NewBasicDetail(err), false)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	} else if collateral.Uint64()+payload.Amount < min.Uint64() {
		customError := NewServerError("Amount is lower than min collateral", make(Details), true)
		ResponseError(w, customError, http.StatusConflict)
		return
	}

	opts := &bind.TransactOpts{
		Value:  big.NewInt(int64(payload.Amount)),
		From:   common.HexToAddress(addrStr),
		Signer: lp.SignTx,
	}

	err = s.rsk.AddCollateral(opts)
	if err != nil {
		log.Error(err)
		customError := NewServerError(GetCollateralError, *NewBasicDetail(err), false)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}

	collateral, _, err = s.rsk.GetCollateral(addrStr)
	if err != nil {
		log.Error(err)
		customError := NewServerError(GetCollateralError, *NewBasicDetail(err), false)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}

	response := &AddCollateralResponse{
		NewCollateralBalance: collateral.Uint64(),
	}

	JsonResponse(w, http.StatusOK, &response)
}

type WithdrawCollateralRequest struct {
	LpRskAddress string `json:"lpRskAddress" validate:"required,eth_addr"`
}

// @Title Withdraw Collateral
// @Description Withdraw Collateral of a resigned LP
// @Param  WithdrawCollateralRequest  body WithdrawCollateralRequest true "Withdraw Collateral Request"
// @Route /withdrawCollateral [post]
// @Success 204 object
func (s *Server) withdrawCollateral(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	payload := WithdrawCollateralRequest{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}
	if isValid := Validate(payload)(w); !isValid {
		return
	}

	opts, err := pegin.GetPeginProviderTransactOpts(s.providers, payload.LpRskAddress)
	if err != nil {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusNotFound)
		return
	}

	if err := s.rsk.WithdrawCollateral(opts); err != nil && errors.Is(err, connectors.WithdrawCollateralError) {
		customError := NewServerError(fmt.Sprintf("%s, please complete resign proccess first", err.Error()), make(Details), true)
		ResponseError(w, customError, http.StatusConflict)
	} else if err != nil {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

type GetCollateralResponse struct {
	Collateral uint64 `json:"collateral"`
}

// @Title Get Collateral
// @Description Get Collateral
// @Param address path  string  true  "Liquidity provider address"
// @Success  200  object GetCollateralResponse
// @Route /collateral/{address} [get]
func (s *Server) getCollateralHandler(w http.ResponseWriter, request *http.Request) {
	address := request.URL.Query().Get("address")
	collateral, _, err := s.rsk.GetCollateral(address)

	var e *connectors.AddressError
	if err != nil && errors.As(err, &e) {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusBadRequest)
	} else if err != nil {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusInternalServerError)
	} else if collateral.Uint64() == 0 {
		customError := NewServerError("no collateral found", make(Details), true)
		ResponseError(w, customError, http.StatusNotFound)
	} else {
		response := &GetCollateralResponse{Collateral: collateral.Uint64()}
		JsonResponse(w, http.StatusOK, response)
	}
}

type ProviderResignRequest struct {
	LpRskAddress string `json:"lpRskAddress" validate:"required,eth_addr"`
}

// @Title Provider resignation
// @Description Provider stops being a liquidity provider
// @Param  ProviderResignRequest  body ProviderResignRequest true "Provider Resignation Request"
// @Route /provider/resignation [post]
// @Success 204 object
func (s *Server) providerResignHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	payload := ProviderResignRequest{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}
	if isValid := Validate(payload)(w); !isValid {
		return
	}

	opts, err := pegin.GetPeginProviderTransactOpts(s.providers, payload.LpRskAddress)
	if err != nil {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusNotFound)
		return
	}

	err = s.rsk.Resign(opts)
	if err != nil && errors.Is(err, connectors.ProviderResignError) {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusConflict)
	} else if err != nil {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
