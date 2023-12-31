package main

import (
	"log"
	"os"
	"time"

	// "fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"

	"go-fiber-auth/apis/account"
	"go-fiber-auth/apis/auth"
	"go-fiber-auth/apis/chat"
	"go-fiber-auth/apis/index"
	"go-fiber-auth/apis/proxy"
	"go-fiber-auth/apis/websocket"

	"go-fiber-auth/configuration"
	"go-fiber-auth/database"
	"go-fiber-auth/utilities"
)

func main() {
	// load environment variables via the .env file
	envError := godotenv.Load()
	if envError != nil {
		log.Fatal(envError)
		return
	}

	// connect to the database
	dbError := database.Connect()
	if dbError != nil {
		log.Fatal(dbError)
		return
	}

	app := fiber.New()

	// middlewares
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
	}))
	app.Use(helmet.New())
	app.Use(limiter.New(limiter.Config{
		Max:      60,
		Duration: 60 * time.Second,
		LimitReached: func(ctx *fiber.Ctx) error {
			return utilities.Response(utilities.ResponseParams{
				Ctx:    ctx,
				Info:   configuration.ResponseMessages.TooManyRequests,
				Status: fiber.StatusTooManyRequests,
			})
		},
	}))
	app.Use(logger.New())

	// available APIs
	account.Setup(app)
	auth.Setup(app)
	chat.Setup(app)
	index.Setup(app)
	websocket.Setup(app)
	proxy.Setup(app)

	// app.Post("/post", func(c *fiber.Ctx) error {
	// 	payload := struct {
	// 		Message  string `json:"message"`
	// 	}{}
	// 	fmt.Println("Raw Body111: %s\n", c.Body())
	// 	if err := c.BodyParser(&payload); err != nil {
	// 		return err
	// 	}
	// 	fmt.Println("Raw Body111: %s\n", c.JSON(payload))

	// 	 return c.JSON(payload)
	// });

	// handle 404
	app.Use(func(ctx *fiber.Ctx) error {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.NotFound,
			Status: fiber.StatusNotFound,
		})
	})

	// get the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "9119"
	}

	// launch the app
	launchError := app.Listen("0.0.0.0:" + port)
	if launchError != nil {
		panic(launchError)
	}
}
