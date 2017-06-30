package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// TODO: Better name for this
type errorMessage struct {
	ErrorMessage string
}

const (
	staticPath  = "/static"
	giftListURL = "https://booking.kuoni.co.uk/ob/X1ROOT?TRNTPD=GF04&TRNNRD=1&GFTGID=34182"
)

var (
	listen = ":" + os.Getenv("PORT")
)

func runServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", serveTemplate("index"))
	r.HandleFunc("/rsvp", serveTemplate("rsvp_form")).Methods("GET")
	r.HandleFunc("/rsvp", serveRSVP).Methods("POST")
	r.HandleFunc("/invitation/{id}", serveInvitationForm).Methods("GET")
	r.HandleFunc("/invitation/{id}", serveInvitation).Methods("POST")
	r.HandleFunc("/invitation/{id}/success", serveTemplate("rsvp_success")).Methods("GET")
	r.Handle("/gifts", http.RedirectHandler(giftListURL, http.StatusFound)).Methods("GET")

	// TODO: Cache assets
	r.PathPrefix(staticPath).
		Handler(http.StripPrefix(staticPath, http.FileServer(http.Dir("./static"))))

	log.Fatal(http.ListenAndServe(listen, r))
}

func serveRSVP(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	if code == "" {
		renderTemplate(w, "rsvp_form", errorMessage{"Please enter the code written on your RSVP card."})
		return
	}

	id, err := getInvitationId(code)
	if err == ErrNotFound {
		renderTemplate(w, "rsvp_form", errorMessage{"Uh-oh, we don't recognise that code. Please make sure you typed it correctly."})
		return
	}
	if err != nil {
		handleErr(w, err)
		return
	}

	http.Redirect(w, r, "/invitation/"+id, http.StatusFound)
}

func serveInvitationForm(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	invitation, err := getInvitation(id)
	if err == ErrNotFound {
		http.Redirect(w, r, "/rsvp", http.StatusFound)
		return
	}
	if err != nil {
		handleErr(w, err)
		return
	}

	renderTemplate(w, "invitation", struct {
		InvitationForm *InvitationForm
		ErrorMessage   string
	}{NewInvitationForm(invitation), ""})
}

func serveInvitation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	invitation, err := getInvitation(id)
	if err != nil {
		handleErr(w, err)
		return
	}

	form := NewInvitationForm(invitation)
	form.Parse(r)

	if valid, errors := form.Validate(); !valid {
		renderTemplate(w, "invitation", struct {
			InvitationForm *InvitationForm
			ErrorMessage   string
		}{form, errors[0]})
		return
	}

	for _, guestForm := range form.GuestForms {
		if err := updateGuest(guestForm.Guest.Id, guestForm.UpdateParams()); err != nil {
			handleErr(w, err)
			return
		}
	}
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
