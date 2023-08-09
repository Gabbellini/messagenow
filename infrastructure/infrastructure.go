package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/repositories"
	"messagenow/settings"
	"messagenow/usecases"
	http_pkg "messagenow/views"
	"net/http"
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
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(apiMiddleware)
	apiRouter.Use(authorizationMiddleware)

	setupAuthorizationModule(apiRouter, db)
	setupAPIModule(apiRouter, db)
	http_pkg.NewPWAWebModule("/static/", "./static").Setup(router)

	return nil
}

func setupAuthorizationModule(router *mux.Router, db *sql.DB) {
	loginRepository := repositories.NewLoginRepository(db)
	loginUseCase := usecases.NewLoginUseCase(loginRepository)
	http_pkg.NewAuthorizationHTTPModule(loginUseCase).Setup(router)
}

func setupAPIModule(router *mux.Router, db *sql.DB) {
	createMessageRepository := repositories.NewCreateMessageRepository(db)
	getMessagesRepository := repositories.NewGetMessagesRepository(db)
	createRoomRepository := repositories.NewCreateRoomRepository(db)
	getRoomsRepository := repositories.NewGetRoomsRepository(db)
	getRoomRepository := repositories.NewGetRoomRepository(db)
	joinRoomRepository := repositories.NewJoinRoomRepository(db)

	createRoomUseCase := usecases.NewCreateRoomUseCase(createRoomRepository)
	joinRoomUseCase := usecases.NewJoinRoomUseCase(joinRoomRepository)
	createMessageUseCase := usecases.NewCreateMessageUseCase(createMessageRepository, getRoomRepository)
	getMessagesUseCase := usecases.NewGetMessagesUseCase(getMessagesRepository)
	getRoomsUseCase := usecases.NewGetRoomsUseCase(getRoomsRepository)

	http_pkg.NewMessageHTTPModule(
		createMessageUseCase,
		getMessagesUseCase,
		createRoomUseCase,
		joinRoomUseCase,
		getRoomsUseCase,
	).Setup(router)

	loginRepository := repositories.NewGetUserRepository(db)
	loginUseCase := usecases.NewGetUserUseCase(loginRepository)
	http_pkg.NewUserHTTPModule(loginUseCase).Setup(router)
}

func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathTemplate, err := mux.CurrentRoute(r).GetPathTemplate()
		if err != nil {
			log.Println("[userAuthorizationMiddleware] Error", err)
			exceptions.HandleError(w, exceptions.NewForbiddenError(exceptions.ForbiddenMessage))
			return
		}

		if pathTemplate == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}

		//Check if the user has the cookie with the token
		cookie, err := r.Cookie("cookie")
		if err != nil {
			if err == http.ErrNoCookie {
				//If the user doesn't have the cookie, return an error
				log.Println("[authorizationMiddleware] Error views.ErrNoCookie", err)
				exceptions.HandleError(w, exceptions.NewUnauthorizedError(exceptions.UnauthorizedMessage))
				return
			}
			//If there is an error, return an error
			log.Println("[authorizationMiddleware] Error r.Cookie", err)
			exceptions.HandleError(w, exceptions.NewUnauthorizedError(exceptions.UnauthorizedMessage))
			return
		}

		token, err := getTokenFromCookie(cookie)
		if err != nil {
			log.Println("[authorizationMiddleware] Error getTokenFromCookie", err)
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

func apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Set content type to JSON
		w.Header().Set("Content-Type", "application/json")

		//Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func getTokenFromCookie(cookie *http.Cookie) (*jwt.Token, error) {
	secureCookie := securecookie.New([]byte("MESSAGE_NOW_SECRET_KEY"), nil)
	var tokenString string
	err := secureCookie.Decode("cookie", cookie.Value, &tokenString)
	if err != nil {
		log.Println("[login] Error Decode", err)
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Println("[login] token.Method.(*jwt.SigningMethodHMAC) !ok", err)
			return nil, errors.New("error parsing token")
		}
		return []byte("MESSAGE_NOW_SECRET_KEY"), nil
	})
	if err != nil {
		log.Println("[isCookieValid] Error parsing token", err)
		return nil, err
	}

	return token, nil
}
