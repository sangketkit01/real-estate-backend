package api

import (
	"fmt"

	"github.com/bytedance/sonic"
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
	server.router = fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	server.router.Static("/static", "/uploads")

	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	server.setupPublicRoutes(server.router)
	server.setupProtectedRoutes(server.router)
	server.setupAdminRoute(server.router)

	return nil
}

func (server *Server) setupPublicRoutes(router *fiber.App) {
	router.Get("/", server.GetAllAssets)
	router.Post("/create-user", server.CreateUser)
	router.Post("/login-user", server.LoginUser)

	router.Get("/watch/:asset_id", server.GetAssetById)
	router.Get("/user/:username", server.GetAssetsByUsername)
}

func (server *Server) setupProtectedRoutes(router *fiber.App) {
	authGroup := router.Group("", server.AuthMiddleware())
	authGroup.Post("/create-asset", server.CreateAsset)
	authGroup.Get("/logout", server.Logout)

	authGroup.Get("/me", server.GetUserData)
	authGroup.Post("/update-profile", server.UpdateUser)
	authGroup.Post("/update-password", server.UpdateUserPassword)

	authGroup.Get("/my-asset", server.AllMyAssets)

	assetGroup := authGroup.Group("/asset")

	assetGroup.Get("/my-asset-detail/:asset_id", server.AssetMiddleware(), server.EditAsset)

	assetGroup.Put("/:asset_id", server.AssetMiddleware(), server.UpdateAsset)
	assetGroup.Delete("/:asset_id", server.AssetMiddleware(), server.DeleteAsset)

	assetGroup.Post("/:asset_id/add-contact", server.AssetMiddleware(), server.AddNewContact)
	assetGroup.Put("/:asset_id/:contact_id", server.AssetMiddleware(), server.UpdateContact)
	assetGroup.Delete("/:asset_id/:contact_id", server.AssetMiddleware(), server.DeleteContact)

	assetGroup.Post("/:asset_id/add-image", server.AssetMiddleware(), server.AddNewImage)
	assetGroup.Delete("/:asset_id/:image_id", server.AssetMiddleware(), server.DeleteImage)
}

func (server *Server) setupAdminRoute(router *fiber.App) {
	adminGroup := router.Group("/admin", server.AuthMiddleware(), server.AdminMiddleware())
	adminGroup.Get("/users", server.GetAllUsers)
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": msg,
	})
}
