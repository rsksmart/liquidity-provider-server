package regtest

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bridgeBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/bridge"
	flyoverConfigurationsBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/flyover_configurations"
	pauseRegistryBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pause_registry"
	peginAddressRegistryBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_address_registry"
	peginCommitFirstBinding "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	entityRootstock "github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test/integration"
	"github.com/stretchr/testify/require"
)

const (
	RskHTTPEndpoint = "http://localhost:4444"
	RskChainID      = 33
	RskBridgeHex    = "0x0000000000000000000000000000000001000006"
	DepositWei      = 5_000_000_000_000_000 // 0.005 ether
	// Account 0 of DEPLOYER_MNEMONIC "test test test test test test test test test test test junk".
	DeployerPrivateKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	DeployerAddressHex    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	LPSHotWalletHex       = "0x9D93929A9099be4355fC2389FbF253982F9dF47c"
	localKeyPassword      = "test"
)

func localKeyFile() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("docker-compose", "local", "localstack", "local-key.json")
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "docker-compose", "local", "localstack", "local-key.json")
}

// noRetry matches cmd/utils/scripts: one RPC attempt, no 1-minute production backoff.
var noRetry = rootstock.RetryParams{Retries: 0, Sleep: 0}

type RootstockStack struct {
	client          rootstock.RpcClientBinding
	rpc             blockchain.RootstockRpcServer
	ChainID         *big.Int
	Deployer        *ecdsa.PrivateKey
	DeployerAddress common.Address
	RegistryAddr    common.Address
	PegInAddr       common.Address
	PauseAddr       common.Address
	registry        blockchain.PegInAddressRegistryContract
	configurations  blockchain.FlyoverConfigurationsContract
	registryBinding *peginAddressRegistryBinding.PegInAddressRegistryContract
	registryBound   *bind.BoundContract
	pegInBinding    *peginCommitFirstBinding.PeginCommitFirstContract
	pegInBound      *bind.BoundContract
	pauseBinding    *pauseRegistryBinding.PauseRegistryContract
	pauseBound      *bind.BoundContract
	bridgeBinding   *bridgeBinding.RskBridge
	bridgeBound     *bind.BoundContract
}

func OpenRootstock(t *testing.T, cfg *integration.Config) *RootstockStack {
	t.Helper()
	ctx := context.Background()
	rskClient, err := bootstrap.Rootstock(ctx, environment.Environment{
		Rsk: environment.RskEnv{Endpoint: RskHTTPEndpoint, ChainId: RskChainID},
	})
	require.NoError(t, err)
	client := rskClient.Rpc()
	t.Cleanup(client.Close)

	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(RskChainID), chainID.Uint64())

	deployer, err := crypto.HexToECDSA(DeployerPrivateKeyHex)
	require.NoError(t, err)
	deployerAddress := crypto.PubkeyToAddress(deployer.PublicKey)
	require.Equal(t, common.HexToAddress(DeployerAddressHex), deployerAddress)

	registryAddr := common.HexToAddress(cfg.Rsk.PeginAddressRegistry)
	configurationsAddr := common.HexToAddress(cfg.Rsk.FlyoverConfigurations)
	pegInAddr := common.HexToAddress(cfg.Rsk.PeginContract)
	pauseAddr := common.HexToAddress(cfg.Rsk.PauseRegistry)
	require.False(t, isZeroAddress(registryAddr), "peginAddressRegistry is missing from integration-test.config.json")
	require.False(t, isZeroAddress(configurationsAddr), "flyoverConfigurations is missing from integration-test.config.json")
	require.False(t, isZeroAddress(pegInAddr), "peginContract is missing from integration-test.config.json")
	require.False(t, isZeroAddress(pauseAddr), "pauseRegistry is missing from integration-test.config.json")

	abis := rootstock.MustLoadFlyoverABIs()
	registryBinding := peginAddressRegistryBinding.NewPegInAddressRegistryContract()
	configurationsBinding := flyoverConfigurationsBinding.NewFlyoverConfigurationsContract()
	pegInBinding := peginCommitFirstBinding.NewPeginCommitFirstContract()
	pauseBinding := pauseRegistryBinding.NewPauseRegistryContract()
	rskBridgeBinding := bridgeBinding.NewRskBridge()
	registryBound := registryBinding.Instance(client, registryAddr)
	configurationsBound := configurationsBinding.Instance(client, configurationsAddr)
	pegInBound := pegInBinding.Instance(client, pegInAddr)
	pauseBound := pauseBinding.Instance(client, pauseAddr)
	bridgeAddr := common.HexToAddress(RskBridgeHex)
	bridgeBound := rskBridgeBinding.Instance(client, bridgeAddr)

	return &RootstockStack{
		client:          client,
		rpc:             rootstock.NewRskjRpcServer(rskClient, noRetry),
		ChainID:         chainID,
		Deployer:        deployer,
		DeployerAddress: deployerAddress,
		RegistryAddr:    registryAddr,
		PegInAddr:       pegInAddr,
		PauseAddr:       pauseAddr,
		registry: rootstock.NewPegInAddressRegistryContractImpl(
			rskClient,
			registryAddr.Hex(),
			registryBound,
			noRetry,
			registryBinding,
			abis,
		),
		configurations: rootstock.NewFlyoverConfigurationsContractImpl(
			rskClient,
			configurationsAddr.Hex(),
			configurationsBound,
			noRetry,
			configurationsBinding,
			abis,
		),
		registryBinding: registryBinding,
		registryBound:   registryBound,
		pegInBinding:    pegInBinding,
		pegInBound:      pegInBound,
		pauseBinding:    pauseBinding,
		pauseBound:      pauseBound,
		bridgeBinding:   rskBridgeBinding,
		bridgeBound:     bridgeBound,
	}
}

