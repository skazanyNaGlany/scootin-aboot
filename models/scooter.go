package models

import (
	"time"

	"github.com/google/uuid"
)

// Scooter represents a scooter entity.
type Scooter struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;"             json:"id"         doc:"ID of the scooter (UUID)"`
	CreatedAt time.Time `                                         json:"created_at" doc:"Time when the scooter was created"`
	UpdatedAt time.Time `                                         json:"updated_at" doc:"Time when the scooter was last updated"`
	Status    string    `gorm:"type:varchar(50);index:,type:hash" json:"status"     doc:"Status of the scooter"                            enum:"occupied,free"`
	UserID    uuid.UUID `gorm:"type:uuid;index"                   json:"user_id"    doc:"ID of the user who is using the scooter (UUID)"`
	ETag      uuid.UUID `gorm:"type:uuid;"                        json:"etag"       doc:"ETag of the scooter, used for optimistic locking"`
}
