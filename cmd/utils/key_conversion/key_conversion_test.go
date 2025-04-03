package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/account"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"sync"
	"testing"
)

type keyConversionTestDataset struct {
	HexKeyBytes       []byte
	TestnetWifKey     string
	MainnetWifKey     string
	RskAddress        string
	BtcMainnetAddress string
	BtcTestnetAddress string
	Keystore          string
	KeystorePassword  string
}

type keystoreContent struct {
	Address string `json:"address"`
}

var keystoreTestLock = sync.Mutex{}

// IMPORTANT: these keys were generated for this test specifically, please DON'T use them in any network
var testDataset = []keyConversionTestDataset{
	{
		HexKeyBytes:       []byte{0xba, 0x65, 0xb6, 0x56, 0x5f, 0xdc, 0xd2, 0x63, 0x81, 0x8f, 0xd1, 0x0c, 0xf6, 0xc1, 0x1e, 0xf7, 0x9a, 0x20, 0x9f, 0x30, 0x49, 0x3c, 0x2a, 0xe4, 0x7e, 0xcc, 0xf9, 0x5e, 0x88, 0x1a, 0xc3, 0x83},
		TestnetWifKey:     "cTq2uVE9nuHJKSw6ntC5PPL8wpjatWSqjQJddrYQMzfrGuGBTK3d",
		MainnetWifKey:     "L3U3SaEJMqb3A1TqQUNx24q5KbSBE4M9fNAAXS5trt1r2AB6jjqE",
		RskAddress:        "0x68b62b92fea356a59a44497a36a2a42cae28c2dd",
		BtcMainnetAddress: "1NzrNq3RMNuc2qqgVy8dRLBbDQEFtWUQiW",
		BtcTestnetAddress: "n3Woft8QAQLroxKJDY71FFPv5PpxqHw8qw",
		Keystore:          "{\"address\":\"68b62b92fea356a59a44497a36a2a42cae28c2dd\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"2ae62c36fb88cd694618bdce60d69b0e64cfea3858835028cad64b1d34b50806\",\"cipherparams\":{\"iv\":\"25bf9053d0c4618dc72597313622e52a\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"0750e94f308e6ee9c5f488231e5305802fa923041ddd3d17ce333107794dcb83\"},\"mac\":\"b4d19295767c071d423c1fd5034e6a6f42d77d89a0101674a6dfe125a9b78b45\"},\"id\":\"3c5351fc-e99f-4700-b7a0-2a57065c6338\",\"version\":3}",
		KeystorePassword:  "unit-test-1",
	},
	{
		HexKeyBytes:       []byte{0x19, 0x44, 0xf5, 0x71, 0xe7, 0x08, 0x84, 0xed, 0xbb, 0xab, 0x4b, 0x40, 0xb4, 0xf0, 0xd5, 0x0a, 0x86, 0x19, 0xf9, 0x28, 0xfe, 0xe7, 0xa3, 0x11, 0x01, 0x6e, 0xf0, 0x43, 0x60, 0x8b, 0xc5, 0xb0},
		TestnetWifKey:     "cNRpdSYzhvU1DRwcnKfY2Fgv1Khg2vWktKAv6c4yNxidkqZ9t4ax",
		MainnetWifKey:     "Kx4qAXZ9Grmk3zUMPurQewBrP6QGNUR4pH2SzBcTsr4dW6YuTT58",
		RskAddress:        "0x390b06eac314cdf661e158835c73dc1e1a552e90",
		BtcMainnetAddress: "19St9boSRFdH7jsmm7parkoeBb9Gq9KTjU",
		BtcTestnetAddress: "moxqSetREH4XtrMPUgnxgg1y3ajyp5s2bx",
		Keystore:          "{\"address\":\"390b06eac314cdf661e158835c73dc1e1a552e90\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"15db3c05252cc44d0ccb7092378803eec564a4d784f5331905d005cf724616cc\",\"cipherparams\":{\"iv\":\"47fd1361f57733d4ae909876d989d77a\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"744dfeaabc59e9af3882d66bfed90cf1ea7c88b53aed3c1d6f8a8835b829c2b5\"},\"mac\":\"39a58286b8ff12df2eedad0f0c9267f33037f2f8221d32961a8548b03323de25\"},\"id\":\"3a586bf3-5a81-4ae1-9e09-46b8d67983ff\",\"version\":3}",
		KeystorePassword:  "random-pwd-test",
	},
	{
		HexKeyBytes:       []byte{0x91, 0xfc, 0xfc, 0x2d, 0xee, 0x42, 0xcc, 0x4d, 0x9d, 0x18, 0x14, 0x4b, 0x54, 0xde, 0x88, 0x4d, 0x86, 0xa1, 0x71, 0x5f, 0x48, 0xd2, 0xab, 0x5d, 0x71, 0xee, 0x87, 0x81, 0xec, 0x0b, 0x57, 0x7b},
		TestnetWifKey:     "cSUV1KGWtP1mwhnrLVpm685DHDqyeeh5KAV2L4H1oMWtXh7tyQ28",
		MainnetWifKey:     "L27VYQGfTKKWnGKax61dioa9ezYZzCbPF8LZDdpWJErtGx6bJaXx",
		RskAddress:        "0x169dc328e420329b61ee83df392a2e22f3deb292",
		BtcMainnetAddress: "1KFSDfxFpN23ETNCq7TpJzVyGusqLC3kUo",
		BtcTestnetAddress: "mymPWj3EdPTJ1ZqpYgSC8uiJ8uUYJp4DHj",
		Keystore:          "{\"address\":\"169dc328e420329b61ee83df392a2e22f3deb292\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"fcf05b77507d531d3f6ff6f9d28227ec971bfbc42890013598863092a4f818d1\",\"cipherparams\":{\"iv\":\"8f8502743df2454609e40934d0a34ebe\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"8e581fe703a4ddac73fd0af8de4b1413ff6f1d2ef572f1684657de45aa30033c\"},\"mac\":\"800658292680b741901acd098d9e4aa502af138180f01611e2f168677022f738\"},\"id\":\"f1656990-92b0-4479-a88a-7b173e00ed5f\",\"version\":3}",
		KeystorePassword:  "test-pwd-123",
	},
}

