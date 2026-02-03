package monitoring

// Currency labels for metrics
const (
	MetricLabelRbtc = "rbtc"
	MetricLabelBtc  = "btc"
)

// Transfer reason labels for cold wallet metrics
const (
	MetricLabelThreshold   = "threshold"
	MetricLabelTimeForcing = "time_forcing"
)

// Asset metric type labels
const (
	MetricLabelTotal                       = "total"
	MetricLabelLocationRskWallet           = "location_rsk_wallet"
	MetricLabelLocationBtcWallet           = "location_btc_wallet"
	MetricLabelLocationLbc                 = "location_lbc"
	MetricLabelLocationFederation          = "location_federation"
	MetricLabelAllocationReservedForUsers  = "allocation_reserved_for_users"
	MetricLabelAllocationWaitingRefund     = "allocation_waiting_refund"
	MetricLabelAllocationAvailable         = "allocation_available"
)
