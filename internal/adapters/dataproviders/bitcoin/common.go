package bitcoin

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/bloom"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
	"math/big"
)

const (
	BtcToSatoshi = 100000000
)

type Connection struct {
	NetworkParams *chaincfg.Params
	client        *rpcclient.Client
}

func NewConnection(networkParams *chaincfg.Params, client *rpcclient.Client) *Connection {
	return &Connection{NetworkParams: networkParams, client: client}
}

func (c *Connection) Shutdown(endChannel chan<- bool) {
	c.client.Disconnect()
	endChannel <- true
	log.Debug("Disconnected from BTC node")
}

func (c *Connection) CheckConnection(ctx context.Context) bool {
	err := c.client.Ping()
	if err != nil {
		log.Error("Error checking BTC node connection: ", err)
	}
	return err == nil
}

func DecodeAddressBase58(address string, keepVersion bool) ([]byte, error) {
	var buff bytes.Buffer
	addressBytes, version, err := base58.CheckDecode(address)
	if err != nil {
		return nil, err
	} else if len(addressBytes) != 20 {
		return nil, fmt.Errorf("decoded address exceeds 20 bytes (%s)", address)
	}
	if keepVersion {
		buff.WriteByte(version)
	}
	buff.Write(addressBytes)
	return buff.Bytes(), nil
}

func DecodeAddressBase58OnlyLegacy(address string, keepVersion bool) ([]byte, error) {
	if !blockchain.IsLegacyBtcAddress(address) {
		return nil, fmt.Errorf("only legacy address allowed (%s)", address)
	}
	return DecodeAddressBase58(address, keepVersion)
}

func toSwappedBytes32(hash *chainhash.Hash) [32]byte {
	var result [32]byte
	for i := 0; i < chainhash.HashSize/2; i++ {
		result[i], result[chainhash.HashSize-1-i] = hash[chainhash.HashSize-1-i], hash[i]
	}
	return result
}

func buildMerkleBranch(merkleTree []*chainhash.Hash, txCount uint32, txIndex uint32) blockchain.MerkleBranch {
	hashes := make([][32]byte, 0)
	path := 0
	pathIndex := 0
	var levelOffset uint32 = 0
	currentNodeOffset := txIndex

	for levelSize := txCount; levelSize > 1; levelSize = (levelSize + 1) / 2 {
		var targetOffset uint32
		if currentNodeOffset%2 == 0 {
			// Target is left hand side, use right hand side
			targetOffset = min(currentNodeOffset+1, levelSize-1)
		} else {
			// Target is right hand side, use left hand side
			targetOffset = currentNodeOffset - 1
			path = path + (1 << pathIndex)
		}
		hashes = append(hashes, toSwappedBytes32(merkleTree[levelOffset+targetOffset]))

		levelOffset += levelSize
		currentNodeOffset = currentNodeOffset / 2
		pathIndex++
	}

	return blockchain.MerkleBranch{
		Hashes: hashes,
		Path:   big.NewInt(int64(path)),
	}
}

func serializePartialMerkleTree(txHash *chainhash.Hash, block *btcutil.Block) ([]byte, error) {
	var err error
	filter := bloom.NewFilter(1, 0, 0, wire.BloomUpdateAll)
	filter.AddHash(txHash)

	msg, indices := bloom.NewMerkleBlock(block, filter)
	if len(indices) > 1 {
		return nil, fmt.Errorf("block matches more than one transaction (%v)", len(indices))
	}

	var buf bytes.Buffer
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(block.Transactions())))
	buf.Write(b)

	err = wire.WriteVarInt(&buf, wire.ProtocolVersion, uint64(len(msg.Hashes)))
	if err != nil {
		return nil, err
	}

	for _, h := range msg.Hashes {
		buf.Write(h[:])
	}
	err = wire.WriteVarInt(&buf, wire.ProtocolVersion, uint64(len(msg.Flags)))
	if err != nil {
		return nil, err
	}

	buf.Write(msg.Flags)
	return buf.Bytes(), nil
}
