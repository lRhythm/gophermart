package rest

import "errors"

type userRegisterRequest struct {
	userCredentialsRequest
}

type userAuthRequest struct {
	userCredentialsRequest
}

type userCredentialsRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *userCredentialsRequest) Validate() error {
	var errLogin, errPassword error
	if r.Login == "" {
		errLogin = errors.New("login required")
	}
	if len(r.Password) < 6 {
		errPassword = errors.New("password too short")
	}
	return errors.Join(errLogin, errPassword)
}

type userBalanceWithdrawRequest struct {
	OrderID string  `json:"order"`
	Sum     float32 `json:"sum"`
}

func (r *userBalanceWithdrawRequest) Validate() error {
	var errOrderID, errSum error
	if !luhnValid(r.OrderID) {
		errOrderID = errors.New("order wrong format")
	}
	// Not required.
	if r.Sum == 0 {
		errSum = errors.New("sum zero")
	}
	return errors.Join(errOrderID, errSum)
}
