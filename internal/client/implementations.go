package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lRhythm/gophermart/internal/models"
	"io"
	"net/http"
)

func (c *Client) AccrualOrderSend(address, orderID string) error {
	b, err := json.Marshal(newOrderSendRequest(orderID))
	if err != nil {
		return err
	}
	resp, err := c.client.Post(address, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusAccepted:
		return nil // Success.

	case http.StatusConflict:
		return models.ErrorClientOrderAlreadyAccepted

	default:
		return fmt.Errorf("unexpected status: %d", resp.StatusCode) // 400, 500.
	}
}

func (c *Client) AccrualOrderGet(address string) (models.AccrualOrderDTO, uint, error) {
	var DTO models.AccrualOrderDTO

	resp, err := c.client.Get(address)
	if err != nil {
		return DTO, 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return DTO, 0, err
		}
		var o orderGetResponse
		err = json.Unmarshal(body, &o)
		if err != nil {
			return DTO, 0, err
		}
		DTO = o.toDTO()
		return DTO, 0, nil // Success.

	case http.StatusNoContent:
		return DTO, 0, models.ErrorClientOrderNotRegistered

	case http.StatusTooManyRequests:
		ra, err := respRetryAfter(resp)
		return DTO, ra, err

	default:
		return DTO, 0, fmt.Errorf("unexpected status: %d", resp.StatusCode) // 500.
	}
}
