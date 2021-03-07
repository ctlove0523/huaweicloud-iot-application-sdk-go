package iot

import "encoding/json"

type ApplicationError struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}


func (e *ApplicationError) Error() string {
	if e == nil {
		return ""
	}

	jsonString, err := json.Marshal(e)
	if err != nil {
		return "marshal error"
	}

	return string(jsonString)
}
