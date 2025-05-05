package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
	"github.com/sangketkit01/real-estate-backend/util"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, newErrorResponse("Invalid request."))
		return
	}

	user, err := server.store.LoginUser(ctx, req.Username)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusUnauthorized, newErrorResponse("Invalid credentials"))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = util.CheckPassword(user.Password, req.Password) ; err != nil{
		if err == bcrypt.ErrMismatchedHashAndPassword{
			ctx.JSON(http.StatusUnauthorized, newErrorResponse("Invalid password"))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	token, _, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.SetCookie("token",token, 7 * 24 * 60 * 60, "/" , "localhost", server.isSecure, true)
	ctx.JSON(http.StatusOK, messageResponse("Login successfully."))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid request.")))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Name:     req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				ctx.JSON(http.StatusForbidden, errorResponse(errors.New("Username or Email already exist.")))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, _, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.SetCookie("token",token, 7 * 24 * 60 * 60 , "/", "localhost", server.isSecure, true)
	ctx.JSON(http.StatusOK, messageResponse("Create account successfully."))
}
