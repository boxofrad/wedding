package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
	httpClient  = http.DefaultClient

	ErrNotFound = errors.New("airtable: not found")
)

type Guest struct {
	Id                            string
	FirstName                     string
	PartOfDay                     string
	RSVPReceived                  bool
	AttendingService              bool
	AttendingReception            bool
	AttendingEvening              bool
	DietaryPreferences            string
	AdditionalDietaryRequirements string
	BingoFact                     string
}

type Invitation struct {
	Id         string
	Addressees string
	GuestIds   []string
	Guests     []*Guest
}

func getInvitationId(code string) (string, error) {
	q := make(url.Values)
	q.Add("filterByFormula", fmt.Sprintf(`{Code}="%s"`, code))

	req, err := http.NewRequest("GET", airtableBaseUrl+"Invitations?"+q.Encode(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("airtable: unexpected response code %d", resp.StatusCode)
	}

	response := struct {
		Records []struct {
			Id string `json:"id"`
		} `json:"records"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Records) == 0 {
		return "", ErrNotFound
	}
	return response.Records[0].Id, nil
}

func getInvitation(id string) (*Invitation, error) {
	req, err := http.NewRequest("GET", airtableBaseUrl+"Invitations/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))

	resp, err := httpClient.Do(req)
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

	guests := make([]*Guest, len(response.Fields.Guests))
	for i, guestId := range response.Fields.Guests {
		guest, err := getGuest(guestId)
		if err != nil {
			return nil, err
		}
		guests[i] = guest
	}

	return &Invitation{
		Id:         response.Id,
		Addressees: response.Fields.Addressees,
		GuestIds:   response.Fields.Guests,
		Guests:     guests,
	}, nil
}

func getGuest(id string) (*Guest, error) {
	req, err := http.NewRequest("GET", airtableBaseUrl+"Guests/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))

	resp, err := httpClient.Do(req)
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
			RSVPReceived                  bool   `json:"RSVP Received?"`
			AttendingService              bool   `json:"Attending Service?"`
			AttendingReception            bool   `json:"Attending Reception?"`
			AttendingEvening              bool   `json:"Attending Evening?"`
			DietaryPreferences            string `json:"Dietary Preferences"`
			AdditionalDietaryRequirements string `json:"Additional Dietary Requirements"`
			BingoFact                     string `json:"Bingo Fact"`
		} `json:"fields"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &Guest{
		Id:                            response.Id,
		FirstName:                     response.Fields.FirstName,
		PartOfDay:                     response.Fields.PartOfDay,
		RSVPReceived:                  response.Fields.RSVPReceived,
		AttendingService:              response.Fields.AttendingService,
		AttendingReception:            response.Fields.AttendingReception,
		AttendingEvening:              response.Fields.AttendingEvening,
		DietaryPreferences:            response.Fields.DietaryPreferences,
		AdditionalDietaryRequirements: response.Fields.AdditionalDietaryRequirements,
		BingoFact:                     response.Fields.BingoFact,
	}, nil
}

type UpdateGuestParams struct {
	RSVPReceived                  bool   `json:"RSVP Received?"`
	AttendingService              bool   `json:"Attending Service?"`
	AttendingReception            bool   `json:"Attending Reception?"`
	AttendingEvening              bool   `json:"Attending Evening?"`
	DietaryPreferences            string `json:"Dietary Preferences,omitempty"`
	AdditionalDietaryRequirements string `json:"Additional Dietary Requirements"`
	BingoFact                     string `json:"Bingo Fact"`
}

func updateGuest(id string, fields UpdateGuestParams) error {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(struct {
		Fields UpdateGuestParams `json:"fields"`
	}{fields}); err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", airtableBaseUrl+"Guests/"+id, b)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("airtable: unexpected response code %d", resp.StatusCode)
	}
	return nil
}

func getInvitationsWithoutCodes() ([]string, error) {
	q := make(url.Values)
	q.Add("filterByFormula", `{Code} = ""`)

	req, err := http.NewRequest("GET", airtableBaseUrl+"Invitations?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airtable: unexpected response code %d", resp.StatusCode)
	}

	response := struct {
		Records []struct {
			Id string `json:"id"`
		} `json:"records"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	ids := make([]string, len(response.Records))
	for i, record := range response.Records {
		ids[i] = record.Id
	}
	return ids, nil
}

func updateInvitationCode(id, code string) error {
	body := struct {
		Fields struct {
			Code string `json:"Code"`
		} `json:"fields"`
	}{}
	body.Fields.Code = code

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", airtableBaseUrl+"Invitations/"+id, b)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", airtableKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("airtable: unexpected response code %d", resp.StatusCode)
	}
	return nil
}
