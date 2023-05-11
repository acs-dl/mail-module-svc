package processor

import (
	"regexp"

	"github.com/acs-dl/mail-module-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) validateDeleteUser(msg data.ModulePayload) error {
	return validation.Errors{
		"email": validation.Validate(msg.Email, validation.Required, validation.Match(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))),
	}.Filter()
}

func (p *processor) HandleDeleteUserAction(msg data.ModulePayload) error {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	err := p.validateDeleteUser(msg)
	if err != nil {
		p.log.WithError(err).Errorf("failed to validate fields for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to validate fields")
	}

	user, err := p.googleClient.GetDomainUserFromApi(msg.Email)
	if err != nil {
		p.log.WithError(err).Errorf("failed to get user from API for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to get user from api")
	}

	if user == nil {
		p.log.WithError(err).Errorf("user is not in domain for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "user is not in domain")
	}

	dbUser, err := p.usersQ.FilterByMailIds(user.MailId).Get()
	if err != nil {
		p.log.WithError(err).Errorf("failed to get user for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to get user")
	}
	if dbUser == nil {
		p.log.Errorf("no such user in module for message action with id `%s`", msg.RequestId)
		return errors.Errorf("no such user in module")
	}

	permissions, err := p.permissionsQ.FilterByMailIds(user.MailId).Select()
	if err != nil {
		p.log.WithError(err).Errorf("failed to select permissions by mail id `%d` for message action with id `%s`", user.MailId, msg.RequestId)
		return errors.Wrap(err, "failed to select permissions")
	}

	for _, permission := range permissions {
		domainUser, err := p.googleClient.GetDomainUserFromApi(msg.Email)
		if err != nil {
			p.log.WithError(err).Errorf("failed to get domain user from API for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "some error while checking user from api")
		}

		if domainUser != nil {
			err = p.googleClient.DeleteUserInDomainFromApi(domainUser.Email)
			if err != nil {
				p.log.WithError(err).Errorf("failed to remove user from API for message action with id `%s`", msg.RequestId)
				return errors.Wrap(err, "some error while removing user from api")
			}
		}
		if err = p.permissionsQ.FilterByMailIds(user.MailId).FilterByLinks(permission.Link).Delete(); err != nil {
			p.log.WithError(err).Errorf("failed to delete permission from db for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to delete permission")
		}
	}

	err = p.usersQ.FilterByMailIds(user.MailId).Delete()
	if err != nil {
		p.log.WithError(err).Errorf("failed to delete user by mail id `%d` for message action with id `%s`", user.MailId, msg.RequestId)
		return errors.Wrap(err, "failed to delete user")
	}

	if dbUser.Id == nil {
		err = p.SendDeleteUser(msg.RequestId, *dbUser)
		if err != nil {
			p.log.WithError(err).Errorf("failed to publish delete user for message action with id `%s`", msg.RequestId)
			return errors.Wrap(err, "failed to publish delete user")
		}
	}
	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	p.resetFilters()
	return nil
}
