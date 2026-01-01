package setup

import (
"github.com/gofiber/fiber/v2"
"github.com/spf13/viper"
)

// NewFiber creates a new Fiber application instance based on configuration
func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
AppName:      config.GetString("app.name"),
ErrorHandler: NewErrorHandler(),
		Prefork:      config.GetBool("web.prefork"),
	})

	return app
}

// NewErrorHandler creates a custom error handler for Fiber
func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
"errors": err.Error(),
		})
	}
}
