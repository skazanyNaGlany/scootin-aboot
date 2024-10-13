package handlers

import (
	"context"
	"log"
	"scootin-aboot/consts"
	"scootin-aboot/formats/hal"
	"scootin-aboot/lerrors"
	"scootin-aboot/models"
	"strings"

	"github.com/google/uuid"
)

type POST_Users_Input struct{}

type POST_Users_Output struct {
	Body hal.User
}

// POST_Users is a handler function that creates a new user.
func POST_Users(ctx context.Context, input *POST_Users_Input) (*POST_Users_Output, error) {
	log.Println("POST_Users called")

	user := models.User{ID: uuid.New()}

	if err := UserRepository.Create(&user); err != nil {
		log.Println("Error creating user:", err)

		return nil, lerrors.ErrResInternalServerError
	}

	response := POST_Users_Output{}
	response.Body.User = &user
	response.Body.Links.Self.Href = strings.ReplaceAll(consts.USERS_ITEM, "{id}", user.ID.String())

	return &response, nil
}
