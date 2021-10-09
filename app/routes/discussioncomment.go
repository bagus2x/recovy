package routes

import (
	"strconv"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/app/middleware"
	"github.com/bagus2x/recovy/discussioncomment"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func DiscussionCommentRoutes(r fiber.Router, mw *middleware.Middleware, service discussioncomment.Service) {
	discussionComment := r.Group("/api/v1/discussion-comment")
	discussionComments := r.Group("/api/v1/discussion-comments")

	discussionComment.Post("/", mw.Auth(), createDiscussionComment(service))
	discussionComment.Get("/:discussionCommentID", getDiscussionCommentByID(service))
	discussionComment.Delete("/:discussionCommentID", mw.Auth(), deleteDiscussionComment(service))
	discussionComments.Get("/", getDiscussionComments(service))
	discussionComments.Get("/:discussionID", getDiscussionCommentsByDiscussionID(service))
}

func createDiscussionComment(service discussioncomment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req discussioncomment.CreateDiscussionCommentReq
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
		req.CommentatorID = userID

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

func getDiscussionComments(service discussioncomment.Service) fiber.Handler {
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

func getDiscussionCommentByID(service discussioncomment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		discussionID, _ := strconv.ParseInt(c.Params("discussionCommentID"), 10, 64)

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

func getDiscussionCommentsByDiscussionID(service discussioncomment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		discussionCommentID, _ := strconv.ParseInt(c.Params("discussionID"), 10, 64)
		res, err := service.GetByDiscussionID(c.Context(), discussionCommentID)
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

func deleteDiscussionComment(service discussioncomment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		webinarID, _ := strconv.ParseInt(c.Params("discussionCommentID"), 10, 64)
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
