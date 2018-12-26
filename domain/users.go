package domain

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/kea-accounting/kea-backend/datastructures"
	"github.com/kea-accounting/kea-backend/util"
	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func Validate(user *datastructures.User) (*datastructures.User, error) {
	var err error
	user.Id = util.NewID()
	user.Created, err = ptypes.TimestampProto(time.Now())
	user.LastUpdated, err = ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	user.Password = hashAndSalt([]byte(user.GetPassword()))

	return user, nil
}
