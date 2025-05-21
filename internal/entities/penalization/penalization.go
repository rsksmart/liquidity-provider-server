package penalization

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
)

type PenalizedEvent struct {
	LiquidityProvider string        `json:"liquidityProvider" bson:"liquidity_provider" validate:"required"`
	Penalty           *entities.Wei `json:"penalty" bson:"penalty" validate:"required"`
	QuoteHash         string        `json:"quoteHash" bson:"quote_hash" validate:"required"`
}

type PenalizedEventRepository interface {
	InsertPenalization(ctx context.Context, event PenalizedEvent) error
	GetPenalizationsByQuoteHashes(ctx context.Context, quoteHashes []string) ([]PenalizedEvent, error)
}
