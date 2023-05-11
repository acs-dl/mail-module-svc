package processor

import (
	"strconv"

	"github.com/acs-dl/mail-module-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) validateVerifyUser(msg data.ModulePayload) error {
	return validation.Errors{
		"user_id": validation.Validate(msg.UserId, validation.Required),
		"email":   validation.Validate(msg.Email, validation.Required),
	}.Filter()
}

func (p *processor) HandleVerifyUserAction(msg data.ModulePayload) error {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	err := p.validateVerifyUser(msg)
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
		return errors.Wrap(err, "some error while getting user from api")
	}
	if user == nil {
		p.log.Errorf("no user was found in domain")
		return errors.Errorf("no user was found in domain")
	}
	user.Id = &userId

	if err = p.usersQ.Upsert(*user); err != nil {
		p.log.WithError(err).Errorf("failed to upsert user in db for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to upsert user in db")
	}

	err = p.SendDeleteUser(msg.RequestId, *user)
	if err != nil {
		p.log.WithError(err).Errorf("failed to publish delete user for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to publish delete user")
	}

	p.resetFilters()
	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return nil
}
