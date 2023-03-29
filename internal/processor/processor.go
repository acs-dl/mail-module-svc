package processor

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/acs/mail-module/internal/config"
	"gitlab.com/distributed_lab/acs/mail-module/internal/data"
	"gitlab.com/distributed_lab/acs/mail-module/internal/data/manager"
	"gitlab.com/distributed_lab/acs/mail-module/internal/data/postgres"
	"gitlab.com/distributed_lab/acs/mail-module/internal/googleApi"
	"gitlab.com/distributed_lab/acs/mail-module/internal/sender"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	serviceName = data.ModuleName + "-processor"

	//add needed actions for module
	GetUsersAction   = "get_users"
	AddUserAction    = "add_user"
	RemoveUserAction = "remove_user"
	VerifyUserAction = "verify_user"
	DeleteUserAction = "delete_user"

	SetUsersAction    = "set_users"
	DeleteUsersAction = "delete_users"
)

type Processor interface {
	HandleNewMessage(msg data.ModulePayload) error
	SendDeleteUser(uuid string, user data.User) error
}

type processor struct {
	log          *logan.Entry
	googleClient googleApi.GoogleClient
	permissionsQ data.Permissions
	usersQ       data.Users
	managerQ     *manager.Manager
	sender       *sender.Sender
}

var handleActions = map[string]func(proc *processor, msg data.ModulePayload) error{
	GetUsersAction:   (*processor).handleGetUsersAction,
	AddUserAction:    (*processor).handleAddUserAction,
	RemoveUserAction: (*processor).handleRemoveUserAction,
	VerifyUserAction: (*processor).handleVerifyUserAction,
	DeleteUserAction: (*processor).handleDeleteUserAction,
}

func NewProcessor(cfg config.Config) Processor {
	return &processor{
		log:          cfg.Log().WithField("service", serviceName),
		googleClient: googleApi.NewGoogle(cfg.Log()),
		permissionsQ: postgres.NewPermissionsQ(cfg.DB()),
		usersQ:       postgres.NewUsersQ(cfg.DB()),
		managerQ:     manager.NewManager(cfg.DB()),
		sender:       sender.NewSender(cfg),
	}
}

func (p *processor) HandleNewMessage(msg data.ModulePayload) error {
	p.log.Infof("handling message with id `%s`", msg.RequestId)

	err := validation.Errors{
		"action": validation.Validate(msg.Action, validation.Required, validation.In(GetUsersAction, AddUserAction, RemoveUserAction, DeleteUserAction, VerifyUserAction)),
	}.Filter()
	if err != nil {
		p.log.WithError(err).Errorf("no such action `%s` to handle for message with id `%s`", msg.Action, msg.RequestId)
		return errors.Wrap(err, fmt.Sprintf("no such action `%s` to handle for message with id `%s`", msg.Action, msg.RequestId))
	}

	requestHandler := handleActions[msg.Action]
	if err = requestHandler(p, msg); err != nil {
		p.log.WithError(err).Errorf("failed to handle message with id `%s`", msg.RequestId)
		return err
	}

	p.log.Infof("finish handling message with id `%s`", msg.RequestId)
	return nil
}
