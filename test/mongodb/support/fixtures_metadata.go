package support

import "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"

// FixtureCollection pairs a MongoDB collection name with its fixture JSON file name.
type FixtureCollection struct {
	Collection string
	FileName   string
}

// FixtureCollections is the single source of truth for which collections have fixtures
// and how they map to JSON files. Both the generate-fixtures tool and the test suite
// derive their collection lists from here.
var FixtureCollections = []FixtureCollection{
	{Collection: mongo.PeginQuoteCollection, FileName: "peginQuote.json"},
	{Collection: mongo.RetainedPeginQuoteCollection, FileName: "retainedPeginQuote.json"},
	{Collection: mongo.PeginCreationDataCollection, FileName: "peginQuoteCreationData.json"},
	{Collection: mongo.PegoutQuoteCollection, FileName: "pegoutQuote.json"},
	{Collection: mongo.RetainedPegoutQuoteCollection, FileName: "retainedPegoutQuote.json"},
	{Collection: mongo.PegoutCreationDataCollection, FileName: "pegoutQuoteCreationData.json"},
	{Collection: mongo.DepositEventsCollection, FileName: "depositEvents.json"},
	{Collection: mongo.LiquidityProviderCollection, FileName: "liquidityProvider.json"},
	{Collection: mongo.TrustedAccountCollection, FileName: "trustedAccounts.json"},
	{Collection: mongo.PenalizedEventCollection, FileName: "penalizedEvent.json"},
	{Collection: mongo.BatchPegOutEventsCollection, FileName: "batchPegOutEvents.json"},
}
