package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/rsksmart/liquidity-provider-server/cmd/utils/scripts"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

const keystoreDir = "keystore"

const (
	rskOriginBlockchain = "rsk"
	btcOriginBlockchain = "btc"
)

const (
	terminalKeySource = "terminal"
	fileKeySource     = "file"
	keystoreKeySource = "keystore"
)

type KeyConversionScriptInput struct {
	KeySource        string `validate:"required,oneof=keystore file terminal"`
	OriginBlockchain string `validate:"required,oneof=btc rsk"`
	Network          string `validate:"required,oneof=regtest testnet mainnet"`
	InputFile        string `validate:"omitempty,filepath"`
}

func main() {
	scripts.SetUsageMessage(
		"This script can be used to get the corresponding BTC address, RSK address, hex private key and WIF " +
			"private key from a specific private key. Also, it generates a keystore file to contain the provided key.",
	)
	scriptInput := new(KeyConversionScriptInput)
	ReadKeyConversionScriptInput(scriptInput)
	rskAccount, err := ParseKeyConversionScriptInput(flag.Parse, term.ReadPassword, scriptInput)
	if err != nil {
		scripts.ExitWithError(2, "Error reading input", err)
	}
	err = ShowKeys(*rskAccount)
	if err != nil {
		scripts.ExitWithError(2, "Error showing keys", err)
	}
}

func ReadKeyConversionScriptInput(scriptInput *KeyConversionScriptInput) {
	flag.StringVar(&scriptInput.KeySource, "key-src", "", "The source to obtain the private key from. Must be one of:\n"+
		"- keystore: if the key is going to be provided through an existing keystore file.\n"+
		"- file: if the key will be provided through a plain text file.\n"+
		"- terminal: if the key will be provided through the terminal.")
	flag.StringVar(&scriptInput.OriginBlockchain, "origin-blockchain", "", "The blockchain where the key was "+
		"originally generated. The format to interpret the private key will depend on this value. Must be:\n"+
		"- btc: if the private key should be interpreted in WIF.\n"+
		"- rsk: if the private key should be interpreted in hex format.")
	flag.StringVar(&scriptInput.Network, "network", "", "The network to generate the addresses to. Must be one of: regtest, testnet, mainnet.")
	flag.StringVar(&scriptInput.InputFile, "input-file", "", "The input file to obtain the private key from. Only required if the key source is 'file' or 'keystore'.")
}

func ParseKeyConversionScriptInput(parse scripts.ParseFunc, pwdReader scripts.PasswordReader, scriptInput *KeyConversionScriptInput) (*account.RskAccount, error) {
	var keystoreBytes []byte
	var password string
	var err error

	defer func() {
		for i := range keystoreBytes {
			keystoreBytes[i] = 0
		}
		password = ""
	}()

	parse()
	btcEnv := environment.BtcEnv{Network: scriptInput.Network}
	btcNetwork, err := btcEnv.GetNetworkParams()
	if err != nil {
		return nil, err
	}
	keystoreBytes, password, err = GetKeystoreAndPassword(scriptInput, pwdReader)
	if err != nil {
		return nil, err
	}
	rskAccount, err := account.GetRskAccountWithDerivation(account.CreationWithDerivationArgs{
		CreationArgs: account.CreationArgs{
			KeyDir:        keystoreDir,
			AccountNum:    0,
			EncryptedJson: string(keystoreBytes),
			Password:      password,
		},
		BtcParams: btcNetwork,
	})
	if err != nil {
		return nil, err
	}
	return rskAccount, nil
}

func ShowKeys(rskAccount account.RskAccount) error {
	btcAddress, err := rskAccount.BtcAddress()
	if err != nil {
		return err
	}
	fmt.Println("BTC Address:", btcAddress.EncodeAddress())
	fmt.Println("RSK Address:", strings.ToLower(rskAccount.Account.Address.Hex()))
	return rskAccount.UsePrivateKeyWif(func(wif *btcutil.WIF) error {
		fmt.Println("BTC Private Key WIF:", wif)
		fmt.Println("RSK Private Key:", hex.EncodeToString(wif.PrivKey.Serialize()))
		return nil
	})
}

