package connectors

import (
	"os"
	"testing"
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

func testNewRSKWithInvalidAddresses(t *testing.T) {
	abiFile, err := os.Open("abi_test.json")
	if err != nil {
		t.Errorf("Unexpected error while opening abi mock file %v: %v", "abi_test.json", err)
	}
	for _, tt := range invalidAddresses {
		res, err := NewRSK(tt.input, abiFile, nil)

		if res != nil {
			t.Errorf("Unexpected value for input %v: %v", tt.input, res)
		}
		if err == nil {
			t.Errorf("Unexpected success for input %v: %v", tt.input, err)
		}
		if err.Error() != "invalid contract address" {
			t.Errorf("Unexpected error for input %v: %v", tt.input, err)
		}
	}
}

func testNewRSKWithValidAddresses(t *testing.T) {
	abiFile, err := os.Open("abi_test.json")
	if err != nil {
		t.Errorf("Unexpected error while opening abi mock file %v: %v", "abi_test.json", err)
	}
	for _, tt := range validTests {
		res, err := NewRSK(tt.input, abiFile, nil)
		if err != nil {
			t.Errorf("Unexpected error for input %v: %v", tt.input, err)
		}
		if res == nil {
			t.Errorf("Unexpected error for input %v: %v", tt.input, err)
		}
	}
}

func TestRSKCreate(t *testing.T) {
	t.Run("new", testNewRSKWithInvalidAddresses)
	t.Run("new", testNewRSKWithValidAddresses)
}
