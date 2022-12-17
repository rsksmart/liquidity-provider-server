package connectors

import (
	"encoding/hex"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/rsksmart/liquidity-provider-server/connectors/testmocks"
	"github.com/stretchr/testify/mock"

	"github.com/btcsuite/btcutil"

	"github.com/stretchr/testify/assert"
)

type testPmt struct {
	h   string
	pmt string
}

var (
	expectedPmts = [2]testPmt{
		{
			h: "07f8b22fa9a3b32e20b59bb90727de05fb634749519ebcb6a887aeaf2c7eb041",
			pmt: "f3080000" + "0d" +
				"41b07e2cafae87a8b6bc9e51494763fb05de2707b99bb5202eb3a3a92fb2f807" +
				"731c671fafb5d234834c726657f29c9af030ccf7068f1ef732af4efd8e146da0" +
				"a9d6075f4758821ceeef2c230cfd2497df2d1d1d02dd19e653d22b3dc271b393" +
				"94c2d2e51eae99af800c39575a11f3b4eb0fdf5d5deff0b9f5ff592566f4f173" +
				"2b3f6d41bded4344899a348d3053ccd68e922626e589da71d9a583ccfe9e3be6" +
				"393a24e14d18006b54a967963c56da58b18ce770cdb3b32e56d88c138c473a1a" +
				"37acc29a7788a88404fb9a05c416c2ad8f340a61c1ea331528a91ae6210db4f0" +
				"e22d21c55b3f6806387ca34aba54522a7fe15e593ab0d0ff89c6d826cf4e7455" +
				"99bdfed66a4ef1a366f056363158b2907388a5bd4013643d83a016469d392aab" +
				"2319910a0ac801d2c9793c661f1bb02863f672288ce3b4c5e365cff81932fbe4" +
				"3d1eba4b027f80b1240ab5b8c677167602ee63c5a6dad213777b6fbbd3dd9778" +
				"71f30b975d1a8f6cd62535985ea4d11f5c7f80549eab1b18a5cc011872b5403d" +
				"ee666302831a3c64d1604c5c0bec9c796d8dcace974ae97e5837ff0d446d060c" +
				"04ff1f0000",
		},
		{
			h: "ddf5061f9707f0c959bf24278d557b264716672c1b601ec50112d6dfe160d9d3",
			pmt: "f3080000" + "0d" +
				"c0746a357444e9948a18a612e02df5a99240e77f1ff75dd949d5b4038dcf3667" +
				"3a03c716cf722cff7d264c763088ceeb1665f26c6fdd5835d841eeee2f3ece4a" +
				"203a24db8b7a51e4ab0e35a6b4151f6d7f1eef96f32e4fceaac6127521911618" +
				"6efb7fdb763e821f99bd6af8d044cc6feadd7b4716e6938335a3e08548f5a077" +
				"5dd364971faab5cd089cd1fa713e8be658a67a704d39952218f6518e5045d269" +
				"d3d960e1dfd61201c51e601b2c671647267b558d2724bf59c9f007971f06f5dd" +
				"0eab2677f52c996a3f941bef3ec57ebdf22429c37dee5ae68892df30f8acfc22" +
				"5c6fed56bdff34686135e68fda4b716713e60258b6971c03091f25115c008eec" +
				"48a828c75ad7340fadbc368636b4014f6e8386c3990a35620cbddca933a72b02" +
				"d990fb8a602fcda9e1e41120c25f4981362a9dfc7f7ed1f5188482b8ee3f532f" +
				"0ee6234e44af99351ee430f4ac0fa7b71fe9c601c78480b9a97fea305d3abca2" +
				"35a6668846093803e07c48dc9a75be90ed6edd4debb0b7b49bc057e093ad395e" +
				"ee666302831a3c64d1604c5c0bec9c796d8dcace974ae97e5837ff0d446d060c" +
				"04dbd50300",
		},
	}
)

func testSerializeTx(t *testing.T) {
	expected := "01000000010000000000000000000000000000000000000000000000000000000000000000" +
		"ffffffff5f034aa00a1c2f5669614254432f4d696e656420627920797a33313936303538372f2cfabe" +
		"6d6dc0a3751203a336deb817199448996ebcb2a0e537b1ce9254fa3e9c3295ca196b10000000000000" +
		"0010c56e6700d262d24bd851bb829f9f0000ffffffff0401b3cc25000000001976a914536ffa992491" +
		"508dca0354e52f32a3a7a679a53a88ac00000000000000002b6a2952534b424c4f434b3a040c866ad2" +
		"fdb8b59b32dd17059edaeef11d295e279a74ab97125d2500371ce90000000000000000266a24b9e11b" +
		"6dab3e2ca50c1a6b01cf80eccb9d291aab8b095d653e348aa9d94a73964ff5cf1b0000000000000000" +
		"266a24aa21a9ed04f0bac0104f4fa47bec8058f2ebddd292dd85027ab0d6d95288d31f12c5a4b800000000"

	// this is block 0000000000000000000aca0460feaf0661f173b75d4cc824b57233aa7c6b7bc3
	f, err := os.Open("./testdata/test_block")
	if err != nil {
		t.Fatalf("error opening test block file: %v", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("error reading test file: %v", err)
	}
	s := string(b)
	h, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("error decoding test file: %v", err)
	}

	block, err := btcutil.NewBlockFromBytes(h)
	if err != nil {
		t.Fatalf("error parsing test block: %v", err)
	}
	tx, err := block.Tx(0)
	if err != nil {
		t.Fatalf("error reading transaction from test block: %v", err)
	}
	result, err := serializeTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(result) != expected {
		t.Errorf("serialized tx does not match expected: %x \n----\n %v", result, expected)
	}
}

