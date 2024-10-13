package handlers

import (
	"context"
	"log"
	"scootin-aboot/consts"
	"scootin-aboot/formats/hal"
	"scootin-aboot/lerrors"
	"scootin-aboot/models"
	"scootin-aboot/params"
	"strconv"
	"strings"
)

type GET_Events_Input struct {
	params.AuthorizationParam
}

type GET_Events_Output struct {
	Body struct {
		EmbeddedEvents hal.EmbeddedEvents `json:"_embedded" doc:"Embedded resources"`
		Links          hal.Links          `json:"_links" doc:"List of links"`
	}
}

// GET_Events retrieves all events.
// It returns a list of events along with any error encountered.
func GET_Events(ctx context.Context, input *GET_Events_Input) (*GET_Events_Output, error) {
	log.Println("GET_Events called")

	items, err := EventRepository.FindAll()
	if err != nil {
		log.Println("Error retrieving events:", err)

		return nil, lerrors.ErrResInternalServerError
	}

	response := GET_Events_Output{}
	response.Body.Links.Self.Href = consts.EVENTS
	response.Body.EmbeddedEvents.Events = make([]hal.Event, 0)

	mergeEventItems(items, &response)

	return &response, nil
}

// mergeEventItems merges the given list of event items with the provided response.
// It creates a HAL event for each item and appends it to the embedded events in the response.
// The HAL event contains the event item and a self link with the corresponding item ID.
func mergeEventItems(items []*models.Event, response *GET_Events_Output) {
	for _, item := range items {
		jsonEvent := hal.Event{
			Event: item,
			Links: hal.Links{
				Self: hal.Self{
					Href: strings.ReplaceAll(consts.EVENTS_ITEM, "{id}", strconv.FormatUint(uint64(item.ID), 10)),
				},
			},
		}

		response.Body.EmbeddedEvents.Events = append(response.Body.EmbeddedEvents.Events, jsonEvent)
	}
}
