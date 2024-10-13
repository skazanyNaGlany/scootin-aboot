package models

import (
	"time"

	"github.com/google/uuid"
)

// Event represents an scooter event.
type Event struct {
	ID        int64     `gorm:"primaryKey;"              json:"id"         doc:"ID of the event"`
	CreatedAt time.Time `gorm:"index:,sort:asc"          json:"created_at" doc:"Time when the event was created"`
	UpdatedAt time.Time `                                json:"updated_at" doc:"Time when the event was last updated"`
	ScooterID uuid.UUID `gorm:"type:uuid;not null;index" json:"scooter_id" doc:"ID of the scooter"`
	UserID    uuid.UUID `gorm:"type:uuid;index"          json:"user_id"    doc:"ID of the user who is using the scooter (UUID)"`
	EventType string    `gorm:"type:varchar(50);"        json:"event_type" doc:"Type of the event"                              enum:"start,stop,location_update"`
	Latitude  float64   `gorm:"index"                    json:"latitude"   doc:"Latitude of the event"`
	Longitude float64   `gorm:"index"                    json:"longitude"  doc:"Longitude of the event"`
}
