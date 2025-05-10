package apifiber

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