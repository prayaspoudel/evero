package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// NewFiberApp creates a new Fiber application instance based on configuration
func NewFiberApp(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewFiberErrorHandler(),
		Prefork:      config.GetBool("web.prefork"),
	})

	return app
}

// NewFiberErrorHandler creates a custom error handler for Fiber
func NewFiberErrorHandler() fiber.ErrorHandler {
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

// NewFiber is an alias for NewFiberApp for backwards compatibility
func NewFiber(config *viper.Viper) *fiber.App {
	return NewFiberApp(config)
}

// NewErrorHandler is an alias for NewFiberErrorHandler for backwards compatibility
func NewErrorHandler() fiber.ErrorHandler {
	return NewFiberErrorHandler()
}
