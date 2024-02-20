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
		{Value: "any value", Result: false},
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
		{Value: "any value", Result: false},
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
		{Value: "any value", Result: false},
	}
	test.RunTable(t, cases, func(value liquidity_provider.ProviderType) bool {
		return value.AcceptsPegout()
	})
}

func TestToProviderType(t *testing.T) {
	var err error
	var result liquidity_provider.ProviderType

	errorCases := test.Table[string, error]{
		{Value: "pegin", Result: nil},
		{Value: "pegout", Result: nil},
		{Value: "both", Result: nil},
		{Value: "", Result: liquidity_provider.InvalidProviderTypeError},
		{Value: "any value", Result: liquidity_provider.InvalidProviderTypeError},
	}

	valueCases := test.Table[string, liquidity_provider.ProviderType]{
		{Value: "pegin", Result: liquidity_provider.PeginProvider},
		{Value: "pegout", Result: liquidity_provider.PegoutProvider},
		{Value: "both", Result: liquidity_provider.FullProvider},
		{Value: "", Result: ""},
		{Value: "any value", Result: ""},
	}

	test.RunTable(t, errorCases, func(value string) error {
		_, err = liquidity_provider.ToProviderType(value)
		return err
	})

	test.RunTable(t, valueCases, func(value string) liquidity_provider.ProviderType {
		result, _ = liquidity_provider.ToProviderType(value)
		return result
	})
}
