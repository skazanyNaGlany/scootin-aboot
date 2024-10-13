package handlers

import (
	"context"
	"errors"
	"log"
	"scootin-aboot/consts"
	"scootin-aboot/enums"
	"scootin-aboot/formats/hal"
	"scootin-aboot/lerrors"
	"scootin-aboot/params"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type PATCH_Scooters_Input struct {
	params.ETagUUIDParam
	params.AuthorizationParam

	ID uuid.UUID `path:"id" doc:"Scooter ID"`

	Body struct {
		Status string `json:"status" doc:"Scooter status" enum:"occupied,free"`
	}
}

type PATCH_Scooters_Output struct {
	Body hal.Scooter
}

// PATCH_Scooters updates the status of a scooter based on the provided input.
func PATCH_Scooters(ctx context.Context, input *PATCH_Scooters_Input) (*PATCH_Scooters_Output, error) {
	log.Println("PATCH_Scooters called", input.ID, input.Body.Status, input.ETag)

	// find the scooter by ID
	scooter, err := ScooterRepository.FindByID(input.ID)

	if err != nil {
		log.Println("Error while looking for scooter: ", err)

		return nil, huma.Error404NotFound("Scooter not found")
	}

	// check if the ETag matches
	if input.ETag != scooter.ETag {
		log.Println("ETag does not match", input.ETag, scooter.ETag)

		return nil, huma.Error412PreconditionFailed("ETag does not match")
	}

	if input.Body.Status == string(enums.ScooterStatusOccupied) {
		// want to set scooter to occupied
		if scooter.Status == string(enums.ScooterStatusOccupied) {
			log.Println("Scooter is already occupied")

			return nil, huma.Error400BadRequest("Scooter is already occupied")
		}
	} else if input.Body.Status == string(enums.ScooterStatusFree) {
		// want to set scooter to free
		if scooter.Status != string(enums.ScooterStatusOccupied) {
			log.Println("Cannot free a scooter that is not occupied")

			return nil, huma.Error400BadRequest("Cannot free a scooter that is not occupied")
		}

		scooter.UserID = uuid.Nil
	}

	// update the editable fields
	scooter.Status = input.Body.Status
	scooter.UserID = input.Authorization
	scooter.ETag = uuid.New()

	// save the updated scooter
	err = ScooterRepository.UpdateWithETag(scooter, input.ETag)

	if err != nil {
		if errors.Is(err, lerrors.ErrDBNoRowsAffected) {
			log.Println("Error while updating scooter (wrong etag?):", err)

			return nil, huma.Error412PreconditionFailed("ETag does not match")
		}

		log.Println("Error while updating scooter: ", err)

		return nil, lerrors.ErrResInternalServerError
	}

	response := PATCH_Scooters_Output{}
	response.Body.Scooter = scooter
	response.Body.Links.Self.Href = strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", scooter.ID.String())

	return &response, nil
}
