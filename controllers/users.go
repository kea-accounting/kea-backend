package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/kea-accounting/kea-backend/data"
	"github.com/kea-accounting/kea-backend/datastructures"
	"github.com/kea-accounting/kea-backend/errors"
	ohttp "github.com/kea-accounting/kea-backend/http"
)

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func Login(signingKey []byte, w http.ResponseWriter, r *http.Request) {

	var login datastructures.User
	jsonParser := json.NewDecoder(r.Body)
	err := jsonParser.Decode(&login)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	user, err := data.GetUser(login.GetEmail())
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
	} else {

		if !comparePasswords(user.GetPassword(), []byte(login.GetPassword())) {
			ohttp.WriteError(w, errors.BadRequest(fmt.Errorf("Invalid Login")))
		} else {
			// Create a new token object, specifying signing method and the claims
			// you would like it to contain.
			jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user": user.GetId(),
			})
			ss, err := jwttoken.SignedString(signingKey)
			if err != nil {
				ohttp.WriteError(w, errors.WrapError(err))
				return
			}
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "jwttoken", Value: ss, Expires: expiration}
			http.SetCookie(w, &cookie)

			ohttp.WriteJSON(w, http.StatusOK, user)
		}
	}
}

func Signup(signingKey []byte, w http.ResponseWriter, r *http.Request) {

	var user datastructures.User
	jsonParser := json.NewDecoder(r.Body)
	err := jsonParser.Decode(&user)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	_, err = data.NewUser(&user)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
	} else {

		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": user.GetId(),
		})
		ss, err := jwttoken.SignedString(signingKey)
		if err != nil {
			ohttp.WriteError(w, errors.WrapError(err))
			return
		}
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "jwttoken", Value: ss, Expires: expiration}
		http.SetCookie(w, &cookie)

		ohttp.WriteJSON(w, http.StatusOK, user)
	}
}
