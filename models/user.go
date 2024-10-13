package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"         doc:"ID of the user (UUID)"`
	CreatedAt time.Time `                             json:"created_at" doc:"Time when the user was created"`
	UpdatedAt time.Time `                             json:"updated_at" doc:"Time when the user was last updated"`
}
