package environment

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"sync"
)

type applicationMutexesImpl struct {
	rskWalletMutex       *sync.Mutex
	btcWalletMutex       *sync.Mutex
	peginLiquidityMutex  *sync.Mutex
	pegoutLiquidityMutex *sync.Mutex
}

func NewApplicationMutexes() entities.ApplicationMutexes {
	rskWalletMutex := sync.Mutex{}
	btcWalletMutex := sync.Mutex{}
	peginLiquidityMutex := sync.Mutex{}
	pegoutLiquidityMutex := sync.Mutex{}
	return &applicationMutexesImpl{
		rskWalletMutex:       &rskWalletMutex,
		btcWalletMutex:       &btcWalletMutex,
		peginLiquidityMutex:  &peginLiquidityMutex,
		pegoutLiquidityMutex: &pegoutLiquidityMutex,
	}
}

func (a *applicationMutexesImpl) RskWalletMutex() *sync.Mutex {
	return a.rskWalletMutex
}

func (a *applicationMutexesImpl) PeginLiquidityMutex() *sync.Mutex {
	return a.peginLiquidityMutex
}

func (a *applicationMutexesImpl) PegoutLiquidityMutex() *sync.Mutex {
	return a.pegoutLiquidityMutex
}

func (a *applicationMutexesImpl) BtcWalletMutex() *sync.Mutex {
	return a.btcWalletMutex
}
