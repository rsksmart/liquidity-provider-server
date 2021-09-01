package main

type config struct {
	LogFile              string
	Debug                bool
	IrisActivationHeight int
	IsTestNet            bool
	ErpKeys              []string

	Server struct {
		Port uint
	}
	DB struct {
		Path string
	}
	RSK struct {
		Endpoint   string
		LBCAddr    string
		BridgeAddr string
	}
	BTC struct {
		Endpoint string
	}
	Provider struct {
		Keystore      string
		RskAccountNum int
		PwdFilePath   string
		BtcAddress    string
	}
}
