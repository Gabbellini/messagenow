package http

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/usecases"
	"net/http"
	"os"
)

type authorizationHttpModule struct {
	loginUseCase usecases.LoginUseCase
}

func NewAuthorizationHTTPModule(loginUseCase usecases.LoginUseCase) ModuleHTTP {
	return &authorizationHttpModule{
		loginUseCase: loginUseCase,
	}
}

func (m authorizationHttpModule) Setup(router *mux.Router) {
	router.HandleFunc("/login", m.login).Methods(http.MethodPost)
}

func (m authorizationHttpModule) login(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[login] Error ReadAll", err)
		exceptions.HandleError(w, err)
		return
	}

	var credentials entities.Credential
	if err = json.Unmarshal(b, &credentials); err != nil {
		log.Println("[login] Error Unmarshal", err)
		exceptions.HandleError(w, exceptions.NewForbiddenError(exceptions.ForbiddenMessage))
		return
	}

	user, err := m.loginUseCase.Execute(r.Context(), credentials)
	if err != nil {
		log.Println("[login] Error Execute", err)
		exceptions.HandleError(w, err)
		return
	}

	userByte, err := json.Marshal(*user)
	if err != nil {
		log.Println("[login] Error Marshal", err)
		exceptions.HandleError(w, err)
		return
	}

	// TODO: store secret key in a safe place.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": string(userByte),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("MESSAGE_NOW_SECRET_KEY")))
	if err != nil {
		log.Println("[login] Error SignedString", err)
		exceptions.HandleError(w, err)
		return
	}

	_, err = w.Write([]byte(tokenString))
	if err != nil {
		log.Println("[login] Error Write", err)
		exceptions.HandleError(w, err)
		return
	}
}
