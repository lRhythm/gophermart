package rest

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lRhythm/gophermart/internal/models"
)

func (s *Server) setupHandlers() *Server {
	jwtmw := s.JWTMiddleware()
	router := s.app.Group("/api/user")
	router.Post("/register", s.userRegisterHandler)
	router.Post("/login", s.userAuthHandler)
	router.Use(jwtmw)
	router.Post("/orders", s.orderCreateHandler)
	router.Get("/orders", s.orderListHandler)
	router.Get("/balance", s.userBalanceGetHandler)
	router.Post("/balance/withdraw", s.userWithdrawGetHandler)
	router.Get("/withdrawals", s.userWithdrawalListHandler)
	return s
}

func (s *Server) userRegisterHandler(c *fiber.Ctx) error {
	var req userRegisterRequest

	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		// 400 - decode.
		return s.err(c, fiber.StatusBadRequest)
	}

	err = req.Validate()
	if err != nil {
		// 400 - validate.
		return s.err(c, fiber.StatusBadRequest)
	}

	userID, err := s.service.UserRegisterAndIDGet(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrorStorageConflict) {
			// 409 - user.login already exists.
			return s.err(c, fiber.StatusConflict)
		}
		// 500 - error.
		return s.err(c, fiber.StatusInternalServerError)
	}

	token, err := s.token(userID)
	if err != nil {
		// 500 - token signature error.
		return s.err(c, fiber.StatusInternalServerError)
	}

	s.headerAuth(c, token)
	// 200 - success.
	return s.resp(c, fiber.StatusOK, nil)
}

func (s *Server) userAuthHandler(c *fiber.Ctx) error {
	var req userAuthRequest

	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		// 400 - decode.
		return s.err(c, fiber.StatusBadRequest)
	}

	err = req.Validate()
	if err != nil {
		// 400 - validate.
		return s.err(c, fiber.StatusBadRequest)
	}

	userID, err := s.service.UserIDByCredentialsGet(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrorStorageNoRows) || errors.Is(err, models.ErrorServiceWrongPassword) {
			// 401 - user not exists or wrong user.password for user.login.
			return s.err(c, fiber.StatusUnauthorized)
		}
		// 500 - error.
		return s.err(c, fiber.StatusInternalServerError)
	}

	token, err := s.token(userID)
	if err != nil {
		// 500 - token signature error.
		return s.err(c, fiber.StatusInternalServerError)
	}

	s.headerAuth(c, token)
	// 200 - success.
	return s.resp(c, fiber.StatusOK, nil)
	//return s.resp(c, fiber.StatusOK, []string{token}) // For local testing.
}

func (s *Server) orderCreateHandler(c *fiber.Ctx) error {
	orderID := string(c.Body())

	if orderID == "" {
		// 400 - bad format.
		return s.err(c, fiber.StatusBadRequest)
	}

	if !luhnValid(orderID) {
		// 422 - validate.
		return s.resp(c, fiber.StatusUnprocessableEntity, nil)
	}

	err := s.service.OrderCreate(orderID, s.userID(c))
	if err != nil {
		if errors.Is(err, models.ErrorServiceOrderAlreadyExistsCurrentUser) {
			// 200 - already created by current user.
			return s.resp(c, fiber.StatusOK, nil)
		}
		if errors.Is(err, models.ErrorServiceOrderAlreadyExistsOtherUser) {
			// 409 - already created by other user.
			return s.resp(c, fiber.StatusConflict, nil)
		}
		// 500 - error.
		return s.err(c, fiber.StatusInternalServerError)
	}

	// 202 - success.
	return s.resp(c, fiber.StatusAccepted, nil)
}

func (s *Server) orderListHandler(c *fiber.Ctx) error {
	orders, err := s.service.OrderList(s.userID(c))
	if err != nil {
		// 500 - error.
		return s.err(c, fiber.StatusInternalServerError)
	}

	if len(orders) == 0 {
		// 204 - no orders.
		return s.resp(c, fiber.StatusNoContent, nil)
	}

	// 200 - success.
	return s.resp(c, fiber.StatusOK, newOrdersListResponse(orders))
}

func (s *Server) userBalanceGetHandler(c *fiber.Ctx) error {
	current, withdrawn, err := s.service.UserBalanceGet(s.userID(c))
	if err != nil {
		// 500 - error.
		return s.err(c, fiber.StatusInternalServerError)
	}

	// 200 - success.
	return s.resp(c, fiber.StatusOK, newUserBalanceResponse(current, withdrawn))
}

func (s *Server) userWithdrawGetHandler(c *fiber.Ctx) error {
	var req userBalanceWithdrawRequest

	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		// 400 - decode (not required).
		return s.err(c, fiber.StatusBadRequest)
	}

	err = req.Validate()
	if err != nil {
		// 422 - validate.
		return s.resp(c, fiber.StatusUnprocessableEntity, nil)
	}

	err = s.service.UserBalanceWithdraw(req.OrderID, req.Sum, s.userID(c))
	if err != nil {
		if errors.Is(err, models.ErrorServiceUserBalanceInsufficientFunds) {
			// 402 - insufficient funds.
			return s.err(c, fiber.StatusPaymentRequired)
		}
		// 500 - error.
		// В т.ч. повторная попытка списания по заказу.
		return s.err(c, fiber.StatusInternalServerError)
	}

	// 200 - success.
	return s.resp(c, fiber.StatusOK, nil)
}

func (s *Server) userWithdrawalListHandler(c *fiber.Ctx) error {
	withdrawals, err := s.service.UserWithdrawalList(s.userID(c))
	if err != nil {
		// 500 - error.
		return s.err(c, fiber.StatusInternalServerError)
	}

	if len(withdrawals) == 0 {
		// 204 - no withdrawals.
		return s.resp(c, fiber.StatusNoContent, nil)
	}

	// 200 - success.
	return s.resp(c, fiber.StatusOK, newUserWithdrawalListResponse(withdrawals))
}
