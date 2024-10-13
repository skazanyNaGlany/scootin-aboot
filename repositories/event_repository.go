package repositories

import (
	"gorm.io/gorm"

	"scootin-aboot/lerrors"
	"scootin-aboot/models"
)

// EventRepository represents a repository for managing events in the application.
type EventRepository struct {
	DB *gorm.DB
}

// Create inserts a new event into the database.
// It takes a pointer to a models.Event object as a parameter and returns an error, if any.
func (r *EventRepository) Create(event *models.Event) error {
	result := r.DB.Create(event)

	if result.RowsAffected == 0 {
		result.Error = lerrors.ErrDBNoRowsAffected
	}

	if result.RowsAffected > 1 {
		result.Error = lerrors.ErrDBMoreThan1RowsAffected
	}

	return result.Error
}

// CreateBatch inserts multiple events into the database.
// It takes a slice of events as input and inserts each event into the database using batch inserts.
// If any error occurs during the insertion, it returns the error.
func (r *EventRepository) CreateBatch(events []*models.Event) error {
	// batch inserts not working with latest GORM version?
	for _, ievent := range events {
		if err := r.Create(ievent); err != nil {
			return err
		}
	}

	return nil
}

// DeleteBatch deletes multiple events from the database.
// It takes a slice of events as input and deletes each event from the database.
// If any error occurs during the deletion process, it returns the error.
// Otherwise, it returns nil.
func (r *EventRepository) DeleteBatch(events []*models.Event) error {
	for _, ievent := range events {
		result := r.DB.Delete(&ievent)

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

// DeleteBatchByIDs deletes multiple events from the database based on their IDs.
// It takes a slice of int64 IDs as input and returns an error if any occurs.
func (r *EventRepository) DeleteBatchByIDs(ids []int64) error {
	result := r.DB.Where("id IN ?", ids).Delete(&models.Event{})

	if result.RowsAffected == 0 {
		result.Error = lerrors.ErrDBNoRowsAffected
	}

	return result.Error
}

// FindAll returns all events from the database.
func (r *EventRepository) FindAll() ([]*models.Event, error) {
	var events []*models.Event

	if err := r.DB.Order("created_at").Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}
