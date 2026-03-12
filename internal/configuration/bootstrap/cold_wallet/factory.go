package cold_wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment/secrets"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
)

func Create(rpc blockchain.Rpc, config secrets.ColdWalletConfiguration) (cold_wallet.ColdWallet, error) {
	switch config.Type {
	case "static":
		var args cold_wallet.StaticColdWalletArgs
		if err := json.Unmarshal(config.Configuration, &args); err != nil {
			return nil, fmt.Errorf("invalid %s cold wallet configuration: %w", config.Type, err)
		}
		return cold_wallet.NewStaticColdWallet(rpc, args), nil
	default:
		return nil, errors.New("unknown cold wallet type")
	}
}
