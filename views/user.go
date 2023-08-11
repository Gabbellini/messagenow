package views

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/usecases"
	"net/http"
	"strconv"
)

type userHttpModule struct {
	createUserUseCase  usecases.CreateUserUseCase
	getUserByIDUseCase usecases.GetUserByIDUseCase
}

func NewUserHTTPModule(createUserUseCase usecases.CreateUserUseCase, getUserByIDUseCase usecases.GetUserByIDUseCase) ModuleHTTP {
	return &userHttpModule{
		getUserByIDUseCase: getUserByIDUseCase,
		createUserUseCase:  createUserUseCase,
	}
}

func (m userHttpModule) Setup(router *mux.Router) {
	router.HandleFunc("/me", m.me).Methods(http.MethodGet)
	router.HandleFunc("/users", m.createUser).Methods(http.MethodPost)
}

func (m userHttpModule) me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := ctx.Value("user").(entities.User)

	user, err := m.getUserByIDUseCase.Execute(ctx, userCtx.ID)
	if err != nil {
		log.Println("[me] Error Execute", err)
		exceptions.HandleError(w, err)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Println("[me] Error Marshal", err)
		exceptions.HandleError(w, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[me] Error Write", err)
		exceptions.HandleError(w, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage))
		return
	}
}

func (m userHttpModule) createUser(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[createUser] Error ReadAll", err)
		exceptions.HandleError(w, exceptions.NewBadRequestError("Corpo da requisição não é valido"))
		return
	}

	var user entities.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Println("[createUser] Error Unmarshal", err)
		exceptions.HandleError(w, exceptions.NewBadRequestError("Corpo da requisição não é valido"))
		return
	}

	userID, err := m.createUserUseCase.Execute(r.Context(), user)
	if err != nil {
		log.Println("[me] Error Execute", err)
		exceptions.HandleError(w, err)
		return
	}

	_, err = w.Write([]byte(strconv.FormatInt(userID, 10)))
	if err != nil {
		log.Println("[me] Error Write", err)
		exceptions.HandleError(w, exceptions.NewUnexpectedError(exceptions.UnexpectedErrorMessage))
		return
	}
}
