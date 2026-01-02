package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prayaspoudel/modules/access/delivery/http"
	"github.com/prayaspoudel/modules/access/middleware"
)

type RouteConfig struct {
	App            *fiber.App
	AuthController *http.AuthController
	AuthMiddleware *middleware.AuthMiddleware
}

func (c *RouteConfig) Setup() {
	// API group
	api := c.App.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", c.AuthController.Register)
	auth.Post("/login", c.AuthController.Login)
	auth.Post("/refresh", c.AuthController.RefreshToken)

	// Protected routes
	auth.Post("/logout", c.AuthMiddleware.Authenticate, c.AuthController.Logout)

	// Health check
	c.App.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":  "ok",
			"service": "sso-access",
		})
	})
}
