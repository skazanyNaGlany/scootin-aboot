package hal

import "scootin-aboot/models"

// Scooter represents a scooter resource in HAL format.
type Scooter struct {
	*models.Scooter `json:",inline" doc:"Scooter resource"`

	EmbeddedEvents *EmbeddedEvents `json:"_embedded,omitempty" doc:"Embedded resources"`
	Links          Links           `json:"_links"              doc:"List of links"`
}