const keyFileTemplate = "text-key-%d.txt"

func TestWifToHex(t *testing.T) {
	t.Run("should convert WIF key to hex key", func(t *testing.T) {
		for _, dataset := range testDataset {
			hexKey, err := WifToHex(dataset.TestnetWifKey)
			require.NoError(t, err)
			assert.Equal(t, dataset.HexKeyBytes, hexKey)
		}
		for _, dataset := range testDataset {
			hexKey, err := WifToHex(dataset.MainnetWifKey)
			require.NoError(t, err)
			assert.Equal(t, dataset.HexKeyBytes, hexKey)
		}
	})
	t.Run("should return error when WIF key is invalid", func(t *testing.T) {
		hexKey, err := WifToHex("invalid-wif")
		require.Error(t, err)
		assert.Nil(t, hexKey)
	})
}

func TestParseRawKeyInput(t *testing.T) {
	t.Run("should parse raw key input with rsk origin", func(t *testing.T) {
		for _, dataset := range testDataset {
			privKey, err := ParseRawKeyInput("rsk", []byte(hex.EncodeToString(dataset.HexKeyBytes)))
			require.NoError(t, err)
			assert.Equal(t, dataset.HexKeyBytes, privKey.Serialize())
		}
	})
	t.Run("should parse raw key input with btc origin", func(t *testing.T) {
		for _, dataset := range testDataset {
			privKey, err := ParseRawKeyInput("btc", []byte(dataset.MainnetWifKey))
			require.NoError(t, err)
			assert.Equal(t, dataset.HexKeyBytes, privKey.Serialize())
		}
		for _, dataset := range testDataset {
			privKey, err := ParseRawKeyInput("btc", []byte(dataset.TestnetWifKey))
			require.NoError(t, err)
			assert.Equal(t, dataset.HexKeyBytes, privKey.Serialize())
		}
	})
	t.Run("should return error when key is invalid", func(t *testing.T) {
		rskKey, rskErr := ParseRawKeyInput("rsk", []byte("invalid-hex-key"))
		require.Error(t, rskErr)
		assert.Empty(t, rskKey)
		btcKey, btcErr := ParseRawKeyInput("btc", []byte("invalid-wif-key"))
		require.Error(t, btcErr)
		assert.Empty(t, btcKey)
	})
	t.Run("should return error when origin is invalid", func(t *testing.T) {
		privKey, err := ParseRawKeyInput("invalid-origin", []byte("invalid-key"))
		require.ErrorContains(t, err, "invalid origin network")
		assert.Empty(t, privKey)
	})
}

