package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	staticPath = "/static"
)

var (
	listen = ":" + os.Getenv("PORT")
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", serveRoot)

	// TODO: Cache assets
	r.PathPrefix(staticPath).
		Handler(http.StripPrefix(staticPath, http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(listen, r))
}

func serveRoot(w http.ResponseWriter, _ *http.Request) {
	assetFn, err := assetPathHelper()
	if err != nil {
		handleErr(w, err)
		return
	}

	tmpl, err := template.New("").
		Funcs(map[string]interface{}{"asset_path": assetFn}).
		ParseFiles("templates/index.html.tmpl")
	if err != nil {
		handleErr(w, err)
		return
	}

	if err = tmpl.ExecuteTemplate(w, "index.html.tmpl", struct{}{}); err != nil {
		handleErr(w, err)
	}
}

func handleErr(w http.ResponseWriter, err error) {
	// TODO: Do something with the error
}
