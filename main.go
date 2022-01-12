package main

import (
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/rs/zerolog/log"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpan("bingo.try.request", tracer.ResourceName("/"))
	defer span.Finish()
	span.SetTag("http.url", r.URL.Path)

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}
func main() {
	tracer.Start(
		tracer.WithEnv("prod"),
		tracer.WithService("bingo-frontend"),
		tracer.WithServiceVersion("v1.0"),
	)
	defer tracer.Stop()
	r := muxtrace.NewRouter()
	r.Handle("/", loggingHandler(&templateHandler{filename: "bingo.html"}))
	r.PathPrefix("/src/").Handler(loggingHandler(http.StripPrefix("/src/", http.FileServer(http.Dir("./src")))))
	log.Print("Start listening on :8000...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("hostname", r.Host).
			Str("method", r.Method).
			Str("proto", r.Proto).
			Str("remote_ip", r.RemoteAddr).
			Str("path", r.RequestURI).
			Str("user-agent", r.UserAgent()).
			Int("status", http.StatusOK).
			Msg("")
		next.ServeHTTP(w, r)
	})
}
