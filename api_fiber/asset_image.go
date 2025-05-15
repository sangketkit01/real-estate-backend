package apifiber

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
)

func (server *Server) AddNewImage(c *fiber.Ctx) error{
	assetId := c.Locals("asset_id").(int)
	
	form, err := c.MultipartForm()
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "failed to read multipart form.")
	}

	files := form.File["images"]
	for _, file := range files{
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		dst := fmt.Sprintf("../uploads/%s", uniqueName)
		if err := c.SaveFile(file, dst) ; err != nil{
			return fiber.NewError(fiber.StatusInternalServerError, "upload failed.")
		}

		arg := db.InsertAssetImageParams{
			AssetID: int64(assetId),
			ImageUrl: "uploads/" + uniqueName,
		} 

		_, err := server.store.InsertAssetImage(c.Context(), arg)
		if err != nil{
			return fiber.NewError(fiber.StatusInternalServerError, "add image failed.")
		}
	}


	return okResponse(c, "add image successfully.")
}

func (server *Server) DeleteImage(c *fiber.Ctx) error{
	imageId, err := strconv.Atoi(c.Params("image_id", "no image id"))
	if err !=  nil{
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} 

	imageData, err := server.store.GetImageById(c.Context(), int64(imageId))
	if err != nil{
		if err == sql.ErrNoRows{
			return fiber.NewError(fiber.StatusNotFound, "image not found.")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "cannot get image.")
	}

	filePath := fmt.Sprintf("../%s", imageData.ImageUrl)
	if err = os.Remove(filePath) ; err != nil{
		fmt.Printf("failed to delete file: %v\n", err)
	}

	if err = server.store.DeleteImage(c.Context(), int64(imageId)); err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return okResponse(c, "delete image successfully.")
}