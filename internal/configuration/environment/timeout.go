package environment

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"time"
)

type Timeout uint64

func (t Timeout) Seconds() time.Duration {
	return time.Duration(t) * time.Second
}

type ApplicationTimeouts struct {
	Bootstrap           Timeout `validate:"required"`
	WatcherPreparation  Timeout `validate:"required"`
	WatcherValidation   Timeout `validate:"required"`
	DatabaseInteraction Timeout `validate:"required"`
	MiningWait          Timeout `validate:"required"`
	DatabaseConnection  Timeout `validate:"required"`
	ServerReadHeader    Timeout `validate:"required"`
	ServerWrite         Timeout `validate:"required"`
	ServerIdle          Timeout `validate:"required"`
	PegoutDepositCheck  Timeout `validate:"required"`
}

func DefaultTimeouts() ApplicationTimeouts {
	return ApplicationTimeouts{
		Bootstrap:           240,
		WatcherPreparation:  15,
		WatcherValidation:   15,
		DatabaseInteraction: 3,
		MiningWait:          300,
		DatabaseConnection:  10,
		ServerReadHeader:    5,
		ServerWrite:         60,
		ServerIdle:          10,
		PegoutDepositCheck:  60,
	}
}

func TimeoutsFromEnv(env TimeoutEnv) (ApplicationTimeouts, error) {
	defaultTimeouts := DefaultTimeouts()
	timeouts := ApplicationTimeouts{}
	timeouts.Bootstrap = utils.FirstNonZero(Timeout(env.Bootstrap), defaultTimeouts.Bootstrap)
	timeouts.WatcherPreparation = utils.FirstNonZero(Timeout(env.WatcherPreparation), defaultTimeouts.WatcherPreparation)
	timeouts.WatcherValidation = utils.FirstNonZero(Timeout(env.WatcherValidation), defaultTimeouts.WatcherValidation)
	timeouts.DatabaseInteraction = utils.FirstNonZero(Timeout(env.DatabaseInteraction), defaultTimeouts.DatabaseInteraction)
	timeouts.MiningWait = utils.FirstNonZero(Timeout(env.MiningWait), defaultTimeouts.MiningWait)
	timeouts.DatabaseConnection = utils.FirstNonZero(Timeout(env.DatabaseConnection), defaultTimeouts.DatabaseConnection)
	timeouts.ServerReadHeader = utils.FirstNonZero(Timeout(env.ServerReadHeader), defaultTimeouts.ServerReadHeader)
	timeouts.ServerWrite = utils.FirstNonZero(Timeout(env.ServerWrite), defaultTimeouts.ServerWrite)
	timeouts.ServerIdle = utils.FirstNonZero(Timeout(env.ServerIdle), defaultTimeouts.ServerIdle)
	timeouts.PegoutDepositCheck = utils.FirstNonZero(Timeout(env.PegoutDepositCheck), defaultTimeouts.PegoutDepositCheck)
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(timeouts); err != nil {
		return ApplicationTimeouts{}, fmt.Errorf("error validating timeouts: %w", err)
	}
	return timeouts, nil
}
