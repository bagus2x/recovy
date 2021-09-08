package middleware

import (
	"strings"

	"github.com/bagus2x/recovy/app"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		bearer := strings.Split(authHeader, " ")
		if len(bearer) != 2 {
			return c.Status(401).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EUnauthorized,
					Messages: []string{"Invalid authorization header format"},
				},
			})
		}

		claims, err := m.authService.ExtractAccessToken(bearer[1])
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}
