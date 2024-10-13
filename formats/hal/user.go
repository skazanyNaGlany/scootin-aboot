package hal

import "scootin-aboot/models"

// User represents a user resource in HAL format.
type User struct {
	*models.User `json:",inline" doc:"User resource"`

	Links Links `json:"_links" doc:"List of links"`
}
