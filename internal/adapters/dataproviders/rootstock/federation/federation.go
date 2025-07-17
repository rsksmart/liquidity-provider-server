package federation

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

const flyoverPrefix byte = 0x20

func GetFedRedeemScript(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params) ([]byte, error) {
	var buf *bytes.Buffer
	var err error

	buf, err = getFedRedeemScript(fedInfo, btcParams)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getFedRedeemScript(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params) (*bytes.Buffer, error) {
	buf, err := GetErpRedeemScriptBuf(fedInfo, btcParams)
	if err != nil {
		return nil, err
	}

	err = ValidateRedeemScript(fedInfo, btcParams, buf.Bytes())
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func GetRedeemScriptBuf(fedInfo blockchain.FederationInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := addStdNToMScriptPart(fedInfo, sb)
	if err != nil {
		return nil, err
	}
	if addMultiSig {
		sb.AddOp(txscript.OP_CHECKMULTISIG)
	}

	sbuf, err := sb.Script()
	if err != nil {
		return nil, err
	}

	buf.Write(sbuf)
	return &buf, nil
}

func addStdNToMScriptPart(fedInfo blockchain.FederationInfo, builder *txscript.ScriptBuilder) error {
	builder.AddOp(getOpCodeFromInt(int(fedInfo.FedThreshold)))

	for _, pubKey := range fedInfo.PubKeys {
		pkBuffer, err := hex.DecodeString(pubKey)
		if err != nil {
			return err
		}
		builder.AddData(pkBuffer)
	}

	builder.AddOp(getOpCodeFromInt(int(fedInfo.FedSize)))
	return nil
}

func ValidateRedeemScript(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params, script []byte) error {
	segwitScript, err := bitcoin.ScriptToP2shP2wsh(script)
	if err != nil {
		return err
	}

	addr, err := btcutil.NewAddressScriptHash(segwitScript, &btcParams)
	if err != nil {
		return err
	}

	fedAddress, err := btcutil.DecodeAddress(fedInfo.FedAddress, &btcParams)
	if err != nil {
		return err
	}
	if !bytes.Equal(addr.ScriptAddress(), fedAddress.ScriptAddress()) {
		return errors.New("the generated redeem script does not match with the federation redeem script")
	}

	return nil
}

func GetErpRedeemScriptBuf(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params) (*bytes.Buffer, error) {
	erpRedeemScriptBuf, err := p2ms(fedInfo, false)
	if err != nil {
		return nil, err
	}

	redeemScriptBuf, err := GetRedeemScriptBuf(fedInfo, true)
	if err != nil {
		return nil, err
	}
	scriptsA := txscript.NewScriptBuilder()
	scriptsA.AddOp(txscript.OP_NOTIF)
	var erpRedeemScriptBuffer bytes.Buffer
	scrA, err := scriptsA.Script()
	if err != nil {
		return nil, err
	}
	erpRedeemScriptBuffer.Write(scrA)
	erpRedeemScriptBuffer.Write(redeemScriptBuf.Bytes())
	erpRedeemScriptBuffer.WriteByte(txscript.OP_ELSE)

	csv, err := hex.DecodeString(getCsvValueFromNetwork(btcParams))
	if err != nil {
		return nil, err
	}
	erpRedeemScriptBuffer.WriteByte(byte(len(csv)))
	erpRedeemScriptBuffer.Write(csv)
	erpRedeemScriptBuffer.WriteByte(txscript.OP_CHECKSEQUENCEVERIFY)
	erpRedeemScriptBuffer.WriteByte(txscript.OP_DROP)
	erpRedeemScriptBuffer.Write(erpRedeemScriptBuf.Bytes())
	erpRedeemScriptBuffer.WriteByte(txscript.OP_CHECKMULTISIG)
	erpRedeemScriptBuffer.WriteByte(txscript.OP_ENDIF)

	return &erpRedeemScriptBuffer, nil
}

func p2ms(fedInfo blockchain.FederationInfo, addMultiSig bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	sb := txscript.NewScriptBuilder()
	err := addErpNToMScriptPart(fedInfo, sb)
	if err != nil {
		return nil, err
	}

	if addMultiSig {
		sb.AddOp(txscript.OP_CHECKMULTISIG)
	}

	sbuf, err := sb.Script()
	if err != nil {
		return nil, err
	}

	buf.Write(sbuf)
	return &buf, nil
}

func addErpNToMScriptPart(fedInfo blockchain.FederationInfo, builder *txscript.ScriptBuilder) error {
	size := len(fedInfo.ErpKeys)
	minimum := size/2 + 1
	builder.AddOp(getOpCodeFromInt(minimum))

	for _, pubKey := range fedInfo.ErpKeys {
		pkBuffer, err := hex.DecodeString(pubKey)
		if err != nil {
			return err
		}
		builder.AddData(pkBuffer)
	}

	builder.AddOp(getOpCodeFromInt(len(fedInfo.ErpKeys)))
	return nil
}

func getFlyoverPrefix(derivationValue []byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(flyoverPrefix)
	buf.Write(derivationValue)
	buf.WriteByte(txscript.OP_DROP)
	return buf.Bytes()
}

func GetFlyoverRedeemScript(derivationValue []byte, fedRedeemScript []byte) []byte {
	var buf bytes.Buffer
	buf.Write(getFlyoverPrefix(derivationValue))
	buf.Write(fedRedeemScript)
	return buf.Bytes()
}

// nolint:cyclop
func getOpCodeFromInt(val int) byte {
	switch val {
	case 1:
		return txscript.OP_1
	case 2:
		return txscript.OP_2
	case 3:
		return txscript.OP_3
	case 4:
		return txscript.OP_4
	case 5:
		return txscript.OP_5
	case 6:
		return txscript.OP_6
	case 7:
		return txscript.OP_7
	case 8:
		return txscript.OP_8
	case 9:
		return txscript.OP_9
	case 10:
		return txscript.OP_10
	case 11:
		return txscript.OP_11
	case 12:
		return txscript.OP_12
	case 13:
		return txscript.OP_13
	case 14:
		return txscript.OP_14
	case 15:
		return txscript.OP_15
	default:
		return txscript.OP_16
	}
}

func getCsvValueFromNetwork(btcParams chaincfg.Params) string {
	switch btcParams.Name {
	case chaincfg.MainNetParams.Name:
		return "50CD00"
	case chaincfg.TestNet3Params.Name:
		return "50CD00"
	default:
		return "F401"
	}
}

func GetDerivationValueHash(args blockchain.FlyoverDerivationArgs) []byte {
	var buf bytes.Buffer
	buf.Write(args.QuoteHash)
	buf.Write(args.UserBtcRefundAddress)
	buf.Write(args.LbcAdress)
	buf.Write(args.LpBtcAddress)

	derivationValueHash := crypto.Keccak256(buf.Bytes())

	return derivationValueHash
}

func CalculateFlyoverDerivationAddress(fedInfo blockchain.FederationInfo, btcParams chaincfg.Params, fedRedeemScript, derivationValue []byte) (blockchain.FlyoverDerivation, error) {
	var err error
	var addressScriptHash *btcutil.AddressScriptHash

	if len(fedRedeemScript) == 0 {
		if fedRedeemScript, err = GetFedRedeemScript(fedInfo, btcParams); err != nil {
			return blockchain.FlyoverDerivation{}, fmt.Errorf("error generating fed redeem script: %w", err)
		}
	} else {
		if err = ValidateRedeemScript(fedInfo, btcParams, fedRedeemScript); err != nil {
			return blockchain.FlyoverDerivation{}, fmt.Errorf("error validating fed redeem script: %w", err)
		}
	}

	flyoverScript := GetFlyoverRedeemScript(derivationValue, fedRedeemScript)

	addressScriptHash, err = bitcoin.ScriptToAddressP2shP2wsh(flyoverScript, &btcParams)
	if err != nil {
		return blockchain.FlyoverDerivation{}, err
	}

	return blockchain.FlyoverDerivation{
		Address:      addressScriptHash.EncodeAddress(),
		RedeemScript: hex.EncodeToString(flyoverScript),
	}, nil
}
