package googleApi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"google.golang.org/api/googleapi"
)

func (g *googleApi) GetDomainUserFromApi(email string) (*data.User, error) {
	response, err := g.client.Users.Get(email).Do()
	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok {
			if apiErr.Code == http.StatusTooManyRequests {
				durationInSeconds, err := strconv.ParseInt(apiErr.Header.Get("Retry-After"), 10, 64)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse `retry-after` header")
				}
				timeoutDuration := time.Second * time.Duration(durationInSeconds)
				time.Sleep(timeoutDuration)
				return g.GetDomainUserFromApi(email)
			}
			if apiErr.Code == http.StatusNotFound {
				return nil, nil
			}
		}

		g.log.WithError(err).Errorf("failed to get user by email `%s`", email)
		return nil, errors.Wrap(err, "failed to get users by email")
	}

	return &data.User{
		Email:  response.PrimaryEmail,
		Name:   response.Name.FullName,
		MailId: response.Id,
	}, nil
}
