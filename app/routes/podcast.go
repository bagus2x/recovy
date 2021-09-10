package routes

import (
	"log"
	"strconv"

	"github.com/bagus2x/recovy/app"
	"github.com/bagus2x/recovy/app/middleware"
	"github.com/bagus2x/recovy/podcast"
	"github.com/gofiber/fiber/v2"
)

func PodcastRoutes(r fiber.Router, mw *middleware.Middleware, service podcast.Service) {
	v1 := r.Group("/api/v1/podcasts")

	v1.Post("/", mw.Auth(), createPodcast(service))
	v1.Get("/", mw.NullableAuth(), getPodcasts(service))
	v1.Get("/:podcastID", mw.NullableAuth(), getPodcast(service))
	v1.Delete("/:podcastID", mw.Auth(), deletePodcast(service))
	v1.Patch("/:podcastID/star", mw.Auth(), starPodcast(service))
	v1.Delete("/:podcastID/star", mw.Auth(), unstarPodcast(service))
}

func createPodcast(service podcast.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req podcast.CreatePodcastReq
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

		userID, _ := c.Locals("userID").(int64)
		req.AuthorID = userID

		res, err := service.Create(c.Context(), &req)
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

func starPodcast(service podcast.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := c.Locals("userID").(int64)
		podcastID, err := strconv.ParseInt(c.Params("podcastID"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		req := podcast.StarPodcastReq{
			PodcastID: podcastID,
			UserID:    userID,
		}

		err = service.StarPodcast(c.Context(), &req)
		log.Println(err)
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
		})
	}
}

func unstarPodcast(service podcast.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := c.Locals("userID").(int64)
		podcastID, err := strconv.ParseInt(c.Params("podcastID"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		req := podcast.StarPodcastReq{
			PodcastID: podcastID,
			UserID:    userID,
		}

		err = service.UnstarPodcast(c.Context(), &req)
		if err != nil {
			log.Println(err)
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

func getPodcast(service podcast.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		podcastID, err := strconv.ParseInt(c.Params("podcastID"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		res, err := service.GetByID(c.Context(), podcastID)
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

func getPodcasts(service podcast.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params, err := getPodcastParams(c.Query)
		if err != nil {
			return c.Status(400).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		res, err := service.GetByParams(c.Context(), &params)
		if err != nil {
			log.Println(err)
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

func deletePodcast(service podcast.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := c.Locals("userID").(int64)
		podcastID, err := strconv.ParseInt(c.Params("podcastID"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: app.ErrorMessage(err),
				},
			})
		}

		err = service.DeleteByPodcastIDAndAthorID(c.Context(), podcastID, userID)
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
		})
	}
}

func getPodcastParams(query func(key string, defaultValue ...string) string) (podcast.Params, error) {
	var params podcast.Params

	cursorStr := query("cursor")
	if len(cursorStr) > 0 {
		cursor, err := strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			return podcast.Params{}, err
		}
		params.Cursor = cursor
	}

	limitStr := query("limit")
	if len(limitStr) > 0 {
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return podcast.Params{}, err
		}
		params.Limit = limit
	}

	params.Direction = query("direction")

	log.Println(cursorStr)
	return params, nil
}
