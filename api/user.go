package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
	"github.com/sangketkit01/real-estate-backend/util"
)

func (server *Server) GetUserData(c *fiber.Ctx) error {
	u := c.Locals("user")
	user, ok := u.(db.User)

	if !ok {
		return fiber.NewError(fiber.StatusForbidden, "invalid user type")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

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

	validator := validator.New()
	if err := validator.Struct(req) ; err != nil{
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := server.store.LoginUser(ctx.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := util.CheckPassword(user.Password, req.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
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
	return ctx.Status(401).JSON(fiber.Map{"message": "Login successfully"})
}

func (server *Server) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "lax",
		Secure:   server.isSecure,
	})

	c.ClearCookie("token")
	return okResponse(c, "logout successfully.")
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

	validator := validator.New()
	if err := validator.Struct(req) ; err != nil{
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

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Create account successfully"})
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"omitempty,email"`
	Phone string `json:"phone" validate:"omitempty,min=10,max=10"`
}

func (server *Server) UpdateUser(c *fiber.Ctx) error {
	user := c.Locals("user").(db.User)

	data := c.FormValue("data")
	var req UpdateUserRequest
	if err := sonic.Unmarshal([]byte(data), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validator := validator.New()
	if err := validator.Struct(req) ; err != nil{
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var imageUrl string
	var valid bool

	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["images"]
		for _, file := range files {
			if user.ProfileUrl.Valid {
				oldPath := fmt.Sprintf("../%s", user.ProfileUrl.String)
				_ = os.Remove(oldPath)
			}

			newProfile := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			dst := fmt.Sprintf("../uploads/%s", newProfile)
			if err := c.SaveFile(file, dst); err == nil {
				imageUrl = "uploads/" + newProfile
				valid = true
			}
		}
	}

	arg := db.UpdateUserParams{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		ProfileUrl: sql.NullString{String: imageUrl, Valid: valid},
		Username:   user.Username,
	}

	if err := server.store.UpdateUser(c.Context(), arg); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "update failed.")
	}

	return okResponse(c, "update profile successfully.")
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqField=NewPassword"`
}

func (server *Server) UpdateUserPassword(c *fiber.Ctx) error {
	user := c.Locals("user").(db.User)

	var req UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validator := validator.New()
	if err := validator.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	hashedPassword, err := server.store.GetUserPassword(c.Context(), user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "no data found, why are you here ?")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = util.CheckPassword(hashedPassword, req.CurrentPassword); err != nil {
		return fiber.NewError(fiber.StatusForbidden, "wrong password.")
	}

	newHashedPassword, err := util.HashedPassword(req.NewPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to hash password.")
	}

	arg := db.UpdateUserPasswordParams{
		Password: newHashedPassword,
		Username: user.Username,
	}

	if err = server.store.UpdateUserPassword(c.Context(), arg); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "update password failed.")
	}

	return okResponse(c, "update password successfully.")
}
