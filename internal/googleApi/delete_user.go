package googleApi

import (
	"net/http"
	"strconv"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"google.golang.org/api/googleapi"
)

func (g *googleApi) DeleteUserInDomainFromApi(email string) error {
	err := g.client.Users.Delete(email).Do()
	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok {
			if apiErr.Code == http.StatusTooManyRequests {
				durationInSeconds, err := strconv.ParseInt(apiErr.Header.Get("Retry-After"), 10, 64)
				if err != nil {
					return errors.Wrap(err, "failed to parse `retry-after` header")
				}
				timeoutDuration := time.Second * time.Duration(durationInSeconds)
				time.Sleep(timeoutDuration)
				return g.DeleteUserInDomainFromApi(email)
			}
		}

		g.log.WithError(err).Errorf("failed to delete user in domain with email `%s`", email)
		return errors.Wrap(err, "failed to delete user in domain")
	}

	return nil
}
