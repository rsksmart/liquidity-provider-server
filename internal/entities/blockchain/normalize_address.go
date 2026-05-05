package blockchain

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// NormalizeEthereumAddress returns the canonical lowercase 0x-prefixed 40-hex form
// for Ethereum/RSK-style account addresses. Input may use any EIP-55 casing.
func NormalizeEthereumAddress(addr string) (string, error) {
	trimmed := strings.TrimSpace(addr)
	if !common.IsHexAddress(trimmed) {
		return "", fmt.Errorf("invalid hex address %q", addr)
	}
	return strings.ToLower(common.HexToAddress(trimmed).Hex()), nil
}
