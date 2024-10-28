package pg

import (
	"github.com/jmoiron/sqlx"
	"github.com/lRhythm/gophermart/internal/models"
	_ "github.com/lib/pq"
)

type DB struct {
	con *sqlx.DB
}

func (db *DB) Close() {
	_ = db.con.Close()
}

type user struct {
	ID           string `db:"id"`
	PasswordHash string `db:"password_hash"`
	Balance      uint   `db:"balance"`
	Withdrawn    uint   `db:"withdrawn"`
}

type order struct {
	ID         string `db:"id"`
	UserID     string `db:"user_id"`
	Status     string `db:"status"`
	Accrual    uint   `db:"accrual"`
	UploadedAt string `db:"uploaded_at"`
}

type orders []order

func (items *orders) toDTO() []models.OrderDTO {
	DTO := make([]models.OrderDTO, 0, len(*items))
	for _, item := range *items {
		DTO = append(DTO, models.OrderDTO{
			ID:         item.ID,
			UserID:     item.UserID,
			Status:     item.Status,
			Accrual:    s2r(item.Accrual),
			UploadedAt: item.UploadedAt,
		})
	}
	return DTO
}

type withdrawal struct {
	UserID      string `db:"user_id"`
	OrderID     string `db:"order_id"`
	Sum         uint   `db:"sum"`
	ProcessedAt string `db:"processed_at"`
}

type withdrawals []withdrawal

func (items *withdrawals) toDTO() []models.WithdrawalDTO {
	DTO := make([]models.WithdrawalDTO, 0, len(*items))
	for _, item := range *items {
		DTO = append(DTO, models.WithdrawalDTO{
			UserID:      item.UserID,
			OrderID:     item.OrderID,
			Sum:         s2r(item.Sum),
			ProcessedAt: item.ProcessedAt,
		})
	}
	return DTO
}
