package rootstock

import "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"

type EventIteratorAdapter[T any] interface {
	Next() bool
	Close() error
	Event() *T
	Error() error
}

// -------------- Event iterator adapters -------------------
type depositEventIteratorAdapter struct {
	*bindings.IPegOutPegOutDepositIterator
}

func (i *depositEventIteratorAdapter) Event() *bindings.IPegOutPegOutDeposit {
	return i.IPegOutPegOutDepositIterator.Event
}

type penalizedEventIteratorAdapter struct {
	*bindings.ICollateralManagementPenalizedIterator
}

func (i *penalizedEventIteratorAdapter) Event() *bindings.ICollateralManagementPenalized {
	return i.ICollateralManagementPenalizedIterator.Event
}

type batchPegOutCreatedEventIteratorAdapter struct {
	*bindings.RskBridgeBatchPegoutCreatedIterator
}

func (i *batchPegOutCreatedEventIteratorAdapter) Event() *bindings.RskBridgeBatchPegoutCreated {
	return i.RskBridgeBatchPegoutCreatedIterator.Event
}

// ---------------------------------------------------------
// ------------------ Contract  adapters -------------------

type rskBridgeAdapter struct {
	*bindings.RskBridge
}

func NewRskBridgeAdapter(rskBridge *bindings.RskBridge) RskBridgeAdapter {
	return &rskBridgeAdapter{RskBridge: rskBridge}
}

func (a *rskBridgeAdapter) BatchPegOutCreatedIteratorAdapter(rawIterator *bindings.RskBridgeBatchPegoutCreatedIterator) EventIteratorAdapter[bindings.RskBridgeBatchPegoutCreated] {
	return &batchPegOutCreatedEventIteratorAdapter{RskBridgeBatchPegoutCreatedIterator: rawIterator}
}

type peginContractAdapter struct {
	*bindings.IPegIn
}

func NewPeginContractAdapter(peginContract *bindings.IPegIn) PeginContractAdapter {
	return &peginContractAdapter{IPegIn: peginContract}
}

func (a *peginContractAdapter) Caller() ContractCallerBinding {
	return &bindings.IPegInCallerRaw{Contract: &a.IPegIn.IPegInCaller}
}

type pegoutContractAdapter struct {
	*bindings.IPegOut
}

func NewPegoutContractAdapter(pegoutContract *bindings.IPegOut) PegoutContractAdapter {
	return &pegoutContractAdapter{IPegOut: pegoutContract}
}

func (a *pegoutContractAdapter) DepositEventIteratorAdapter(rawIterator *bindings.IPegOutPegOutDepositIterator) EventIteratorAdapter[bindings.IPegOutPegOutDeposit] {
	return &depositEventIteratorAdapter{IPegOutPegOutDepositIterator: rawIterator}
}

func (a *pegoutContractAdapter) Caller() ContractCallerBinding {
	return &bindings.IPegOutCallerRaw{Contract: &a.IPegOut.IPegOutCaller}
}

type collateralManagementAdapter struct {
	*bindings.ICollateralManagement
}

func NewCollateralManagementAdapter(collateralManagement *bindings.ICollateralManagement) CollateralManagementAdapter {
	return &collateralManagementAdapter{ICollateralManagement: collateralManagement}
}

func (a *collateralManagementAdapter) Caller() ContractCallerBinding {
	return &bindings.ICollateralManagementCallerRaw{Contract: &a.ICollateralManagement.ICollateralManagementCaller}
}

func (a *collateralManagementAdapter) PenalizedEventIteratorAdapter(rawIterator *bindings.ICollateralManagementPenalizedIterator) EventIteratorAdapter[bindings.ICollateralManagementPenalized] {
	return &penalizedEventIteratorAdapter{ICollateralManagementPenalizedIterator: rawIterator}
}

// ---------------------------------------------------------
