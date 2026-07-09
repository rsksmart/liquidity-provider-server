package logging_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/logging"
	"github.com/stretchr/testify/require"
)

func TestCaptureStack_IncludesCallerOutsideLoggingPackage(t *testing.T) {
	stack := logging.CaptureStackForTest(0)
	require.NotEmpty(t, stack)
	require.Contains(t, stack, "stack_test.go")
	require.NotContains(t, stack, "internal/logging/logging.go")
}