func testPMTSerialization(t *testing.T) {
	// this is block 0000000000000000000aca0460feaf0661f173b75d4cc824b57233aa7c6b7bc3
	f, err := os.Open("./testdata/test_block")
	if err != nil {
		t.Fatalf("error opening test block file: %v", err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("error reading test file: %v", err)
	}
	s := string(b)
	h, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("error decoding test file: %v", err)
	}

	block, err := btcutil.NewBlockFromBytes(h)
	if err != nil {
		t.Fatalf("error parsing test block: %v", err)
	}

	for _, p := range expectedPmts {
		serializedPMT, err := serializePMT(p.h, block)
		if err != nil {
			t.Errorf("error serializing PMT:\n %v", p.h)
			continue
		}
		result := hex.EncodeToString(serializedPMT)
		if result != p.pmt {
			t.Errorf("expected PMT:\n%v\n is different from serialized PMT:\n%v\n", p.pmt, result)
		}
	}
}

func testCheckBtcAddr(t *testing.T) {
	btcClientMock := new(testmocks.BTCClientMock)
	addrWatcherMock := new(testmocks.AddressWatcherMock)
	amountInBtc := float64(1)
	amount, err := btcutil.NewAmount(amountInBtc)
	assert.Nil(t, err)
	var confirmations int64

	btc, err := NewBTC("mainnet")
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}
	btc.c = btcClientMock

	btcAddr, err := btcutil.DecodeAddress("38r8PQdgw5vdebE9h12Eum6saVnWEXxbve", &btc.params)
	if err != nil {
		t.Fatalf("error initializing BTC: %v", err)
	}

	// check error when retrieving unspent outputs for address
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{}, errors.New("ListUnspentMinMaxAddresses failed")).Once()
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, btcutil.Amount(0), time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(0, 0) })
	assert.NotNil(t, err)
	assert.EqualValues(t, "error retrieving unspent outputs for address 38r8PQdgw5vdebE9h12Eum6saVnWEXxbve: ListUnspentMinMaxAddresses failed", err.Error())
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check happy flow
	confirmations = 0
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{{TxID: "0xabc", Confirmations: 1, Amount: amountInBtc}}, nil).Once()
	addrWatcherMock.On("OnNewConfirmation", "0xabc", int64(1), amount).Once()
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(0, 0) })
	assert.Nil(t, err)
	assert.EqualValues(t, 1, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check happy flow #2: case when agreed amount has been deposited in the second tx (with two UTXOs)
	confirmations = 0
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{
		{TxID: "0xabc", Confirmations: 1, Amount: float64(0.98)},
		{TxID: "0xdef", Confirmations: 1, Amount: float64(0.4)}, // \
		{TxID: "0xdef", Confirmations: 1, Amount: float64(0.6)}, // -- these two txs with hash 0xdef are going to be selected
		{TxID: "0xghi", Confirmations: 1, Amount: float64(0.99)},
		{TxID: "0xjkl", Confirmations: 1, Amount: float64(1.1)},
	}, nil).Once()
	addrWatcherMock.On("OnNewConfirmation", "0xdef", int64(1), amount).Once()
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(0, 0) })
	assert.Nil(t, err)
	assert.EqualValues(t, 1, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check case when time for depositing has elapsed
	confirmations = 0
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{{TxID: "0xabc", Confirmations: 1, Amount: float64(0.98)}}, nil).Once()
	addrWatcherMock.On("OnExpire").Once()
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(1, 0) })
	assert.NotNil(t, err)
	assert.EqualValues(t, "time for depositing 1 BTC has elapsed; addr: 38r8PQdgw5vdebE9h12Eum6saVnWEXxbve", err.Error())
	assert.EqualValues(t, 0, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check case when time for depositing has elapsed, but agreed amount has been deposited
	confirmations = 0
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{{TxID: "0xabc", Confirmations: 1, Amount: amountInBtc}}, nil).Times(1)
	addrWatcherMock.On("OnNewConfirmation", "0xabc", int64(1), amount).Times(0)
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(1, 0) })
	assert.Nil(t, err)
	assert.EqualValues(t, 1, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)

	// check case when number of confirmations has not advanced after the previous check
	confirmations = 1
	btcClientMock.On("ListUnspentMinMaxAddresses", 0, 9999, mock.AnythingOfType("[]btcutil.Address")).Return([]btcjson.ListUnspentResult{{TxID: "0xabc", Confirmations: 1, Amount: amountInBtc}}, nil).Times(1)
	addrWatcherMock.On("OnNewConfirmation", "0xabc", int64(1), amount).Times(0)
	err = btc.checkBtcAddr(addrWatcherMock, btcAddr, amount, time.Unix(0, 0), &confirmations, func() time.Time { return time.Unix(0, 0) })
	assert.NotNil(t, err)
	assert.EqualValues(t, "num of confirmations has not advanced; conf: 1", err.Error())
	assert.EqualValues(t, 1, confirmations)
	btcClientMock.AssertExpectations(t)
	addrWatcherMock.AssertExpectations(t)
}

func TestBitcoinConnector(t *testing.T) {
	t.Run("test pmt serialization", testPMTSerialization)
	t.Run("test tx serialization", testSerializeTx)
	t.Run("test check btc addr", testCheckBtcAddr)
}
