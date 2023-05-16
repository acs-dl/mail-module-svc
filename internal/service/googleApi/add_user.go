package googleApi

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

func (g *googleApi) AddUserInDomainFromApi(firstName, lastName, email string) (*data.User, error) {
	response, err := g.client.Users.Insert(&admin.User{
		ChangePasswordAtNextLogin: true,
		Name: &admin.UserName{
			FamilyName: lastName,
			GivenName:  firstName,
		},
		Emails: []admin.UserEmail{
			{
				Type:    "work",
				Address: email,
			},
		},
		PrimaryEmail: fmt.Sprintf("%s.%s@centrilisedgym.online", firstName, lastName),
		Password:     fmt.Sprintf("%s.%s@centrilisedgym.online", firstName, lastName),
	}).Do()
	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok {
			if apiErr.Code == http.StatusTooManyRequests {
				durationInSeconds, err := strconv.ParseInt(apiErr.Header.Get("Retry-After"), 10, 64)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse `retry-after` header")
				}
				timeoutDuration := time.Second * time.Duration(durationInSeconds)
				time.Sleep(timeoutDuration)
				return g.AddUserInDomainFromApi(firstName, lastName, email)
			}
		}

		g.log.WithError(err).Errorf("failed to add user in domain with email `%s`", email)
		return nil, errors.Wrap(err, "failed to add user in domain")
	}

	return &data.User{
		Email:  response.PrimaryEmail,
		Name:   firstName + " " + lastName,
		MailId: response.Id,
	}, nil
}
