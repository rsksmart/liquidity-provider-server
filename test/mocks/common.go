package mocks

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
)

func GetFakeFedInfo() rootstock.FederationInfo {
	var keys []string
	keys = append(keys, "02cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1")
	keys = append(keys, "0362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a124")
	keys = append(keys, "03c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db")

	var erpPubKeys []string
	erpPubKeys = append(erpPubKeys, "0257c293086c4d4fe8943deda5f890a37d11bebd140e220faa76258a41d077b4d4")
	erpPubKeys = append(erpPubKeys, "03c2660a46aa73078ee6016dee953488566426cf55fc8011edd0085634d75395f9")
	erpPubKeys = append(erpPubKeys, "03cd3e383ec6e12719a6c69515e5559bcbe037d0aa24c187e1e26ce932e22ad7b3")
	erpPubKeys = append(erpPubKeys, "02370a9838e4d15708ad14a104ee5606b36caaaaf739d833e67770ce9fd9b3ec80")

	return rootstock.FederationInfo{
		ActiveFedBlockHeight: 0,
		ErpKeys:              erpPubKeys,
		FedSize:              int64(len(keys)),
		FedThreshold:         int64(len(keys)/2 + 1),
		PubKeys:              keys,
		IrisActivationHeight: 0,
	}
}
