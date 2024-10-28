package service

import (
	"github.com/lRhythm/gophermart/internal/config"
	"github.com/lRhythm/gophermart/internal/models"
	"github.com/sirupsen/logrus"
	"time"
)

type storager interface {
	storageUser
	storageOrder
	storageExchanger
}

type storageUser interface {
	UserCreateAndIDGet(login, passwordHash string) (userID string, err error)
	UserIDAndPasswordHashByLoginGet(login string) (userID, passwordHash string, err error)
	storageBalancer
}

type storageOrder interface {
	OrderUserIDGetAndInsert(orderID, userID string, uploadedAt time.Time) (exists bool, orderUserID string, err error)
	OrderList(userID string) (orders []models.OrderDTO, err error)
}

type storageBalancer interface {
	UserBalanceAndWithdrawnGet(userID string) (current, withdrawn float32, err error)
	UserBalanceGet(userID string) (balance float32, err error)
	UserBalanceWithdraw(orderID string, sum float32, userID string, ProcessedAt time.Time) (err error)
	UserWithdrawalList(userID string) (withdrawals []models.WithdrawalDTO, err error)
}

type storageExchanger interface {
	OrdersIDForProcessing(statuses []string) (ordersID []string, err error)
	OrderUnlockProcessing(orderID string) (err error)
	OrderStatusAndAccrualSet(orderID string, status string, accrual float32) (err error)
}

type clienter interface {
	AccrualOrderSend(address, orderID string) (err error)
	AccrualOrderGet(address string) (DTO models.AccrualOrderDTO, retryAfter uint, err error)
}

type Client struct {
	logs       *logrus.Logger
	config     *config.Accrual
	storage    storager
	client     clienter
	ordersIDCh chan string
}

func (c *Client) Close() {
	close(c.ordersIDCh)
}
