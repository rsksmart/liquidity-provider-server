// Package regtest drives the local regtest stack (docker-compose/local) from integration tests.
// It reads the env file lps-local.sh writes, so contract addresses always match the last deploy.
package regtest

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/stretchr/testify/require"
)

const (
	envFileOverride = "REGTEST_ENV_FILE"
	defaultEnvFile  = "docker-compose/local/.env.regtest"
	// The env file holds the compose service names; the tests run on the host.
	hostAddress = "127.0.0.1"
)

type Contracts struct {
	PegIn                 common.Address
	PegOut                common.Address
	CollateralManagement  common.Address
	Discovery             common.Address
	PegInAddressRegistry  common.Address
	FlyoverConfigurations common.Address
	PauseRegistry         common.Address
	Bridge                common.Address
}

type Mongo struct {
	Host     string
	User     string
	Password string
}

// Stack is the local regtest stack as seen from the host.
type Stack struct {
	LpsURL            string
	RskURL            string
	ChainID           string
	Btc               environment.BtcEnv
	Mongo             Mongo
	Contracts         Contracts
	LiquidityProvider common.Address
}

// LoadStack reads the env file of the local stack. REGTEST_ENV_FILE overrides its path, and
// REGTEST_LPS_URL, REGTEST_RSK_URL, REGTEST_BTC_ENDPOINT and REGTEST_MONGODB_HOST override the
// endpoints derived from it.
func LoadStack(t *testing.T) Stack {
	t.Helper()
	path := os.Getenv(envFileOverride)
	if path == "" {
		path = filepath.Join(moduleRoot(t), defaultEnvFile)
	}
	values := readEnvFile(t, path)
	required := func(key string) string {
		value := values[key]
		require.NotEmptyf(t, value, "%s is missing in %s. Run lps-local.sh first", key, path)
		return value
	}
	address := func(key string) common.Address {
		value := required(key)
		require.Truef(t, common.IsHexAddress(value), "%s=%s in %s is not an address", key, value, path)
		parsed := common.HexToAddress(value)
		require.NotEqualf(t, common.Address{}, parsed, "%s is the zero address in %s. Did the deploy run?", key, path)
		return parsed
	}

	return Stack{
		LpsURL:  overrideOr("REGTEST_LPS_URL", "http://"+net.JoinHostPort(hostAddress, required("SERVER_PORT"))),
		RskURL:  overrideOr("REGTEST_RSK_URL", onHost(t, required("RSK_ENDPOINT"))),
		ChainID: required("CHAIN_ID"),
		Btc: environment.BtcEnv{
			Network:  required("BTC_NETWORK"),
			Username: required("BTC_USERNAME"),
			Password: required("BTC_PASSWORD"),
			Endpoint: overrideOr("REGTEST_BTC_ENDPOINT", hostPortOnHost(t, required("BTC_ENDPOINT"))),
		},
		Mongo: Mongo{
			Host:     overrideOr("REGTEST_MONGODB_HOST", net.JoinHostPort(hostAddress, required("MONGODB_PORT"))),
			User:     required("MONGODB_USER"),
			Password: required("MONGODB_PASSWORD"),
		},
		Contracts: Contracts{
			PegIn:                 address("PEGIN_CONTRACT_ADDRESS"),
			PegOut:                address("PEGOUT_CONTRACT_ADDRESS"),
			CollateralManagement:  address("COLLATERAL_MANAGEMENT_ADDRESS"),
			Discovery:             address("DISCOVERY_ADDRESS"),
			PegInAddressRegistry:  address("PEGIN_ADDRESS_REGISTRY_ADDRESS"),
			FlyoverConfigurations: address("FLYOVER_CONFIGURATIONS_ADDRESS"),
			PauseRegistry:         address("PAUSE_REGISTRY_ADDRESS"),
			Bridge:                address("RSK_BRIDGE_ADDR"),
		},
		LiquidityProvider: address("LIQUIDITY_PROVIDER_RSK_ADDR"),
	}
}

// RequireLpsHealthy fails unless the LPS reports itself and all its services as ok.
func (stack Stack) RequireLpsHealthy(t *testing.T) {
	t.Helper()
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, stack.LpsURL+"/health", nil)
	require.NoError(t, err)
	response, err := http.DefaultClient.Do(request)
	require.NoErrorf(t, err, "the LPS is not reachable at %s. Run lps-local.sh first", stack.LpsURL)
	defer func() { require.NoError(t, response.Body.Close()) }()
	require.Equal(t, http.StatusOK, response.StatusCode)

	var health pkg.HealthResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&health))
	const ok = "ok"
	require.Equalf(t, pkg.HealthResponse{Status: ok, Services: pkg.Services{Db: ok, Rsk: ok, Btc: ok}}, health,
		"the LPS at %s is not healthy", stack.LpsURL)
}

// onHost keeps the scheme and port of a compose URL and points it at the host.
func onHost(t *testing.T, rawURL string) string {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	require.NoError(t, err)
	parsed.Host = net.JoinHostPort(hostAddress, parsed.Port())
	return parsed.String()
}

// hostPortOnHost keeps the port of a compose host:port and points it at the host.
func hostPortOnHost(t *testing.T, hostPort string) string {
	t.Helper()
	_, port, err := net.SplitHostPort(hostPort)
	require.NoErrorf(t, err, "%s is not a host:port", hostPort)
	return net.JoinHostPort(hostAddress, port)
}

func overrideOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		require.NotEqualf(t, dir, parent, "cannot find go.mod above %s", dir)
		dir = parent
	}
}

// readEnvFile reads the KEY=VALUE lines lps-local.sh sources: quoted values keep their content,
// and unquoted values drop a trailing " # comment".
func readEnvFile(t *testing.T, path string) map[string]string {
	t.Helper()
	file, err := os.Open(path)
	require.NoErrorf(t, err, "cannot read %s. Run lps-local.sh first", path)
	defer func() { require.NoError(t, file.Close()) }()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		values[strings.TrimSpace(key)] = envValue(strings.TrimSpace(value))
	}
	require.NoError(t, scanner.Err())
	return values
}

func envValue(value string) string {
	for _, quote := range []string{`"`, `'`} {
		if strings.HasPrefix(value, quote) {
			if end := strings.Index(value[1:], quote); end >= 0 {
				return value[1 : end+1]
			}
		}
	}
	if comment := strings.Index(value, " #"); comment >= 0 {
		value = value[:comment]
	}
	return strings.TrimSpace(value)
}
