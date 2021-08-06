package http

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"context"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/connectors"
	"github.com/rsksmart/liquidity-provider-server/storage"
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	srv       http.Server
	providers []providers.LiquidityProvider
	rsk       *connectors.RSK
	db        *storage.DB
}

func New(rsk *connectors.RSK, db *storage.DB) Server {
	provs := []providers.LiquidityProvider{}
	return Server{rsk: rsk, db: db, providers: provs}
}

func (s *Server) AddProvider(lp providers.LiquidityProvider) {
	s.providers = []providers.LiquidityProvider{lp}
}

func (s *Server) Start(port uint) error {
	r := mux.NewRouter()
	r.Path("/getQuote").Methods(http.MethodPost).HandlerFunc(s.getQuoteHandler)
	r.Path("/acceptQuote").Methods(http.MethodPost).HandlerFunc(s.acceptQuoteHandler)
	w := log.StandardLogger().WriterLevel(log.DebugLevel)
	h := handlers.LoggingHandler(w, r)
	defer w.Close()

	s.srv = http.Server{
		Addr:    ":" + fmt.Sprint(port),
		Handler: h,
	}
	log.Info("starting server at localhost", s.srv.Addr)

	err := s.srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	log.Info("started server at localhost", s.srv.Addr)
	return nil
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

func (s *Server) getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	q := types.Quote{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Debug("received quote request: ", fmt.Sprintf("%+v", q))

	gas, err := s.rsk.EstimateGas(q.ContractAddr, q.Value, []byte(q.Data))
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

	quotes := []*types.Quote{}
	// TODO: fill in LBC and Fed address with existing info and prevent receiving it from the request payload

	for _, p := range s.providers {
		pq := p.GetQuote(q, gas, *price)

		// TODO: validate that the received quote matches the expected params

		if pq != nil {
			err = s.storeQuote(pq)

			if err != nil {
				log.Error(err)
			} else {
				quotes = append(quotes, pq)
			}
		}
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&quotes)
	if err != nil {
		log.Error("error encoding quote list: ", err.Error())
		http.Error(w, "error processing quotes", http.StatusInternalServerError)
		return
	}
}

func (s *Server) acceptQuoteHandler(w http.ResponseWriter, r *http.Request) {
	type acceptReq struct {
		QuoteHash string
	}

	type acceptRes struct {
		Signature                 string `json:"signature"`
		BitcoinDepositAddressHash string `json:"bitcoinDepositAddressHash"`
	}
	req := acceptReq{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := s.providers[0]

	response := acceptRes{}
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

	signature, err := p.SignHash(hashBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.Signature = hex.EncodeToString(signature)

	derivedFedAddress, err := getDerivedBitcoinAddressHash(
		[]byte(quote.BTCRefundAddr), []byte(quote.LBCAddr), []byte(quote.LPBTCAddr), []byte(req.QuoteHash))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.BitcoinDepositAddressHash = hex.EncodeToString(derivedFedAddress)

	enc := json.NewEncoder(w)
	err = enc.Encode(response)

	// TODO: ensure that the quote is not processed if there is any kind of error in the communication with the client
	if err != nil {
		log.Error("error encoding response: ", err.Error())
		http.Error(w, "error processing request", http.StatusInternalServerError)
		return
	}
}

func getDerivedBitcoinAddressHash(userBtcRefundAddress []byte, lbcAddress []byte, lpBtcAddress []byte, derivationArgumentsHash []byte) ([]byte, error) {
	var resultData []byte
	resultData = append(resultData, derivationArgumentsHash...)
	resultData = append(resultData, userBtcRefundAddress...)
	resultData = append(resultData, lpBtcAddress...)
	resultData = append(resultData, lbcAddress...)

	derivationValueHash := crypto.Keccak256(resultData)
	btcAddressHash := deriveFastBridgeFederationAddress(derivationValueHash)
	return btcAddressHash, nil
}

func deriveFastBridgeFederationAddress(derivationValue []byte) []byte {
	return derivationValue // TODO: implement derivation
}

func (s *Server) storeQuote(q *types.Quote) error {
	h, err := s.rsk.HashQuote(q)
	if err != nil {
		return fmt.Errorf("error hashing quote: %v", err)
	}

	err = s.db.InsertQuote(h, q)
	if err != nil {
		log.Fatal("error inserting quote: %v", err)
	}
	return nil
}
