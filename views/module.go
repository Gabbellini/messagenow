package views

import "github.com/gorilla/mux"

type ModuleHTTP interface {
	Setup(router *mux.Router)
}
