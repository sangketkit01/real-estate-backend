package apifiber

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
)

func (server *Server) GetAllAssets(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page","1"))
	if err != nil || page < 1{
		page = 1
	}

	limit := 10
	offset := (page - 1) * limit

	arg := db.GetAllAssetsParams{
		Limit: int32(limit),
		Offset: int32(offset),
	}

	assets, err := server.store.GetAllAssets(c.Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "no asset found.")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "cannot get assets")
	}

	total, err := server.store.GetAssetCount(c.Context())
	if err !=  nil{
		return fiber.NewError(fiber.StatusInternalServerError, "cannot count asset")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assets" : assets,
		"page" : page,
		"limit" : limit,
		"total" : total,
	})
}

func (server *Server) GetAssetById(c *fiber.Ctx) error {
	assetIdString := c.Params("asset_id", "no asset id")
	assetId, err := strconv.Atoi(assetIdString)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	asset, err := server.store.GetAssetById(c.Context(), int64(assetId))
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "asset not found.")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "cannot not get asset.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"asset": asset})
}

func (server *Server) GetAssetsByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	if strings.TrimSpace(username) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "username is not provided.")
	}

	page, err := strconv.Atoi(c.Query("page","1"))
	if err != nil || page < 1{
		page = 1
	}

	limit := 10
	offset := (page - 1) * limit

	arg := db.GetAssetsByUsernameParams{
		Limit: int32(limit),
		Offset: int32(offset),
		Owner: username,
	}

	assets, err := server.store.GetAssetsByUsername(c.Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "no asset found for this person.")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "cannot get assets.")
	}

	total, err := server.store.GetAssetCountByUsername(c.Context(), username)
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "cannot count assets.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assets" : assets,
		"page" : page,
		"limit" : limit,
		"total" : total,
	})
}

func (server *Server) AllMyAssets(c *fiber.Ctx) error {
	user := c.Locals("user").(db.User)

	page, err := strconv.Atoi(c.Query("page","1"))
	if err != nil || page < 1{
		page = 1
	}

	limit := 10
	offset := (page - 1) * limit

	arg := db.GetAssetsByUsernameParams{
		Limit: int32(limit),
		Offset: int32(offset),
		Owner: user.Username,
	}

	assets, err := server.store.GetAssetsByUsername(c.Context(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "you have no asset.")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "cannot get assets.")
	}

	total, err := server.store.GetAssetCountByUsername(c.Context(), user.Username)
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "cannot count assets.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assets" : assets,
		"page" : page,
		"limit" : limit,
		"total" : total,
	})
}

func (server *Server) EditAsset(c *fiber.Ctx) error {
	assetIdString := c.Params("asset_id", "no asset id")
	assetId, err := strconv.Atoi(assetIdString)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	asset, err := server.store.GetAssetById(c.Context(), int64(assetId))
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "asset not found.")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "cannot not get asset.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"asset": asset})
}

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

type CreateAssetRequest struct {
	Asset    AssetRequest          `json:"asset" validate:"required"`
	Contacts []AssetContactRequest `json:"asset_contacts"`
}

func (server *Server) CreateAsset(c *fiber.Ctx) error {
	userData := c.Locals("user").(db.User)

	data := c.FormValue("data")
	var req CreateAssetRequest
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON data.")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assetArg := db.InsertAssetParams{
		Owner:  userData.Username,
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

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to read multipart form")
	}
	files := form.File["images"]

	for _, file := range files {
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		dst := fmt.Sprintf("../uploads/%s", uniqueName)

		if err := c.SaveFile(file, dst); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "upload failed")
		}

		server.store.InsertAssetImage(c.Context(), db.InsertAssetImageParams{
			AssetID:  asset.ID,
			ImageUrl: "uploads/" + uniqueName,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Create asset successfully."})
}

type UpdateAssetRequest struct {
	Price  *int64  `json:"price" validate:"omitempty,min=0"`
	Detail *string `json:"detail"`
}

func (server *Server) UpdateAsset(c *fiber.Ctx) error {
	assetId := c.Locals("asset_id").(int)

	var req UpdateAssetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	arg := db.UpdateAssetParams{
		ID:     int64(assetId),
		Price:  *req.Price,
		Detail: *req.Detail,
	}
	err := server.store.UpdateAsset(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return okResponse(c, "update asset successfully.")
}

func (server *Server) DeleteAsset(c *fiber.Ctx) error {
	assetId := c.Locals("asset_id").(int)
	err := server.store.DeleteAsset(c.Context(), int64(assetId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	return okResponse(c, "delete asset successfully.")
}
