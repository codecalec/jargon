package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"jargon/pkg/api"
	"jargon/pkg/db"
)

type indexHandler struct {
	database *db.Database
}

func (h indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	jargons, err := h.database.GetAllJargons()
	if err != nil {
		log.Fatal(err)
	}

	t := template.Must(template.ParseFiles(
		"static/html/index.html",
		"static/html/home.html",
		"static/html/jargon_list.html"))

	if err := t.Execute(w, jargons); err != nil {
		log.Fatal(err)
	}
}

type jargonPageHandler struct {
	database *db.Database
}

func (h jargonPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	label := strings.TrimPrefix(r.URL.Path, "/page/")

	i, err := strconv.ParseUint(label, 10, 32)
	if err != nil {
	}

	var j *api.Jargon
	j, err = h.database.GetJargon(uint32(i))
	if err != nil {
		http.NotFoundHandler().ServeHTTP(w, r)
	}

	t := template.Must(template.ParseFiles(
		"static/html/index.html",
		"static/html/page.html"))

	if err := t.Execute(w, *j); err != nil {
		log.Fatal(err)
	}
}

func muxWithLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got a %s request from %v for: %v", r.Method, r.RemoteAddr, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func StartServer(database *db.Database, port uint) {
	mux := http.NewServeMux()
	mux.Handle("/", indexHandler{database})
	mux.Handle("/page/", jargonPageHandler{database})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	handler := muxWithLogging(mux)

	log.Printf("Starting server on port %v\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), handler); err != nil {
		log.Fatal(err)
	}
}