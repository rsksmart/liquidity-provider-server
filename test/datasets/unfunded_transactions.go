package datasets

type UnfundedTransactionTestCase struct {
	Description        string
	Address            string
	ValueSatoshis      int64
	OpReturnContent    []byte
	RawTxPaymentOnly   string // Transaction with only payment output (from CreateRawTransaction)
	ExpectedCompleteTx string // Complete transaction with payment + OP_RETURN outputs
}

var UnfundedTxTestCases = []UnfundedTransactionTestCase{
	{
		Description:        "P2PKH address, 0.6 BTC, quote hash",
		Address:            "mxqk28jvEtvjxRN8k7W9hFEJfWz5VcUgHW",
		ValueSatoshis:      60000000,
		OpReturnContent:    []byte{0xa1, 0xb2, 0xc3, 0xd4, 0xe5, 0xf6, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22},
		RawTxPaymentOnly:   "02000000000100879303000000001976a914be07cb9dfdc7dfa88436fa4128410e2126d6979688ac00000000",
		ExpectedCompleteTx: "02000000000200879303000000001976a914be07cb9dfdc7dfa88436fa4128410e2126d6979688ac0000000000000000226a20a1b2c3d4e5f60123456789abcdef112233445566778899aabbccddeeff00112200000000",
	},
	{
		Description:        "P2SH address, 0.5 BTC, quote hash",
		Address:            "2N4DTeBWDF9yaF9TJVGcgcZDM7EQtsGwFjX",
		ValueSatoshis:      50000000,
		OpReturnContent:    []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11},
		RawTxPaymentOnly:   "02000000000180f0fa020000000017a9147853f2f139767d6548f38193afbdc136bfc9a9628700000000",
		ExpectedCompleteTx: "02000000000280f0fa020000000017a9147853f2f139767d6548f38193afbdc136bfc9a962870000000000000000226a201122334455667788aabbccddeeff00112233445566778899aabbccddeeff001100000000",
	},
	{
		Description:        "P2WPKH address, 0.3 BTC, quote hash",
		Address:            "tb1qlh84gv84mf7e28lsk3m75sgy7rx2lmvpr77rmw",
		ValueSatoshis:      30000000,
		OpReturnContent:    []byte{0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89},
		RawTxPaymentOnly:   "02000000000180c3c90100000000160014fdcf5430f5da7d951ff0b477ea4104f0ccafed8100000000",
		ExpectedCompleteTx: "02000000000280c3c90100000000160014fdcf5430f5da7d951ff0b477ea4104f0ccafed810000000000000000226a20abcdef0123456789abcdef0123456789abcdef0123456789abcdef012345678900000000",
	},
	{
		Description:        "P2WSH address, 0.7 BTC, quote hash",
		Address:            "tb1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3q0sl5k7",
		ValueSatoshis:      70000000,
		OpReturnContent:    []byte{0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11},
		RawTxPaymentOnly:   "020000000001801d2c04000000002200201863143c14c5166804bd19203356da136c985678cd4d27a1b8c632960490326200000000",
		ExpectedCompleteTx: "020000000002801d2c04000000002200201863143c14c5166804bd19203356da136c985678cd4d27a1b8c63296049032620000000000000000226a202233445566778899aabbccddeeff00112233445566778899aabbccddeeff001100000000",
	},
	{
		Description:        "P2TR address, 0.4 BTC, quote hash",
		Address:            "tb1ptt2hnzgzfhrfdyfz02l02wam6exd0mzuunfdgqg3ttt9yagp6daslx6grp",
		ValueSatoshis:      40000000,
		OpReturnContent:    []byte{0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		RawTxPaymentOnly:   "020000000001005a6202000000002251205ad57989024dc69691227abef53bbbd64cd7ec5ce4d2d401115ad6527501d37b00000000",
		ExpectedCompleteTx: "020000000002005a6202000000002251205ad57989024dc69691227abef53bbbd64cd7ec5ce4d2d401115ad6527501d37b0000000000000000226a2099887766554433221100ffeeddccbbaa99887766554433221100aabbccddeeff00000000",
	},
}
