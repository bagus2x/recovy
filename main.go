package main

import (
	"log"

	"github.com/bagus2x/recovy/app/middleware"
	"github.com/bagus2x/recovy/app/routes"
	"github.com/bagus2x/recovy/auth"
	"github.com/bagus2x/recovy/config"
	"github.com/bagus2x/recovy/db"
	"github.com/bagus2x/recovy/podcast"
	"github.com/bagus2x/recovy/starredpodcast"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()
	cache := db.Cache(cfg)
	db, err := db.OpenPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("openapi.html")
	})

	app.Use(cors.New(cors.ConfigDefault))
	app.Use(logger.New())

	authRepo := auth.NewRepository(db)
	authCacheRepo := auth.NewCacheRepository(cache, cfg)
	podcastRepo := podcast.NewRepository(db)
	starredPodcastRepo := starredpodcast.NewRepository(db)

	authService := auth.NewService(authRepo, authCacheRepo, cfg)
	podcastService := podcast.NewService(podcastRepo, starredPodcastRepo)

	mw := middleware.NewMiddleware(authService)

	routes.AuthRoutes(app, mw, authService)
	routes.PodcastRoutes(app, mw, podcastService)

	app.Listen(cfg.AppPort())
}
