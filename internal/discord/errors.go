package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OAuthError struct {
	Status      int
	Code        string
	Description string
}

func (e *OAuthError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf(
			"discord oauth: %s: %s (http %d)",
			e.Code,
			e.Description,
			e.Status,
		)
	}

	return fmt.Sprintf("discord oauth: %s (http %d)", e.Code, e.Status)
}

func parseOAuthError(status int, body []byte) error {
	var payload struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	_ = json.Unmarshal(body, &payload)

	return &OAuthError{
		Status:      status,
		Code:        payload.Error,
		Description: payload.Description,
	}
}

type APIError struct {
	Status  int
	Code    int
	Message string
	Body    string
}

func (e *APIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf(
			"discord api: %s (http %d, code %d)",
			e.Message,
			e.Status,
			e.Code,
		)
	}

	return fmt.Sprintf("discord api: http %d: %s (body: %.200q)", e.Status, e.Message, e.Body)
}

func parseAPIError(status int, body []byte) error {
	var payload struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	_ = json.Unmarshal(body, &payload)

	msg := payload.Message
	if msg == "" {
		msg = http.StatusText(status)
	}

	return &APIError{
		Status:  status,
		Code:    payload.Code,
		Message: msg,
		Body:    string(body),
	}
}
