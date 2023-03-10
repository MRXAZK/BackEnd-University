package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MRXAZK/BackEnd-University/controllers"
	"github.com/MRXAZK/BackEnd-University/initializers"
	"github.com/MRXAZK/BackEnd-University/middleware"

	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	initializers.ConnectDB(&config)
	initializers.ConnectRedis(&config)
}

func main() {
	app := fiber.New()
	micro := fiber.New()

	app.Mount("/api", micro)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
	}))

	app.All("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"Message": "Career Assistent BackEnd - University",
			"Author":  "Farhan Aulianda",
		})
	})

	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/register", controllers.SignUpUniversity)
		router.Post("/login", controllers.SignInUniversity)
		router.Get("/logout", middleware.DeserializeUniversity, controllers.LogoutUniversity)
		router.Get("/refresh", controllers.RefreshAccessToken)
	})

	micro.Get("/university/me", middleware.DeserializeUniversity, controllers.GetMe)

	ctx := context.TODO()
	value, err := initializers.RedisClient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": value,
		})
	})

	micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})

	log.Fatal(app.Listen(":8000"))
}
