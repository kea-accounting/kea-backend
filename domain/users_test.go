package domain

import (
	"testing"

	"github.com/kea-accounting/kea-backend/datastructures"
)

func TestValidate(t *testing.T) {
	user := &datastructures.User{
		Email:    "test@test.com",
		Password: "test",
	}
	Validate(user)
	assert(t, user.GetCreated() != nil)
	assert(t, user.GetEmail() == "test@test.com")
	assert(t, user.GetPassword() != "test")

}
