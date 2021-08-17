package federation

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

var testQuotes = []struct {
	BTCRefundAddr string
	LBCAddr string
	LPBTCAddr string
	QuoteHash string
	Expected string
}{
	{
			LPBTCAddr: "746231713634706b667a306368773065753936776c6a677663683779357633397270643465786a337276",
			LBCAddr:"0xD2244D24FDE5353e4b3ba3b6e05821b456e04d95",
			BTCRefundAddr: "74623171336b3463726832367936367335747575617333386733347775756a7074793466647579726e76",
			QuoteHash: "4a3eca107f22707e5dbc79964f3e6c21ec5e354e0903391245d9fdbe6bd2b2f0",
			Expected: "bb0afe83ab039e7ac51c581f0c9909a2c3c02e956c9c0f975baf5baefb39a715",
	},
}

func TestGetDerivationValueHash(t *testing.T) {
	for _, tt := range testQuotes {
		btcRefAddr, err := hex.DecodeString(tt.BTCRefundAddr)
		if err != nil || len(btcRefAddr) == 0 {
			t.Errorf("Cannot parse BTCRefundAddr correctly. value: %v, error: %v", tt.BTCRefundAddr, err)
		}

		if !common.IsHexAddress(tt.LBCAddr) {
			t.Errorf("invalid address: %v", tt.LBCAddr)
		}

		lbcAddr := common.FromHex(tt.LBCAddr)
		if err != nil || len(lbcAddr) == 0 {
			t.Errorf("Cannot parse LBCAddr correctly. value: %v, error: %v", tt.LBCAddr, err)
		}

		lpBTCAdrr, err := hex.DecodeString(tt.LPBTCAddr)
		if err != nil || len(lpBTCAdrr) == 0 {
			t.Errorf("Cannot parse LPBTCAddr correctly. value: %v, error: %v", tt.LPBTCAddr, err)
		}
		hashBytes, err := hex.DecodeString(tt.QuoteHash)
		if err != nil || len(hashBytes) == 0 {
			t.Errorf("Cannot parse QuoteHash correctly. value: %v, error: %v", tt.QuoteHash, err)
		}

		value, err := GetDerivationValueHash(btcRefAddr, lbcAddr, lpBTCAdrr, hashBytes)

		if err != nil {
			t.Errorf("Unexpected error for quotehash %v: %v", tt.QuoteHash, err)
		}

		result := hex.EncodeToString(value)
		if result != tt.Expected {
			t.Errorf("Unexpected derivation value. value: %v, expected: %v, error: %v", result, tt.Expected, err)
		}
	}
}
func TestRSKCreate(t *testing.T) {
	t.Run("new", TestGetDerivationValueHash)
}