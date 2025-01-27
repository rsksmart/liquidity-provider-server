package liquidity_provider_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"testing"
)

func TestProviderType_IsValid(t *testing.T) {
	cases := test.Table[liquidity_provider.ProviderType, bool]{
		{Value: "pegin", Result: true},
		{Value: "pegout", Result: true},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: test.AnyString, Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.IsValid()
	})
}

func TestProviderType_AcceptsPegin(t *testing.T) {
	cases := test.Table[liquidity_provider.ProviderType, bool]{
		{Value: "pegin", Result: true},
		{Value: "pegout", Result: false},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: test.AnyString, Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.AcceptsPegin()
	})
}

func TestProviderType_AcceptsPegout(t *testing.T) {
	cases := test.Table[liquidity_provider.ProviderType, bool]{
		{Value: "pegin", Result: false},
		{Value: "pegout", Result: true},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: test.AnyString, Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.AcceptsPegout()
	})
}

func TestToProviderType(t *testing.T) {
	type testResult struct {
		Result liquidity_provider.ProviderType
		Error  error
	}

	cases := test.Table[string, testResult]{
		{Value: "pegin", Result: testResult{Result: liquidity_provider.PeginProvider, Error: nil}},
		{Value: "pegout", Result: testResult{Result: liquidity_provider.PegoutProvider, Error: nil}},
		{Value: "both", Result: testResult{Result: liquidity_provider.FullProvider, Error: nil}},
		{Value: "", Result: testResult{Result: "", Error: liquidity_provider.InvalidProviderTypeError}},
		{Value: test.AnyString, Result: testResult{Result: "", Error: liquidity_provider.InvalidProviderTypeError}},
	}

	test.RunTable(t, cases, func(value string) testResult {
		result, err := liquidity_provider.ToProviderType(value)
		return testResult{
			Result: result,
			Error:  err,
		}
	})
}
