package googleApi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"google.golang.org/api/googleapi"
)

func (g *googleApi) SearchByUsersFromApi(email string) ([]data.User, error) {
	response, err := g.client.Users.List().Query("email:" + email).Do()

	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok {
			if apiErr.Code == http.StatusTooManyRequests {
				durationInSeconds, err := strconv.ParseInt(apiErr.Header.Get("Retry-After"), 10, 64)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse `retry-after` header")
				}
				timeoutDuration := time.Second * time.Duration(durationInSeconds)
				time.Sleep(timeoutDuration)
				return g.SearchByUsersFromApi(email)
			}
		}

		g.log.WithError(err).Errorf("failed to search users by email `%s`", email)
		return nil, errors.Wrap(err, "failed to search users by email")
	}

	result := make([]data.User, 0)

	for _, user := range response.Users {
		result = append(result, data.User{
			Email:  user.PrimaryEmail,
			Name:   user.Name.FullName,
			MailId: user.Id,
		})
	}

	return result, nil
}
