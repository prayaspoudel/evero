package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/prayaspoudel/modules/access/features/auth"
	"github.com/prayaspoudel/modules/access/middleware"
	"github.com/prayaspoudel/modules/access/model"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log         *logrus.Logger
	AuthUseCase *auth.AuthUseCase
	Validator   *validator.Validate
}

func NewAuthController(log *logrus.Logger, authUseCase *auth.AuthUseCase, validator *validator.Validate) *AuthController {
	return &AuthController{
		Log:         log,
		AuthUseCase: authUseCase,
		Validator:   validator,
	}
}

// WebResponse generic response wrapper
type WebResponse[T any] struct {
	Data   T      `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req model.RegisterUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.Validator.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := c.AuthUseCase.Register(&req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(WebResponse[*model.UserResponse]{
		Status: "success",
		Data:   user,
	})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req model.LoginUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.Validator.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	ipAddress := ctx.IP()
	response, err := c.AuthUseCase.Login(&req, ipAddress)
	if err != nil {
		return err
	}

	return ctx.JSON(WebResponse[*model.LoginResponse]{
		Status: "success",
		Data:   response,
	})
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	authCtx := middleware.GetAuth(ctx)
	if authCtx == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	token := ctx.Get("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := c.AuthUseCase.Logout(authCtx.UserID, token); err != nil {
		return err
	}

	return ctx.JSON(WebResponse[any]{
		Status: "success",
		Data:   fiber.Map{"message": "logged out successfully"},
	})
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	var req model.RefreshTokenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.Validator.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response, err := c.AuthUseCase.RefreshToken(&req)
	if err != nil {
		return err
	}

	return ctx.JSON(WebResponse[*model.LoginResponse]{
		Status: "success",
		Data:   response,
	})
}
