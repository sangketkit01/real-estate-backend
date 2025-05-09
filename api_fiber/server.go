package apifiber

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
	"github.com/sangketkit01/real-estate-backend/util"
)

type Server struct {
	router     *fiber.App
	store      *db.Store
	config     util.Config
	tokenMaker util.Maker
	isSecure   bool
}

func NewServer(store *db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := util.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	server.isSecure = config.Environment == "production"

	err = server.setUpRoute()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (server *Server) Start() error {
	return server.router.Listen(":8080")
}

func (server *Server) setUpRoute() error {
	router := fiber.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	router.Get("/",server.HomePage)
	router.Post("/create-user",server.CreateUser)
	router.Post("/login-user",server.LoginUser)

	server.router = router

	return nil
}
