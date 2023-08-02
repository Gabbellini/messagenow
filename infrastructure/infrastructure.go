package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	http_pkg "messagenow/http"
	"messagenow/repositories"
	"messagenow/settings"
	"messagenow/usecases"
	"net/http"
	"os"
	"strings"
	"time"
)

func Setup(settings settings.Settings, router *mux.Router) error {
	db, err := setupDataBase(settings)
	if err != nil {
		log.Println("[Setup] Error setupDataBase", err)
		return err
	}

	err = setupModules(router, db)
	if err != nil {
		log.Println("[Setup] Error setupModules", err)
		return err
	}

	return nil
}

// SetupDataBase set the connection to the database and set connection entities.go.
func setupDataBase(settings settings.Settings) (*sql.DB, error) {
	db, err := sql.Open("mysql", settings.GetDBSource())
	if err != nil {
		log.Println("[Setup] Error connecting to database", err)
		return nil, err
	}

	// Limit the amount of time the connections are kept in the pool
	db.SetConnMaxLifetime(time.Minute * 10)

	// Limit the number of connections stored in the pool
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)

	return db, nil
}

// setupModules set the MVC structure for the application.
func setupModules(router *mux.Router, db *sql.DB) error {
	router.Use(rootMiddleware)
	setupAuthorizationModule(router, db)
	setupAPIModule(router, db)
	return nil
}

func setupAuthorizationModule(router *mux.Router, db *sql.DB) {
	loginRepository := repositories.NewLoginRepository(db)
	loginUseCase := usecases.NewLoginUseCase(loginRepository)
	http_pkg.NewAuthorizationHTTPModule(loginUseCase).Setup(router)
}

func setupAPIModule(router *mux.Router, db *sql.DB) {
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(authorizationMiddleware)

	createTextMessageRepository := repositories.NewCreateTextMessageRepository(db)
	createTextMessageUseCase := usecases.CreateTextMessageUseCase(createTextMessageRepository)
	http_pkg.NewMessageHTTPModule(createTextMessageUseCase).Setup(apiRouter)
}

// rootMiddleware set the response content type for the api as json.
func rootMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Set the origin to allow all.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		//Set the valid methods to all.
		w.Header().Set("Access-Control-Allow-Methods", "*")

		//Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathTemplate, err := mux.CurrentRoute(r).GetPathTemplate()
		if err != nil {
			log.Println("[userAuthorizationMiddleware] Error", err)
			exceptions.HandleError(w, exceptions.NewForbiddenError(exceptions.ForbiddenMessage))
			return
		}

		if pathTemplate == "/user/login" {
			next.ServeHTTP(w, r)
			return
		}

		token, err := getTokenFromRequest(r)
		if err != nil {
			log.Println("[authorizationMiddleware] Error getTokenFromRequest", err)
			exceptions.HandleError(w, exceptions.NewForbiddenError(exceptions.ForbiddenMessage))
			return
		}

		//Check if the token is valid
		if !token.Valid {
			//If the token is not valid, return an error
			log.Println("[authorizationMiddleware] Error !token.Valid", err)
			exceptions.HandleError(w, exceptions.NewForbiddenError(exceptions.ForbiddenMessage))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			log.Println("[authorizationMiddleware] Error !ok && !token.Valid", err)
			exceptions.HandleError(w, exceptions.NewForbiddenError(exceptions.ForbiddenMessage))
			return
		}

		userString, ok := claims["user"]
		if !ok {
			log.Println("[authorizationMiddleware] Error !ok", err)
			exceptions.HandleError(w, exceptions.NewForbiddenError(exceptions.ForbiddenMessage))
			return
		}

		var user entities.User
		err = json.Unmarshal([]byte(userString.(string)), &user)
		if err != nil {
			log.Println("[authorizationMiddleware] Error strconv.Atoi", err)
			exceptions.HandleError(w, exceptions.NewForbiddenError(exceptions.ForbiddenMessage))
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		//Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// getTokenFromRequest get the token from the cookie.
func getTokenFromRequest(r *http.Request) (*jwt.Token, error) {
	splitToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(splitToken) != 2 {
		log.Println("[getTokenFromRequest] Error len(splitToken) == 0")
		return nil, errors.New("error Authorization Bearer not valid")
	}
	tokenString := splitToken[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Println("[getTokenFromRequest] token.Method.(*jwt.SigningMethodHMAC) !ok")
			return nil, errors.New("error parsing token")
		}
		return []byte(os.Getenv("MESSAGE_NOW_SECRET_KEY")), nil
	})
	if err != nil {
		log.Println("[getTokenFromRequest] Error parsing token", err)
		return nil, err
	}

	return token, nil
}
