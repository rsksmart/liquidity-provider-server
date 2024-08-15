package bitcoin_test

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type decodedAddress struct {
	address  string
	expected []byte
}

var base58Addresses = []decodedAddress{
	{"n1BE7ioVukYS2GC88hT2K6cUvRiKwMwio7", []byte{111, 215, 167, 103, 99, 62, 208, 72, 131, 184, 122, 185, 112, 220, 93, 130, 94, 43, 74, 67, 67}},
	{"n2aSettzgmgwxMoaNQfNL58BqQtEAqCRp1", []byte{111, 231, 3, 144, 1, 201, 195, 77, 230, 3, 157, 153, 171, 189, 171, 215, 212, 180, 203, 182, 41}},
	{"mugStBVzw3ZRafzrYqzSrq5Qq9HH9cZJst", []byte{111, 155, 93, 94, 120, 235, 41, 106, 157, 161, 118, 38, 245, 188, 189, 62, 16, 242, 77, 50, 225}},
	{"mfng2KVQdxd1wmSEk5TRKVLwMarsX9sD7D", []byte{111, 2, 249, 18, 163, 98, 241, 57, 132, 12, 133, 190, 167, 51, 135, 39, 8, 238, 37, 238, 22}},
	{"miba3yyWiomi5HMXboKZ4YmEnj1T6m3RLx", []byte{111, 33, 199, 227, 213, 45, 69, 31, 199, 87, 161, 249, 199, 169, 25, 147, 105, 219, 96, 219, 46}},
	{"1Ab1Jfe6xQHzL8RHoHDukDQBEks35KFWHC", []byte{0, 105, 39, 135, 109, 109, 187, 191, 255, 241, 12, 238, 207, 1, 71, 119, 140, 83, 84, 176, 148}},
	{"19hiJTQpZyT3C7Hu29dJE2YYToCeKp6cGu", []byte{0, 95, 116, 30, 252, 161, 93, 178, 87, 245, 15, 180, 19, 235, 18, 247, 75, 185, 227, 86, 164}},
	{"1LQdpgVCY2nYzsoRNRHWhuCLxMpYzb6zzg", []byte{0, 212, 226, 177, 205, 89, 105, 145, 127, 135, 118, 204, 8, 84, 231, 158, 29, 254, 239, 62, 150}},
	{"17oKLsbZsd2BZdCDn1dbrbk2TT9HSzw2aM", []byte{0, 74, 147, 51, 90, 109, 233, 79, 47, 123, 200, 20, 172, 19, 242, 57, 107, 24, 145, 19, 190}},
	{"1NbJonAytRKfCFkvGcQNEUCXAFnf17bYQG", []byte{0, 236, 215, 165, 247, 66, 139, 143, 251, 27, 232, 167, 211, 106, 236, 65, 215, 144, 67, 50, 188}},
	{"1MwygkmvJHwwG934EbtkjhRUFyfMHLEPi9", []byte{0, 229, 200, 75, 213, 72, 142, 34, 190, 217, 45, 198, 157, 186, 72, 245, 47, 208, 252, 113, 119}},
	{"2MyKtiQcyAgQQdPDmBJcT4UMn6jyVdKVwxg", []byte{196, 66, 178, 190, 224, 218, 119, 248, 26, 82, 40, 48, 190, 49, 159, 125, 114, 243, 42, 175, 30}},
	{"2N9TujzoGXb4LkzCQscMuuRSwaTVPhFfYSy", []byte{196, 177, 232, 48, 84, 2, 170, 51, 48, 251, 57, 152, 125, 51, 230, 130, 144, 98, 27, 231, 156}},
	{"2NAvNk9t8fKwFRhMJGgY2wMLRboN8DHXeBB", []byte{196, 193, 225, 171, 152, 53, 142, 203, 39, 152, 199, 170, 199, 208, 233, 205, 65, 249, 239, 208, 91}},
	{"2MzqFJDg3QcbzeWX87XpxRpHZm6SSNoGdoF", []byte{196, 83, 56, 29, 183, 70, 201, 79, 21, 103, 129, 237, 73, 8, 214, 82, 67, 212, 124, 164, 213}},
	{"2MuBrh282woZkhycbpKAZ8zEptTAEtRSM62", []byte{196, 21, 77, 57, 173, 29, 183, 149, 178, 0, 34, 217, 55, 64, 184, 30, 90, 13, 91, 214, 241}},
	{"33BTakydSPnJfSfR13foniEsPCB2nuHiCb", []byte{5, 16, 89, 49, 17, 100, 144, 216, 162, 233, 218, 176, 118, 166, 174, 244, 82, 223, 195, 106, 191}},
	{"3Fi3ywSD7eBEtnoUuiW7zFSBmpATd3YDLs", []byte{5, 153, 195, 229, 51, 197, 179, 77, 218, 70, 180, 220, 57, 200, 231, 161, 90, 119, 51, 54, 31}},
	{"3DCT2YtzwZZYdr3pPEhqVjA1Amak89YrHf", []byte{5, 126, 58, 102, 228, 17, 153, 252, 64, 101, 122, 39, 140, 215, 39, 222, 50, 192, 152, 23, 248}},
	{"3EuQdJ651cN7Cv9jJk2EJPdRhKT9JJFpt8", []byte{5, 144, 241, 147, 242, 132, 35, 204, 174, 111, 172, 89, 117, 131, 75, 45, 4, 128, 176, 86, 42}},
	{"3BcLTSd24JRtJhLcKqkeF83rFmFxxY5qH9", []byte{5, 108, 206, 165, 252, 49, 214, 45, 206, 182, 122, 231, 101, 130, 27, 245, 51, 13, 14, 60, 174}},
}

