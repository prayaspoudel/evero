package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/prayaspoudel/modules/access/features/auth"
)

type AuthMiddleware struct {
	AuthUseCase *auth.AuthUseCase
}

func NewAuthMiddleware(authUseCase *auth.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		AuthUseCase: authUseCase,
	}
}

type AuthContext struct {
	UserID string
	Email  string
}

func (m *AuthMiddleware) Authenticate(ctx *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	// Extract token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format")
	}

	token := parts[1]

	// Verify token
	claims, err := m.AuthUseCase.VerifyAccessToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token")
	}

	// Set user context
	ctx.Locals("auth", &AuthContext{
		UserID: claims.UserID,
		Email:  claims.Email,
	})

	return ctx.Next()
}

// GetAuth retrieves auth context from fiber context
func GetAuth(ctx *fiber.Ctx) *AuthContext {
	auth, ok := ctx.Locals("auth").(*AuthContext)
	if !ok {
		return nil
	}
	return auth
}
