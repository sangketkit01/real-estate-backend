package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
)


type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}

	log.Println(req.Username, req.Password)
	ctx.JSON(http.StatusOK, req)
}

type CreateUserRequst struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequst
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldErr := range validationErrors {
			log.Println("Field:", fieldErr.Field(), "Error:", fieldErr.Tag())
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Name:     req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Username or email already exists"})
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, _, err = server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
	}

	// ctx.SetCookie("token",token, 7 * 24 * 60 * 60 , "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, "Create account seccessfully")
}
