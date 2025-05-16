package api

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	db "github.com/sangketkit01/real-estate-backend/db/sqlc"
)

type NewContactRequest struct{
	ContactName string `json:"contact_name" validate:"required"`
	ContactDetail string `json:"contact_detail" validate:"required"`
}

func (server *Server) AddNewContact(c *fiber.Ctx) error{
	assetId := c.Locals("asset_id").(int)
	var req NewContactRequest
	if err := c.BodyParser(&req) ; err != nil{
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validator := validator.New()
	if err := validator.Struct(req) ; err != nil{
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	arg := db.InsertAssetContactParams{
		AssetID: int64(assetId),
		ContactName: req.ContactName,
		ContactDetail: req.ContactDetail,
	}
	_, err := server.store.InsertAssetContact(c.Context(), arg)
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "add contact failed.")
	}

	return okResponse(c, "add contact successfully.")
}

type UpdateContactRequest struct{
	ContactName *string `json:"contact_name"`
	ContactDetail *string `json:"contact_detail"`
}

func (server *Server) UpdateContact(c *fiber.Ctx) error{
	contactId, err := strconv.Atoi(c.Params("contact_id"))
	if err != nil{
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	var req UpdateContactRequest
	if err := c.BodyParser(&req) ; err != nil{
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	validator := validator.New()
	if err := validator.Struct(req) ; err != nil{
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	arg := db.UpdateContactParams{
		ID: int64(contactId),
		ContactName: *req.ContactName,
		ContactDetail: *req.ContactDetail,
	}

	err = server.store.UpdateContact(c.Context(), arg)
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "update contact failed.")
	} 

	
	return okResponse(c, "update contact successfully.")
}


func (server *Server) DeleteContact(c *fiber.Ctx) error{
	contactId, err := strconv.Atoi(c.Params("contact_id"))
	if err != nil{
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	err = server.store.RemoveContact(c.Context(), int64(contactId))
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "update contact failed.")
	} 

	
	return okResponse(c, "delete contact successfully.")
}

func (server *Server) GetAssetContacts(c *fiber.Ctx) error{
	assetId, err := strconv.Atoi(c.Params("asset_id", "no asset id"))
	if err != nil{
		return fiber.NewError(fiber.StatusBadRequest, "invalid asset id.")
	}

	contacts, err := server.store.GetAssetContacts(c.Context(), int64(assetId))
	if err != nil{
		return fiber.NewError(fiber.StatusInternalServerError, "cannot get asset's contacts.")
	}

	if len(contacts) == 0{
		return fiber.NewError(fiber.StatusNotFound, "no contact found.")
	}

	return okAndJsonResponse(c, contacts)
}