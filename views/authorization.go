package views

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"io"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/usecases"
	"net/http"
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

	tokenString, err := token.SignedString([]byte("MESSAGE_NOW_SECRET_KEY"))
	if err != nil {
		log.Println("[login] Error SignedString", err)
		exceptions.HandleError(w, err)
		return
	}

	secureCookie := securecookie.New([]byte("MESSAGE_NOW_SECRET_KEY"), nil)
	encodedTokenString, err := secureCookie.Encode("cookie", tokenString)
	if err != nil {
		log.Println("[login] Error Encode", err)
		exceptions.HandleError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:  "cookie",
		Value: encodedTokenString,
	}

	http.SetCookie(w, cookie)

	b, err = json.Marshal(user)
	if err != nil {
		log.Println("[login] Error Marshal", err)
		exceptions.HandleError(w, exceptions.NewInternalServerError(exceptions.InternalErrorMessage))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[login] Error Write", err)
		exceptions.HandleError(w, exceptions.NewInternalServerError(exceptions.InternalErrorMessage))
		return
	}
}
