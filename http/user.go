package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"messagenow/domain/entities"
	"messagenow/exceptions"
	"messagenow/usecases"
	"net/http"
)

type userHttpModule struct {
	getUserUseCase usecases.GetUserUseCase
}

func NewUserHTTPModule(getUserUseCase usecases.GetUserUseCase) ModuleHTTP {
	return &userHttpModule{
		getUserUseCase: getUserUseCase,
	}
}

func (m userHttpModule) Setup(router *mux.Router) {
	router.HandleFunc("/me", m.me).Methods(http.MethodGet)
}

func (m userHttpModule) me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCtx := ctx.Value("user").(entities.User)

	user, err := m.getUserUseCase.Execute(ctx, userCtx.ID)
	if err != nil {
		log.Println("[me] Error Execute", err)
		exceptions.HandleError(w, err)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Println("[me] Error Marshal", err)
		exceptions.HandleError(w, exceptions.NewInternalServerError(exceptions.InternalErrorMessage))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[me] Error Write", err)
		exceptions.HandleError(w, exceptions.NewInternalServerError(exceptions.InternalErrorMessage))
		return
	}
}
