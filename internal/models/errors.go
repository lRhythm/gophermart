package models

import "errors"

var (
	// storage

	ErrorStorageConflict = errors.New("conflict")
	ErrorStorageNoRows   = errors.New("no rows")

	// service

	ErrorServiceWrongPassword                 = errors.New("wrong password")
	ErrorServiceOrderAlreadyExistsCurrentUser = errors.New("order create: order already exists for current user")
	ErrorServiceOrderAlreadyExistsOtherUser   = errors.New("order create: order already exists for other user")
	ErrorServiceUserBalanceInsufficientFunds  = errors.New("user balance withdraw: insufficient funds")

	// client

	ErrorClientOrderAlreadyAccepted = errors.New("accrual: order already accepted")
	ErrorClientOrderNotRegistered   = errors.New("accrual: order not registered")
)
