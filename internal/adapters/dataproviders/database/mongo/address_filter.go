package mongo

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func rskAddressEqualsFilter(field, address string) (bson.M, error) {
	normalized, err := blockchain.NormalizeRskAddress(address)
	if err != nil {
		return nil, err
	}
	return bson.M{
		"$expr": bson.M{
			"$eq": []interface{}{
				bson.M{"$toLower": "$" + field},
				normalized,
			},
		},
	}, nil
}

func retainedQuotesOwnerFilter(address string, states any) (bson.M, error) {
	filter, err := rskAddressEqualsFilter("owner_account_address", address)
	if err != nil {
		return nil, err
	}
	filter["state"] = bson.D{
		primitive.E{Key: "$in", Value: states},
	}
	return filter, nil
}
