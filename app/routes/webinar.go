package routes

import (
	"strconv"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/app/middleware"
	"github.com/bagus2x/recovy/webinar"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func WebinarRoutes(r fiber.Router, mw *middleware.Middleware, service webinar.Service) {
	webinar := r.Group("/api/v1/webinar")
	webinars := r.Group("/api/v1/webinars")

	webinar.Post("/", mw.Auth(), createWebinar(service))
	webinar.Get("/:webinarID", getWebinarByID(service))
	webinar.Delete("/:webinarID", mw.Auth(), deleteWebinar(service))
	webinars.Get("/", getWebinars(service))
	webinars.Get("/:category", getWebinarsByCategory(service))
}

func createWebinar(service webinar.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req webinar.CreateWebinarReq
		err := c.BodyParser(&req)
		if err != nil {
			logrus.Error(err)
			return c.Status(app.Status(err)).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
		}

		userID, _ := c.Locals("userID").(int64)
		req.AuthorID = userID

		res, err := service.Create(c.Context(), &req)
		if err != nil {
			logrus.Error(err)
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

func getWebinars(service webinar.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		res, err := service.Get(c.Context())
		if err != nil {
			logrus.Error(err)
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

func getWebinarByID(service webinar.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		webinarID, _ := strconv.ParseInt(c.Params("webinarID"), 10, 64)

		res, err := service.GetByID(c.Context(), webinarID)
		if err != nil {
			logrus.Error(err)
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

func getWebinarsByCategory(service webinar.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		res, err := service.GetByCategory(c.Context(), c.Params("category"))
		if err != nil {
			logrus.Error(err)
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

func deleteWebinar(service webinar.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		webinarID, _ := strconv.ParseInt(c.Params("webinarID"), 10, 64)
		userID, _ := c.Locals("userID").(int64)

		err := service.Delete(c.Context(), webinarID, userID)
		if err != nil {
			logrus.Error(err)
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
		})
	}
}
