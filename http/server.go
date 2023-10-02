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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rsksmart/liquidity-provider-server/account"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"io"
	"math"
	"math/big"

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
const PegoutDepositCheckInterval = 2 * time.Minute

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
const ErrorNotLiquidity = "Not enough liquidity"

type LiquidityProviderList struct {
	Endpoint                    string   `env:"RSK_ENDPOINT"`
	LBCAddr                     string   `env:"LBC_ADDR"`
	BridgeAddr                  string   `env:"RSK_BRIDGE_ADDR"`
	RequiredBridgeConfirmations int64    `env:"RSK_REQUIRED_BRIDGE_CONFIRMATONS"`
	LpsAddress                  string   `env:"LIQUIDITY_PROVIDER_RSK_ADDR"`
	ChainId                     *big.Int `env:"CHAIN_ID"`
}

type ConfigData struct {
	RSK                  LiquidityProviderList
	QuoteCacheStartBlock uint64
	CaptchaSecretKey     string
	CaptchaSiteKey       string
	CaptchaThreshold     float32
}

type Server struct {
	srv                  http.Server
	provider             pegin.LiquidityProvider
	pegoutProvider       pegout.LiquidityProvider
	rsk                  connectors.RSKConnector
	btc                  connectors.BTCConnector
	dbMongo              mongoDB.DBConnector
	now                  func() time.Time
	watchers             map[string]*BTCAddressWatcher
	pegOutWatchers       map[string]*BTCAddressPegOutWatcher
	pegOutDepositWatcher DepositEventWatcher
	lpFundsEventtWatcher LpFundsEventWatcher
	addWatcherMu         sync.Mutex
	sharedPeginMutex     sync.Mutex
	sharedPegoutMutex    sync.Mutex
	cfgData              ConfigData
	ProviderRespository  *storage.LPRepository
	ProviderConfig       pegin.ProviderConfig
	PegoutConfig         pegout.ProviderConfig
	AccountProvider      account.AccountProvider
	awsConfig            aws.Config
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

type acceptRes struct {
	Signature                 string `json:"signature" required:"" example:"0x0" description:"Signature of the quote"`
	BitcoinDepositAddressHash string `json:"bitcoinDepositAddressHash" required:"" example:"0x0" description:"Hash of the deposit BTC address"`
	FlyoverRedeemScript       string `json:"-"`
}
type acceptResPegOut struct {
	Signature  string `json:"signature" required:"" example:"0x0" description:"Signature of the quote"`
	LbcAddress string `json:"lbcAddress" required:"" example:"0x0" description:"LBC address to execute depositPegout function"`
}

type AcceptResPegOut struct {
	Signature string `json:"signature" required:"" example:"0x0" description:"Signature"`
}

func New(rsk connectors.RSKConnector, btc connectors.BTCConnector, dbMongo mongoDB.DBConnector, cfgData ConfigData,
	LPRep *storage.LPRepository, ProviderConfig pegin.ProviderConfig, pegoutConfig pegout.ProviderConfig, accountProvider account.AccountProvider, awsConfig aws.Config) Server {
	return newServer(rsk, btc, dbMongo, time.Now, cfgData, LPRep, ProviderConfig, pegoutConfig, accountProvider, awsConfig)
}

func newServer(rsk connectors.RSKConnector, btc connectors.BTCConnector, dbMongo mongoDB.DBConnector, now func() time.Time,
	cfgData ConfigData, LPRep *storage.LPRepository, ProviderConfig pegin.ProviderConfig, pegoutConfig pegout.ProviderConfig, accountProvider account.AccountProvider,
	awsConfig aws.Config) Server {
	return Server{
		rsk:                 rsk,
		btc:                 btc,
		dbMongo:             dbMongo,
		provider:            nil,
		pegoutProvider:      nil,
		now:                 now,
		watchers:            make(map[string]*BTCAddressWatcher),
		pegOutWatchers:      make(map[string]*BTCAddressPegOutWatcher),
		cfgData:             cfgData,
		ProviderRespository: LPRep,
		ProviderConfig:      ProviderConfig,
		PegoutConfig:        pegoutConfig,
		AccountProvider:     accountProvider,
		awsConfig:           awsConfig,
	}
}

func (s *Server) AddProvider(peginProvider pegin.LiquidityProvider, pegoutProvider pegout.LiquidityProvider, providerDetails types.ProviderRegisterRequest) error {
	var peginCollateral, pegoutCollateral, minCollateral *big.Int
	var operationalForPegin, operationalForPegout bool
	var err error

	s.provider = peginProvider
	s.pegoutProvider = pegoutProvider

	if providerDetails.ProviderType != "pegin" && providerDetails.ProviderType != "pegout" && providerDetails.ProviderType != "both" {
		return errors.New("invalid provider type")
	}

	if peginCollateral, minCollateral, err = s.rsk.GetCollateral(peginProvider.Address()); err != nil {
		return err
	}

	if pegoutCollateral, _, err = s.rsk.GetPegoutCollateral(pegoutProvider.Address()); err != nil {
		return err
	}

	if operationalForPegin, err = s.rsk.IsOperational(&bind.CallOpts{}, common.HexToAddress(peginProvider.Address())); err != nil {
		return err
	}

	if operationalForPegout, err = s.rsk.IsOperationalForPegout(&bind.CallOpts{}, common.HexToAddress(peginProvider.Address())); err != nil {
		return err
	}

	if isProviderRegistered(providerDetails.ProviderType, operationalForPegin, operationalForPegout) {
		log.Debug("Already registered")
		return nil
	}

	if (providerDetails.ProviderType == "pegin" || providerDetails.ProviderType == "both") && !operationalForPegin && peginCollateral.Cmp(big.NewInt(0)) != 0 {
		return s.addPeginCollateral(peginProvider, peginCollateral, minCollateral)
	}
	if (providerDetails.ProviderType == "pegout" || providerDetails.ProviderType == "both") && !operationalForPegout && pegoutCollateral.Cmp(big.NewInt(0)) != 0 {
		return s.addPegoutCollateral(pegoutProvider, pegoutCollateral, minCollateral)
	}

	return s.performRegisterProvider(peginProvider, pegoutProvider, providerDetails, minCollateral)
}

func isProviderRegistered(providerType string, isOperationalForPegin, isOperationalForPegout bool) bool {
	return (providerType == "both" && isOperationalForPegin && isOperationalForPegout) ||
		(providerType == "pegin" && isOperationalForPegin) ||
		(providerType == "pegout" && isOperationalForPegout)
}

func (s *Server) addPeginCollateral(peginProvider pegin.LiquidityProvider, peginCollateral, minCollateral *big.Int) error {
	if peginCollateral.Cmp(minCollateral) >= 0 {
		return nil
	}
	opts := &bind.TransactOpts{
		Value:  minCollateral.Sub(minCollateral, peginCollateral),
		From:   common.HexToAddress(peginProvider.Address()),
		Signer: peginProvider.SignTx,
	}
	return s.rsk.AddCollateral(opts)
}

func (s *Server) addPegoutCollateral(pegoutProvider pegout.LiquidityProvider, pegoutCollateral, minCollateral *big.Int) error {
	if pegoutCollateral.Cmp(minCollateral) >= 0 {
		return nil
	}
	opts := &bind.TransactOpts{
		Value:  minCollateral.Sub(minCollateral, pegoutCollateral),
		From:   common.HexToAddress(pegoutProvider.Address()),
		Signer: pegoutProvider.SignTx,
	}
	return s.rsk.AddPegoutCollateral(opts)
}

func (s *Server) performRegisterProvider(peginProvider pegin.LiquidityProvider, pegoutProvider pegout.LiquidityProvider,
	providerDetails types.ProviderRegisterRequest, minCollateral *big.Int) error {
	log.Debug("Registering new provider...")
	var signer bind.SignerFn
	var address string
	if providerDetails.ProviderType == "pegin" || providerDetails.ProviderType == "both" {
		address = peginProvider.Address()
		signer = peginProvider.SignTx
	} else {
		address = pegoutProvider.Address()
		signer = pegoutProvider.SignTx
	}
	opts := &bind.TransactOpts{
		Value:  new(big.Int).Mul(minCollateral, big.NewInt(2)),
		From:   common.HexToAddress(address),
		Signer: signer,
	}

	providerID, err := s.rsk.RegisterProvider(opts, providerDetails.Name, providerDetails.ApiBaseUrl, providerDetails.Status, providerDetails.ProviderType)
	if err != nil {
		return err
	}
	err = s.dbMongo.InsertProvider(providerID, providerDetails, address, providerDetails.ProviderType)
	if err != nil {
		return err
	}
	return nil
}

type RegistrationStatus struct {
	Status string `json:"Status" example:"Provider Created Successfully" description:"Returned Status"`
}

// @Title Register Pegin Provider
// @Description Registers New Pegin Provider
// @Param  RegisterRequest  body types.ProviderRegisterRequest true "Provider Register Request"
// @Success  200 object RegistrationStatus
// @Route /provider/pegin/register [post]
func (s *Server) registerPeginProviderHandler(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)
	payload := types.ProviderRegisterRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Errorf(UnableToDeserializePayloadError, err)
		http.Error(w, UnableToDeserializePayloadError, http.StatusBadRequest)
		return
	}
	lp, err := pegin.NewLocalProvider(s.ProviderConfig, s.ProviderRespository, s.AccountProvider, s.cfgData.RSK.ChainId)
	if err != nil {
		log.Error(ErrorCreatingLocalProvider, err)
		http.Error(w, ErrorCreatingLocalProvider, http.StatusBadRequest)
		return
	}
	payload.ProviderType = "pegin"
	err = s.AddProvider(lp, nil, payload)
	if err != nil {
		log.Errorf(ErrorAddingProvider, err)
		http.Error(w, ErrorAddingProvider, http.StatusBadRequest)
		return
	}
	response := RegistrationStatus{Status: "Pegin Provider Created Successfully"}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&response)
	if err != nil {
		http.Error(w, UnableToBuildResponse, http.StatusInternalServerError)
		return
	}
}

