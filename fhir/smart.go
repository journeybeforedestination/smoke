package fhir

import "net/url"

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
