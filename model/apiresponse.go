package model

type APIResponse struct {

	// code
	Code int `json:"code"`

	// message
	Message string `json:"message,omitempty"`
}
