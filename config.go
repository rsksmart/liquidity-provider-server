package main

type config struct {
	LogFile              string
	Debug                bool
	IrisActivationHeight int
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
		Username string
		Password string
		Network  string
	}
	Provider struct {
		Keystore      string
		RskAccountNum int
		PwdFilePath   string
		BtcAddress    string
	}
}
