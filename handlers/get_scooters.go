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

	"github.com/danielgtaylor/huma/v2"
)

type GET_Scooters_Input struct {
	params.AuthorizationParam

	Status       string  `query:"status"        doc:"Filter by status"  enum:"free,occupied"`
	MinLatitude  float64 `query:"min_latitude"  doc:"Minimum latitude"`
	MinLongitude float64 `query:"min_longitude" doc:"Minimum longitude"`
	MaxLatitude  float64 `query:"max_latitude"  doc:"Maximum latitude"`
	MaxLongitude float64 `query:"max_longitude" doc:"Maximum longitude"`

	hasStatus   bool
	hasLocation bool
}

type GET_Scooters_Output struct {
	Body struct {
		EmbeddedScooters hal.EmbeddedScooters `json:"_embedded" doc:"Embedded resources"`
		Links            hal.Links            `json:"_links" doc:"List of links"`
	}
}

// Resolve resolves the GET_Scooters_Input by checking the location parameters and setting the status and location flags.
// It returns a list of errors encountered during the resolution process.
// The goal of this method is to check if the location and status parameters are provided and set the status and location flags accordingly.
func (i *GET_Scooters_Input) Resolve(ctx huma.Context) []error {
	log.Println("Resolving GET_Scooters_Input")

	if err := i.checkLocationParams(ctx); err != nil {
		log.Println("Error checking location params:", err)

		return []error{err}
	}

	if i.hasStatusAndLocationParams(ctx) {
		i.hasStatus = true
		i.hasLocation = true

		log.Println("Status and location parameters are provided.")
	} else if ctx.Query("status") != "" {
		i.hasStatus = true

		log.Println("Only status parameter is provided.")
	}

	return nil
}

// GET_Scooters retrieves a list of scooters based on the provided input parameters.
// If both status and location are provided, it queries the scooters by location and status.
// If only status is provided, it queries the scooters by status.
// If neither status nor location is provided, it retrieves all scooters.
// The function returns a list of scooters along with any error that occurred during the retrieval process.
func GET_Scooters(ctx context.Context, input *GET_Scooters_Input) (*GET_Scooters_Output, error) {
	var scooters []*models.Scooter
	var err error

	log.Println(
		"GET_Scooters called",
		input.Status,
		input.MinLatitude,
		input.MinLongitude,
		input.MaxLatitude,
		input.MaxLongitude,
	)

	if input.hasStatus && input.hasLocation {
		log.Println(
			"Querying scooters by location and status",
			input.Status,
			input.MinLatitude,
			input.MinLongitude,
			input.MaxLatitude,
			input.MaxLongitude,
		)

		scooterEvent, err := ScooterRepository.QueryScootersByLocationAndStatus(
			input.MinLatitude, input.MinLongitude, input.MaxLatitude, input.MaxLongitude, input.Status,
		)

		if err != nil {
			log.Println("Error querying scooters by location and status:", err)

			return nil, lerrors.ErrResInternalServerError
		}

		response := GET_Scooters_Output{}
		response.Body.Links.Self.Href = consts.SCOOTERS
		response.Body.EmbeddedScooters.Scooters = make([]hal.Scooter, 0)

		mergeScooterEventItems(scooterEvent, &response)

		return &response, nil
	} else if input.hasStatus {
		log.Println("Querying scooters by status", input.Status)

		scooters, err = ScooterRepository.QueryScootersByStatus(input.Status)

		if err != nil {
			log.Println("Error querying scooters by status:", err)

			return nil, lerrors.ErrResInternalServerError
		}
	} else {
		log.Println("Querying all scooters")

		scooters, err = ScooterRepository.FindAll()

		if err != nil {
			log.Println("Error querying all scooters:", err)

			return nil, lerrors.ErrResInternalServerError
		}
	}

	response := GET_Scooters_Output{}
	response.Body.Links.Self.Href = consts.SCOOTERS
	response.Body.EmbeddedScooters.Scooters = make([]hal.Scooter, 0)

	mergeScooterItems(scooters, &response)

	return &response, nil
}

// mergeScooterEventItems merges the scooter events with the response body.
// It takes a slice of scooter events and a pointer to the GET_Scooters_Output struct as input.
// For each scooter event, it creates a HALScooter object with the scooter details and links.
// If the scooter event has an associated event, it creates a HALEvent object with the event details and links.
// The HALScooter object is then appended to the EmbeddedScooters slice in the response body.
func mergeScooterEventItems(scooterEvent []*models.ScooterEvent, response *GET_Scooters_Output) {
	for _, item := range scooterEvent {
		jsonScooter := hal.Scooter{
			Scooter: item.Scooter,
			Links: hal.Links{
				Self: hal.Self{
					Href: strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", item.Scooter.ID.String()),
				},
			},
		}

		jsonScooter.EmbeddedEvents = new(hal.EmbeddedEvents)

		if item.Event != nil {
			halEvent := hal.Event{
				Event: item.Event,
				Links: hal.Links{
					Self: hal.Self{
						Href: strings.ReplaceAll(
							consts.EVENTS_ITEM,
							"{id}",
							strconv.FormatUint(uint64(item.Event.ID), 10),
						),
					},
				},
			}

			jsonScooter.EmbeddedEvents.Events = []hal.Event{halEvent}
		}

		response.Body.EmbeddedScooters.Scooters = append(response.Body.EmbeddedScooters.Scooters, jsonScooter)
	}
}

// mergeScooterItems merges the given list of scooters with the provided response object.
// It creates a JSON representation of each scooter and appends it to the embedded scooters list in the response.
// The JSON representation includes the scooter details and a self link with the corresponding scooter ID.
func mergeScooterItems(scooters []*models.Scooter, response *GET_Scooters_Output) {
	for _, item := range scooters {
		jsonScooter := hal.Scooter{
			Scooter: item,
			Links: hal.Links{
				Self: hal.Self{
					Href: strings.ReplaceAll(consts.SCOOTERS_ITEM, "{id}", item.ID.String()),
				},
			},
		}

		response.Body.EmbeddedScooters.Scooters = append(response.Body.EmbeddedScooters.Scooters, jsonScooter)
	}
}

// hasStatusAndLocationParams checks if the given context has all the status and location parameters.
// It returns true if all the parameters are present, otherwise false.
func (i *GET_Scooters_Input) hasStatusAndLocationParams(ctx huma.Context) bool {
	return ctx.Query("status") != "" && ctx.Query("min_latitude") != "" && ctx.Query("min_longitude") != "" &&
		ctx.Query("max_latitude") != "" && ctx.Query("max_longitude") != ""
}

// checkLocationParams checks if all the location parameters are provided in the query.
// It returns an error if any of the required parameters are missing.
func (i *GET_Scooters_Input) checkLocationParams(ctx huma.Context) error {
	if ctx.Query("min_latitude") != "" || ctx.Query("min_longitude") != "" || ctx.Query("max_latitude") != "" ||
		ctx.Query("max_longitude") != "" {
		if ctx.Query("min_latitude") == "" || ctx.Query("min_longitude") == "" || ctx.Query("max_latitude") == "" ||
			ctx.Query("max_longitude") == "" {
			return &huma.ErrorDetail{
				Message:  "When querying for the location you need to provide all of the following parameters: min_latitude, min_longitude, max_latitude, max_longitude",
				Location: "query",
			}
		}
	}

	return nil
}