func (stack *RootstockStack) NewUser(t *testing.T) common.Address {
	t.Helper()
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	return crypto.PubkeyToAddress(key.PublicKey)
}

func (stack *RootstockStack) GetPegInAddress(t *testing.T, user common.Address) string {
	t.Helper()
	out, err := stack.registry.GetPegInAddress(user.Hex())
	require.NoError(t, err)
	require.Equal(t, blockchain.PegInAddressRegistryEncodingBase58, out.Encoding)
	encoded, err := bitcoin.EncodeAddressBase58(out.Payload)
	require.NoError(t, err)
	return encoded
}

func (stack *RootstockStack) CalculatePegInFee(t *testing.T, amount *entities.Wei) *entities.Wei {
	t.Helper()
	fee, err := stack.configurations.CalculatePegInFee(amount)
	require.NoError(t, err)
	return fee
}

func (stack *RootstockStack) RequiredConfirmations(t *testing.T, amount *entities.Wei) uint64 {
	t.Helper()
	confirmations, err := stack.configurations.GetRequiredPegInBtcConfirmations(amount)
	require.NoError(t, err)
	return confirmations
}

func (stack *RootstockStack) MinAmount(t *testing.T) *entities.Wei {
	t.Helper()
	min, err := stack.configurations.MinAmount()
	require.NoError(t, err)
	return min
}

func (stack *RootstockStack) SetPauseLevel(t *testing.T, level uint8, reason string) {
	t.Helper()
	callData, err := stack.pauseBinding.TryPackSetPauseLevel(level, reason)
	require.NoError(t, err)
	stack.send(t, stack.Deployer, stack.PauseAddr, big.NewInt(0), callData)
}

func (stack *RootstockStack) UnpackPegInRequested(t *testing.T, txHash string) blockchain.PegInRequestedEvent {
	t.Helper()
	receipt, err := stack.rpc.GetTransactionReceipt(context.Background(), txHash)
	require.NoError(t, err)
	event, err := unpackPegInRequestedForTest(stack.pegInBinding, stack.PegInAddr, receipt)
	require.NoError(t, err)
	return event
}

