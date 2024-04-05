package liquidity_provider

import "errors"

var BadLoginError = errors.New("incorrect username or credentials")

type LoginUseCase struct{}

func NewLoginUseCase() *LoginUseCase {
	return &LoginUseCase{}
}

func (useCase *LoginUseCase) Run() error {
	// TODO add here the logic to login when implementing the login task
	return nil
}
