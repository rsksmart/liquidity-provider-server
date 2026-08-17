package bindings_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	quotepegin "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitfirst "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// requestPegInSelector is keccak256("requestPegIn(address,bytes,bytes,bytes32,uint256,bytes32[])")[:4]
// from the pinned IPegInCommitFirst ABI, not from the packer under test.
var requestPegInSelector = []byte{0xa3, 0x55, 0xe9, 0x35}

func TestPeginCommitFirstBindingPacksRequestPegIn(t *testing.T) {
	contract := commitfirst.NewPeginCommitFirstContract()
	rskAddr := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	rawTx := []byte{0x01, 0x02, 0x03}
	blockHash := [32]byte{0x11}
	path := big.NewInt(1)
	hashes := [][32]byte{{0x22}}

	calldata := contract.PackRequestPegIn(rskAddr, rawTx, []byte{}, blockHash, path, hashes)

	require.GreaterOrEqual(t, len(calldata), 4)
	assert.Equal(t, requestPegInSelector, calldata[:4])
	assert.Contains(t, commitfirst.PeginCommitFirstContractMetaData.ABI, `"name":"requestPegIn"`)
	assert.Contains(t, commitfirst.PeginCommitFirstContractMetaData.ABI, `"name":"PegInRequested"`)
	assert.Contains(t, commitfirst.PeginCommitFirstContractMetaData.ABI, `"name":"PegInAlreadyProcessed"`)
	assert.Contains(t, commitfirst.PeginCommitFirstContractMetaData.ABI, `"name":"AddressNotRegistered"`)
	assert.Contains(t, commitfirst.PeginCommitFirstContractMetaData.ABI, `"name":"DepositOutputNotFound"`)
	assert.Contains(t, commitfirst.PeginCommitFirstContractMetaData.ABI, `"name":"InsufficientConfirmations"`)
	assert.Contains(t, commitfirst.PeginCommitFirstContractMetaData.ABI, `"name":"IncorrectFronting"`)

	var (
		_ = contract.UnpackError
		_ = contract.UnpackPegInRequestedEvent
		_ = contract.UnpackPegInAlreadyProcessedError
		_ *commitfirst.PeginCommitFirstContractPegInAlreadyProcessed
		_ *commitfirst.PeginCommitFirstContractAddressNotRegistered
		_ *commitfirst.PeginCommitFirstContractDepositOutputNotFound
		_ *commitfirst.PeginCommitFirstContractInsufficientConfirmations
		_ *commitfirst.PeginCommitFirstContractIncorrectFronting
		_ *commitfirst.PeginCommitFirstContractPegInRequested
	)
}

func TestQuotePeginBindingStillExposesRegisterPegIn(t *testing.T) {
	quote := quotepegin.NewPeginContract()
	assert.Contains(t, quotepegin.PeginContractMetaData.ABI, `"name":"registerPegIn"`)
	assert.Contains(t, quotepegin.PeginContractMetaData.ABI, `"name":"callForUser"`)
	assert.NotContains(t, quotepegin.PeginContractMetaData.ABI, `"name":"requestPegIn"`)

	var (
		_ = quote.PackRegisterPegIn
		_ = quote.PackCallForUser
	)
}
