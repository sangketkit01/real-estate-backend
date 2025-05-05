package api

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
	"github.com/sangketkit01/real-estate-backend/util"
)

type Server struct {
	router     *gin.Engine
	store      *db.Store
	config     util.Config
	tokenMaker util.Maker
}

func NewServer(store *db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := util.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:  store,
		config: config,
		tokenMaker: tokenMaker,
	}

	err = server.setUpRoute()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (server *Server) Start() error {
	return server.router.Run(":8080")
}

func (server *Server) setUpRoute() error {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/", server.HomePage)
	router.POST("/login", server.LoginUser)
	router.POST("/create-account",server.CreateUser)

	server.router = router

	return nil
}
