package processor

import (
	"regexp"
	"strconv"
	"time"

	"github.com/acs-dl/mail-module-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) validateAddUser(msg data.ModulePayload) error {
	return validation.Errors{
		"user_id":    validation.Validate(msg.UserId, validation.Required),
		"first_name": validation.Validate(msg.FirstName, validation.Required),
		"last_name":  validation.Validate(msg.LastName, validation.Required),
		"email":      validation.Validate(msg.Email, validation.Required, validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))),
		"domain":     validation.Validate(msg.Link, validation.Required),
	}.Filter()
}

func (p *processor) HandleAddUserAction(msg data.ModulePayload) error {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	err := p.validateAddUser(msg)
	if err != nil {
		p.log.WithError(err).Errorf("failed to validate fields for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to validate fields")
	}

	userId, err := strconv.ParseInt(msg.UserId, 10, 64)
	if err != nil {
		p.log.WithError(err).Errorf("failed to parse user id `%s` for message action with id `%s`", msg.UserId, msg.RequestId)
		return errors.Wrap(err, "failed to parse user id")
	}

	user, err := p.googleClient.GetDomainUserFromApi(msg.Email)
	if err != nil {
		p.log.WithError(err).Errorf("failed to get user from API for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to get user from api")
	}

	if user != nil {
		p.log.WithError(err).Errorf("user is already in domain for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "user is already in domain")
	}

	user, err = p.googleClient.AddUserInDomainFromApi(msg.FirstName, msg.LastName, msg.Email)
	if err != nil {
		p.log.WithError(err).Errorf("failed to add user from API for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to add user from api")
	}

	user.CreatedAt = time.Now()
	user.Id = &userId

	err = p.managerQ.Transaction(func() error {
		if err = p.usersQ.Upsert(*user); err != nil {
			p.log.WithError(err).Errorf("failed to upsert user in db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to upsert user in db")
		}

		if err = p.permissionsQ.Upsert(data.Permission{
			RequestId: msg.RequestId,
			MailId:    user.MailId,
			Link:      msg.Link,
			CreatedAt: user.CreatedAt,
		}); err != nil {
			p.log.WithError(err).Errorf("failed to upsert permission in db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to upsert permission in db")
		}

		return nil
	})
	if err != nil {
		p.log.WithError(err).Errorf("failed to make add user transaction for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to make add user transaction")
	}

	err = p.SendDeleteUser(msg.RequestId, *user)
	if err != nil {
		p.log.WithError(err).Errorf("failed to publish users for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to publish users")
	}

	p.resetFilters()
	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return nil
}

func (p *processor) resetFilters() {
	p.usersQ = p.usersQ.New()
	p.permissionsQ = p.permissionsQ.New()
}
