package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lRhythm/gophermart/internal/models"
)

func (s *Server) userID(c *fiber.Ctx) string {
	return c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)[models.JWTKeyUserID].(string)
}

func (s *Server) token(userID string) (string, error) {
	return s.service.TokenGet(userID, s.config.JWT.TTL, s.config.JWT.Secret)
}

func (s *Server) headerAuth(c *fiber.Ctx, token string) {
	c.Set(models.JWTHeaderKey, models.JWTHeaderValue(token))
}