var bech32Addresses = []decodedAddress{
	{"tb1q22cm3qarlpj3gnf5h03kpdhaftdvf98q58dp75", []byte{0, 10, 10, 24, 27, 17, 0, 29, 3, 31, 1, 18, 17, 8, 19, 9, 20, 23, 15, 17, 22, 1, 13, 23, 29, 9, 11, 13, 12, 9, 5, 7, 0}},
	{"tb1qja7532egus56jkjnu6xgf9nh96q9up7gq5473m", []byte{0, 18, 29, 30, 20, 17, 10, 25, 8, 28, 16, 20, 26, 18, 22, 18, 19, 28, 26, 6, 8, 9, 5, 19, 23, 5, 26, 0, 5, 28, 1, 30, 8}},
	{"tb1qug3kle73ze6wcstdc4wunkjxapqnaeetprqjql", []byte{0, 28, 8, 17, 22, 31, 25, 30, 17, 2, 25, 26, 14, 24, 16, 11, 13, 24, 21, 14, 28, 19, 22, 18, 6, 29, 1, 0, 19, 29, 25, 25, 11}},
	{"tb1q66e97gspk233et7k24334zm2femvf5tpsq8ggm", []byte{0, 26, 26, 25, 5, 30, 8, 16, 1, 22, 10, 17, 17, 25, 11, 30, 22, 10, 21, 17, 17, 21, 2, 27, 10, 9, 25, 27, 12, 9, 20, 11, 1}},
	{"tb1qjtl57d37ccadme30hv3jhytt9gc9p4dq9zrz49", []byte{0, 18, 11, 31, 20, 30, 13, 17, 30, 24, 24, 29, 13, 27, 25, 17, 15, 23, 12, 17, 18, 23, 4, 11, 11, 5, 8, 24, 5, 1, 21, 13, 0}},
	{"bc1qzulaxy8fmvk8a92sec8s8u0xcqwcxw4fx037d8", []byte{0, 2, 28, 31, 29, 6, 4, 7, 9, 27, 12, 22, 7, 29, 5, 10, 16, 25, 24, 7, 16, 7, 28, 15, 6, 24, 0, 14, 24, 6, 14, 21, 9}},
	{"bc1q8sr9tv9ng4yd8s6s9eenfs7mh24jv64vnwzl0p", []byte{0, 7, 16, 3, 5, 11, 12, 5, 19, 8, 21, 4, 13, 7, 16, 26, 16, 5, 25, 25, 19, 9, 16, 30, 27, 23, 10, 21, 18, 12, 26, 21, 12}},
	{"bc1q5pfzfxmtx3kn7j8wqwe6336tmg0n5lmpqss9kx", []byte{0, 20, 1, 9, 2, 9, 6, 27, 11, 6, 17, 22, 19, 30, 18, 7, 14, 0, 14, 25, 26, 17, 17, 26, 11, 27, 8, 15, 19, 20, 31, 27, 1}},
	{"bc1qgq506g46u2dnua70k3dypu6r7xu3kfqeee3c38", []byte{0, 8, 0, 20, 15, 26, 8, 21, 26, 28, 10, 13, 19, 28, 29, 30, 15, 22, 17, 13, 4, 1, 28, 26, 3, 30, 6, 28, 17, 22, 9, 0, 25}},
	{"bc1qk2r5qt94fluyehjhr6neka0agpxung28pndjly", []byte{0, 22, 10, 3, 20, 0, 11, 5, 21, 9, 31, 28, 4, 25, 23, 18, 23, 3, 26, 19, 25, 22, 29, 15, 29, 8, 1, 6, 28, 19, 8, 10, 7}},
	{"tb1qzda4qlkdpjgmwxt9zr29pphhzqf2ku09p7dj33qyugqn80kg5muq8x0wyv", []byte{0, 2, 13, 29, 21, 0, 31, 22, 13, 1, 18, 8, 27, 14, 6, 11, 5, 2, 3, 10, 5, 1, 1, 23, 23, 2, 0, 9, 10, 22, 28, 15, 5, 1, 30, 13, 18, 17, 17, 0, 4, 28, 8, 0, 19, 7, 15, 22, 8, 20, 27, 28, 0}},
	{"tb1qgpgtqj68zwsdz7xmvqxxxaan7dcfgu76jz0cfzynqgrtvdsxlyqsf7dfz8", []byte{0, 8, 1, 8, 11, 0, 18, 26, 7, 2, 14, 16, 13, 2, 30, 6, 27, 12, 0, 6, 6, 6, 29, 29, 19, 30, 13, 24, 9, 8, 28, 30, 26, 18, 2, 15, 24, 9, 2, 4, 19, 0, 8, 3, 11, 12, 13, 16, 6, 31, 4, 0, 16}},
	{"tb1qkp4lxc09e34cc5vw383j42rgacurp7wrpnwjmvazv6g23c2ydz3qx5tfhl", []byte{0, 22, 1, 21, 31, 6, 24, 15, 5, 25, 17, 21, 24, 24, 20, 12, 14, 17, 7, 17, 18, 21, 10, 3, 8, 29, 24, 28, 3, 1, 30, 14, 3, 1, 19, 14, 18, 27, 12, 29, 2, 12, 26, 8, 10, 17, 24, 10, 4, 13, 2, 17, 0}},
	{"tb1qzhu8fjgw5aaqgv0q2jey4dnc3pgcr4cks858d6eaf97ljxywe70qwwsdku", []byte{0, 2, 23, 28, 7, 9, 18, 8, 14, 20, 29, 29, 0, 8, 12, 15, 0, 10, 18, 25, 4, 21, 13, 19, 24, 17, 1, 8, 24, 3, 21, 24, 22, 16, 7, 20, 7, 13, 26, 25, 29, 9, 5, 30, 31, 18, 6, 4, 14, 25, 30, 15, 0}},
	{"tb1qzda4qlkdpjgmwxt9zr29pphhzqf2ku09p7dj33qyugqn80kg5muq8x0wyv", []byte{0, 2, 13, 29, 21, 0, 31, 22, 13, 1, 18, 8, 27, 14, 6, 11, 5, 2, 3, 10, 5, 1, 1, 23, 23, 2, 0, 9, 10, 22, 28, 15, 5, 1, 30, 13, 18, 17, 17, 0, 4, 28, 8, 0, 19, 7, 15, 22, 8, 20, 27, 28, 0}},
	{"bc1qhnumvtg3c9xj2q7jmt8xnk4p5kmk52ffqwax8crfn4hqtry6qseq8vahua", []byte{0, 23, 19, 28, 27, 12, 11, 8, 17, 24, 5, 6, 18, 10, 0, 30, 18, 27, 11, 7, 6, 19, 22, 21, 1, 20, 22, 27, 22, 20, 10, 9, 9, 0, 14, 29, 6, 7, 24, 3, 9, 19, 21, 23, 0, 11, 3, 4, 26, 0, 16, 25, 0}},
	{"bc1qv47nn097m6hujqadw6kgt5hsk9h06k7tgq05empl3nn3mska8cfqpkjl36", []byte{0, 12, 21, 30, 19, 19, 15, 5, 30, 27, 26, 23, 28, 18, 0, 29, 13, 14, 26, 22, 8, 11, 20, 23, 16, 22, 5, 23, 15, 26, 22, 30, 11, 8, 0, 15, 20, 25, 27, 1, 31, 17, 19, 19, 17, 27, 16, 22, 29, 7, 24, 9, 0}},
	{"bc1qj8pqhwkv0k6h2tm3wtqu793njkvfd66dva04zldpdcey4sak5h3qx3n8nz", []byte{0, 18, 7, 1, 0, 23, 14, 22, 12, 15, 22, 26, 23, 10, 11, 27, 17, 14, 11, 0, 28, 30, 5, 17, 19, 18, 22, 12, 9, 13, 26, 26, 13, 12, 29, 15, 21, 2, 31, 13, 1, 13, 24, 25, 4, 21, 16, 29, 22, 20, 23, 17, 0}},
	{"bc1qem2ta6uk98rfr779t4wftq4qjtr3xtja9vf9yy3rgtczapc78j3sxa6570", []byte{0, 25, 27, 10, 11, 29, 26, 28, 22, 5, 7, 3, 9, 3, 30, 30, 5, 11, 21, 14, 9, 11, 0, 21, 0, 18, 11, 3, 17, 6, 11, 18, 29, 5, 12, 9, 5, 4, 4, 17, 3, 8, 11, 24, 2, 29, 1, 24, 30, 7, 18, 17, 16}},
	{"bc1qazm8jprsdjxn0qq77yrzw7m2340ys0kuuylg05vul4t5ll2lhduquuhngw", []byte{0, 29, 2, 27, 7, 18, 1, 3, 16, 13, 18, 6, 19, 15, 0, 0, 30, 30, 4, 3, 2, 14, 30, 27, 10, 17, 21, 15, 4, 16, 15, 22, 28, 28, 4, 31, 8, 15, 20, 12, 28, 31, 21, 11, 20, 31, 31, 10, 31, 23, 13, 28, 0}},
}

