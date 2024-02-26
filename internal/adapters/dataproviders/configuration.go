package dataproviders

type Configuration struct {
	RskConfig    RskConfig
	BtcConfig    BitcoinConfig
	PeginConfig  PeginConfig
	PegoutConfig PegoutConfig
}

type RskConfig struct {
	ChainId       uint64
	Account       uint64
	Confirmations map[int]uint16
}

type BitcoinConfig struct {
	BtcAddress    string
	Confirmations map[int]uint16
}

// This structures were kept just in case, right now all the parameters are manipulated through management API

type PeginConfig struct{}

type PegoutConfig struct{}
