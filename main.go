package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	staticPath  = "/static"
	giftListURL = "https://booking.kuoni.co.uk/ob/X1ROOT?TRNTPD=GF04&TRNNRD=1&GFTGID=34182"
)

var (
	listen = ":" + os.Getenv("PORT")
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", serveTemplate("index"))
	r.HandleFunc("/rsvp", serveRSVP).Methods("GET")
	r.HandleFunc("/rsvp/{id}", serveRSVPResponse).Methods("POST")
	r.HandleFunc("/rsvp/{id}/success", serveTemplate("rsvp_success")).Methods("GET")
	r.Handle("/gifts", http.RedirectHandler(giftListURL, http.StatusFound)).Methods("GET")

	// TODO: Cache assets
	r.PathPrefix(staticPath).
		Handler(http.StripPrefix(staticPath, http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(listen, r))
}

// TODO: Handle the error
func serveRSVP(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		renderTemplate(w, "rsvp", nil)
		return
	}

	invitation, err := invitationWithCode(code)
	if err != nil {
		handleErr(w, err)
		return
	}

	renderTemplate(w, "invitation", invitation)
}

func serveRSVPResponse(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	i, err := invitationWithId(id)
	if err != nil {
		handleErr(w, err)
		return
	}

	for _, guestId := range i.GuestIds {
		formVal := func(k string) string {
			return r.FormValue(fmt.Sprintf("guest[%s][%s]", guestId, k))
		}

		err = updateGuest(guestId, UpdateGuestFields{
			RSVPReceived:                  true,
			AttendingService:              formVal("attending_service") == "1",
			AttendingReception:            formVal("attending_reception") == "1",
			AttendingEvening:              formVal("attending_evening") == "1",
			MealType:                      formVal("meal_type"),
			AdditionalDietaryRequirements: formVal("additional_dietary_requirements"),
		})

		if err != nil {
			handleErr(w, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/rsvp/%s/success", id), http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	assetFn, err := assetPathHelper()
	if err != nil {
		handleErr(w, err)
		return
	}

	tmpl, err := template.New("").
		Funcs(map[string]interface{}{"asset_path": assetFn}).
		ParseFiles("templates/"+name+".html.tmpl", "templates/layout.html.tmpl")
	if err != nil {
		handleErr(w, err)
		return
	}

	if err = tmpl.ExecuteTemplate(w, "base", data); err != nil {
		handleErr(w, err)
	}
}

func serveTemplate(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		renderTemplate(w, name, nil)
	}
}

func handleErr(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, "ERROR: %s", err)
}