func unpackPegInRequestedForTest(
	binding *peginCommitFirstBinding.PeginCommitFirstContract,
	contract common.Address,
	receipt blockchain.TransactionReceipt,
) (blockchain.PegInRequestedEvent, error) {
	parsed, err := peginCommitFirstBinding.PeginCommitFirstContractMetaData.ParseABI()
	if err != nil {
		return blockchain.PegInRequestedEvent{}, err
	}
	eventID := parsed.Events["PegInRequested"].ID
	for _, eventLog := range receipt.Logs {
		if !eventLog.Removed &&
			common.HexToAddress(eventLog.Address) == contract &&
			len(eventLog.Topics) > 0 &&
			common.Hash(eventLog.Topics[0]) == eventID &&
			len(eventLog.Data) > 0 {
			topics := make([]common.Hash, len(eventLog.Topics))
			for i, topic := range eventLog.Topics {
				topics[i] = topic
			}
			gethLog := &types.Log{
				Address: common.HexToAddress(eventLog.Address),
				Topics:  topics,
				Data:    eventLog.Data,
			}
			unpacked, unpackErr := binding.UnpackPegInRequestedEvent(gethLog)
			if unpackErr != nil {
				return blockchain.PegInRequestedEvent{}, unpackErr
			}
			return blockchain.PegInRequestedEvent{
				PegInId:     unpacked.PegInId,
				Claimer:     unpacked.Claimer.Hex(),
				RskAddress:  unpacked.RskAddr.Hex(),
				Amount:      entities.NewBigWei(unpacked.Amount),
				NetToUser:   entities.NewBigWei(unpacked.NetToUser),
				CallSuccess: unpacked.CallSuccess,
			}, nil
		}
	}
	return blockchain.PegInRequestedEvent{}, fmt.Errorf("PegInRequested event not found in %s", receipt.TransactionHash)
}

func (stack *RootstockStack) GetRegistrationRoot(t *testing.T, blockNumber uint64) [32]byte {
	t.Helper()
	root, err := stack.registry.GetRegistrationRoot(context.Background(), blockNumber)
	require.NoError(t, err)
	return root
}

func (stack *RootstockStack) Balance(t *testing.T, addr common.Address) *big.Int {
	t.Helper()
	balance, err := stack.rpc.GetBalance(context.Background(), addr.Hex())
	require.NoError(t, err)
	return balance.AsBigInt()
}

func (stack *RootstockStack) Head(t *testing.T) uint64 {
	t.Helper()
	head, err := stack.rpc.GetHeight(context.Background())
	require.NoError(t, err)
	return head
}

func (stack *RootstockStack) Advance(t *testing.T, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		stack.send(t, stack.Deployer, stack.DeployerAddress, big.NewInt(0), nil)
	}
}

func (stack *RootstockStack) BridgeBtcHeight(t *testing.T) int64 {
	t.Helper()
	callData, err := stack.bridgeBinding.TryPackGetBtcBlockchainBestChainHeight()
	require.NoError(t, err)
	height, err := bind.Call(stack.bridgeBound, &bind.CallOpts{}, callData, stack.bridgeBinding.UnpackGetBtcBlockchainBestChainHeight)
	require.NoError(t, err)
	return height.Int64()
}

func (stack *RootstockStack) RegisterAddress(t *testing.T, user common.Address, btc blockchain.BitcoinNetwork, depositTxID string) {
	t.Helper()
	rawTx, err := btc.GetRawTransaction(depositTxID)
	require.NoError(t, err)
	block, err := btc.GetTransactionBlockInfo(depositTxID)
	require.NoError(t, err)
	merkle, err := btc.BuildMerkleBranch(depositTxID)
	require.NoError(t, err)
	stack.WaitBridgeBtcAtLeast(t, block.Height.Int64())
	callData, err := stack.registryBinding.TryPackRegisterAddress(
		user,
		rawTx,
		block.Hash,
		merkle.Path,
		merkle.Hashes,
	)
	require.NoError(t, err)
	stack.send(t, stack.Deployer, stack.RegistryAddr, big.NewInt(0), callData)
}

