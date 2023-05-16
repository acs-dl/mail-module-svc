package googleApi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"google.golang.org/api/googleapi"
)

func (g *googleApi) GetUsersFromApi(domain string) ([]data.User, error) {
	result := make([]data.User, 0)
	nextPageToken := ""

	for {
		response, err := g.client.Users.List().Domain(domain).PageToken(nextPageToken).Do()
		if err != nil {
			if apiErr, ok := err.(*googleapi.Error); ok {
				if apiErr.Code == http.StatusTooManyRequests {
					durationInSeconds, err := strconv.ParseInt(apiErr.Header.Get("Retry-After"), 10, 64)
					if err != nil {
						return nil, errors.Wrap(err, "failed to parse `retry-after` header")
					}
					timeoutDuration := time.Second * time.Duration(durationInSeconds)
					time.Sleep(timeoutDuration)
					continue
				}
			}

			g.log.WithError(err).Errorf("failed to get users in domain `%s`", domain)
			return nil, errors.Wrap(err, "failed to get users in domain")
		}

		users := make([]data.User, 0)
		for _, user := range response.Users {
			users = append(users, data.User{
				Email:  user.PrimaryEmail,
				Name:   user.Name.FullName,
				MailId: user.Id,
			})
		}

		result = append(result, users...)

		if response.NextPageToken == "" {
			break
		}
		nextPageToken = response.NextPageToken
	}

	return result, nil
}
