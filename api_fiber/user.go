package apifiber

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
	"github.com/sangketkit01/real-estate-backend/util"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (server *Server) LoginUser(ctx *fiber.Ctx) error {
	var req LoginUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	user, err := server.store.LoginUser(ctx.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = util.CheckPassword(user.Password, req.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	token, _, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(7 * 24 * 60 * 60),
		Secure:   server.isSecure,
		HTTPOnly: true,
		SameSite: "lax",
	})
	return ctx.JSON(fiber.Map{"message": "Login successfully"})
}

type CreateUserRequst struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

func (server *Server) CreateUser(ctx *fiber.Ctx) error {
	var req CreateUserRequst
	if err := ctx.BodyParser(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldErr := range validationErrors {
			log.Println("Field:", fieldErr.Field(), "Error:", fieldErr.Tag())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Name:     req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return fiber.NewError(fiber.StatusForbidden, "Username or email already exist")
			}
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	token, _, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(7 * 24 * 60 * 60),
		Secure:   server.isSecure,
		HTTPOnly: true,
		SameSite: "lax",
	})

	return ctx.JSON(fiber.Map{"message": "Create account successfully"})
}