func GetKeystoreAndPassword(scriptInput *KeyConversionScriptInput, pwdReader scripts.PasswordReader) ([]byte, string, error) {
	var rawPrivateKey []byte
	var parsedPrivateKey secp256k1.PrivateKey
	var err error
	defer func() {
		for i := range rawPrivateKey {
			rawPrivateKey[i] = 0
		}
		parsedPrivateKey.Zero()
	}()
	switch scriptInput.KeySource {
	case keystoreKeySource:
		return CreateKeystore(scriptInput, parsedPrivateKey, pwdReader)
	case fileKeySource:
		if rawPrivateKey, err = os.ReadFile(scriptInput.InputFile); err != nil {
			return nil, "", err
		}
		if parsedPrivateKey, err = ParseRawKeyInput(scriptInput.OriginBlockchain, rawPrivateKey); err != nil {
			return nil, "", err
		}
		return CreateKeystore(scriptInput, parsedPrivateKey, pwdReader)
	case terminalKeySource:
		fmt.Println("Insert private key:")
		if rawPrivateKey, err = pwdReader(syscall.Stdin); err != nil {
			return nil, "", err
		}
		if parsedPrivateKey, err = ParseRawKeyInput(scriptInput.OriginBlockchain, rawPrivateKey); err != nil {
			return nil, "", err
		}
		return CreateKeystore(scriptInput, parsedPrivateKey, pwdReader)
	default:
		return nil, "", errors.New("invalid source")
	}
}

func CreateKeystore(scriptInput *KeyConversionScriptInput, privateKey secp256k1.PrivateKey, pwdReader scripts.PasswordReader) ([]byte, string, error) {
	var passwordBytes, keystoreBytes []byte
	var importedAccount accounts.Account
	var err error

	fmt.Println("Insert keystore password:")
	passwordBytes, err = pwdReader(syscall.Stdin)
	if err != nil {
		return nil, "", err
	}
	password := string(passwordBytes)
	if scriptInput.KeySource == keystoreKeySource {
		keystoreBytes, err = os.ReadFile(scriptInput.InputFile)
		return keystoreBytes, password, err
	} else if scriptInput.KeySource == fileKeySource || scriptInput.KeySource == terminalKeySource {
		ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
		if importedAccount, err = ks.ImportECDSA(privateKey.ToECDSA(), password); err != nil {
			return nil, "", err
		}
		if keystoreBytes, err = ks.Export(importedAccount, password, password); err != nil {
			return nil, "", err
		}
		return keystoreBytes, password, nil
	} else {
		return nil, "", errors.New("invalid source")
	}
}

func WifToHex(wifString string) ([]byte, error) {
	wif, err := btcutil.DecodeWIF(wifString)
	if err != nil {
		return nil, err
	}
	defer func(w *btcutil.WIF) { w.PrivKey.Zero() }(wif)
	return wif.PrivKey.Serialize(), nil
}

func ParseRawKeyInput(originNetwork string, rawKey []byte) (secp256k1.PrivateKey, error) {
	var err error
	var pkBytes []byte
	if originNetwork == btcOriginBlockchain {
		pkBytes, err = WifToHex(strings.TrimSpace(string(rawKey)))
	} else if originNetwork == rskOriginBlockchain {
		pkBytes, err = hex.DecodeString(strings.TrimSpace(string(rawKey)))
	} else {
		return secp256k1.PrivateKey{}, errors.New("invalid origin network")
	}
	if err != nil {
		return secp256k1.PrivateKey{}, err
	}
	parsedPk := secp256k1.PrivKeyFromBytes(pkBytes)
	defer func() {
		parsedPk.Zero()
		for i := range pkBytes {
			pkBytes[i] = 0
		}
		for i := range rawKey {
			rawKey[i] = 0
		}
	}()
	return *parsedPk, nil
}
