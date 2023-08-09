package views

import (
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func NewPWAWebModule(
	path string,
	rootDir string,
) ModuleHTTP {
	return studentModule{
		path:    path,
		rootDir: rootDir,
	}
}

type studentModule struct {
	path    string
	rootDir string
}

func (s studentModule) Path() string {
	return s.path
}

func (s studentModule) Setup(router *mux.Router) {
	subRouter := router.PathPrefix(s.Path()).Subrouter()
	subRouter.PathPrefix("").Handler(s)
}

func (s studentModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	indexName := filepath.Join(s.rootDir, "index.html")
	fileName := filepath.Join(s.rootDir, strings.TrimPrefix(r.URL.Path, s.Path()))

	var f *os.File
	var err error

	if r.URL.Path == s.Path() || r.URL.Path == s.Path() {
		f, err = os.Open(indexName)
		if err != nil {
			http.Error(w, r.RequestURI, http.StatusNotFound)
			return
		}
	} else {
		f, err = os.Open(fileName)
		if err != nil {
			log.Println(err)
			if !errors.Is(err, os.ErrNotExist) {
				http.Error(w, r.RequestURI, http.StatusNotFound)
				return
			} else {
				log.Println(err)
				f, err = os.Open(indexName)
				if err != nil {
					http.Error(w, r.RequestURI, http.StatusNotFound)
					return
				}
			}
		}
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, r.RequestURI, http.StatusNotFound)
		return
	}

	http.ServeContent(w, r, fileName, fi.ModTime(), f)
}