func TestCreateKeystore(t *testing.T) {
	keystoreTestLock.Lock()
	defer keystoreTestLock.Unlock()
	const testPassword = "unit-test-password"
	t.Run("should fail if password is empty", func(t *testing.T) {
		keystoreBytes, password, err := CreateKeystore(&KeyConversionScriptInput{}, secp256k1.PrivateKey{}, func(i int) ([]byte, error) { return []byte{}, nil })
		require.Error(t, err)
		assert.Nil(t, keystoreBytes)
		assert.Empty(t, password)
	})
	t.Run("should read existing keystore", func(t *testing.T) {
		for i, dataset := range testDataset {
			keystorePath := test.WriteTestFile(t, fmt.Sprintf("existing-keystore-%d.json", i), []byte(dataset.Keystore))
			pwdReaderFunc := func(i int) ([]byte, error) { return []byte(dataset.KeystorePassword), nil }
			keystoreBytes, password, err := CreateKeystore(&KeyConversionScriptInput{KeySource: keystoreKeySource, InputFile: keystorePath}, secp256k1.PrivateKey{}, pwdReaderFunc)
			require.NoError(t, err)
			assert.Equal(t, dataset.Keystore, string(keystoreBytes))
			assert.Equal(t, dataset.KeystorePassword, password)
			require.NoError(t, os.RemoveAll(keystoreDir))
		}
	})
	t.Run("should create new keystore if key was obtained from terminal", func(t *testing.T) {
		for _, dataset := range testDataset {
			privateKeyParsed := secp256k1.PrivKeyFromBytes(dataset.HexKeyBytes)
			pwdReaderFunc := func(i int) ([]byte, error) { return []byte(dataset.KeystorePassword), nil }
			keystoreBytes, password, err := CreateKeystore(&KeyConversionScriptInput{KeySource: terminalKeySource}, *privateKeyParsed, pwdReaderFunc)
			require.NoError(t, err)
			content := &keystoreContent{}
			require.NoError(t, json.Unmarshal(keystoreBytes, content))
			assert.Equal(t, dataset.RskAddress[2:], content.Address)
			assert.Equal(t, dataset.KeystorePassword, password)
			// ensure created keystore can be opened
			_, err = account.GetRskAccount(account.CreationArgs{KeyDir: keystoreDir, AccountNum: 0, EncryptedJson: string(keystoreBytes), Password: password})
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(keystoreDir))
		}
	})
	t.Run("should create new keystore if key was obtained from a text file", func(t *testing.T) {
		for i, dataset := range testDataset {
			privateKeyParsed := secp256k1.PrivKeyFromBytes(dataset.HexKeyBytes)
			keyPath := test.WriteTestFile(t, fmt.Sprintf(keyFileTemplate, i), []byte(hex.EncodeToString(dataset.HexKeyBytes)))
			pwdReaderFunc := func(i int) ([]byte, error) { return []byte(dataset.KeystorePassword), nil }
			keystoreBytes, password, err := CreateKeystore(&KeyConversionScriptInput{KeySource: fileKeySource, InputFile: keyPath}, *privateKeyParsed, pwdReaderFunc)
			require.NoError(t, err)
			content := &keystoreContent{}
			require.NoError(t, json.Unmarshal(keystoreBytes, content))
			assert.Equal(t, dataset.RskAddress[2:], content.Address)
			assert.Equal(t, dataset.KeystorePassword, password)
			// ensure created keystore can be opened
			_, err = account.GetRskAccount(account.CreationArgs{KeyDir: keystoreDir, AccountNum: 0, EncryptedJson: string(keystoreBytes), Password: password})
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(keystoreDir))
		}
	})
	t.Run("should fail if key source is invalid", func(t *testing.T) {
		scriptInput := &KeyConversionScriptInput{KeySource: "invalid-source"}
		keystoreBytes, password, err := CreateKeystore(scriptInput, secp256k1.PrivateKey{}, func(i int) ([]byte, error) { return []byte(testPassword), nil })
		require.Error(t, err)
		assert.Nil(t, keystoreBytes)
		assert.Empty(t, password)
	})
}

