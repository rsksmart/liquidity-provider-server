package bitcoin

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

type bitcoinNetworkAddresses struct {
	P2PKH  string
	P2SH   string
	P2WPKH string
	P2WSH  string
	P2TR   string
}

type bitcoinAllNetworksAddresses struct {
	Mainnet bitcoinNetworkAddresses
	Testnet bitcoinNetworkAddresses
	Regtest bitcoinNetworkAddresses
}

func (networks *bitcoinAllNetworksAddresses) Network(network string) (bitcoinNetworkAddresses, error) {
	switch network {
	case "mainnet":
		return networks.Mainnet, nil
	case "testnet3":
		return networks.Testnet, nil
	case "regtest":
		return networks.Regtest, nil
	default:
		return bitcoinNetworkAddresses{}, errors.New("unknown network")
	}
}

func (addresses *bitcoinNetworkAddresses) Address(address blockchain.BtcAddressType) (string, error) {
	switch address {
	case "p2pkh":
		return addresses.P2PKH, nil
	case "p2sh":
		return addresses.P2SH, nil
	case "p2wpkh":
		return addresses.P2WPKH, nil
	case "p2wsh":
		return addresses.P2WSH, nil
	case "p2tr":
		return addresses.P2TR, nil
	default:
		return "", errors.New("unknown address type")
	}
}

var bitcoinZeroAddresses = bitcoinAllNetworksAddresses{
	Mainnet: bitcoinNetworkAddresses{
		P2PKH:  blockchain.BitcoinMainnetP2PKHZeroAddress,
		P2SH:   blockchain.BitcoinMainnetP2SHZeroAddress,
		P2WPKH: blockchain.BitcoinMainnetP2WPKHZeroAddress,
		P2WSH:  blockchain.BitcoinMainnetP2WSHZeroAddress,
		P2TR:   blockchain.BitcoinMainnetP2TRZeroAddress,
	},
	Testnet: bitcoinNetworkAddresses{
		P2PKH:  blockchain.BitcoinTestnetP2PKHZeroAddress,
		P2SH:   blockchain.BitcoinTestnetP2SHZeroAddress,
		P2WPKH: blockchain.BitcoinTestnetP2WPKHZeroAddress,
		P2WSH:  blockchain.BitcoinTestnetP2WSHZeroAddress,
		P2TR:   blockchain.BitcoinTestnetP2TRZeroAddress,
	},
	Regtest: bitcoinNetworkAddresses{
		P2PKH:  blockchain.BitcoinTestnetP2PKHZeroAddress,
		P2SH:   blockchain.BitcoinTestnetP2SHZeroAddress,
		P2WPKH: blockchain.BitcoinRegtestP2WPKHZeroAddress,
		P2WSH:  blockchain.BitcoinRegtestP2WSHZeroAddress,
		P2TR:   blockchain.BitcoinRegtestP2TRZeroAddress,
	},
}
