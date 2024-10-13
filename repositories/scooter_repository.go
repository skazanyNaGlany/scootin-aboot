package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	lerrors "scootin-aboot/lerrors"
	"scootin-aboot/models"
)

// ScooterRepository represents a repository for managing scooter data.
type ScooterRepository struct {
	DB *gorm.DB
}

// Create inserts a new scooter into the database.
// It takes a pointer to a models.Scooter object as a parameter and returns an error, if any.
func (r *ScooterRepository) Create(scooter *models.Scooter) error {
	result := r.DB.Create(scooter)

	if result.RowsAffected == 0 {
		result.Error = lerrors.ErrDBNoRowsAffected
	}

	if result.RowsAffected > 1 {
		result.Error = lerrors.ErrDBMoreThan1RowsAffected
	}

	return result.Error
}

// QueryScootersByLocationAndStatus queries scooters based on their location and status.
// It takes in the minimum and maximum latitude and longitude values to define the location range,
// and the status of the scooters to filter the results.
// It returns a slice of ScooterEvent pointers and an error if any occurred.
func (r *ScooterRepository) QueryScootersByLocationAndStatus(
	minLatitude, minLongitude, maxLatitude, maxLongitude float64,
	status string,
) ([]*models.ScooterEvent, error) {
	var maps []map[string]any
	var result []*models.ScooterEvent

	err := r.DB.Table("scooters").
		Select("scooters.id AS scooter__id, scooters.created_at AS scooter__created_at, scooters.updated_at AS scooter__updated_at, scooters.status AS scooter__status, scooters.user_id AS scooter__user_id, scooters.e_tag AS scooter__e_tag, events.id AS event__id, events.created_at AS event__created_at, events.updated_at AS event__updated_at, events.scooter_id AS event__scooter_id, events.user_id AS event__user_id, events.event_type AS event__event_type, events.latitude AS event__latitude, events.longitude AS event__longitude").
		Joins("JOIN events ON (scooters.id = events.scooter_id AND events.id = (SELECT MAX(events.id) FROM events WHERE events.scooter_id = scooters.id))").
		Where("events.latitude >= ? AND events.latitude <= ? AND events.longitude >= ? AND events.longitude <= ? AND scooters.status = ?", minLatitude, maxLatitude, minLongitude, maxLongitude, status).
		Order("events.id").
		Find(&maps).Error

	if err != nil {
		return nil, err
	}

	for _, imap := range maps {
		scooter := models.Scooter{}
		scooter.ID = uuid.MustParse(imap["scooter__id"].(string))
		scooter.CreatedAt = imap["scooter__created_at"].(time.Time)
		scooter.UpdatedAt = imap["scooter__updated_at"].(time.Time)
		scooter.Status = imap["scooter__status"].(string)
		scooter.UserID = uuid.MustParse(imap["scooter__user_id"].(string))
		scooter.ETag = uuid.MustParse(imap["scooter__e_tag"].(string))

		event := models.Event{}
		event.ID = imap["event__id"].(int64)
		event.CreatedAt = imap["event__created_at"].(time.Time)
		event.UpdatedAt = imap["event__updated_at"].(time.Time)
		event.ScooterID = uuid.MustParse(imap["event__scooter_id"].(string))
		event.UserID = uuid.MustParse(imap["event__user_id"].(string))
		event.EventType = imap["event__event_type"].(string)
		event.Latitude = imap["event__latitude"].(float64)
		event.Longitude = imap["event__longitude"].(float64)

		scooterEvent := models.ScooterEvent{}
		scooterEvent.Scooter = &scooter
		scooterEvent.Event = &event

		result = append(result, &scooterEvent)
	}

	return result, nil
}

// QueryScootersByStatus retrieves a list of scooters with the specified status.
// It takes a status string as a parameter and returns a slice of Scooter pointers
// and an error if any occurred during the query.
func (r *ScooterRepository) QueryScootersByStatus(status string) ([]*models.Scooter, error) {
	var scooters []*models.Scooter

	err := r.DB.Table("scooters").
		Select("scooters.*").
		Where("scooters.status = ?", status).
		Group("scooters.id").
		Find(&scooters).Error

	if err != nil {
		return nil, err
	}

	return scooters, nil
}

// CreateBatch inserts a batch of scooters into the database.
// It takes a slice of scooter models as input and inserts each scooter into the database using GORM.
// If any error occurs during the insertion, it returns the error.
// Otherwise, it returns nil.
func (r *ScooterRepository) CreateBatch(scooters []*models.Scooter) error {
	// batch inserts not working with latest GORM version?
	for _, scooter := range scooters {
		result := r.DB.Create(scooter)

		if result.RowsAffected == 0 {
			result.Error = lerrors.ErrDBNoRowsAffected
		}

		if result.RowsAffected > 1 {
			result.Error = lerrors.ErrDBMoreThan1RowsAffected
		}

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// DeleteBatch deletes a batch of scooters from the database.
// It takes a slice of scooter objects as input and deletes each scooter individually.
// If any error occurs during the deletion process, it returns the error.
// Otherwise, it returns nil.
func (r *ScooterRepository) DeleteBatch(scooters []*models.Scooter) error {
	for _, scooter := range scooters {
		result := r.DB.Delete(scooter)

		if result.RowsAffected == 0 {
			result.Error = lerrors.ErrDBNoRowsAffected
		}

		if result.RowsAffected > 1 {
			result.Error = lerrors.ErrDBMoreThan1RowsAffected
		}

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// Update updates the given scooter in the database.
// It returns an error if there was an issue updating the scooter.
func (r *ScooterRepository) Update(scooter *models.Scooter) error {
	result := r.DB.Save(scooter)

	if result.RowsAffected == 0 {
		result.Error = lerrors.ErrDBNoRowsAffected
	}

	if result.RowsAffected > 1 {
		result.Error = lerrors.ErrDBMoreThan1RowsAffected
	}

	return result.Error
}

func (r *ScooterRepository) UpdateWithETag(scooter *models.Scooter, etag uuid.UUID) error {
	result := r.DB.Model(scooter).Where("id = ? AND e_tag = ?", scooter.ID, etag).Updates(scooter)

	if result.RowsAffected == 0 {
		result.Error = lerrors.ErrDBNoRowsAffected
	}

	if result.RowsAffected > 1 {
		result.Error = lerrors.ErrDBMoreThan1RowsAffected
	}

	return result.Error
}

// FindByID retrieves a scooter from the database based on its ID.
// It returns a pointer to the found scooter and an error, if any.
func (r *ScooterRepository) FindByID(id uuid.UUID) (*models.Scooter, error) {
	var scooter models.Scooter
	if err := r.DB.First(&scooter, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &scooter, nil
}

// FindAll returns all the scooters from the database.
func (r *ScooterRepository) FindAll() ([]*models.Scooter, error) {
	var scooters []*models.Scooter

	if err := r.DB.Find(&scooters).Error; err != nil {
		return nil, err
	}
	return scooters, nil
}

// Count returns the total number of scooters in the database.
func (r *ScooterRepository) Count() (int64, error) {
	var count int64
	if err := r.DB.Model(&models.Scooter{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
