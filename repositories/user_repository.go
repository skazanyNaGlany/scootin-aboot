package repositories

import (
	"scootin-aboot/lerrors"
	"scootin-aboot/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository represents a repository for managing user data.
type UserRepository struct {
	DB *gorm.DB
}

// Create inserts a new user record into the database.
// It takes a pointer to a User model as a parameter and returns an error, if any.
func (r *UserRepository) Create(user *models.User) error {
	result := r.DB.Create(user)

	if result.RowsAffected == 0 {
		result.Error = lerrors.ErrDBNoRowsAffected
	}

	if result.RowsAffected > 1 {
		result.Error = lerrors.ErrDBMoreThan1RowsAffected
	}

	return result.Error
}

// CreateBatch inserts multiple users into the database in a batch.
// It takes a slice of user models as input and returns an error if any
// error occurs during the insertion process.
func (r *UserRepository) CreateBatch(users []*models.User) error {
	// batch inserts not working with latest GORM version?
	for _, user := range users {
		err := r.Create(user)

		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteBatch deletes multiple users from the database.
// It takes a slice of user models as input and deletes each user from the database.
// If any error occurs during the deletion process, it returns the error.
// Otherwise, it returns nil.
func (r *UserRepository) DeleteBatch(users []*models.User) error {
	for _, user := range users {
		result := r.DB.Delete(&user)

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

// FindByID retrieves a user from the database based on the provided ID.
// It returns the user if found, otherwise it returns an error.
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteByID deletes a user from the database by their ID.
// It takes the ID of the user to be deleted as a parameter and returns an error if any.
func (r *UserRepository) DeleteByID(id uuid.UUID) error {
	result := r.DB.Delete(&models.User{}, "id = ?", id)

	if result.RowsAffected == 0 {
		result.Error = lerrors.ErrDBNoRowsAffected
	}

	if result.RowsAffected > 1 {
		result.Error = lerrors.ErrDBMoreThan1RowsAffected
	}

	return result.Error
}