func TestGetKeystoreAndPassword(t *testing.T) {
	keystoreTestLock.Lock()
	defer keystoreTestLock.Unlock()
	t.Run("should get keystore and password from existing keystore file", func(t *testing.T) {
		for i, dataset := range testDataset {
			keystorePath := test.WriteTestFile(t, fmt.Sprintf("existing-keystore-%d.json", i), []byte(dataset.Keystore))
			pwdReaderFunc := func(i int) ([]byte, error) { return []byte(dataset.KeystorePassword), nil }
			keystoreBytes, password, err := GetKeystoreAndPassword(&KeyConversionScriptInput{KeySource: keystoreKeySource, InputFile: keystorePath}, pwdReaderFunc)
			require.NoError(t, err)
			assert.Equal(t, dataset.Keystore, string(keystoreBytes))
			assert.Equal(t, dataset.KeystorePassword, password)
			_, err = account.GetRskAccount(account.CreationArgs{KeyDir: keystoreDir, AccountNum: 0, EncryptedJson: string(keystoreBytes), Password: password})
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(keystoreDir))
		}
	})
	t.Run("should get keystore and password from a private key provided by the terminal", func(t *testing.T) {
		for _, dataset := range testDataset {
			createPwdReader := func(d keyConversionTestDataset) func(i int) ([]byte, error) {
				called := false
				return func(i int) ([]byte, error) {
					if !called {
						called = true
						return []byte(hex.EncodeToString(d.HexKeyBytes)), nil
					}
					return []byte(d.KeystorePassword), nil
				}
			}
			input := &KeyConversionScriptInput{KeySource: terminalKeySource, OriginBlockchain: rskOriginBlockchain}
			keystoreBytes, password, err := GetKeystoreAndPassword(input, createPwdReader(dataset))
			require.NoError(t, err)
			content := &keystoreContent{}
			require.NoError(t, json.Unmarshal(keystoreBytes, content))
			assert.Equal(t, dataset.RskAddress[2:], content.Address)
			assert.Equal(t, dataset.KeystorePassword, password)
			_, err = account.GetRskAccount(account.CreationArgs{KeyDir: keystoreDir, AccountNum: 0, EncryptedJson: string(keystoreBytes), Password: password})
			require.NoError(t, err) // ensure created keystore can be opened
			require.NoError(t, os.RemoveAll(keystoreDir))
		}
	})
	t.Run("should get keystore and password from a private key provided by a text file", func(t *testing.T) {
		for i, dataset := range testDataset {
			keyPath := test.WriteTestFile(t, fmt.Sprintf(keyFileTemplate, i), []byte(hex.EncodeToString(dataset.HexKeyBytes)))
			pwdReaderFunc := func(i int) ([]byte, error) { return []byte(dataset.KeystorePassword), nil }
			keystoreBytes, password, err := GetKeystoreAndPassword(&KeyConversionScriptInput{KeySource: fileKeySource, InputFile: keyPath, OriginBlockchain: rskOriginBlockchain}, pwdReaderFunc)
			require.NoError(t, err)
			content := &keystoreContent{}
			require.NoError(t, json.Unmarshal(keystoreBytes, content))
			assert.Equal(t, dataset.RskAddress[2:], content.Address)
			assert.Equal(t, dataset.KeystorePassword, password)
			_, err = account.GetRskAccount(account.CreationArgs{KeyDir: keystoreDir, AccountNum: 0, EncryptedJson: string(keystoreBytes), Password: password})
			require.NoError(t, err) // ensure created keystore can be opened
			require.NoError(t, os.RemoveAll(keystoreDir))
		}
	})
	t.Run("should fail if key source is invalid", func(t *testing.T) {
		keystoreBytes, password, err := GetKeystoreAndPassword(&KeyConversionScriptInput{KeySource: "invalid-source"}, func(i int) ([]byte, error) { return nil, nil })
		require.Error(t, err)
		assert.Nil(t, keystoreBytes)
		assert.Empty(t, password)
	})
}

func TestShowKeys(t *testing.T) {
	keystoreTestLock.Lock()
	defer keystoreTestLock.Unlock()
	showKeysTest := func(t *testing.T, network string) {
		var expectedBtcAddress, expectedWifKey string
		for i, dataset := range testDataset {
			if network == "mainnet" {
				expectedBtcAddress = dataset.BtcMainnetAddress
				expectedWifKey = dataset.MainnetWifKey
			} else {
				expectedBtcAddress = dataset.BtcTestnetAddress
				expectedWifKey = dataset.TestnetWifKey
			}
			fileName := fmt.Sprintf(keyFileTemplate, i)
			keyPath := test.WriteTestFile(t, fileName, []byte(hex.EncodeToString(dataset.HexKeyBytes)))
			input := &KeyConversionScriptInput{
				KeySource:        fileKeySource,
				OriginBlockchain: rskOriginBlockchain,
				Network:          network,
				InputFile:        keyPath,
			}
			rskAccount, err := ParseKeyConversionScriptInput(flag.Parse, func(i int) ([]byte, error) {
				return []byte(dataset.KeystorePassword), nil
			}, input)
			require.NoError(t, err)

			r, w, err := os.Pipe()
			require.NoError(t, err)
			originalStdOut := os.Stdout
			os.Stdout = w

			err = ShowKeys(*rskAccount)
			require.NoError(t, err)

			os.Stdout = originalStdOut
			require.NoError(t, w.Close())
			result, err := io.ReadAll(r)
			require.NoError(t, err)
			assert.Contains(t, string(result), "BTC Address: "+expectedBtcAddress)
			assert.Contains(t, string(result), "RSK Address: "+dataset.RskAddress)
			assert.Contains(t, string(result), "BTC Private Key WIF: "+expectedWifKey)
			assert.Contains(t, string(result), "RSK Private Key: "+hex.EncodeToString(dataset.HexKeyBytes))
			require.NoError(t, os.RemoveAll(keystoreDir))
		}
	}
	t.Run("should show keys for testnet", func(t *testing.T) {
		showKeysTest(t, "testnet")
	})
	t.Run("should show keys for mainnet", func(t *testing.T) {
		showKeysTest(t, "mainnet")
	})
	t.Run("should show keys for regtest", func(t *testing.T) {
		showKeysTest(t, "regtest")
	})
}
