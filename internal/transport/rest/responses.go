package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lRhythm/gophermart/internal/models"
)

func (s *Server) resp(c *fiber.Ctx, status int, body any) error {
	if body == nil {
		return c.Status(status).Send(nil)
	}
	return c.Status(status).JSON(body)
}

func (s *Server) err(c *fiber.Ctx, status int) error {
	return s.resp(c, status, nil)
}

type orderListItemResponse struct {
	ID         string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

func newOrdersListResponse(DTOs []models.OrderDTO) []orderListItemResponse {
	items := make([]orderListItemResponse, 0, len(DTOs))
	for _, DTO := range DTOs {
		items = append(items, orderListItemResponse{
			ID:         DTO.ID,
			Status:     DTO.Status,
			Accrual:    DTO.Accrual,
			UploadedAt: DTO.UploadedAt,
		})
	}
	return items
}

type userBalanceResponse struct {
	Balance   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

func newUserBalanceResponse(current, withdrawn float32) userBalanceResponse {
	return userBalanceResponse{
		Balance:   current,
		Withdrawn: withdrawn,
	}
}

type userWithdrawalListItemResponse struct {
	OrderID     string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func newUserWithdrawalListResponse(DTOs []models.WithdrawalDTO) []userWithdrawalListItemResponse {
	items := make([]userWithdrawalListItemResponse, 0, len(DTOs))
	for _, DTO := range DTOs {
		items = append(items, userWithdrawalListItemResponse{
			OrderID:     DTO.OrderID,
			Sum:         DTO.Sum,
			ProcessedAt: DTO.ProcessedAt,
		})
	}
	return items
}
