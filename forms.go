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
		errors = append(errors, guestForm.Validate()...)
	}

	return len(errors) == 0, errors
}

type GuestForm struct {
	Guest *Guest

	RawAttendingService              string
	RawAttendingReception            string
	RawAttendingEvening              string
	RawMealType                      string
	RawAdditionalDietaryRequirements string
	RawBingoFact                     string
}

func (g *GuestForm) Parse(r *http.Request) {
	g.RawAttendingService = r.FormValue(g.FieldName("attending_service"))
	g.RawAttendingReception = r.FormValue(g.FieldName("attending_reception"))
	g.RawAttendingEvening = r.FormValue(g.FieldName("attending_evening"))
	g.RawMealType = r.FormValue(g.FieldName("meal_type"))
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

func (g GuestForm) MealType() string {
	if g.RawMealType == "" {
		return g.Guest.MealType
	}

	return g.RawMealType
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

func (g GuestForm) Validate() []string {
	errors := make([]string, 0)

	if g.RawAttendingService == "" {
		errors = append(errors, fmt.Sprintf("Please indicate if %s will be attending the service.", g.Guest.FirstName))
	}

	if g.Guest.PartOfDay == "Service & Reception" {
		if g.RawAttendingReception == "" {
			errors = append(errors, fmt.Sprintf("Please indicate if %s will be attending the reception.", g.Guest.FirstName))
		}

		if g.RawAttendingReception == "1" {
			if g.RawMealType == "" {
				errors = append(errors, fmt.Sprintf("Please indicate what kind of meal %s would like.", g.Guest.FirstName))
			}

			if g.RawBingoFact == "" {
				errors = append(errors, fmt.Sprintf("Please tell us something random about %s.", g.Guest.FirstName))
			}
		}
	}

	if g.Guest.PartOfDay == "Service & Reception" || g.Guest.PartOfDay == "Service & Evening" {
		if g.RawAttendingEvening == "" {
			errors = append(errors, fmt.Sprintf("Please indicate if %s will be attending the evening reception.", g.Guest.FirstName))
		}
	}

	return errors
}

func (g GuestForm) UpdateParams() UpdateGuestParams {
	return UpdateGuestParams{
		RSVPReceived:                  true,
		AttendingService:              g.AttendingService(),
		AttendingReception:            g.AttendingReception(),
		AttendingEvening:              g.AttendingEvening(),
		MealType:                      g.MealType(),
		AdditionalDietaryRequirements: g.AdditionalDietaryRequirements(),
		BingoFact:                     g.BingoFact(),
	}
}
