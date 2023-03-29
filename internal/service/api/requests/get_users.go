package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetUsersRequest struct {
	Email *string `filter:"username"`
}

func NewGetUsersRequest(r *http.Request) (GetUsersRequest, error) {
	var request GetUsersRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, err
}
