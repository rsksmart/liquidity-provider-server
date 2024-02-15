package entities_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/test"
	"testing"
)

func TestProviderType_IsValid(t *testing.T) {
	cases := test.Table[entities.ProviderType, bool]{
		{Value: "pegin", Result: true},
		{Value: "pegout", Result: true},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: "any value", Result: false},
	}
	test.RunTable(t, cases, func(value entities.ProviderType) bool {
		return value.IsValid()
	})
}

func TestProviderType_AcceptsPegin(t *testing.T) {
	cases := test.Table[entities.ProviderType, bool]{
		{Value: "pegin", Result: true},
		{Value: "pegout", Result: false},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: "any value", Result: false},
	}
	test.RunTable(t, cases, func(value entities.ProviderType) bool {
		return value.AcceptsPegin()
	})
}

func TestProviderType_AcceptsPegout(t *testing.T) {
	cases := test.Table[entities.ProviderType, bool]{
		{Value: "pegin", Result: false},
		{Value: "pegout", Result: true},
		{Value: "both", Result: true},
		{Value: "", Result: false},
		{Value: "any value", Result: false},
	}
	test.RunTable(t, cases, func(value entities.ProviderType) bool {
		return value.AcceptsPegout()
	})
}

func TestToProviderType(t *testing.T) {
	var err error
	var result entities.ProviderType

	errorCases := test.Table[string, error]{
		{Value: "pegin", Result: nil},
		{Value: "pegout", Result: nil},
		{Value: "both", Result: nil},
		{Value: "", Result: entities.InvalidProviderTypeError},
		{Value: "any value", Result: entities.InvalidProviderTypeError},
	}

	valueCases := test.Table[string, entities.ProviderType]{
		{Value: "pegin", Result: entities.PeginProvider},
		{Value: "pegout", Result: entities.PegoutProvider},
		{Value: "both", Result: entities.FullProvider},
		{Value: "", Result: ""},
		{Value: "any value", Result: ""},
	}

	test.RunTable(t, errorCases, func(value string) error {
		_, err = entities.ToProviderType(value)
		return err
	})

	test.RunTable(t, valueCases, func(value string) entities.ProviderType {
		result, _ = entities.ToProviderType(value)
		return result
	})
}
