package hal

import "scootin-aboot/models"

// Event represents a HAL (Hypertext Application Language) event resource.
type Event struct {
	*models.Event `json:",inline" doc:"Event resource"`

	Links Links `json:"_links" doc:"List of links"`
}
