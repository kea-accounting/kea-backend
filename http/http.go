// Copyright 2018 Mike Jarmy. All rights reserved.  Use of this
// source statusCode is governed by a MIT-style
// license that can be found in the LICENSE file.

package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kea-accounting/kea-backend/errors"
)

const ContentType = "Content-Type"
const ApplicationJson = "application/json; charset=utf-8"

// WriteJSON writes the provide value as json
func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}) {

	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		WriteError(w, errors.WrapError(err))
		return
	}

	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s\n", json)
}

// WriteError writes the error as a json value
func WriteError(w http.ResponseWriter, err *errors.Error) {

	log.Printf("%d %s", err.StatusCode, err.Message)

	w.WriteHeader(err.StatusCode)
	fmt.Fprintf(w, "%s\n", err.Error())
}
