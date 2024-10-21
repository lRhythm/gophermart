package rest

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(s.config.JWT.Secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return s.err(c, fiber.StatusUnauthorized)
		},
	})
}
