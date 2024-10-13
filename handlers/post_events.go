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
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type POST_Events_Input struct {
	params.AuthorizationParam

	Body struct {
		ScooterID uuid.UUID `json:"scooter_id" doc:"ID of the scooter"`
		EventType string    `json:"event_type" doc:"Type of the event"      enum:"start,stop,location_update"`
		Latitude  float64   `json:"latitude"   doc:"Latitude of the event"`
		Longitude float64   `json:"longitude"  doc:"Longitude of the event"`
	}
}

type POST_Events_Output struct {
	Body hal.Event
}

// POST_Events handles the HTTP POST request for creating events.
func POST_Events(ctx context.Context, input *POST_Events_Input) (*POST_Events_Output, error) {
	log.Println(
		"POST_Events called",
		input.Body.ScooterID,
		input.Body.EventType,
		input.Body.Latitude,
		input.Body.Longitude,
	)

	scooter, err := ScooterRepository.FindByID(input.Body.ScooterID)

	if err != nil {
		log.Println("Error finding scooter:", err)

		return nil, huma.Error404NotFound("Scooter not found")
	}

	if scooter.Status != string(enums.ScooterStatusOccupied) {
		log.Println("Scooter is not occupied")

		return nil, huma.Error400BadRequest("Scooter is not occupied")
	}

	if scooter.UserID != input.Authorization {
		log.Println("Scooter is occupied by another user")

		return nil, huma.Error409Conflict("Scooter is occupied by another user")
	}

	_, err = UserRepository.FindByID(input.Authorization)

	if err != nil {
		log.Println("Error finding user:", err)

		return nil, huma.Error404NotFound("User not found")
	}

	event := models.Event{}
	event.ScooterID = input.Body.ScooterID
	event.UserID = input.Authorization
	event.EventType = input.Body.EventType
	event.Latitude = input.Body.Latitude
	event.Longitude = input.Body.Longitude

	if err := EventRepository.Create(&event); err != nil {
		log.Println("Error creating event:", err)

		return nil, lerrors.ErrResInternalServerError
	}

	response := POST_Events_Output{}
	response.Body.Event = &event
	response.Body.Links.Self.Href = consts.EVENTS_ITEM + strconv.FormatUint(uint64(event.ID), 10)

	return &response, nil
}
