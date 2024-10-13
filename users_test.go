package main

import (
	"encoding/json"
	"net/http"
	"scootin-aboot/consts"
	"scootin-aboot/handlers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var usersTest = UsersTest{}

// UsersTest represents a test suite for user-related functionality.
type UsersTest struct {
	BaseTest
}

// TestPostUsers is a unit test function that tests the POST request for creating users.
// It sends a POST request to the "/users" endpoint with the provided authorization header and empty request body.
func (st *UsersTest) TestPostUsers(t *testing.T) {
	var responseMap map[string]any

	st.setup(t)
	defer st.teardown(t)

	response := st.wrappedAPI.Post(consts.USERS, struct{}{})

	assert.Equal(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &responseMap)

	st.Test200OKResponseMapItem(t, responseMap)
	st.TestUserMap(t, responseMap)

	userId := responseMap["id"].(string)

	err := handlers.UserRepository.DeleteByID(uuid.MustParse(userId))
	assert.Nil(t, err)
}

func TestPostUsers(t *testing.T) {
	usersTest.TestPostUsers(t)
}