var bech32mAddresses = []decodedAddress{
	{"tb1pqaas5xm75dny58s452949c9ak5qd53shfkln490ju4ny2afs2ldsput844", []byte{1, 0, 29, 29, 16, 20, 6, 27, 30, 20, 13, 19, 4, 20, 7, 16, 21, 20, 10, 5, 21, 5, 24, 5, 29, 22, 20, 0, 13, 20, 17, 16, 23, 9, 22, 31, 19, 21, 5, 15, 18, 28, 21, 19, 4, 10, 29, 9, 16, 10, 31, 13, 16}},
	{"tb1p25h0xs3840q7aex3kl9dshd8q99qzaxkh8r5p70z54r4ykmn2rtsgcsj34", []byte{1, 10, 20, 23, 15, 6, 16, 17, 7, 21, 15, 0, 30, 29, 25, 6, 17, 22, 31, 5, 13, 16, 23, 13, 7, 0, 5, 5, 0, 2, 29, 6, 22, 23, 7, 3, 20, 1, 30, 15, 2, 20, 21, 3, 21, 4, 22, 27, 19, 10, 3, 11, 16}},
	{"tb1p7hvw8mnqlrtp7ffa8wzmhq7vddegffdeus4sl0yj6fw54zjda36qhc5q8y", []byte{1, 30, 23, 12, 14, 7, 27, 19, 0, 31, 3, 11, 1, 30, 9, 9, 29, 7, 14, 2, 27, 23, 0, 30, 12, 13, 13, 25, 8, 9, 9, 13, 25, 28, 16, 21, 16, 31, 15, 4, 18, 26, 9, 14, 20, 21, 2, 18, 13, 29, 17, 26, 0}},
	{"tb1pnqdr56lugmtrcxtae8k9cfe7hve8986ud0daktljsh93wf8q7u4qhc2q3c", []byte{1, 19, 0, 13, 3, 20, 26, 31, 28, 8, 27, 11, 3, 24, 6, 11, 29, 25, 7, 22, 5, 24, 9, 25, 30, 23, 12, 25, 7, 5, 7, 26, 28, 13, 15, 13, 29, 22, 11, 31, 18, 16, 23, 5, 17, 14, 9, 7, 0, 30, 28, 21, 0}},
	{"tb1pa54gmj3dzr9g5p7qx6kupqg9xkvtv2cdcty78wgyaycxtqc72h5qlqgz2c", []byte{1, 29, 20, 21, 8, 27, 18, 17, 13, 2, 3, 5, 8, 20, 1, 30, 0, 6, 26, 22, 28, 1, 0, 8, 5, 6, 22, 12, 11, 12, 10, 24, 13, 24, 11, 4, 30, 7, 14, 8, 4, 29, 4, 24, 6, 11, 0, 24, 30, 10, 23, 20, 0}},
	{"tb1p8lkxfnps5wd6rsrusvytp8zllrmxz05e0ttessnhyzwl0kusc2as4s72wz", []byte{1, 7, 31, 22, 6, 9, 19, 1, 16, 20, 14, 13, 26, 3, 16, 3, 28, 16, 12, 4, 11, 1, 7, 2, 31, 31, 3, 27, 6, 2, 15, 20, 25, 15, 11, 11, 25, 16, 16, 19, 23, 4, 2, 14, 31, 15, 22, 28, 16, 24, 10, 29, 16}},
	{"bc1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297", []byte{1, 20, 13, 30, 3, 18, 0, 30, 8, 26, 3, 13, 22, 10, 4, 23, 2, 22, 16, 5, 16, 27, 31, 29, 0, 11, 25, 13, 3, 21, 13, 25, 22, 0, 15, 7, 8, 25, 7, 2, 11, 14, 29, 24, 30, 10, 16, 9, 3, 5, 3, 28, 16}},
	{"bc1py8g4v4ehll399qlpaxyxykg37pszhad9yg0dphxvjhdmhy7f08vsn43s6p", []byte{1, 4, 7, 8, 21, 12, 21, 25, 23, 31, 31, 17, 5, 5, 0, 31, 1, 29, 6, 4, 6, 4, 22, 8, 17, 30, 1, 16, 2, 23, 29, 13, 5, 4, 8, 15, 13, 1, 23, 6, 12, 18, 23, 13, 27, 23, 4, 30, 9, 15, 7, 12, 16}},
	{"bc1pc09cafvlgu5ykmxyyzr4gu5qwx9a2zz6fz3lljeyddc9z7n75n9qfz7ckr", []byte{1, 24, 15, 5, 24, 29, 9, 12, 31, 8, 28, 20, 4, 22, 27, 6, 4, 4, 2, 3, 21, 8, 28, 20, 0, 14, 6, 5, 29, 10, 2, 2, 26, 9, 2, 17, 31, 31, 18, 25, 4, 13, 13, 24, 5, 2, 30, 19, 30, 20, 19, 5, 0}},
	{"bc1p74k39706fe6n0qv5k30z4xpqd8gcf8apyzn9s5rujkz4jln2u3fqwwta94", []byte{1, 30, 21, 22, 17, 5, 30, 15, 26, 9, 25, 26, 19, 15, 0, 12, 20, 22, 17, 15, 2, 21, 6, 1, 0, 13, 7, 8, 24, 9, 7, 29, 1, 4, 2, 19, 5, 16, 20, 3, 28, 18, 22, 2, 21, 18, 31, 19, 10, 28, 17, 9, 0}},
	{"bc1petgnkphl82md05d84gwee0alkuzpphfjy8ycxs932ngvdx8z8u0s3dwj5t", []byte{1, 25, 11, 8, 19, 22, 1, 23, 31, 7, 10, 27, 13, 15, 20, 13, 7, 21, 8, 14, 25, 25, 15, 29, 31, 22, 28, 2, 1, 1, 23, 9, 18, 4, 7, 4, 24, 6, 16, 5, 17, 10, 19, 8, 12, 13, 6, 7, 2, 7, 28, 15, 16}},
}

