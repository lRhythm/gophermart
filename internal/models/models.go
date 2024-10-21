package models

type OrderDTO struct {
	ID         string
	UserID     string
	Status     string
	Accrual    float32
	UploadedAt string
}

type WithdrawalDTO struct {
	UserID      string
	OrderID     string
	Sum         float32
	ProcessedAt string
}

type AccrualOrderDTO struct {
	ID      string
	Status  string
	Accrual float32
}
