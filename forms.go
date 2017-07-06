package main

import (
	"fmt"
	"net/http"
)

type InvitationForm struct {
	Invitation *Invitation
	GuestForms []*GuestForm
}

func NewInvitationForm(invitation *Invitation) *InvitationForm {
	form := &InvitationForm{
		Invitation: invitation,
		GuestForms: make([]*GuestForm, len(invitation.Guests)),
	}
	for i, guest := range invitation.Guests {
		form.GuestForms[i] = &GuestForm{Guest: guest}
	}
	return form
}

func (i *InvitationForm) Parse(r *http.Request) {
	for _, guestForm := range i.GuestForms {
		guestForm.Parse(r)
	}
}

func (i InvitationForm) Validate() (bool, []string) {
	errors := make([]string, 0)

	for _, guestForm := range i.GuestForms {
		errors = append(errors, guestForm.Validate(len(i.GuestForms) == 1)...)
	}

	return len(errors) == 0, errors
}

type GuestForm struct {
	Guest *Guest

	RawAttendingService              string
	RawAttendingReception            string
	RawAttendingEvening              string
	RawDietaryPreferences            string
	RawAdditionalDietaryRequirements string
	RawBingoFact                     string
}

func (g *GuestForm) Parse(r *http.Request) {
	g.RawAttendingService = r.FormValue(g.FieldName("attending_service"))
	g.RawAttendingReception = r.FormValue(g.FieldName("attending_reception"))
	g.RawAttendingEvening = r.FormValue(g.FieldName("attending_evening"))
	g.RawDietaryPreferences = r.FormValue(g.FieldName("dietary_preferences"))
	g.RawAdditionalDietaryRequirements = r.FormValue(g.FieldName("additional_dietary_requirements"))
	g.RawBingoFact = r.FormValue(g.FieldName("bingo_fact"))
}

func (g GuestForm) FieldName(name string) string {
	return "guest[" + g.Guest.Id + "][" + name + "]"
}

func (g GuestForm) AttendingService() bool {
	if g.RawAttendingService == "" {
		return g.Guest.AttendingService
	}
	return g.RawAttendingService == "1"
}

func (g GuestForm) NotAttendingService() bool {
	if g.RawAttendingService == "" && g.Guest.RSVPReceived {
		return !g.Guest.AttendingService
	}
	return g.RawAttendingService == "0"
}

func (g GuestForm) AttendingReception() bool {
	if g.RawAttendingReception == "" {
		return g.Guest.AttendingReception
	}
	return g.RawAttendingReception == "1"
}

func (g GuestForm) NotAttendingReception() bool {
	if g.RawAttendingReception == "" && g.Guest.RSVPReceived {
		return !g.Guest.AttendingReception
	}
	return g.RawAttendingReception == "0"
}

func (g GuestForm) AttendingEvening() bool {
	if g.RawAttendingEvening == "" {
		return g.Guest.AttendingEvening
	}
	return g.RawAttendingEvening == "1"
}

func (g GuestForm) NotAttendingEvening() bool {
	if g.RawAttendingEvening == "" && g.Guest.RSVPReceived {
		return !g.Guest.AttendingEvening
	}
	return g.RawAttendingEvening == "0"
}

func (g GuestForm) DietaryPreferences() string {
	if g.RawDietaryPreferences == "" {
		return g.Guest.DietaryPreferences
	}

	return g.RawDietaryPreferences
}

func (g GuestForm) AdditionalDietaryRequirements() string {
	if g.RawAdditionalDietaryRequirements == "" {
		return g.Guest.AdditionalDietaryRequirements
	}

	return g.RawAdditionalDietaryRequirements
}

func (g GuestForm) BingoFact() string {
	if g.RawBingoFact == "" {
		return g.Guest.BingoFact
	}

	return g.RawBingoFact
}

func (g GuestForm) Validate(single bool) []string {
	errors := make([]string, 0)

	var salutation string
	if single {
		salutation = "you"
	} else {
		salutation = g.Guest.FirstName
	}

	if g.RawAttendingService == "" {
		errors = append(errors, fmt.Sprintf("Please indicate if %s will be attending the service.", salutation))
	}

	if g.Guest.PartOfDay == "Service & Reception" {
		if g.RawAttendingReception == "" {
			errors = append(errors, fmt.Sprintf("Please indicate if %s will be attending the reception.", salutation))
		}

		if g.RawAttendingReception == "1" && g.Guest.Age != "Baby" {
			if g.RawDietaryPreferences == "" {
				errors = append(errors, fmt.Sprintf("Please indicate %s's dietary preferences.", salutation))
			}

			if g.RawBingoFact == "" {
				var who string
				if single {
					who = "yourself"
				} else {
					who = g.Guest.FirstName
				}

				errors = append(errors, fmt.Sprintf("Please tell us something random about %s.", who))
			}
		}
	}

	if g.Guest.PartOfDay == "Service & Evening" {
		if g.RawAttendingEvening == "" {
			errors = append(errors, fmt.Sprintf("Please indicate if %s will be attending the evening reception.", salutation))
		}
	}

	return errors
}

func (g GuestForm) UpdateParams() UpdateGuestParams {
	var attendingEvening bool
	if g.Guest.PartOfDay == "Service & Reception" {
		attendingEvening = g.AttendingReception()
	} else {
		attendingEvening = g.AttendingEvening()
	}

	return UpdateGuestParams{
		RSVPReceived:                  true,
		AttendingService:              g.AttendingService(),
		AttendingReception:            g.AttendingReception(),
		AttendingEvening:              attendingEvening,
		DietaryPreferences:            g.DietaryPreferences(),
		AdditionalDietaryRequirements: g.AdditionalDietaryRequirements(),
		BingoFact:                     g.BingoFact(),
	}
}