func TestDecodeAddressBase58(t *testing.T) {
	var decodedAddresses []decodedAddress
	decodedAddresses = append(decodedAddresses, base58Addresses...)
	cases := decodedAddresses
	for _, c := range cases {
		withVersion, err := bitcoin.DecodeAddressBase58(c.address, true)
		require.NoError(t, err)
		withoutVersion, err := bitcoin.DecodeAddressBase58(c.address, false)
		require.NoError(t, err)
		assert.Equal(t, c.expected, withVersion)
		assert.Equal(t, c.expected[1:], withoutVersion)
	}
}

func TestDecodeAddressBase58_ErrorHandling(t *testing.T) {
	var errorCases []string = []string{
		"5Hwgr3u458GLafKBgxtssHSPqJnYoGrSzgQsPwLFhLNYskDPyyA",
		"not in bas58",
		"A",
		"0x79568c2989232dCa1840087D73d403602364c0D4",
	}
	var result []byte
	var err error
	for _, c := range errorCases {
		result, err = bitcoin.DecodeAddressBase58(c, true)
		require.Error(t, err)
		assert.Nil(t, result)
		result, err = bitcoin.DecodeAddressBase58(c, false)
		require.Error(t, err)
		assert.Nil(t, result)
	}
}

func TestDecodeAddress(t *testing.T) {
	var decodedAddresses []decodedAddress
	decodedAddresses = append(decodedAddresses, base58Addresses...)
	decodedAddresses = append(decodedAddresses, bech32Addresses...)
	decodedAddresses = append(decodedAddresses, bech32mAddresses...)
	cases := decodedAddresses
	for _, c := range cases {
		decoded, err := bitcoin.DecodeAddress(c.address)
		require.NoError(t, err)
		assert.Equal(t, c.expected, decoded)
	}
}

