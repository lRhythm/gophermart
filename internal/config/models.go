package config

import (
	"net/url"
	"strings"
	"time"
)

const (
	a = "a" // flag: -a; env: RUN_ADDRESS.
	d = "d" // flag: -d; env: DATABASE_URI.
	r = "r" // flag: -r; env: ACCRUAL_SYSTEM_ADDRESS.
)

type Config struct {
	Server
	Database
	Accrual
	JWT
}

type Server struct {
	Address string `env:"RUN_ADDRESS"` // -a
}

type Database struct {
	DSN string `env:"DATABASE_URI"` // -d
}

type Accrual struct {
	Delay         time.Duration `env:"ACCRUAL_EXCHANGE_DELAY" envDefault:"5s"`                           // Периодичность запуска обработки получения информации о расчёте начислений баллов лояльности.
	Address       string        `env:"ACCRUAL_SYSTEM_ADDRESS"`                                           // -r
	OrderSendPath string        `env:"ACCRUAL_SYSTEM_ORDER_SEND_PATH" envDefault:"/api/orders"`          // Путь маршрута отправки информации для расчёта начислений баллов лояльности.
	OrderGetPath  string        `env:"ACCRUAL_SYSTEM_ORDER_GET_PATH" envDefault:"/api/orders/{orderId}"` // Путь маршрута получения информации о расчёте начислений баллов лояльности.
}

func (c *Accrual) OrderSendAddress() string {
	address, _ := url.JoinPath(c.Address, c.OrderSendPath)
	return address
}

func (c *Accrual) OrderGetAddress(orderID string) string {
	address, _ := url.JoinPath(c.Address, strings.Replace(c.OrderGetPath, "{orderId}", orderID, 1))
	return address
}

type JWT struct {
	Secret string        `env:"JWT_SECRET" envDefault:"MOr1oLgf2cSRst+Enpq7lA/5Hq/eGp0SCxz76NR8J+0="` // openssl rand -base64 32
	TTL    time.Duration `env:"JWT_TTL" envDefault:"30m"`
}
