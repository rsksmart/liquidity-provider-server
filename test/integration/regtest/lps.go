package regtest

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/stretchr/testify/require"
)

// Lps calls the public REST API of the LPS, the way a client of the v2.6 quote flow does.
type Lps struct {
	url string
}

// AcceptResult keeps the status and body of a rejected accept, which tests assert on.
type AcceptResult struct {
	StatusCode int
	Body       string
	Accepted   pkg.AcceptPeginRespose
}

func NewLps(stack Stack) *Lps {
	return &Lps{url: stack.LpsURL}
}

// AvailablePegInLiquidity is the liquidity the LPS offers to new quotes, after subtracting the
// quotes accepted and still waiting for their deposit.
func (lps *Lps) AvailablePegInLiquidity(t *testing.T) *big.Int {
	t.Helper()
	status, body := lps.request(t, http.MethodGet, "/providers/liquidity", nil)
	require.Equalf(t, http.StatusOK, status, "GET /providers/liquidity: %s", body)
	var liquidity pkg.AvailableLiquidityDTO
	require.NoError(t, json.Unmarshal(body, &liquidity))
	return liquidity.PeginLiquidityAmount
}

// PegInQuote asks for a v2.6 peg-in quote that transfers value to destination.
func (lps *Lps) PegInQuote(t *testing.T, destination common.Address, value *big.Int) pkg.GetPeginQuoteResponse {
	t.Helper()
	status, body := lps.request(t, http.MethodPost, "/pegin/getQuote", pkg.PeginQuoteRequest{
		CallEoaOrContractAddress: destination.Hex(),
		ValueToTransfer:          value,
		RskRefundAddress:         destination.Hex(),
	})
	require.Equalf(t, http.StatusOK, status, "POST /pegin/getQuote: %s", body)
	var quotes []pkg.GetPeginQuoteResponse
	require.NoError(t, json.Unmarshal(body, &quotes))
	require.Len(t, quotes, 1)
	return quotes[0]
}

// AcceptPegInQuote accepts a v2.6 quote, which reserves its liquidity until the deposit arrives
// or the quote expires. A rejection is returned, not failed, so tests can assert on it.
func (lps *Lps) AcceptPegInQuote(t *testing.T, quoteHash string) AcceptResult {
	t.Helper()
	status, body := lps.request(t, http.MethodPost, "/pegin/acceptQuote", pkg.AcceptQuoteRequest{QuoteHash: quoteHash})
	result := AcceptResult{StatusCode: status, Body: string(body)}
	if status == http.StatusOK {
		require.NoError(t, json.Unmarshal(body, &result.Accepted))
	}
	return result
}

func (lps *Lps) request(t *testing.T, method, path string, payload any) (int, []byte) {
	t.Helper()
	var reader io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		require.NoError(t, err)
		reader = bytes.NewReader(encoded)
	}
	request, err := http.NewRequestWithContext(t.Context(), method, lps.url+path, reader)
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	require.NoErrorf(t, err, "the LPS is not reachable at %s", lps.url)
	defer func() { require.NoError(t, response.Body.Close()) }()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return response.StatusCode, body
}
