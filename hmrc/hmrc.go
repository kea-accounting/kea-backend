package hmrc

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kea-accounting/kea-backend/data"
	"github.com/kea-accounting/kea-backend/util"

	"github.com/kea-accounting/kea-backend/errors"
	"github.com/kea-accounting/kea-backend/globals"
	ohttp "github.com/kea-accounting/kea-backend/http"
	"golang.org/x/oauth2"
)

// Authorize
func Authorize(oauthConfig *oauth2.Config, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	state := r.Form.Get("state")
	if state == "" {
		ohttp.WriteError(w, errors.BadRequest(fmt.Errorf("Missing State")))
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		ohttp.WriteError(w, errors.BadRequest(fmt.Errorf("Missing Code")))
		return
	}

	userID := r.Context().Value(globals.UserIDKey)
	user, err := data.GetUserById(userID.(string))

	if user.GetAccessToken() != state {
		ohttp.WriteError(w, errors.BadRequest(fmt.Errorf("Invalid State")))
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code, oauth2.SetAuthURLParam("client_id", oauthConfig.ClientID), oauth2.SetAuthURLParam("client_secret", oauthConfig.ClientSecret))
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken
	user.TokenType = token.TokenType
	user.TokenExpiry = token.Expiry.Format(time.RFC3339)

	_, err = data.SaveUser(user)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}
	http.Redirect(w, r, "/hello", http.StatusFound)
}

// Login
func Login(oauthConfig *oauth2.Config, w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(globals.UserIDKey)
	user, err := data.GetUserById(userID.(string))

	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}
	user.AccessToken = util.NewID()
	_, err = data.SaveUser(user)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}
	u := oauthConfig.AuthCodeURL(user.GetAccessToken())
	http.Redirect(w, r, u, http.StatusFound)
}

// Hello
func CallHMRC(oauthConfig *oauth2.Config, method string, api string, w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(globals.UserIDKey)
	user, err := data.GetUserById(userID.(string))

	t := new(oauth2.Token)
	t.AccessToken = user.AccessToken
	t.RefreshToken = user.RefreshToken
	expiryString := user.TokenExpiry
	expiryTime, err := time.Parse(time.RFC3339, expiryString)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
	}
	t.Expiry = expiryTime
	t.TokenType = user.TokenType
	client := oauthConfig.Client(r.Context(), t)
	uri := api + strings.Replace(r.RequestURI, "/hmrc", "/", 1)
	var req *http.Request
	if method == "GET" {
		req, err = http.NewRequest(method, uri, nil)
	} else {
		req, err = http.NewRequest(method, uri, r.Body)
	}
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}
	req.Header.Add("Accept", "application/vnd.hmrc.1.0+json")
	test := req.Header.Get("Gov-Test-Scenario")
	if test != "" {
		req.Header.Add("Gov-Test-Scenario", test)
	}
	resp, err := client.Do(req)
	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ohttp.WriteError(w, errors.WrapError(err))
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error %s %s", resp.StatusCode, body)
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

}
