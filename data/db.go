package data

import (
	"fmt"
	"log"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	ss "github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/golang/protobuf/proto"
	"github.com/kea-accounting/kea-backend/errors"
)

var db *fdb.Database
var subspaces *subspaceStruct

type subspaceStruct struct {
	invoice            ss.Subspace
	invoiceByClientID  ss.Subspace
	invoiceByStatus    ss.Subspace
	invoiceByVATPeriod ss.Subspace
	user               ss.Subspace
	userByKey          ss.Subspace
}

func InitDb(dataDb *fdb.Database) {
	db = dataDb
	dir, err := directory.CreateOrOpen(db, []string{"development"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	subspaces = &subspaceStruct{
		invoice:            dir.Sub("invoice"),
		invoiceByClientID:  dir.Sub("invoiceByClientID"),
		invoiceByStatus:    dir.Sub("invoiceByStatus"),
		invoiceByVATPeriod: dir.Sub("invoiceByVATPeriod"),
		user:               dir.Sub("user"),
		userByKey:          dir.Sub("userByKey"),
	}
}

func RunTransaction(object proto.Message, transaction func(out *[]byte, tr *fdb.Transaction) (interface{}, error)) (interface{}, error) {

	return db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		options := tr.Options()
		options.SetTimeout(2000)
		options.SetRetryLimit(3)

		out, err := proto.Marshal(object)
		if err != nil {
			return nil, errors.WrapError(err)
		}

		return transaction(&out, &tr)

	})

}

func GetById(key fdb.KeyConvertible) ([]byte, error) {
	i, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {

		bytes := tr.Get(key).MustGet()
		if bytes == nil {
			return nil, errors.NotFound(
				fmt.Errorf("Couldn't find '%s' does not exist", key))
		}

		return bytes, nil
	})

	if err != nil {
		return nil, err
	}

	return i.([]byte), nil
}

func GetByLookup(key fdb.KeyConvertible) ([]byte, error) {
	i, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {

		id := tr.Get(key).MustGet()
		if id == nil {
			return nil, errors.NotFound(
				fmt.Errorf("Couldn't find '%s' does not exist", key))
		}

		bytes := tr.Get(fdb.Key(id)).MustGet()
		if bytes == nil {
			return nil, errors.NotFound(
				fmt.Errorf("Couldn't find ID '%s' does not exist", id))
		}

		return bytes, nil
	})

	if err != nil {
		return nil, err
	}

	return i.([]byte), nil
}
