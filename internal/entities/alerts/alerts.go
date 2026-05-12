package alerts

import "context"

// This constants are used to standardize the alert subjects so they can be used in the alerting system
// Changing one of these constants impacts the external alerting system and there is not an automatic way
// to identify that error.
const (
	AlertSubjectPenalization         = "LPS has been penalized"
	AlertSubjectPeginOutOfLiquidity  = "PegIn: Out of liquidity"
	AlertSubjectPegoutOutOfLiquidity = "PegOut: Out of liquidity"
	AlertSubjectEclipseAttack        = "Node Eclipse Detected"
)

type AlertSender interface {
	SendAlert(ctx context.Context, subject, body string, recipient []string) error
}
