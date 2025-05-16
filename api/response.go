package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func messageResponse(message string) fiber.Map{
	return fiber.Map{"message" : message}
}

func errorResponse(err error) fiber.Map{
	return fiber.Map{"error" : err.Error()}
}

func newErrorResponse(message string) fiber.Map {
	return fiber.Map{"error" : errors.New(message)}
}

func okResponse(c *fiber.Ctx, message string) error{
	return c.Status(fiber.StatusOK).JSON(messageResponse(message))
} 

func okAndJsonResponse(c *fiber.Ctx, data any) error{
	return c.Status(fiber.StatusOK).JSON(data)
}