func (stack *RootstockStack) RequestPegInAsCompetitor(t *testing.T, user common.Address, btc blockchain.BitcoinNetwork, depositTxID string, net *entities.Wei) {
	t.Helper()
	rawTx, err := btc.GetRawTransaction(depositTxID)
	require.NoError(t, err)
	block, err := btc.GetTransactionBlockInfo(depositTxID)
	require.NoError(t, err)
	merkle, err := btc.BuildMerkleBranch(depositTxID)
	require.NoError(t, err)
	stack.WaitBridgeBtcAtLeast(t, block.Height.Int64())
	callData, err := stack.pegInBinding.TryPackRequestPegIn(
		user,
		rawTx,
		[]byte{},
		block.Hash,
		merkle.Path,
		merkle.Hashes,
	)
	require.NoError(t, err)
	stack.send(t, stack.Deployer, stack.PegInAddr, net.AsBigInt(), callData)
}

func (stack *RootstockStack) FundLPS(t *testing.T, amount *big.Int) {
	t.Helper()
	stack.send(t, stack.Deployer, common.HexToAddress(LPSHotWalletHex), amount, nil)
}

func (stack *RootstockStack) DrainLPS(t *testing.T) {
	t.Helper()
	lpsKey := decryptLPSHotWallet(t)
	lpsAddr := crypto.PubkeyToAddress(lpsKey.PublicKey)
	require.Equal(t, common.HexToAddress(LPSHotWalletHex), lpsAddr)
	balance := stack.Balance(t, lpsAddr)
	gasPrice, err := stack.client.SuggestGasPrice(context.Background())
	require.NoError(t, err)
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(21000))
	value := new(big.Int).Sub(balance, gasCost)
	if value.Sign() <= 0 {
		return
	}
	stack.send(t, lpsKey, stack.DeployerAddress, value, nil)
}

func (stack *RootstockStack) send(
	t *testing.T,
	key *ecdsa.PrivateKey,
	to common.Address,
	value *big.Int,
	data []byte,
) *types.Receipt {
	t.Helper()
	opts := bind.NewKeyedTransactor(key, stack.ChainID)
	opts.Value = value
	var tx *types.Transaction
	var err error
	if len(data) == 0 {
		nonce, nonceErr := stack.client.PendingNonceAt(context.Background(), opts.From)
		require.NoError(t, nonceErr)
		gasPrice, gasErr := stack.client.SuggestGasPrice(context.Background())
		require.NoError(t, gasErr)
		tx = types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &to,
			Value:    value,
			Gas:      21000,
			GasPrice: gasPrice,
		})
		tx, err = opts.Signer(opts.From, tx)
		require.NoError(t, err)
		require.NoError(t, stack.client.SendTransaction(context.Background(), tx))
	} else {
		contract := stack.contractAt(t, to)
		tx, err = bind.Transact(contract, opts, data)
		require.NoError(t, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, stack.client, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	return receipt
}

func (stack *RootstockStack) contractAt(t *testing.T, to common.Address) *bind.BoundContract {
	t.Helper()
	switch to {
	case stack.RegistryAddr:
		return stack.registryBound
	case stack.PegInAddr:
		return stack.pegInBound
	case stack.PauseAddr:
		return stack.pauseBound
	default:
		require.FailNowf(t, "unknown contract", "no BoundContract for %s", to.Hex())
		return nil
	}
}

func isZeroAddress(addr common.Address) bool {
	return addr == common.Address{}
}

func decryptLPSHotWallet(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	raw, err := os.ReadFile(localKeyFile())
	require.NoError(t, err)
	var wrapped struct {
		HotWallet json.RawMessage `json:"hotWallet"`
	}
	require.NoError(t, json.Unmarshal(raw, &wrapped))
	key, err := keystore.DecryptKey(wrapped.HotWallet, localKeyPassword)
	require.NoError(t, err)
	return key.PrivateKey
}

func DepositAmountWei() *big.Int {
	return new(big.Int).SetUint64(DepositWei)
}

func PayableNet(t *testing.T, amount, fee *entities.Wei) *entities.Wei {
	t.Helper()
	net, err := entityRootstock.CalculatePegInClaimPayableValue(amount, fee)
	require.NoError(t, err)
	return net
}
