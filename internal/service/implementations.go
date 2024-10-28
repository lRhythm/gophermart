package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/lRhythm/gophermart/internal/models"
	"strings"
	"time"
)

// Token.

func (c *Client) TokenGet(userID string, TTL time.Duration, secret string) (string, error) {
	token, err := jwt.
		NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.MapClaims{
				models.JWTKeyExp:    time.Now().Add(TTL).Unix(),
				models.JWTKeyUserID: userID,
			},
		).
		SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// User.

func (c *Client) UserRegisterAndIDGet(login, password string) (string, error) {
	passwordHash, err := c.generateFromPassword(password)
	if err != nil {
		return "", err
	}
	id, err := c.storage.UserCreateAndIDGet(strings.ToLower(login), passwordHash)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (c *Client) UserIDByCredentialsGet(login, password string) (string, error) {
	id, passwordHash, err := c.storage.UserIDAndPasswordHashByLoginGet(strings.ToLower(login))
	if err != nil {
		return "", err
	}
	if !c.compareHashAndPassword(passwordHash, password) {
		return "", models.ErrorServiceWrongPassword
	}
	return id, nil
}

func (c *Client) UserBalanceGet(userID string) (float32, float32, error) {
	return c.storage.UserBalanceAndWithdrawnGet(userID)
}

func (c *Client) UserBalanceWithdraw(orderID string, sum float32, userID string) error {
	// Можно сделать транзакцией, но в этом случае логика сравнения sum и user.balance будет перенесена в слой БД.
	balance, err := c.storage.UserBalanceGet(userID)
	if err != nil {
		return err
	}
	if sum > balance {
		return models.ErrorServiceUserBalanceInsufficientFunds
	}
	return c.storage.UserBalanceWithdraw(orderID, sum, userID, time.Now())
}

func (c *Client) UserWithdrawalList(userID string) ([]models.WithdrawalDTO, error) {
	return c.storage.UserWithdrawalList(userID)
}

// Order.

func (c *Client) OrderCreate(orderID, userID string) error {
	// Без транзакции.
	//orderUserID, err := c.storage.OrderUserIDGet(orderID)
	//if err != nil {
	//	if errors.Is(err, models.ErrorStorageNoRows) {
	//		return c.storage.OrderInsert(orderID, userID, time.Now())
	//	}
	//	return err
	//}

	exists, orderUserID, err := c.storage.OrderUserIDGetAndInsert(orderID, userID, time.Now())
	if err != nil {
		return err
	}
	// Заказ не существовал до вставки, вставка успешна → отправка заказа в систему расчёта начислений баллов лояльности.
	if !exists {
		go c.sendOrderToAccrualSystem(orderID)
		return nil
	}
	// Заказ существует → сравнение пользователя заказа для определения статуса состояния HTTP ответа.
	if orderUserID == userID {
		return models.ErrorServiceOrderAlreadyExistsCurrentUser
	}
	return models.ErrorServiceOrderAlreadyExistsOtherUser
}

func (c *Client) OrderList(userID string) ([]models.OrderDTO, error) {
	orders, err := c.storage.OrderList(userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
