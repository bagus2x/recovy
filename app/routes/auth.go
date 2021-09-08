package routes

import (
	"net/http"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/app/middleware"
	"github.com/bagus2x/recovy/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(r fiber.Router, mw *middleware.Middleware, service auth.Service) {
	v1 := r.Group("/api/v1/auth")

	v1.Get("/", mw.Auth(), getAuthenticatedUser(service))
	v1.Post("/signup", signUp(service))
	v1.Post("/signin", signIn(service))
	v1.Delete("/signout", mw.Auth(), signOut(service))
	v1.Post("/refresh", refreshToken(service))
}

func getAuthenticatedUser(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := c.Locals("userID").(int64)

		res, err := service.GetUserByID(c.Context(), userID)
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		return c.Status(200).JSON(app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func signUp(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req auth.SignUpReq
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
		}

		res, err := service.SignUp(c.Context(), &req)
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		return c.Status(200).JSON(app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func signIn(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req auth.SignInReq
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
		}

		res, err := service.SignIn(c.Context(), &req)
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		return c.Status(200).JSON(app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func refreshToken(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req auth.RefreshTokenReq
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
		}

		res, err := service.RefreshToken(c.Context(), &req)
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		return c.Status(200).JSON(app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func signOut(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := c.Locals("userID").(int64)
		err := service.SignOut(c.Context(), userID)
		if err != nil {
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		return c.SendStatus(http.StatusNoContent)
	}
}
