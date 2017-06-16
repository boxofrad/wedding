package main

import (
	"encoding/xml"
	"errors"
	"net/http"

	"golang.org/x/net/html/charset"
	xmlpath "gopkg.in/xmlpath.v2"
)

// Annoyingly Kuoni's gift list site doesn't provide a permalink for our list
// e.g. some-url.com/gift-list?id=1234 because the form is only available via
// a POST request. Meaning we can't redirect to it.
//
// Kind people buying us gifts would have to hunt around for a practically hidden
// iframed form, figure out which is the correct box and type in our list number.
//
// There's luckily no CSRF protection on this form endpoint though so we can
// just render an equivilent form our page and automatically submit it with
// JavaScript. \o/
//
// The only caveat is that the form embeds a unique session id which we can't fake.
// Luckily we can just GET the form, which for some reason returns XML? huh? and pull
// the session id out with XPATH.
func getGiftListSessionId() (string, error) {
	req, err := http.NewRequest("GET", "http://booking.kuoni.co.uk/ob/x1root?TRNTPD=GF03", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	root, err := xmlpath.ParseDecoder(decoder)
	if err != nil {
		return "", err
	}

	xpath, err := xmlpath.Compile(`//page/reqrsp/@sessionId`)
	if err != nil {
		return "", err
	}

	if value, found := xpath.String(root); found {
		return value, nil
	} else {
		return "", errors.New("gifts: unable to scrape session id")
	}
}
