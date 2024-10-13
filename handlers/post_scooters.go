package handlers

import (
	"context"
	"log"
	"scootin-aboot/consts"
	"scootin-aboot/enums"
	"scootin-aboot/formats/hal"
	"scootin-aboot/lerrors"
	"scootin-aboot/models"
	"scootin-aboot/params"
	"strings"

	"github.com/google/uuid"
)

type POST_Scooters_Input struct {
	params.AuthorizationParam

	Body struct{}
}

type POST_Scooters_Output struct {
	Body hal.Scooter
}

// POST_Scooters updates the status of a scooter based on the provided input.
func POST_Scooters(ctx context.Context, input *POST_Scooters_Input) (*POST_Scooters_Output, error) {
	log.Println("POST_Scooters called")

	scooter := models.Scooter{ID: uuid.New()}
	scooter.Status = string(enums.ScooterStatusFree)
	scooter.ETag = uuid.New()

	if err := ScooterRepository.Create(&scooter); err != nil {
		log.Println("Error while creating scooter:", err)

		return nil, lerrors.ErrResInternalServerError
	}

	response := POST_Scooters_Output{}
	response.Body.Scooter = &scooter
	response.Body.Links.Self.Href = strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", scooter.ID.String())

	return &response, nil
}
