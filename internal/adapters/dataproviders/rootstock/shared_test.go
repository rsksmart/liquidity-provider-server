package rootstock_test

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	discoveryBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/discovery"
	pegoutBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout"
	"github.com/stretchr/testify/require"
	"testing"
)

func filterMatchFunc(from uint64, to uint64) func(ethereum.FilterQuery) bool {
	return func(opts ethereum.FilterQuery) bool {
		return from == opts.FromBlock.Uint64() && to == opts.ToBlock.Uint64()
	}
}

func mustPackBool(t *testing.T, v bool) []byte {
	t.Helper()
	boolType, err := abi.NewType("bool", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: boolType}}.Pack(v)
	require.NoError(t, err)
	return out
}

func mustPackUint8(t *testing.T, v uint8) []byte {
	t.Helper()
	uint8Type, err := abi.NewType("uint8", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: uint8Type}}.Pack(v)
	require.NoError(t, err)
	return out
}

func mustPackString(t *testing.T, v string) []byte {
	t.Helper()
	strType, err := abi.NewType("string", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: strType}}.Pack(v)
	require.NoError(t, err)
	return out
}

func mustPackBytes(t *testing.T, v string) []byte {
	t.Helper()
	strType, err := abi.NewType("bytes", "", nil)
	require.NoError(t, err)

	parsedValue, err := hex.DecodeString(v)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: strType}}.Pack(parsedValue)
	require.NoError(t, err)
	return out
}

func mustPackBytes32(t *testing.T, v [32]byte) []byte {
	t.Helper()
	byteType, err := abi.NewType("bytes32", "", nil)
	require.NoError(t, err)

	out, err := abi.Arguments{{Type: byteType}}.Pack(v)
	require.NoError(t, err)
	return out
}

func mustPackLiquidityProvider(t *testing.T, provider discoveryBindings.FlyoverLiquidityProvider) []byte {
	t.Helper()
	lpType, err := abi.NewType("tuple", "structFlyover.LiquidityProvider", []abi.ArgumentMarshaling{
		{Name: "id", Type: "uint256"},
		{Name: "providerAddress", Type: "address"},
		{Name: "status", Type: "bool"},
		{Name: "providerType", Type: "uint8"},
		{Name: "name", Type: "string"},
		{Name: "apiBaseUrl", Type: "string"},
	})
	if err != nil {
		t.Fatal(err)
	}
	args := abi.Arguments{{Type: lpType}}
	encoded, err := args.Pack(provider)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func mustPackLiquidityProviders(t *testing.T, provider []discoveryBindings.FlyoverLiquidityProvider) []byte {
	t.Helper()
	lpType, err := abi.NewType("tuple[]", "structFlyover.LiquidityProvider[]", []abi.ArgumentMarshaling{
		{Name: "id", Type: "uint256"},
		{Name: "providerAddress", Type: "address"},
		{Name: "status", Type: "bool"},
		{Name: "providerType", Type: "uint8"},
		{Name: "name", Type: "string"},
		{Name: "apiBaseUrl", Type: "string"},
	})
	if err != nil {
		t.Fatal(err)
	}
	args := abi.Arguments{{Type: lpType}}
	encoded, err := args.Pack(provider)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func mustPackPegoutQuote(t *testing.T, q pegoutBindings.QuotesPegOutQuote) []byte {
	t.Helper()
	quoteType, err := abi.NewType("tuple", "structQuotes.PegOutQuote", []abi.ArgumentMarshaling{
		{Name: "callFee", Type: "uint256"},
		{Name: "penaltyFee", Type: "uint256"},
		{Name: "value", Type: "uint256"},
		{Name: "gasFee", Type: "uint256"},
		{Name: "lbcAddress", Type: "address"},
		{Name: "lpRskAddress", Type: "address"},
		{Name: "rskRefundAddress", Type: "address"},
		{Name: "nonce", Type: "int64"},
		{Name: "agreementTimestamp", Type: "uint32"},
		{Name: "depositDateLimit", Type: "uint32"},
		{Name: "transferTime", Type: "uint32"},
		{Name: "expireDate", Type: "uint32"},
		{Name: "expireBlock", Type: "uint32"},
		{Name: "depositConfirmations", Type: "uint16"},
		{Name: "transferConfirmations", Type: "uint16"},
		{Name: "depositAddress", Type: "bytes"},
		{Name: "btcRefundAddress", Type: "bytes"},
		{Name: "lpBtcAddress", Type: "bytes"},
	})
	if err != nil {
		t.Fatal(err)
	}
	args := abi.Arguments{{Type: quoteType}}
	encoded, err := args.Pack(q)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}
