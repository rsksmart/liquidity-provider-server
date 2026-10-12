package regtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	lpsMongo "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const databaseTimeout = 5 * time.Second

// LpsDB reads the LPS's peg-in state through the LPS's own repositories. It never writes: the
// test observes the LPS, it does not steer it.
type LpsDB struct {
	claims  rootstock.PegInClaimRepository
	watches interface {
		Get(ctx context.Context, rskAddress string) (*rootstock.PegInWatch, error)
	}
}

func NewLpsDB(t *testing.T, stack Stack) *LpsDB {
	t.Helper()
	// The replica set advertises its member as mongo01, which the host cannot resolve.
	uri := fmt.Sprintf("mongodb://%s:%s@%s/admin?directConnection=true", stack.Mongo.User, stack.Mongo.Password, stack.Mongo.Host)
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, client.Disconnect(context.Background())) })
	require.NoErrorf(t, client.Ping(t.Context(), nil), "MongoDB is not reachable at %s", stack.Mongo.Host)

	connection := lpsMongo.NewConnection(lpsMongo.NewClientWrapper(client), databaseTimeout)
	return &LpsDB{
		claims:  lpsMongo.NewPegInClaimMongoRepository(connection),
		watches: lpsMongo.NewPegInWatchMongoRepository(connection),
	}
}

// Watch is nil until the LPS discovers rskAddr in the registry.
func (db *LpsDB) Watch(t *testing.T, rskAddr common.Address) *rootstock.PegInWatch {
	t.Helper()
	watch, err := db.watches.Get(t.Context(), rskAddressKey(rskAddr))
	require.NoError(t, err)
	return watch
}

// Claim is nil until the LPS starts claiming the deposit.
func (db *LpsDB) Claim(t *testing.T, rskAddr common.Address, depositTxID string) *rootstock.PegInClaim {
	t.Helper()
	claim, err := db.claims.Get(t.Context(), rskAddressKey(rskAddr), depositTxID)
	require.NoError(t, err)
	return claim
}

// rskAddressKey is the checksummed hex the registry watcher stores RSK addresses as.
func rskAddressKey(address common.Address) string {
	return address.Hex()
}
