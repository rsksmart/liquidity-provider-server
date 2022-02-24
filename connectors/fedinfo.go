package connectors

type FedInfo struct {
	FedSize              int
	FedThreshold         int
	PubKeys              []string
	FedAddress           string
	ActiveFedBlockHeight int
	IrisActivationHeight int
	ErpKeys              []string
}
