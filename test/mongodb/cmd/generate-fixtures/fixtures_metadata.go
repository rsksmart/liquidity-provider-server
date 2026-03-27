package main

import "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"

type fixturesCollection struct {
	collection string
	fileName   string
}

var fixtureCollections = []fixturesCollection{
	{collection: mongo.PeginQuoteCollection, fileName: "peginQuote.json"},
	{collection: mongo.RetainedPeginQuoteCollection, fileName: "retainedPeginQuote.json"},
	{collection: mongo.PeginCreationDataCollection, fileName: "peginQuoteCreationData.json"},
	{collection: mongo.PegoutQuoteCollection, fileName: "pegoutQuote.json"},
	{collection: mongo.RetainedPegoutQuoteCollection, fileName: "retainedPegoutQuote.json"},
	{collection: mongo.PegoutCreationDataCollection, fileName: "pegoutQuoteCreationData.json"},
	{collection: mongo.DepositEventsCollection, fileName: "depositEvents.json"},
	{collection: mongo.LiquidityProviderCollection, fileName: "liquidityProvider.json"},
	{collection: mongo.TrustedAccountCollection, fileName: "trustedAccounts.json"},
	{collection: mongo.PenalizedEventCollection, fileName: "penalizedEvent.json"},
	{collection: mongo.BatchPegOutEventsCollection, fileName: "batchPegOutEvents.json"},
}
