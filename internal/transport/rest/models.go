package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lRhythm/gophermart/internal/config"
	"github.com/lRhythm/gophermart/internal/models"
	"github.com/sirupsen/logrus"
	"time"
)

type servicer interface {
	tokener
	user
	order
}

type tokener interface {
	TokenGet(userID string, TTL time.Duration, secret string) (token string, err error)
}

type user interface {
	UserRegisterAndIDGet(login, password string) (userID string, err error)
	UserIDByCredentialsGet(login, password string) (userID string, err error)
	balancer
}

type order interface {
	OrderCreate(orderID, userID string) (err error)
	OrderList(userID string) (orders []models.OrderDTO, err error)
}

type balancer interface {
	UserBalanceGet(userID string) (current, withdrawn float32, err error)
	UserBalanceWithdraw(orderID string, sum float32, userID string) (err error)
	UserWithdrawalList(userID string) (withdrawals []models.WithdrawalDTO, err error)
}

type Server struct {
	app     *fiber.App
	logs    *logrus.Logger
	config  *config.Config
	service servicer
}
