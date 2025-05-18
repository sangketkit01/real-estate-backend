package api

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
)

func (server *Server) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenCookie := c.Cookies("token", "invalid token")
		payload, err := server.tokenMaker.VerifyToken(tokenCookie)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "invalid token"})
		}

		userData, err := server.store.GetUser(c.Context(), payload.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user not found, then why are you here ?"})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		c.Locals("user", userData)
		return c.Next()
	}
}

func (server *Server) AssetMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println(c.AllParams(), c.Route(), c.Path())
		assetIdParam := c.Params("asset_id")
		if assetIdParam == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "asset_id is required"})
		}

		assetId, err := strconv.Atoi(assetIdParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid asset_id"})
		}

		asset, err := server.store.GetAssetById(c.Context(), int64(assetId))
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "asset not found"})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		userData := c.Locals("user")
		user, ok := userData.(db.User)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user context error"})
		}

		if asset.Owner != user.Username {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not the owner, could not edit other's asset."})
		}

		c.Locals("asset_id", assetId)
		return c.Next()
	}
}

func (server *Server) AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userData := c.Locals("user")
		user, ok := userData.(db.User)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, "user context error")
		}

		if user.Roles != "admin" {
			return fiber.NewError(fiber.StatusUnauthorized, "not authorized")
		}

		return c.Next()
	}
}
