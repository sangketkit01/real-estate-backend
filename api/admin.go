package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
)

func (server *Server) GetAllUsers(c *fiber.Ctx) error {

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	users, err := server.store.GetAllUsers(c.Context(), db.GetAllUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch users")
	}

	if len(users) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no users found")
	}

	return c.JSON(users)
}
