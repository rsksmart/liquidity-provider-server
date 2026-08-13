package watcher_test

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rsksmart/liquidity-provider-server/test/datasets"
	"github.com/stretchr/testify/require"
)

type depositAddress struct {
	payload []byte
	address string
}

type pinnedLBCRegistration struct {
	rskAddress string
	root       [32]byte
}

func pinnedLBCRegistrations(t *testing.T) []pinnedLBCRegistration {
	t.Helper()
	return []pinnedLBCRegistration{
		{
			rskAddress: "0x00000000000000000000000000000000000000a1",
			root:       registrationRoot(t, "5a16856e66cb2b1b463f7773c427085d55afdd19d778290b45fb959a6224877e"),
		},
		{
			rskAddress: "0x00000000000000000000000000000000000000a2",
			root:       registrationRoot(t, "95bdcbc864b6248af173e5feb03f5830479e972a492e705c3c2a2bddcb8ca643"),
		},
		{
			rskAddress: "0x00000000000000000000000000000000000000a3",
			root:       registrationRoot(t, "b80604f49f0685bb17dd0f5cc0f611383d724f10f53db8aebcec3e9541f552d8"),
		},
	}
}

func knownDepositAddress(index int) depositAddress {
	const checksumSize = 4
	decoded := datasets.Base58Addresses[index]
	payload := make([]byte, 0, len(decoded.Expected)+checksumSize)
	payload = append(payload, decoded.Expected...)
	payload = append(payload, chainhash.DoubleHashB(decoded.Expected)[:checksumSize]...)
	return depositAddress{payload: payload, address: decoded.Address}
}

func registrationRoot(t *testing.T, encoded string) [32]byte {
	t.Helper()
	decoded, err := hex.DecodeString(encoded)
	require.NoError(t, err)
	require.Len(t, decoded, 32)
	return [32]byte(decoded)
}
