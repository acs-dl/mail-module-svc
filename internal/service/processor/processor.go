package processor

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"

	"github.com/acs-dl/mail-module-svc/internal/config"
	"github.com/acs-dl/mail-module-svc/internal/data"
	"github.com/acs-dl/mail-module-svc/internal/data/manager"
	"github.com/acs-dl/mail-module-svc/internal/data/postgres"
	"github.com/acs-dl/mail-module-svc/internal/service/googleApi"
	"github.com/acs-dl/mail-module-svc/internal/service/sender"
)

const (
	ServiceName = data.ModuleName + "-processor"

	//add needed actions for module
	SetUsersAction    = "set_users"
	DeleteUsersAction = "delete_users"
)

type Processor interface {
	HandleGetUsersAction(msg data.ModulePayload) error
	HandleAddUserAction(msg data.ModulePayload) error
	HandleRemoveUserAction(msg data.ModulePayload) error
	HandleDeleteUserAction(msg data.ModulePayload) error
	HandleVerifyUserAction(msg data.ModulePayload) error
	SendDeleteUser(uuid string, user data.User) error
}

type processor struct {
	log             *logan.Entry
	googleClient    googleApi.GoogleClient
	permissionsQ    data.Permissions
	usersQ          data.Users
	managerQ        *manager.Manager
	sender          *sender.Sender
	unverifiedTopic string
}

func NewProcessorAsInterface(cfg config.Config, ctx context.Context) interface{} {
	return interface{}(&processor{
		log:             cfg.Log().WithField("service", ServiceName),
		googleClient:    googleApi.GoogleClientInstance(ctx),
		permissionsQ:    postgres.NewPermissionsQ(cfg.DB()),
		usersQ:          postgres.NewUsersQ(cfg.DB()),
		managerQ:        manager.NewManager(cfg.DB()),
		sender:          sender.SenderInstance(ctx),
		unverifiedTopic: cfg.Amqp().Unverified,
	})
}

func ProcessorInstance(ctx context.Context) Processor {
	return ctx.Value(ServiceName).(Processor)
}

func CtxProcessorInstance(entry interface{}, ctx context.Context) context.Context {
	return context.WithValue(ctx, ServiceName, entry)
}
