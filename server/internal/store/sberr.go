package store

import (
	"encoding/json"
	"regexp"
	"strconv"
)

type SupabaseAuthError struct {
	StatusCode       int    `json:"status_code"`
	ErrorName        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *SupabaseAuthError) Error() string {
	return e.ErrorDescription
}

func NewSupabaseAuthError(sbErr error) (*SupabaseAuthError, bool) {
	re := regexp.MustCompile(`response status code (\d+): (.+)`)
	matches := re.FindStringSubmatch(sbErr.Error())
	if len(matches) != 3 {
		return nil, false
	}

	statusCode, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, false
	}

	var parsedBody struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	err = json.Unmarshal([]byte(matches[2]), &parsedBody)
	if err != nil {
		return nil, false
	}

	authError := &SupabaseAuthError{
		StatusCode:       statusCode,
		ErrorName:        parsedBody.Error,
		ErrorDescription: parsedBody.ErrorDescription,
	}
	return authError, true
}
