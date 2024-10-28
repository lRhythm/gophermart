package service

import (
	"github.com/lRhythm/gophermart/internal/config"
	"github.com/sirupsen/logrus"
)

func New(opts ...func(*Client)) *Client {
	c := new(Client)
	c.ordersIDCh = make(chan string)
	for _, o := range opts {
		o(c)
	}
	return c
}

func WithLogs(i *logrus.Logger) func(*Client) {
	return func(c *Client) {
		c.logs = i
	}
}

func WithConfig(i *config.Accrual) func(*Client) {
	return func(c *Client) {
		c.config = i
	}
}

func WithStorage(i storager) func(*Client) {
	return func(c *Client) {
		c.storage = i
	}
}

func WithClient(i clienter) func(*Client) {
	return func(c *Client) {
		c.client = i
	}
}
