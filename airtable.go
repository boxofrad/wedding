package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	airtableBaseUrl = "https://api.airtable.com/v0/apppZ6wOjC6cbodTV/"
)

var (
	airtableKey = os.Getenv("AIRTABLE_KEY")
)

type Guest struct {
	Id                            string
	FirstName                     string
	PartOfDay                     string
	AttendingService              bool
	AttendingReception            bool
	AttendingEvening              bool
	MealType                      string
	AdditionalDietaryRequirements string
}

type Invitation struct {
	Id         string
	Addressees string
	GuestIds   []string
	Guests     []*Guest
}

func invitationWithCode(code string) (*Invitation, error) {
	q := make(url.Values)
	q.Add("filterByFormula", fmt.Sprintf(`{Code}="%s"`, code))

	req, err := http.NewRequest("GET", airtableBaseUrl+"Invitations?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airtable: unexpected response code %d", resp.StatusCode)
	}

	response := struct {
		Records []struct {
			Id     string `json:"id"`
			Fields struct {
				Addressees string   `json:"addressees"`
				Guests     []string `json:"guests"`
			} `json:"fields"`
		} `json:"records"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Records) == 0 {
		return nil, fmt.Errorf(`airtable: no invitation found with code "%s"`, code)
	}

	fields := response.Records[0].Fields
	guests := make([]*Guest, len(fields.Guests))
	for i, guestId := range fields.Guests {
		guest, err := getGuest(guestId)
		if err != nil {
			return nil, err
		}
		guests[i] = guest
	}

	return &Invitation{
		Id:         response.Records[0].Id,
		Addressees: fields.Addressees,
		GuestIds:   fields.Guests,
		Guests:     guests,
	}, nil
}

func getGuest(id string) (*Guest, error) {
	req, err := http.NewRequest("GET", airtableBaseUrl+"Guests/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airtable: unexpected response code %d", resp.StatusCode)
	}

	response := struct {
		Id     string `json:"id"`
		Fields struct {
			FirstName                     string `json:"First Name"`
			PartOfDay                     string `json:"Part of Day"`
			AttendingService              bool   `json:"Attending Service?"`
			AttendingReception            bool   `json:"Attending Reception?"`
			AttendingEvening              bool   `json:"Attending Evening?"`
			MealType                      string `json:"Meal Type"`
			AdditionalDietaryRequirements string `json:"Additional Dietary Requirements"`
		} `json:"fields"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &Guest{
		Id:                            response.Id,
		FirstName:                     response.Fields.FirstName,
		PartOfDay:                     response.Fields.PartOfDay,
		AttendingService:              response.Fields.AttendingService,
		AttendingReception:            response.Fields.AttendingReception,
		AttendingEvening:              response.Fields.AttendingEvening,
		MealType:                      response.Fields.MealType,
		AdditionalDietaryRequirements: response.Fields.AdditionalDietaryRequirements,
	}, nil
}

func invitationWithId(id string) (*Invitation, error) {
	req, err := http.NewRequest("GET", airtableBaseUrl+"Invitations/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airtable: unexpected response code %d", resp.StatusCode)
	}

	response := struct {
		Id     string `json:"id"`
		Fields struct {
			Addressees string   `json:"addressees"`
			Guests     []string `json:"guests"`
		} `json:"fields"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &Invitation{
		Id:         response.Id,
		Addressees: response.Fields.Addressees,
		GuestIds:   response.Fields.Guests,
	}, nil
}

type UpdateGuestFields struct {
	RSVPReceived                  bool   `json:"RSVP Received?"`
	AttendingService              bool   `json:"Attending Service?"`
	AttendingReception            bool   `json:"Attending Reception?"`
	AttendingEvening              bool   `json:"Attending Evening?"`
	MealType                      string `json:"Meal Type"`
	AdditionalDietaryRequirements string `json:"Additional Dietary Requirements"`
}

func updateGuest(id string, fields UpdateGuestFields) error {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(struct {
		Fields UpdateGuestFields `json:"fields"`
	}{fields}); err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", airtableBaseUrl+"Guests/"+id, b)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("airtable: unexpected response code %d", resp.StatusCode)
	}
	return nil
}
