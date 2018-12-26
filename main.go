package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"

	"github.com/kea-accounting/kea-backend/config"
	"github.com/kea-accounting/kea-backend/controllers"
	"github.com/kea-accounting/kea-backend/data"
	"github.com/kea-accounting/kea-backend/globals"
	"github.com/kea-accounting/kea-backend/hmrc"
	ohttp "github.com/kea-accounting/kea-backend/http"
)

var db fdb.Database
var oauthConfig = &oauth2.Config{
	Scopes: []string{"hello", "read:vat", "write:vat"},
}

var cfg config.Config

var (
	listenAddr string
	healthy    int32
)

// invoices
func listInvoices(w http.ResponseWriter, r *http.Request) {
	controllers.ListInvoices(w, r)
}
func createInvoice(w http.ResponseWriter, r *http.Request) {
	controllers.CreateInvoice(w, r)
}
func updateInvoice(w http.ResponseWriter, r *http.Request) {
	controllers.UpdateInvoice(w, r)
}
func getInvoice(w http.ResponseWriter, r *http.Request) {
	controllers.GetInvoice(w, r)
}
func getTotal(w http.ResponseWriter, r *http.Request) {
	controllers.TotalInvoicesByPeriod(w, r)
}

// user
func signup(w http.ResponseWriter, r *http.Request) {
	controllers.Signup(cfg.SigningKey, w, r)
}
func login(w http.ResponseWriter, r *http.Request) {
	controllers.Login(cfg.SigningKey, w, r)
}

// hmrc
func authorize(w http.ResponseWriter, r *http.Request) {
	hmrc.Authorize(oauthConfig, w, r)
}
func link(w http.ResponseWriter, r *http.Request) {
	hmrc.Login(oauthConfig, w, r)
}

func callHMRCGet(w http.ResponseWriter, r *http.Request) {
	hmrc.CallHMRC(oauthConfig, "GET", cfg.HMRC.APIBase, w, r)
}

func callHMRCPost(w http.ResponseWriter, r *http.Request) {
	hmrc.CallHMRC(oauthConfig, "POST", cfg.HMRC.APIBase, w, r)
}

func loadConfig() {
	config := config.LoadConfiguration("./config.json")
	_, e := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		log.Printf("Loading Config")
		oauthConfig.ClientID = config.HMRC.ClientID
		oauthConfig.ClientSecret = config.HMRC.ClientSecret
		oauthConfig.Endpoint.AuthURL = config.HMRC.APIBase + "oauth/authorize"
		oauthConfig.Endpoint.TokenURL = config.HMRC.APIBase + "oauth/token"
		oauthConfig.RedirectURL = config.BaseURL + "oauth2"
		cfg = config
		return nil, nil
	})
	if e != nil {
		log.Fatalf("Failed loading config: %s", e)
	}
}

func main() {

	// TODO configure these values
	port := ":8080"
	readTimeout := 15 * time.Second
	writeTimeout := 15 * time.Second
	idleTimeout := 60 * time.Second

	//---------------------------------

	fdb.MustAPIVersion(600)
	db = fdb.MustOpenDefault()
	options := fdb.Options()
	options.SetTraceEnable("")

	data.InitDb(&db)
	loadConfig()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	// set up http routing
	r := mux.NewRouter()
	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	authRouter := r.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		userID := r.Context().Value(globals.UserIDKey)
		return userID != nil
	}).Subrouter()
	authRouter.HandleFunc("/oauth2", authorize).Methods("GET")
	authRouter.HandleFunc("/link", link).Methods("GET")
	authRouter.HandleFunc("/invoice/{status:NEW|SENT}", listInvoices).Methods("GET")
	authRouter.HandleFunc("/invoice/total/{periodKey}", getTotal).Methods("GET")
	authRouter.HandleFunc("/invoice", createInvoice).Methods("POST")
	authRouter.HandleFunc("/invoice/{id}", getInvoice).Methods("GET")
	authRouter.HandleFunc("/invoice/{id}", updateInvoice).Methods("PUT")
	authRouter.HandleFunc("/hmrc/hello/user", callHMRCGet).Methods("GET")
	authRouter.HandleFunc("/hmrc/organisations/vat/{vrn}/obligations", callHMRCGet).Methods("GET")
	authRouter.HandleFunc("/hmrc/organisations/vat/{vrn}/returns", callHMRCGet).Methods("GET")
	authRouter.HandleFunc("/hmrc/organisations/vat/{vrn}/returns/{periodKey}", callHMRCPost).Methods("GET")
	authRouter.HandleFunc("/hmrc/organisations/vat/{vrn}/liabilities", callHMRCGet).Methods("GET")
	authRouter.HandleFunc("/hmrc/organisations/vat/{vrn}/payments", callHMRCGet).Methods("POST")

	r.Handle("/healthz", healthz())

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(ohttp.ContentType, ohttp.ApplicationJson)
			next.ServeHTTP(w, r)
		})
	})

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// start http server
	srv := &http.Server{
		Addr:         port,
		Handler:      authentication()(tracing(nextRequestID)(logging(logger)(r))),
		ErrorLog:     logger,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(globals.RequestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), globals.RequestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func authentication() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tokenCookie, err := r.Cookie("jwttoken")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			tokenString := tokenCookie.Value

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					next.ServeHTTP(w, r)
				}
				return cfg.SigningKey, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := claims["user"].(string)

				ctx := context.WithValue(r.Context(), globals.UserIDKey, userID)
				w.Header().Set("X-User-Id", userID)

				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r)
			}

		})
	}
}
