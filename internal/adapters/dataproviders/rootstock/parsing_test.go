package rootstock_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

// nolint:funlen
func TestParseReceipt(t *testing.T) {
	const (
		txHash      = "0x2c78749f71d5f1e870397eaf401c3683ca752577a7a2fcbf64e3eb2091c817dd"
		blockHash   = "0x9ca0a12f5854335bd78c145c105543b724043277a71dc63fb7b8c5ed9760af02"
		blockNumber = 0x755fa7
		lbcAddress  = "0xAA9cAf1e3967600578727F975F283446A3Da6612"
	)
	to := common.HexToAddress(lbcAddress)
	txData := geth.LegacyTx{
		Nonce:    0x1d0,
		GasPrice: big.NewInt(0x18dbac0),
		Gas:      0x2625a0,
		To:       &to,
		Value:    big.NewInt(0),
		Data:     hexutil.MustDecode("0x2d67952c00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000003c00000000000000000000000000000000000000000000000000000000000000440000000000000000000000000000000000000000000000000000000000000058000000000000000000000000000000000000000000000000000000000000dc3870ae2bb1b3bb7aed09b42aad17f4ebd54fca6f3ca000000000000000000000000000000000000000000000000aa9caf1e3967600578727f975f283446a3da661200000000000000000000000082a06ebdb97776a2da4041df8f2b2ea8d32578520000000000000000000000000000000000000000000000000000000000000280000000000000000000000000ac31a4beedd7ec916b7a48a612230cb85c1aaf5600000000000000000000000000000000000000000000000000000000000002c000000000000000000000000000000000000000000000000000005af3107a4000000000000000000000000000000000000000000000000000000009184e72a000000000000000000000000000ac31a4beedd7ec916b7a48a612230cb85c1aaf560000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000520800000000000000000000000000000000000000000000000074546abbfd23123b00000000000000000000000000000000000000000000000014c2dc12e7fc000000000000000000000000000000000000000000000000000000000000685527e500000000000000000000000000000000000000000000000000000000000015180000000000000000000000000000000000000000000000000000000000001c200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007f723f560000000000000000000000000000000000000000000000000000000000000000150000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001500840098213fec4001cdc4a77cc3340f5bb83d9ed5000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041dc942ceb0e3bdd9815f53e536e70aaac2217fc422c9370a0f07fde7cdab39eed4729a933f1e90e6c8cbd7bbfaff15273de6f82d0f34516f99783f18bf558c7641c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000011601000000053d0ab910f144862c6c9d7049f3cf70a4fa8e3eb2185726f99173caaa8b0f00a10100000000ffffffffc4655c058fabcf3b0ed7a5186d81e04fc97f1ea4b5ffcca19694da67c839a2420000000000ffffffffd870ecc8e2a3cecef46d7244264678277840aa2f6d82a2dd973acf888e0ac8dd0100000000ffffffffcacaee423f27ad6cfc76fe2cfbd4bd20171919f6615e07cb27f2fcca21bb0ed80100000000ffffffff28721107c40671b831a30fbaf2e6e8776a0947229e4bf8fb41c9d47894d59aa30200000000ffffffff0247deea080000000017a9145be63a8b8f3f848981d31fa163a0c89f354a1c50878b7d000000000000160014f011a71128cb7e04fe511e868543b9f6cf81426000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000169070200000bcd88f6365a9cd71dbde59464d2c4b16b2080b631e61d62383a790dc37acfff6e6dfc3733f939bdb9a53b643139f7eddde0927c236b8b99e93fed7a1d27626abf880ca69c03569d61efa7ec9f1039d748212a73a7deeb8efa9654c79c002441867a140bd11767a41bf7c71ace188719e1d2c233fbdb64b4fb223b761daf311982ffffffc6fc50ff2cd0636b80d2d9ac06db90534cb03bfa347c631ecb40849fbade60d5e556b13f08131f047066d5d9040a1e0a3820d12bfae8a4cadd6c4d0400b9fa236eb57fd4cd2614e9c1f38f84e6d35f40800bb8076197596bc9be8db148428dbed1254a3369a99492c6c02bc99b84ca97dcf623cd88d2aa56dfc6f63fb3f3756234efce77e21d9c7d704b243b4ca842ffa1b7aede2d01edbc50717abad92f8b8b68234c8d6cb9084a2c0819f21e237781d341905cbe0d3ea47be8f4ce44d25a9f4860bf57d610f8638dd85202a12beaee7664d922e48184550501cd66c8036f2f000000000000000000000000000000000000000000000000"),
		V:        big.NewInt(0x60),
		R:        new(big.Int).SetBytes(hexutil.MustDecode("0xda4204b9a42bcdc0442d36fd87fee053dc121674e7aaf7639dcdb67631f755ac")),
		S:        new(big.Int).SetBytes(hexutil.MustDecode("0x192902edaeb78e096b555b5884b6c0a62d10809f00a28ad8883e0074a4755f7f")),
	}
	rawTx := geth.NewTx(&txData)
	rawReceipt := &geth.Receipt{
		Type:              0,
		Status:            1,
		CumulativeGasUsed: 0x2314b,
		Bloom:             geth.BytesToBloom(hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000080000000002000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000040000000000000000001000000000000000000000000000000000000100000040000000400000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000400000000")),
		Logs: []*geth.Log{
			{
				BlockNumber: blockNumber,
				BlockHash:   common.HexToHash(blockHash),
				TxHash:      common.HexToHash(txHash),
				Index:       0,
				Address:     common.HexToAddress(lbcAddress),
				Data:        hexutil.MustDecode("0x00000000000000000000000082a06ebdb97776a2da4041df8f2b2ea8d325785200000000000000000000000000000000000000000000000014c337856ab59600"),
				Topics: []common.Hash{
					common.HexToHash("0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53"),
				},
				TxIndex: 0,
				Removed: false,
			},
			{
				BlockNumber: blockNumber,
				BlockHash:   common.HexToHash(blockHash),
				TxHash:      common.HexToHash(txHash),
				Index:       1,
				Address:     common.HexToAddress(lbcAddress),
				Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000014c3378607043c00"),
				Topics: []common.Hash{
					common.HexToHash("0x0629ae9d1dc61501b0ca90670a9a9b88daaf7504b54537b53e1219de794c63d2"),
					common.HexToHash("0xc9dffd38b14139b2fb49ed67bc1a221d3f84e119a638bb424511ba4cd43e4c5e"),
				},
				TxIndex: 0,
				Removed: false,
			},
		},
		TxHash:            common.HexToHash("0x2c78749f71d5f1e870397eaf401c3683ca752577a7a2fcbf64e3eb2091c817dd"),
		GasUsed:           0x2314b,
		EffectiveGasPrice: big.NewInt(0x18dbac0),
		BlockHash:         common.HexToHash("0x9ca0a12f5854335bd78c145c105543b724043277a71dc63fb7b8c5ed9760af02"),
		BlockNumber:       big.NewInt(blockNumber),
		TransactionIndex:  0,
	}

	t.Run("should parse receipt correctly", func(t *testing.T) {
		receipt, err := rootstock.ParseReceipt(rawTx, rawReceipt)
		require.NoError(t, err)
		assert.Equal(t, blockchain.TransactionReceipt{
			TransactionHash:   txHash,
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			From:              "0x82a06eBDB97776a2da4041dF8f2b2ea8D3257852",
			To:                lbcAddress,
			CumulativeGasUsed: big.NewInt(143691),
			GasUsed:           big.NewInt(143691),
			Value:             entities.NewWei(0),
			Logs: []blockchain.TransactionLog{
				{
					Address: lbcAddress,
					Topics: [][32]byte{
						utils.To32Bytes(hexutil.MustDecode("0x42cfb81a915ac5a674852db250bf722637bee705a267633b68cab3a2dde06f53")),
					},
					Data:        hexutil.MustDecode("0x00000000000000000000000082a06ebdb97776a2da4041df8f2b2ea8d325785200000000000000000000000000000000000000000000000014c337856ab59600"),
					BlockNumber: blockNumber,
					TxHash:      txHash,
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       0,
					Removed:     false,
				},
				{
					Address: lbcAddress,
					Topics: [][32]byte{
						utils.To32Bytes(hexutil.MustDecode("0x0629ae9d1dc61501b0ca90670a9a9b88daaf7504b54537b53e1219de794c63d2")),
						utils.To32Bytes(hexutil.MustDecode("0xc9dffd38b14139b2fb49ed67bc1a221d3f84e119a638bb424511ba4cd43e4c5e")),
					},
					Data:        hexutil.MustDecode("0x00000000000000000000000000000000000000000000000014c3378607043c00"),
					BlockNumber: blockNumber,
					TxHash:      txHash,
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       1,
					Removed:     false,
				},
			},
		}, receipt)
	})
	t.Run("should return error if any parameter is nil", func(t *testing.T) {
		result, err := rootstock.ParseReceipt(nil, rawReceipt)
		require.Error(t, err)
		assert.Empty(t, result)

		result, err = rootstock.ParseReceipt(rawTx, nil)
		require.Error(t, err)
		assert.Empty(t, result)
	})
}

// nolint:funlen
func TestParseDepositEvent(t *testing.T) {
	const (
		txHash      = "0xfba869597b09185666429924ee5adc15289e131171c5018b353343e9783236a9"
		blockHash   = "0x450a4391a92630c83798f3814e047556b2129479ae95cce773c220884ea5e006"
		blockNumber = 7709148
		from        = "0xACa43E826BE4d5CbFf195797968A3fcf20cC7813"
		to          = "0xAa9caf1e3967600578727f975F283446a3dA6612"
	)
	var (
		amount = entities.NewWei(38805670000000000)
	)
	t.Run("should parse deposit correctly", func(t *testing.T) {
		receipt := blockchain.TransactionReceipt{
			TransactionHash:   txHash,
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			From:              from,
			To:                to,
			CumulativeGasUsed: big.NewInt(909304),
			GasUsed:           big.NewInt(410230),
			Value:             entities.NewWei(38805670000000000),
			Logs: []blockchain.TransactionLog{
				{
					Address: "0xAa9caf1e3967600578727f975F283446a3dA6612",
					Topics: [][32]byte{
						utils.To32Bytes(hexutil.MustDecode("0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f")),
						utils.To32Bytes(hexutil.MustDecode("0xeb8a4598a3cb0b8a697206316216b791e7b16dd5a8496349a6aad6fac8f190e7")),
						utils.To32Bytes(hexutil.MustDecode("0x000000000000000000000000aca43e826be4d5cbff195797968a3fcf20cc7813")),
						utils.To32Bytes(hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000685c4f0a")),
					},
					Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000089dd8d1f9efc00"),
					BlockNumber: blockNumber,
					TxHash:      txHash,
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       0,
					Removed:     false,
				},
			},
		}
		deposit, err := rootstock.ParseDepositEvent(receipt)
		require.NoError(t, err)
		assert.Equal(t, blockchain.ParsedLog[quote.PegoutDeposit]{
			Log: quote.PegoutDeposit{
				TxHash:      txHash,
				QuoteHash:   "eb8a4598a3cb0b8a697206316216b791e7b16dd5a8496349a6aad6fac8f190e7",
				Amount:      amount,
				Timestamp:   time.Unix(1750880010, 0),
				BlockNumber: blockNumber,
				From:        from,
			},
			RawLog: blockchain.TransactionLog{
				Address: "0xAa9caf1e3967600578727f975F283446a3dA6612",
				Topics: [][32]byte{
					utils.To32Bytes(hexutil.MustDecode("0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f")),
					utils.To32Bytes(hexutil.MustDecode("0xeb8a4598a3cb0b8a697206316216b791e7b16dd5a8496349a6aad6fac8f190e7")),
					utils.To32Bytes(hexutil.MustDecode("0x000000000000000000000000aca43e826be4d5cbff195797968a3fcf20cc7813")),
					utils.To32Bytes(hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000685c4f0a")),
				},
				Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000089dd8d1f9efc00"),
				BlockNumber: blockNumber,
				TxHash:      txHash,
				TxIndex:     0,
				BlockHash:   blockHash,
				Index:       0,
				Removed:     false,
			},
		}, deposit)
	})
	t.Run("should return error when log is not present", func(t *testing.T) {
		receipt := blockchain.TransactionReceipt{
			TransactionHash:   txHash,
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			From:              from,
			To:                to,
			CumulativeGasUsed: big.NewInt(909304),
			GasUsed:           big.NewInt(410230),
			Value:             entities.NewWei(38805670000000000),
			Logs:              []blockchain.TransactionLog{},
		}
		deposit, err := rootstock.ParseDepositEvent(receipt)
		require.ErrorContains(t, err, "deposit event not found in receipt logs")
		assert.Empty(t, deposit)
	})
	t.Run("should return error on malformed log topics", func(t *testing.T) {
		receipt := blockchain.TransactionReceipt{
			TransactionHash:   txHash,
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			From:              from,
			To:                to,
			CumulativeGasUsed: big.NewInt(909304),
			GasUsed:           big.NewInt(410230),
			Value:             entities.NewWei(38805670000000000),
			Logs: []blockchain.TransactionLog{
				{
					Address: "0xAa9caf1e3967600578727f975F283446a3dA6612",
					Topics: [][32]byte{
						utils.To32Bytes(hexutil.MustDecode("0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f")),
						utils.To32Bytes(hexutil.MustDecode("0x000000000000000000000000aca43e826be4d5cbff195797968a3fcf20cc7813")),
					},
					Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000089dd8d1f9efc0000000000000000000000000000000000000000000000000000000000685c4f0a"),
					BlockNumber: blockNumber,
					TxHash:      txHash,
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       0,
					Removed:     false,
				},
			},
		}
		deposit, err := rootstock.ParseDepositEvent(receipt)
		require.ErrorContains(t, err, "invalid number of topics for PegOutDeposit event")
		assert.Empty(t, deposit)
	})
	t.Run("should return error on malformed log data", func(t *testing.T) {
		receipt := blockchain.TransactionReceipt{
			TransactionHash:   txHash,
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			From:              from,
			To:                to,
			CumulativeGasUsed: big.NewInt(909304),
			GasUsed:           big.NewInt(410230),
			Value:             entities.NewWei(38805670000000000),
			Logs: []blockchain.TransactionLog{
				{
					Address: "0xAa9caf1e3967600578727f975F283446a3dA6612",
					Topics: [][32]byte{
						utils.To32Bytes(hexutil.MustDecode("0xb1bc7bfc0dab19777eb03aa0a5643378fc9f186c8fc5a36620d21136fbea570f")),
						utils.To32Bytes(hexutil.MustDecode("0xeb8a4598a3cb0b8a697206316216b791e7b16dd5a8496349a6aad6fac8f190e7")),
						utils.To32Bytes(hexutil.MustDecode("0x000000000000000000000000aca43e826be4d5cbff195797968a3fcf20cc7813")),
					},
					Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000089dd8d1f9efc"),
					BlockNumber: blockNumber,
					TxHash:      txHash,
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       0,
					Removed:     false,
				},
			},
		}
		deposit, err := rootstock.ParseDepositEvent(receipt)
		require.Error(t, err)
		assert.Empty(t, deposit)
	})
}
