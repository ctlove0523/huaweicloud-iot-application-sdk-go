package iot

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

func convertResponseToApplicationError(response *resty.Response) error {
	are := &ApplicationResponseError{}

	err := json.Unmarshal(response.Body(), are)
	if err != nil {
		return err
	}

	ae := &ApplicationError{
		ErrorMsg:  are.ErrorMsg,
		ErrorCode: are.ErrorCode,
	}

	return ae
}

type ApplicationResponseError struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}
