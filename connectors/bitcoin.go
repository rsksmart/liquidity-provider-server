package connectors

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/bloom"
)

// Computes the serialized partial merkle tree of a transaction in a block
// the format of the serialized PMT is:
//
// - uint32     total_transactions (4 bytes)
// - varint     number of hashes   (1-3 bytes)
// - uint256[]  hashes in depth-first order (<= 32*N bytes)
// - varint     number of bytes of flag bits (1-3 bytes)
// - byte[]     flag bits, packed per 8 in a byte, least significant bit first (<= 2*N-1 bits)
//
func SerializePMT(h string, block *btcutil.Block) ([]byte, error) {
	filter := bloom.NewFilter(1, 0, 0, wire.BloomUpdateAll)
	hash, err := chainhash.NewHashFromStr(h)
	if err != nil {
		return nil, fmt.Errorf("error parsing hash: %v", err)
	}
	filter.AddHash(hash)
	msg, indices := bloom.NewMerkleBlock(block, filter)
	if len(indices) > 1 {
		return nil, fmt.Errorf("block matches more than one transaction (%v)", len(indices))
	}

	var buf bytes.Buffer

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(block.Transactions())))
	buf.Write(b)

	wire.WriteVarInt(&buf, wire.ProtocolVersion, uint64(len(msg.Hashes)))

	for _, h := range msg.Hashes {
		buf.Write(h[:])
	}
	wire.WriteVarInt(&buf, wire.ProtocolVersion, uint64(len(msg.Flags)))

	buf.Write(msg.Flags)

	return buf.Bytes(), nil
}
