package processor

import (
	"github.com/acs-dl/mail-module-svc/internal/data"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) validateRemoveUser(msg data.ModulePayload) error {
	return validation.Errors{
		"domain": validation.Validate(msg.Link, validation.Required),
		"email":  validation.Validate(msg.Email, validation.Required),
	}.Filter()
}

func (p *processor) HandleRemoveUserAction(msg data.ModulePayload) error {
	p.log.Infof("start handle message action with id `%s`", msg.RequestId)

	err := p.validateRemoveUser(msg)
	if err != nil {
		p.log.WithError(err).Errorf("failed to validate fields for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to validate fields")
	}

	user, err := p.googleClient.GetDomainUserFromApi(msg.Email)
	if err != nil {
		p.log.WithError(err).Errorf("failed to get user from API for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "some error while getting user from api")
	}
	if user == nil {
		p.log.Errorf("user is not in domain for message action with id `%s`", msg.RequestId)
		return errors.Errorf("user is not in domain")
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

	err = p.googleClient.DeleteUserInDomainFromApi(dbUser.Email)
	if err != nil {
		p.log.WithError(err).Errorf("failed to remove user from API for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to remove user from api")
	}

	err = p.managerQ.Transaction(func() error {
		err := p.permissionsQ.FilterByMailIds(user.MailId).FilterByLinks(msg.Link).Delete()
		if err != nil {
			p.log.WithError(err).Errorf("failed to delete permission by mail id `%d` for message action with id `%s`", user.MailId, msg.RequestId)
			return errors.Wrap(err, "failed to delete permission")
		}

		permissionsCount, err := p.permissionsQ.Count().FilterByMailIds(user.MailId).GetTotalCount()
		if err != nil {
			p.log.WithError(err).Errorf("failed to select permissions by mail id `%d` for message action with id `%s`", user.MailId, msg.RequestId)
			return errors.Wrap(err, "failed to select permissions")
		}

		if permissionsCount == 0 {
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
		}

		return nil
	})
	if err != nil {
		p.log.WithError(err).Errorf("failed to make remove user transaction for message action with id `%s`", msg.RequestId)
		return errors.Wrap(err, "failed to make remove user transaction")
	}

	p.resetFilters()
	p.log.Infof("finish handle message action with id `%s`", msg.RequestId)
	return nil
}
