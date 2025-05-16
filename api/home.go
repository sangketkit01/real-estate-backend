package api

import "github.com/gofiber/fiber/v2"

func (server *Server) HomePage(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message" : "Hello, World!"})
}