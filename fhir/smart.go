package fhir

import (
	"fmt"
	"net/http"
	"net/url"
)

type Launch struct {
	ISS    string `json:"iss"`
	Launch string `json:"launch"`
}

func (l Launch) Params() string {
	params := url.Values{}
	params.Set("iss", url.QueryEscape(l.ISS))
	params.Set("launch", url.QueryEscape(l.Launch))
	return params.Encode()
}

func ParseLaunch(r *http.Request) (Launch, error) {
	iss, err := url.QueryUnescape(r.FormValue("iss"))
	if err != nil {
		return Launch{}, fmt.Errorf("invalid iss: %v", err)
	}

	l, err := url.QueryUnescape(r.FormValue("launch"))
	if err != nil {
		return Launch{}, fmt.Errorf("invalid launch: %v", err)
	}
	return Launch{Launch: l, ISS: iss}, nil
}
