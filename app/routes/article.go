package routes

import (
	"strconv"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/app/middleware"
	"github.com/bagus2x/recovy/article"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ArticleRoutes(r fiber.Router, mw *middleware.Middleware, service article.Service) {
	article := r.Group("/api/v1/article")
	articles := r.Group("/api/v1/articles")

	article.Post("/", mw.Auth(), createArticle(service))
	article.Get("/:articleID", getArticleByID(service))
	article.Delete("/:articleID", mw.Auth(), deleteArticle(service))
	articles.Get("/", getArticles(service))
	articles.Get("/:category", getArticlesByCategory(service))
}

func createArticle(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req article.CreateArticleReq
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

func getArticles(service article.Service) fiber.Handler {
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

func getArticleByID(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		articleID, _ := strconv.ParseInt(c.Params("articleID"), 10, 64)

		res, err := service.GetByID(c.Context(), articleID)
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

func getArticlesByCategory(service article.Service) fiber.Handler {
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

func deleteArticle(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		webinarID, _ := strconv.ParseInt(c.Params("articleID"), 10, 64)
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
