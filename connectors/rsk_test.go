package connectors

import (
	"math/rand"
	"testing"

	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/assert"
)

var validTests = []struct {
	input string
}{
	{"0xD2244D24FDE5353e4b3ba3b6e05821b456e04d95"},
}

var invalidAddresses = []struct {
	input    string
	expected string
}{
	{"123", "invalid contract address"},
	{"b3ba3b6e05821b456e04d95", "invalid contract address"},
}

var quotes = []*types.Quote{
	{
		FedBTCAddr:         "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		LBCAddr:            "2ff74F841b95E000625b3A77fed03714874C4fEa",
		LPRSKAddr:          "0x00d80aA033fb51F191563B08Dc035fA128e942C5",
		BTCRefundAddr:      "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		RSKRefundAddr:      "0x5F3b836CA64DA03e613887B46f71D168FC8B5Bdf",
		LPBTCAddr:          "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		CallFee:            250,
		PenaltyFee:         5000,
		ContractAddr:       "0x87136cf829edaF7c46Eb943063369a1C8D4f9085",
		Data:               "",
		GasLimit:           6000000,
		Nonce:              int64(rand.Int()),
		Value:              250,
		AgreementTimestamp: 0,
		TimeForDeposit:     3600,
		CallTime:           3600,
		Confirmations:      10,
	},
}

func testNewRSKWithInvalidAddresses(t *testing.T) {

	for _, tt := range invalidAddresses {
		res, err := NewRSK(tt.input, tt.input, 10)

		if res != nil {
			t.Errorf("Unexpected value for input %v: %v", tt.input, res)
		}
		if err == nil {
			t.Errorf("Unexpected success for input %v: %v", tt.input, err)
		}
		if err.Error() != "invalid contract address" && err.Error() != "invalid LBC contract address" {
			t.Errorf("Unexpected error for input %v: %v", tt.input, err)
		}
	}
}

func testNewRSKWithValidAddresses(t *testing.T) {
	for _, tt := range validTests {
		res, err := NewRSK(tt.input, tt.input, 10)
		if err != nil {
			t.Errorf("Unexpected error for input %v: %v", tt.input, err)
		}
		if res == nil {
			t.Errorf("Unexpected error for input %v: %v", tt.input, err)
		}
	}
}

func testParseQuote(t *testing.T) {
	for _, tt := range validTests {
		rsk, err := NewRSK(tt.input, tt.input, 10)
		if err != nil {
			t.Errorf("Unexpected error for input %v: %v", tt.input, err)
		}
		for _, quote := range quotes {
			result, err := rsk.ParseQuote(quote)
			if err != nil {
				t.Errorf("Unexpected error parsing quote %v: %v", quote, err)
			}

			assert.EqualValues(t, len(result.BtcRefundAddress), 21, "BtcRefundAddress should not be empty")
			assert.EqualValues(t, len(result.LiquidityProviderBtcAddress), 21, "LiquidityProviderBtcAddress should not be empty")
			assert.NotEmpty(t, len(result.FedBtcAddress), 20, "FedBtcAddress should not be empty")
		}
	}
}

func TestRSKCreate(t *testing.T) {
	t.Run("new", testNewRSKWithInvalidAddresses)
	t.Run("new", testNewRSKWithValidAddresses)
	t.Run("parse quote", testParseQuote)
}
