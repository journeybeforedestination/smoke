package fhir

type Conformance struct {
	AuthEndpoint                  string   `json:"authorization_endpoint"`
	TokenEndpoint                 string   `json:"token_endpoint"`
	GrantTypesSupported           []string `json:"grant_types_supported"`
	ScopesSupported               []string `json:"scopes_supported"`
	ResponseTypesSupported        []string `json:"response_types_supported"`
	CodeChallengeMethodsSupported []string `json:"code_challenge_methods_supported"`
	Capabilities                  []string `json:"capabilities"`
}
