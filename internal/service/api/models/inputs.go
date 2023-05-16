package models

import (
	"github.com/acs-dl/mail-module-svc/resources"
)

func NewInputsModel() resources.Inputs {
	result := resources.Inputs{
		Key: resources.Key{
			ID:   "0",
			Type: resources.INPUTS,
		},
		Attributes: resources.InputsAttributes{
			Email: "string",
			Link:  "string",
		},
	}

	return result
}

func NewInputsResponse() resources.InputsResponse {
	return resources.InputsResponse{
		Data: NewInputsModel(),
	}
}
