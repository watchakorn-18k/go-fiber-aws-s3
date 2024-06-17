package main

import (
	"go-fiber-aws-s3/configuration"
	ds "go-fiber-aws-s3/domain/datasources"
	repo "go-fiber-aws-s3/domain/repositories"
	gw "go-fiber-aws-s3/src/gateways"
	"go-fiber-aws-s3/src/middlewares"
	sv "go-fiber-aws-s3/src/services"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	// // // remove this before deploy ###################
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// /// ############################################

	app := fiber.New(configuration.NewFiberConfiguration())
	middlewares.Logger(app)
	app.Use(recover.New())
	app.Use(cors.New())

	mongodb := ds.NewMongoDB(10)

	userRepo := repo.NewUsersRepository(mongodb)

	sv0 := sv.NewUsersService(userRepo)

	gw.NewHTTPGateway(app, sv0)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}
