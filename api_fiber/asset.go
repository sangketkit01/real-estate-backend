package apifiber

import (
	"database/sql"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
)

type AssetRequest struct {
	Owner  string `json:"owner"`
	Price  int    `json:"price" validate:"required,min=0"`
	Detail string `json:"detail" validate:"required"`
}

type AssetContactRequest struct {
	AssetID       int64  `json:"asset_id" validate:"min=1"`
	ContactName   string `json:"contact_name" validate:"required"`
	ContactDetail string `json:"contact_detail" validate:"required"`
}

type AssetImageRequest struct {
	AssetID  int64  `json:"asset_id" validate:"min=1"`
	ImageUrl string `json:"image_url" validate:"required"`
}

type CreateAssetRequest struct {
	Asset    AssetRequest          `json:"asset" validate:"required"`
	Contacts []AssetContactRequest `json:"asset_contacts"`
	Images   []AssetImageRequest   `json:"asset_images" validate:"required"`
}

func (server *Server) CreateAsset(c *fiber.Ctx) error {
	tokenCookie := c.Cookies("token")
	payload, err := server.tokenMaker.VerifyToken(tokenCookie)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, "invalid token")
	}

	userData, err := server.store.GetUser(c.Context(), payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusForbidden, "user not found, then why are you here ?")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var req CreateAssetRequest
	if err = c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	req.Asset.Owner = userData.Username

	assetArg := db.InsertAssetParams{
		Owner:  req.Asset.Owner,
		Price:  int64(req.Asset.Price),
		Detail: req.Asset.Detail,
	}
	asset, err := server.store.InsertAsset(c.Context(), assetArg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "cannot create asset")
	}

	for _, contact := range req.Contacts {
		contactArg := db.InsertAssetContactParams{
			AssetID:       asset.ID,
			ContactName:   contact.ContactName,
			ContactDetail: contact.ContactDetail,
		}
		if _, err := server.store.InsertAssetContact(c.Context(), contactArg); err != nil {
			fmt.Printf("contact insert error: %v\n", err)
		}
	}

	for _, image := range req.Images {
		imageArg := db.InsertAssetImageParams{
			AssetID:  asset.ID,
			ImageUrl: image.ImageUrl,
		}
		if _, err := server.store.InsertAssetImage(c.Context(), imageArg); err != nil {
			fmt.Printf("image insert error: %v\n", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Create asset successfully."})
}