func TestToSwappedBytes32(t *testing.T) {
	var bytes32 = [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20}
	var hash = chainhash.Hash([32]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40})
	var hashPointer, _ = chainhash.NewHash([]byte{0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60})

	swappedBytes := bitcoin.ToSwappedBytes32(bytes32)
	swappedHash := bitcoin.ToSwappedBytes32(hash)
	swappedHashPointer := bitcoin.ToSwappedBytes32(hashPointer)

	assert.Equal(t, [32]byte{0x20, 0x1F, 0x1E, 0x1D, 0x1C, 0x1B, 0x1A, 0x19, 0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11, 0x10, 0x0F, 0x0E, 0x0D, 0x0C, 0x0B, 0x0A, 0x09, 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}, swappedBytes)
	assert.Equal(t, [32]byte{0x40, 0x3F, 0x3E, 0x3D, 0x3C, 0x3B, 0x3A, 0x39, 0x38, 0x37, 0x36, 0x35, 0x34, 0x33, 0x32, 0x31, 0x30, 0x2F, 0x2E, 0x2D, 0x2C, 0x2B, 0x2A, 0x29, 0x28, 0x27, 0x26, 0x25, 0x24, 0x23, 0x22, 0x21}, swappedHash)
	assert.Equal(t, [32]byte{0x60, 0x5F, 0x5E, 0x5D, 0x5C, 0x5B, 0x5A, 0x59, 0x58, 0x57, 0x56, 0x55, 0x54, 0x53, 0x52, 0x51, 0x50, 0x4F, 0x4E, 0x4D, 0x4C, 0x4B, 0x4A, 0x49, 0x48, 0x47, 0x46, 0x45, 0x44, 0x43, 0x42, 0x41}, swappedHashPointer)
}
