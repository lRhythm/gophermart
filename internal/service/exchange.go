package service

import (
	"context"
	"github.com/lRhythm/gophermart/internal/models"
	"time"
)

// sendOrderToAccrualSystem - отправка заказа в систему расчёта начислений баллов лояльности.
func (c *Client) sendOrderToAccrualSystem(orderID string) {
	err := c.client.AccrualOrderSend(c.config.OrderSendAddress(), orderID)
	if err != nil {
		c.logs.Error(err)
	}
}

// Exchange - запуск при старте сервера процессов получения заказов для получения данных о начислениях:
// получение заказов из хранилища и получения данных из системы расчёта начислений баллов лояльности.
func (c *Client) Exchange(ctx context.Context) {
	go c.receivingOrdersFromStorageToProcessing(ctx)
	go c.getOrderFromAccrualSystem(ctx)
}

// receivingOrdersFromStorageToProcessing - запуск периодического получения id заказов из БД и отправки в канал, который
// вычитывается методом получения данных о начислениях из системы расчёта баллов лояльности.
func (c *Client) receivingOrdersFromStorageToProcessing(ctx context.Context) {
	ticker := time.NewTicker(c.config.Delay)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ordersID, err := c.ordersIDForProcessing()
			if err != nil {
				c.logs.Error(err)
			}
			c.ordersIDToProcessing(ordersID)
		}
	}
}

// ordersIDForProcessing - выборка id заказов из БД для обработки.
func (c *Client) ordersIDForProcessing() ([]string, error) {
	ordersID, err := c.storage.OrdersIDForProcessing(models.ServiceStatusesInProgress)
	if err != nil {
		return nil, err
	}
	return ordersID, nil
}

// ordersIDToProcessing - отправка идентификаторов заказ в канал для вычитки методом getOrderFromAccrualSystem().
func (c *Client) ordersIDToProcessing(ordersID []string) {
	for _, orderID := range ordersID {
		c.ordersIDCh <- orderID
	}
}

// getOrderFromAccrualSystem - вычитка из канала id заказов для получения данных из системы расчёта начислений баллов лояльности.
func (c *Client) getOrderFromAccrualSystem(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			return
		case orderID, ok := <-c.ordersIDCh:
			if !ok {
				break loop
			}
		retry:
			accrualOrderDTO, retryAfter, err := c.client.AccrualOrderGet(c.config.OrderGetAddress(orderID))
			if err != nil {
				c.logs.Error(err)
				err = c.storage.OrderUnlockProcessing(orderID)
				if err != nil {
					c.logs.Error(err)
				}
				break loop
			}
			if retryAfter > 0 {
				time.Sleep(time.Duration(retryAfter) * time.Second)
				goto retry
			}
			go c.orderStatusAndAccrualSet(accrualOrderDTO)
		}
	}
}

// orderStatusAndAccrualSet - вызывается после получения данных из системы расчёта начислений баллов лояльности для
// обновления соответствующих данных (начислений в заказе и изменения баланса у пользователя) в хранилище.
func (c *Client) orderStatusAndAccrualSet(accrualOrderDTO models.AccrualOrderDTO) {
	status := models.StatusesMap[accrualOrderDTO.Status]
	// Для статусов отличных от `PROCESSED` не должно быть начисления баллов.
	if status != models.ServiceStatusProcessed {
		accrualOrderDTO.Accrual = 0
	}
	err := c.storage.OrderStatusAndAccrualSet(accrualOrderDTO.ID, status, accrualOrderDTO.Accrual)
	if err != nil {
		c.logs.Error(err)
		err = c.storage.OrderUnlockProcessing(accrualOrderDTO.ID)
		if err != nil {
			c.logs.Error(err)
		}
	}
}
