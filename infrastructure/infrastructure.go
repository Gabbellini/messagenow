package infrastructure

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	http_pkg "messagenow/http"
	"messagenow/repositories"
	"messagenow/settings"
	"messagenow/usecases"
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
	router.Use(rootMiddleware)
	setupCreateMessageModule(router, db)
	return nil
}

func setupCreateMessageModule(router *mux.Router, db *sql.DB) {
	createTextMessageRepository := repositories.NewCreateTextMessageRepository(db)
	createTextMessageUseCase := usecases.CreateTextMessageUseCase(createTextMessageRepository)
	http_pkg.NewMessageHTTPModule(createTextMessageUseCase).Setup(router)
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
