package routes

import (
	"strconv"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/app/middleware"
	"github.com/bagus2x/recovy/discussion"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func DiscussionRoutes(r fiber.Router, mw *middleware.Middleware, service discussion.Service) {
	discussion := r.Group("/api/v1/discussion")
	discussions := r.Group("/api/v1/discussions")

	discussion.Post("/", mw.Auth(), createDiscussion(service))
	discussion.Get("/:discussionID", getDiscussionByID(service))
	discussion.Delete(":discussionID", mw.Auth(), deleteDiscussion(service))
	discussions.Get("/", getDiscussions(service))
	discussions.Get("/:category", getDiscussionsByCategory(service))
}

func createDiscussion(service discussion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req discussion.CreateDiscussionReq
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

func getDiscussions(service discussion.Service) fiber.Handler {
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

func getDiscussionByID(service discussion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		discussionID, _ := strconv.ParseInt(c.Params("discussionID"), 10, 64)

		res, err := service.GetByID(c.Context(), discussionID)
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

func getDiscussionsByCategory(service discussion.Service) fiber.Handler {
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

func deleteDiscussion(service discussion.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		webinarID, _ := strconv.ParseInt(c.Params("discussionID"), 10, 64)
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
