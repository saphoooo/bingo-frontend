package main

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}
func main() {
	r := mux.NewRouter()
	r.Handle("/", &templateHandler{filename: "bingo.html"})
	log.Print("Start listening on :8000...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
}
