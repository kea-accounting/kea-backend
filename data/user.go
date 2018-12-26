package data

import (
	"fmt"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
	"github.com/golang/protobuf/proto"
	"github.com/kea-accounting/kea-backend/datastructures"
	"github.com/kea-accounting/kea-backend/domain"
	"github.com/kea-accounting/kea-backend/errors"
)

func GetUser(email string) (*datastructures.User, error) {
	bytes, err := GetByLookup(subspaces.userByKey.Pack(tuple.Tuple{email}))
	if err != nil {
		return nil, errors.WrapError(err)
	}
	if bytes == nil {
		return nil, errors.NotFound(
			fmt.Errorf("User with email '%s' does not exist", email))
	}

	user := &datastructures.User{}
	if err := proto.Unmarshal(bytes, user); err != nil {
		return nil, errors.WrapError(err)
	}

	return user, nil
}

func GetUserById(id string) (*datastructures.User, error) {
	bytes, err := GetById(subspaces.user.Pack(tuple.Tuple{id}))
	if err != nil {
		return nil, errors.WrapError(err)
	}
	if bytes == nil {
		return nil, errors.NotFound(
			fmt.Errorf("User with id '%s' does not exist", id))
	}

	user := &datastructures.User{}
	if err := proto.Unmarshal(bytes, user); err != nil {
		return nil, errors.WrapError(err)
	}

	return user, nil
}

func NewUser(user *datastructures.User) (*datastructures.User, error) {

	_, err := domain.Validate(user)
	if err != nil {
		return nil, errors.WrapError(err)
	}

	i, err := RunTransaction(user, func(out *[]byte, tr *fdb.Transaction) (interface{}, error) {

		existing := tr.Get(subspaces.userByKey.Pack(tuple.Tuple{user.GetEmail()})).MustGet()

		if existing != nil {
			return nil, errors.Conflict(
				fmt.Errorf("User with email '%s' already exists", user.GetEmail()))
		}

		tr.Set(subspaces.userByKey.Pack(tuple.Tuple{user.GetEmail()}), subspaces.user.Pack(tuple.Tuple{user.GetId()}))

		tr.Set(subspaces.user.Pack(tuple.Tuple{user.GetId()}), *out)

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return i.(*datastructures.User), nil
}

func SaveUser(user *datastructures.User) (*datastructures.User, error) {

	i, err := RunTransaction(user, func(out *[]byte, tr *fdb.Transaction) (interface{}, error) {

		tr.Set(subspaces.user.Pack(tuple.Tuple{user.GetId()}), *out)

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return i.(*datastructures.User), nil
}
