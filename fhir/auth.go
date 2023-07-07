package fhir

import "net/url"

type AuthCode struct {
	ResponseType        string // response_type
	ClientID            string // client_id
	RedirectURI         string // redirect_uri
	Scope               string // scope
	State               string // state
	Audience            string // aud
	CodeChallenge       string //code_challenge
	CodeChallengeMethod string // code_challenge_method
}

func (a AuthCode) Params() string {

	params := url.Values{}
	params.Set("response_type", a.ResponseType)
	params.Set("client_id", a.ClientID)
	params.Set("redirect_uri", a.RedirectURI)
	params.Set("state", a.State)
	params.Set("aud", a.Audience)
	params.Set("code_challenge", a.CodeChallenge)
	params.Set("code_challenge_method", a.CodeChallengeMethod)
	params.Add("scope", a.Scope)

	return params.Encode()
}

func NewAuthCode(params url.Values) AuthCode {
	authParams := AuthCode{
		ResponseType:        params.Get("response_type"),
		ClientID:            params.Get("client_id"),
		RedirectURI:         params.Get("redirect_uri"),
		Scope:               params.Get("scope"),
		State:               params.Get("state"),
		Audience:            params.Get("aud"),
		CodeChallenge:       params.Get("code_challenge"),
		CodeChallengeMethod: params.Get("code_challenge_method"),
	}

	return authParams
}