// @Title Register Pegout Provider
// @Description Registers New Pegout Provider
// @Param  RegisterRequest  body types.ProviderRegisterRequest true "Provider Register Request"
// @Success  200 object RegistrationStatus
// @Route /provider/pegout/register [post]
func (s *Server) registerPegoutProviderHandler(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)
	payload := types.ProviderRegisterRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Errorf(UnableToDeserializePayloadError, err)
		http.Error(w, UnableToDeserializePayloadError, http.StatusBadRequest)
		return
	}
	lp, err := pegout.NewLocalProvider(&s.PegoutConfig, s.ProviderRespository, s.AccountProvider, s.cfgData.RSK.ChainId)
	if err != nil {
		log.Error(ErrorCreatingLocalProvider, err)
		http.Error(w, ErrorCreatingLocalProvider, http.StatusBadRequest)
		return
	}
	payload.ProviderType = "pegout"
	err = s.AddProvider(nil, lp, payload)
	if err != nil {
		log.Errorf(ErrorAddingProvider, err)
		http.Error(w, ErrorAddingProvider, http.StatusBadRequest)
		return
	}

	response := RegistrationStatus{Status: "Pegout Provider Created Successfully"}
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
	if s.provider.Address() == providerAddress.Provider {
		lp = s.provider
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
	r.Use(s.corsMiddleware)
	r.Path("/health").Methods(http.MethodGet).HandlerFunc(s.checkHealthHandler)
	r.Path("/getProviders").Methods(http.MethodGet).HandlerFunc(s.getProvidersHandler)
	r.Path("/pegin/getQuote").Methods(http.MethodPost).HandlerFunc(s.getQuoteHandler)
	r.Path("/pegin/acceptQuote").Methods(http.MethodPost).Handler(s.captchaMiddleware(http.HandlerFunc(s.acceptQuoteHandler)))
	r.Path("/pegout/getQuotes").Methods(http.MethodPost).HandlerFunc(s.getPegoutQuoteHandler)
	r.Path("/pegout/acceptQuote").Methods(http.MethodPost).Handler(s.captchaMiddleware(http.HandlerFunc(s.acceptQuotePegOutHandler)))
	r.Path("/collateral").Methods(http.MethodGet).HandlerFunc(s.getCollateralHandler)
	r.Path("/addCollateral").Methods(http.MethodPost).HandlerFunc(s.addCollateral)
	r.Path("/withdrawCollateral").Methods(http.MethodPost).HandlerFunc(s.withdrawCollateral)
	r.Path("/providers/pegin/register").Methods(http.MethodPost).HandlerFunc(s.registerPeginProviderHandler)
	r.Path("/providers/pegout/register").Methods(http.MethodPost).HandlerFunc(s.registerPegoutProviderHandler)
	r.Path("/providers/changeStatus").Methods(http.MethodPost).HandlerFunc(s.changeStatusHandler)
	r.Path("/providers/resignation").Methods(http.MethodPost).HandlerFunc(s.providerResignHandler)
	r.Path("/providers/sync").Methods(http.MethodPost).HandlerFunc(s.providerSyncHandler)
	r.Path("/userQuotes").Methods(http.MethodGet).HandlerFunc(s.getUserQuotesHandler)
	r.Path("/providers/details").Methods(http.MethodGet).HandlerFunc(s.providerDetailHandler)

	r.Methods("OPTIONS").HandlerFunc(s.handleOptions)
	w := log.StandardLogger().WriterLevel(log.DebugLevel)
	h := handlers.LoggingHandler(w, r)
	defer func(w *io.PipeWriter) {
		_ = w.Close()
	}(w)

	if err := s.initDepositsCache(); err != nil {
		return err
	}
	err := s.initPeginWatchers()
	if err != nil {
		return err
	}

	provider := s.pegoutProvider
	s.pegOutDepositWatcher = NewDepositEventWatcher(PegoutDepositCheckInterval, provider, &s.addWatcherMu, &s.sharedPegoutMutex, make(chan bool), s.rsk, s.btc, s.dbMongo,
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

	peginProvider := s.provider
	pegoutProvider := s.pegoutProvider
	s.lpFundsEventtWatcher = NewLpFundsEventWatcher(1*time.Minute, make(chan bool), s.rsk, peginProvider, pegoutProvider, s.awsConfig)
	s.lpFundsEventtWatcher.Init()

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
	w.WriteHeader(http.StatusOK)
}

func (s *Server) initDepositsCache() error {
	height, err := s.rsk.GetRskHeight()
	if err != nil {
		return err
	}
	events, err := s.rsk.GetDepositEvents(s.cfgData.QuoteCacheStartBlock, height)
	if err != nil {
		return err
	}

	return s.dbMongo.UpsertDepositEvents(events)
}

func (s *Server) initPegoutWatchers() error {
	quoteStatesToWatch := []types.RQState{types.RQStateCallForUserSucceeded, types.RQStateWaitingForDeposit, types.RQStateWaitingForDepositConfirmations}
	quotes, err := s.dbMongo.GetRetainedPegOutQuoteByState(quoteStatesToWatch)
	if err != nil {
		return err
	}
	waitingForDepositQuotes := make(map[string]*WatchedQuote, 0)
	waitingForConfirmationQuotes := make(map[string]*WatchedQuote, 0)
	for _, entry := range quotes {
		quote, err := s.dbMongo.GetPegOutQuote(entry.QuoteHash)
		if err != nil || quote == nil {
			log.Errorf("initPegoutWatchers: quote not found for hash: %s. Watcher not initialized for address %s", entry.QuoteHash, entry.DepositAddr)
			continue
		}

		p := pegout.GetPegoutProviderByAddress(s.pegoutProvider, quote.LPRSKAddr)
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
			waitingForConfirmationQuotes[entry.QuoteHash] = &WatchedQuote{
				Signature:          entry.Signature,
				Data:               quote,
				DepositTransaction: entry.DepositTransaction,
				QuoteHash:          entry.QuoteHash,
			}
		} else {
			waitingForDepositQuotes[entry.QuoteHash] = &WatchedQuote{
				Signature: entry.Signature,
				Data:      quote,
				QuoteHash: entry.QuoteHash,
			}
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

		p := pegin.GetPeginProviderByAddress(s.provider, quote.LPRSKAddr)
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
		hash:         hash,
		btc:          s.btc,
		rsk:          s.rsk,
		lp:           provider,
		dbMongo:      s.dbMongo,
		quote:        quote,
		state:        state,
		signature:    signB,
		done:         make(chan struct{}),
		sharedLocker: &s.sharedPegoutMutex,
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
	w.Header().Set("Content-Type", "application/json")

	providerList, error := s.dbMongo.GetProviders()
	if error != nil {
		log.Error("Error fetching providers. Error: ", error)
		customError := NewServerError(ErrorFetchingMongoDBProviders+error.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}
	var ids []int64
	for _, address := range providerList {
		ids = append(ids, address.Id)
	}
	providers, error := s.rsk.GetProviders(ids)

	if error != nil {
		log.Error("GetProviders - error encoding response: ", error)
		customError := NewServerError("GetProviders - error encoding response: "+error.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	response := make([]*types.GlobalProvider, 0)
	for _, provider := range providers {
		response = append(response, toGlobalProvider(&provider))
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

	err = s.validateAmountForProvider(new(big.Int).SetUint64(qr.ValueToTransfer), &s.ProviderConfig)
	if err != nil {
		log.Error(err)
		customError := NewServerError(err.Error(), Details{}, true)
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
	pq, err := s.provider.GetQuote(q, gas, types.NewBigWei(price))
	if err != nil {
		log.Error("error getting quote: ", err)
		getQuoteFailed = true
	}
	if pq != nil {
		if new(types.Wei).Add(pq.Value, pq.CallFee).Cmp(minLockTxValueInWei) < 0 {
			log.Error("error getting quote; requested amount below bridge's min pegin tx value: ", qr.ValueToTransfer)
			amountBelowMinLockTxValue = true
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

	err = s.validateAmountForProvider(new(big.Int).SetUint64(qr.ValueToTransfer), &s.PegoutConfig.ProviderConfig)
	if err != nil {
		log.Error(err)
		customError := NewServerError(err.Error(), Details{}, true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	amountInSatoshi, _ := types.NewUWei(qr.ValueToTransfer).ToSatoshi().Uint64()
	feeInSatoshi, err := s.btc.EstimateFees(qr.To, amountInSatoshi)
	if err != nil && strings.Contains(err.Error(), "Insufficient funds") {
		log.Error(ErrorNotLiquidity)
		customError := NewServerError(ErrorNotLiquidity, make(Details), true)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	} else if err != nil {
		log.Error(err.Error())
		customError := NewServerError(err.Error(), make(Details), false)
		ResponseError(w, customError, http.StatusInternalServerError)
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
	q := parseReqToPegOutQuote(qr, s.rsk.GetLBCAddress())
	rskBlockNumber, err := s.rsk.GetRskHeight()
	if err != nil {
		log.Error("Error getting last block", err.Error())
		customError := NewServerError("Error getting last block", make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusInternalServerError)
		return
	}
	pq, err := s.pegoutProvider.GetQuote(q, rskBlockNumber, types.SatoshiToWei(feeInSatoshi))
	if err != nil {
		log.Error("error getting quote: ", err)
		getQuoteFailed = true
	}
	if pq != nil {
		if new(types.Wei).Add(pq.Value, pq.CallFee).Cmp(minLockTxValueInWei) < 0 {
			log.Error("error getting quote; requested amount below bridge's min pegin tx value: ", qr.ValueToTransfer)
			amountBelowMinLockTxValue = true
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
	returnQuoteSignFunc := func(w http.ResponseWriter, signature, depositAddr, flyoverRedeemScript string) {
		enc := json.NewEncoder(w)
		response := acceptRes{
			Signature:                 signature,
			BitcoinDepositAddressHash: depositAddr,
			FlyoverRedeemScript:       flyoverRedeemScript,
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
		returnQuoteSignFunc(w, rq.Signature, rq.DepositAddr, rq.FlyoverRedeemScript)
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

	depositAddress, flyoverRedeemScript, err := s.rsk.GetDerivedBitcoinAddress(fedInfo, s.btc.GetParams(), btcRefAddr, lbcAddr, lpBTCAddr, hashBytes)
	if err != nil {
		log.Error("error getting derived bitcoin address: ", err.Error())
		customError := NewServerError("error getting derived bitcoin address: "+err.Error(), make(map[string]interface{}), true)
		ResponseError(w, customError, http.StatusBadRequest)
		return
	}

	p := pegin.GetPeginProviderByAddress(s.provider, quote.LPRSKAddr)
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
	signB, err := p.SignQuote(hashBytes, depositAddress, flyoverRedeemScript, reqLiq)
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
	returnQuoteSignFunc(w, signature, depositAddress, flyoverRedeemScript)
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

func parseReqToPegOutQuote(qr QuotePegOutRequest, lbcAddr string) *pegout.Quote {
	return &pegout.Quote{
		LBCAddr:       lbcAddr,
		BtcRefundAddr: qr.BitcoinRefundAddress,
		RSKRefundAddr: qr.RskRefundAddress,
		DepositAddr:   qr.To,
		Value:         types.NewWei(int64(qr.ValueToTransfer)),
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

	expTime := quote.GetExpirationTime()
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

	p := pegout.GetPegoutProviderByAddress(s.pegoutProvider, quote.LPRSKAddr)

	reqLiq := quote.CallCost.Uint64() + quote.Value.Uint64()
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
		Signature:  signature,
		LbcAddress: depositAddr,
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

type SenBTCResponse struct {
	TxHash string `json:"txHash" example:"0x0" description:"TxHash of the BTC transaction sent to the address"`
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

	lp := pegin.GetPeginProviderByAddress(s.provider, payload.LpRskAddress)
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
	payload := WithdrawCollateralRequest{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}
	if isValid := Validate(payload)(w); !isValid {
		return
	}

	opts, err := pegin.GetPeginProviderTransactOpts(s.provider, payload.LpRskAddress)
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
	payload := ProviderResignRequest{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		buildErrorDecodingRequest(w, err)
		return
	}
	if isValid := Validate(payload)(w); !isValid {
		return
	}

	opts, err := pegin.GetPeginProviderTransactOpts(s.provider, payload.LpRskAddress)
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

// @Title GetUserQuotes
// @Description Returns user quotes for address.
// @Param   UserQuoteRequest query types.UserQuoteRequest true "User Quote Request Details"
// @Success 200 {array} pegout.DepositEvent "Successfully retrieved the user quotes"
// @Router /userQuotes [get]
func (s *Server) getUserQuotesHandler(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address parameter is required", http.StatusBadRequest)
		return
	}

	events, err := s.dbMongo.GetDepositEvents(address)
	if err != nil {
		log.Error("error getting user quotes: ", err.Error())
	}
	if events == nil {
		events = []*pegout.DepositEvent{}
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(&events)
	if err != nil {
		log.Error("error encoding user events")
		return
	}
}

type ProviderDetail struct {
	Fee                   uint64 `json:"fee"  required:""`
	MinTransactionValue   uint64 `json:"minTransactionValue"  required:""`
	MaxTransactionValue   uint64 `json:"maxTransactionValue"  required:""`
	RequiredConfirmations uint16 `json:"requiredConfirmations"  required:""`
}

type ProviderDetailResponse struct {
	SiteKey string         `json:"siteKey" required:""`
	Pegin   ProviderDetail `json:"pegin" required:""`
	Pegout  ProviderDetail `json:"pegout" required:""`
}

// @Title Provider detail
// @Description Returns the details of the provider that manages this instance of LPS
// @Param   UserQuoteRequest query types.UserQuoteRequest true "User Quote Request Details"
// @Success 200 object ProviderDetailResponse "Detail of the provider that manges this instance"
// @Router /providers/details [get]
func (s *Server) providerDetailHandler(w http.ResponseWriter, r *http.Request) {
	toRestAPI(w)

	detail := ProviderDetailResponse{
		SiteKey: s.cfgData.CaptchaSiteKey,
		Pegin: ProviderDetail{
			Fee:                   s.ProviderConfig.Fee.Uint64(),
			MinTransactionValue:   s.ProviderConfig.MinTransactionValue.Uint64(),
			MaxTransactionValue:   s.ProviderConfig.MaxTransactionValue.Uint64(),
			RequiredConfirmations: s.ProviderConfig.MaxConf,
		},
		Pegout: ProviderDetail{
			Fee:                   s.PegoutConfig.Fee.Uint64(),
			MinTransactionValue:   s.PegoutConfig.MinTransactionValue.Uint64(),
			MaxTransactionValue:   s.PegoutConfig.MaxTransactionValue.Uint64(),
			RequiredConfirmations: s.PegoutConfig.MaxConf,
		},
	}

	err := json.NewEncoder(w).Encode(&detail)
	if err != nil {
		log.Error("error encoding user events")
		return
	}
}

// @Title Provider Synchronization
// @Description Synchronizes providers with MongoDB
// @Route /provider/sync [post]
// @Success 204 object
func (s *Server) providerSyncHandler(w http.ResponseWriter, r *http.Request) {
	providerIds, err := s.rsk.GetProviderIds()
	if err != nil {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusNotFound)
		return
	}
	providersIdList, err := createArrayFromOneToN(providerIds)
	if err != nil {
		customError := NewServerError(err.Error(), make(Details), true)
		ResponseError(w, customError, http.StatusNotFound)
		return
	}
	providers, err := s.rsk.GetProviders(providersIdList)
	var providerDTOs []*types.GlobalProvider
	for _, provider := range providers {
		providerDTOs = append(providerDTOs, toGlobalProvider(&provider))
	}
	filteredProviders := filterProvidersByAddress(s.cfgData.RSK.LpsAddress, providerDTOs)
	if err != nil {
		http.Error(w, UnableToBuildResponse, http.StatusInternalServerError)
		return
	}
	err = s.dbMongo.ResetProviders(filteredProviders)
	response := filteredProviders
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&response)
}

type CaptchaValidationResponse struct {
	Success     bool      `json:"success"`
	Score       *float32  `json:"score"`
	Action      *string   `json:"action"`
	ChallengeTs time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func (s *Server) captchaMiddleware(next http.Handler) http.Handler {
	if s.cfgData.CaptchaThreshold < 0.5 {
		log.Warn("Too low captcha threshold value!")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Captcha-Token")
		disabled, _ := strconv.ParseBool(os.Getenv("DISABLE_CAPTCHA"))
		if disabled {
			log.Warning("IMPORTANT! Handling request with captcha validation disabled")
			next.ServeHTTP(w, r)
			return
		} else if token == "" {
			customError := NewServerError("missing X-Captcha-Token header", make(Details), true)
			ResponseError(w, customError, http.StatusBadRequest)
			return
		}

		form := make(url.Values)
		form.Set("secret", s.cfgData.CaptchaSecretKey)
		form.Set("response", token)
		res, err := http.DefaultClient.PostForm("https://www.google.com/recaptcha/api/siteverify", form)

		if err != nil {
			details := make(Details)
			details["error"] = err.Error()
			customError := NewServerError("error validating captcha", details, false)
			ResponseError(w, customError, http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		var validation CaptchaValidationResponse
		err = json.NewDecoder(res.Body).Decode(&validation)
		if err != nil {
			customError := NewServerError("error validating captcha", make(Details), false)
			ResponseError(w, customError, http.StatusInternalServerError)
			return
		}

		validCaptcha := validation.Success
		if validation.Score != nil { // if is v3 we also use the score
			validCaptcha = validCaptcha && *validation.Score >= s.cfgData.CaptchaThreshold
		}

		if validCaptcha {
			log.Debugf("Valid captcha solved on %s\n", validation.Hostname)
			next.ServeHTTP(w, r)
		} else {
			details := make(Details)
			details["errors"] = validation.ErrorCodes
			customError := NewServerError("error validating captcha", details, true)
			ResponseError(w, customError, http.StatusBadRequest)
		}
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Add("Access-Control-Allow-Origin", "*")
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
		headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, X-Captcha-Token")
		headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) validateAmountForProvider(amount *big.Int, provider *pegin.ProviderConfig) error {
	var min, max = provider.MinTransactionValue, provider.MaxTransactionValue
	if amount.Cmp(min) < 0 || amount.Cmp(max) > 0 {
		return fmt.Errorf("amount out of provider range which is (%d, %d)", min, max)
	}
	return nil
}

func filterProvidersByAddress(address string, providers []*types.GlobalProvider) []*types.GlobalProvider {
	filteredProviders := make([]*types.GlobalProvider, 0)
	lowercaseAddress := strings.ToLower(address)

	for _, provider := range providers {
		if strings.ToLower(provider.Provider) == lowercaseAddress {
			filteredProviders = append(filteredProviders, provider)
		}
	}

	return filteredProviders
}
func createArrayFromOneToN(providerIds *big.Int) ([]int64, error) {
	n := providerIds.Int64()
	if n < 1 {
		return nil, fmt.Errorf("The input number should be greater than 0")
	}

	array := make([]int64, n)
	for i := int64(1); i <= n; i++ {
		array[i-1] = i
	}

	return array, nil
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
