package utils_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFirstNonZero(t *testing.T) {
	assert.Equal(t, 5, utils.FirstNonZero(0, 0, 5, 10))
	assert.Equal(t, "hello", utils.FirstNonZero("", "hello", "world"))
	assert.InEpsilon(t, 3.14, utils.FirstNonZero(0.0, 0.0, 3.14), 0.0000)
	assert.Zero(t, utils.FirstNonZero(0.0, 0.0, 0.0))
	assert.Equal(t, 0, utils.FirstNonZero(0, 0, 0))
	assert.Equal(t, 2, utils.FirstNonZero(2))
	assert.Equal(t, time.Duration(3), utils.FirstNonZero(time.Duration(0), time.Duration(3), time.Duration(9), time.Duration(0)))
}
