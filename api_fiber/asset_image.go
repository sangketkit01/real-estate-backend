package apifiber

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
)

type AddNewImageRequest struct{
	ImageUrl string `json:"image_url" validate:"required"`
}

func (server *Server) AddNewImage(c *fiber.Ctx) error{
	assetId := c.Locals("asset_id").(int)
	
	var req AddNewImageRequest
	if err := c.BodyParser(&req) ; err != nil{
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	arg := db.InsertAssetImageParams{
		AssetID: int64(assetId),
		ImageUrl: req.ImageUrl,
	}
	_, err := server.store.InsertAssetImage(c.Context(), arg)
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "add image failed.")
	}

	return okResponse(c, "add image successfully.")
}

func (server *Server) DeleteImage(c *fiber.Ctx) error{
	imageId, err := strconv.Atoi(c.Params("image_id", "no image id"))
	if err !=  nil{
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} 

	if err = server.store.DeleteImage(c.Context(), int64(imageId)); err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return okResponse(c, "delete image successfully.")
}