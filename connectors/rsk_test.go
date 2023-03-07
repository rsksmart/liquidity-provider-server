package connectors

import (
	"context"
	crand "crypto/rand"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/connectors/testmocks"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/pegin"
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/assert"
)

var contextMatcher = mock.MatchedBy(func(ctx context.Context) bool { return true })

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

var quotes = []*pegin.Quote{
	{
		FedBTCAddr:         "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		LBCAddr:            "2ff74F841b95E000625b3A77fed03714874C4fEa",
		LPRSKAddr:          "0x00d80aA033fb51F191563B08Dc035fA128e942C5",
		BTCRefundAddr:      "mnxKdPFrYqLSUy2oP1eno8n5X8AwkcnPjk",
		RSKRefundAddr:      "0x5F3b836CA64DA03e613887B46f71D168FC8B5Bdf",
		LPBTCAddr:          "2NDjJznHgtH1rzq63eeFG3SiDi5wxE25FSz",
		CallFee:            types.NewWei(250),
		PenaltyFee:         types.NewWei(5000),
		ContractAddr:       "0x87136cf829edaF7c46Eb943063369a1C8D4f9085",
		Data:               "",
		GasLimit:           6000000,
		Nonce:              int64(rand.Int()),
		Value:              types.NewWei(250),
		AgreementTimestamp: 0,
		TimeForDeposit:     3600,
		LpCallTime:         3600,
		Confirmations:      10,
	},
}

func testNewRSKWithInvalidAddresses(t *testing.T) {

	for _, tt := range invalidAddresses {
		res, err := NewRSK(tt.input, tt.input, 10, 0, nil)

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
		res, err := NewRSK(tt.input, tt.input, 10, 0, nil)
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
		rsk, err := NewRSK(tt.input, tt.input, 10, 0, nil)
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

func testCopyBtcAddress(t *testing.T) {
	err := copyBtcAddr("1PRTTaJesdNovgne6Ehcdu1fpEdX7913CK", []byte{})
	assert.Empty(t, err)
}

func testCopyBtcAddressWithAnInvalidAddress(t *testing.T) {
	err := copyBtcAddr("0x895E7D15510C3f77726422669B0Ef768d26F7FEb", []byte{})
	assert.Equal(t, "invalid format: version and/or checksum bytes missing", err.Error())
}

func testIsEOA(t *testing.T) {
	validAddress := "0xD2244D24FDE5353e4b3ba3b6e05821b456e04d95"

	eOAAddress := "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"
	sCAddress := "0xd9145CCE52D386f254917e481eB44e9943F39138"
	invalidAddress := "sfdasfasfaf"

	fakeBytecode := make([]byte, 4)
	crand.Read(fakeBytecode)

	rskClientMock := new(testmocks.RSKClientMock)
	rsk, err := NewRSK(validAddress, validAddress, 10, 0, nil)
	rsk.c = rskClientMock

	if err != nil {
		t.Errorf("Unexpected error creating RSK Client: %v", err)
	}

	rskClientMock.On("CodeAt", contextMatcher, common.HexToAddress(eOAAddress), mock.AnythingOfType("*big.Int")).Return([]byte{}, nil)
	rskClientMock.On("CodeAt", contextMatcher, common.HexToAddress(sCAddress), mock.AnythingOfType("*big.Int")).Return(fakeBytecode, nil)

	testCases := []*struct {
		caseName   string
		address    string
		assertions func(result bool, errorResult error)
	}{
		{
			caseName: "Validates that address is EOA",
			address:  eOAAddress,
			assertions: func(result bool, errorResult error) {
				rskClientMock.AssertNumberOfCalls(t, "CodeAt", 1)
				assert.Nil(t, errorResult)
				assert.True(t, result)
			},
		},
		{
			caseName: "Validates that address is SC",
			address:  sCAddress,
			assertions: func(result bool, errorResult error) {
				rskClientMock.AssertNumberOfCalls(t, "CodeAt", 1)
				assert.Nil(t, errorResult)
				assert.False(t, result)
			},
		},
		{
			caseName: "Returns error on invalid address",
			address:  invalidAddress,
			assertions: func(result bool, errorResult error) {
				rskClientMock.AssertNumberOfCalls(t, "CodeAt", 0)
				assert.Error(t, errorResult, "invalid address")
				assert.False(t, result)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.caseName, func(t *testing.T) {
			result, errorResult := rsk.IsEOA(test.address)
			test.assertions(result, errorResult)
			rskClientMock.Calls = []mock.Call{}
		})
	}

}

func TestRSKCreate(t *testing.T) {
	t.Run("new invalid", testNewRSKWithInvalidAddresses)
	t.Run("new valid", testNewRSKWithValidAddresses)
	t.Run("parse quote", testParseQuote)
	t.Run("test copy btc address", testCopyBtcAddress)
	t.Run("test copy btc address with an invalid address", testCopyBtcAddressWithAnInvalidAddress)
	t.Run("test EOA address validation", testIsEOA)
}
