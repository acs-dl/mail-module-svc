package worker

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"time"

	"github.com/acs-dl/mail-module-svc/internal/config"
	"github.com/acs-dl/mail-module-svc/internal/data"
	"github.com/acs-dl/mail-module-svc/internal/data/postgres"
	"github.com/acs-dl/mail-module-svc/internal/service/processor"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const ServiceName = data.ModuleName + "-worker"

type IWorker interface {
	Run(ctx context.Context)
	ProcessPermissions(_ context.Context) error
	GetEstimatedTime() time.Duration
}

type Worker struct {
	logger        *logan.Entry
	processor     processor.Processor
	linksQ        data.Links
	usersQ        data.Users
	permissionsQ  data.Permissions
	runnerDelay   time.Duration
	estimatedTime time.Duration
}

func NewWorkerAsInterface(cfg config.Config, ctx context.Context) interface{} {
	return interface{}(&Worker{
		logger:        cfg.Log().WithField("runner", ServiceName),
		processor:     processor.ProcessorInstance(ctx),
		linksQ:        postgres.NewLinksQ(cfg.DB()),
		usersQ:        postgres.NewUsersQ(cfg.DB()),
		permissionsQ:  postgres.NewPermissionsQ(cfg.DB()),
		runnerDelay:   cfg.Runners().Worker,
		estimatedTime: time.Duration(0),
	})
}

func (w *Worker) Run(ctx context.Context) {
	running.WithBackOff(
		ctx,
		w.logger,
		ServiceName,
		w.ProcessPermissions,
		w.runnerDelay,
		w.runnerDelay,
		w.runnerDelay,
	)
}

func (w *Worker) ProcessPermissions(_ context.Context) error {
	w.logger.Info("fetching links")

	startTime := time.Now()

	links, err := w.linksQ.Select()
	if err != nil {
		return errors.Wrap(err, "failed to get links")
	}

	reqAmount := len(links)
	if reqAmount == 0 {
		w.logger.Info("no links were found")
		return nil
	}

	w.logger.Infof("found %v links", reqAmount)

	for _, link := range links {
		w.logger.Infof("processing link `%s`", link.Link)

		err = w.createPermissions(link.Link)
		if err != nil {
			w.logger.Infof("failed to create permissions for subs")
			return errors.Wrap(err, "failed to create permissions for subs")
		}

		w.logger.Infof("successfully processed link `%s`", link.Link)
	}

	err = w.removeOldUsers(startTime)
	if err != nil {
		w.logger.WithError(err).Errorf("failed to remove old users")
		return errors.Wrap(err, "failed to remove old users")
	}

	err = w.removeOldPermissions(startTime)
	if err != nil {
		w.logger.WithError(err).Errorf("failed to remove old permissions")
		return errors.Wrap(err, "failed to remove old permissions")
	}

	return nil
}

func (w *Worker) removeOldUsers(borderTime time.Time) error {
	w.logger.Infof("started removing old users")

	users, err := w.usersQ.FilterByLowerTime(borderTime).Select()
	if err != nil {
		w.logger.Infof("failed to select users")
		return errors.Wrap(err, " failed to select users")
	}

	w.logger.Infof("found `%d` users to delete", len(users))

	for _, user := range users {
		if user.Id == nil { //if unverified user we need to remove them from `unverified-svc`
			err = w.processor.SendDeleteUser(uuid.New().String(), user)
			if err != nil {
				w.logger.WithError(err).Errorf("failed to publish delete user")
				return errors.Wrap(err, " failed to publish delete user")
			}
		}

		err = w.usersQ.FilterByMailIds(user.MailId).Delete()
		if err != nil {
			w.logger.Infof("failed to delete user with telegram id `%d`", user.MailId)
			return errors.Wrap(err, " failed to delete user")
		}
	}

	w.logger.Infof("finished removing old users")
	return nil
}

func (w *Worker) removeOldPermissions(borderTime time.Time) error {
	w.logger.Infof("started removing old permissions")

	permissions, err := w.permissionsQ.FilterByLowerTime(borderTime).Select()
	if err != nil {
		w.logger.Infof("failed to select permissions")
		return errors.Wrap(err, " failed to select permissions")
	}

	w.logger.Infof("found `%d` permissions to delete", len(permissions))

	for _, permission := range permissions {
		err = w.permissionsQ.FilterByMailIds(permission.User.MailId).FilterByLinks(permission.Link).Delete()
		if err != nil {
			w.logger.Infof("failed to delete permission")
			return errors.Wrap(err, " failed to delete permission")
		}
	}

	w.logger.Infof("finished removing old permissions")
	return nil
}

func (w *Worker) RefreshSubmodules(msg data.ModulePayload) error {
	w.logger.Infof("started refresh submodules")

	for _, link := range msg.Links {
		w.logger.Infof("started refreshing `%s`", link)

		err := w.createPermissions(link)
		if err != nil {
			w.logger.Infof("failed to create subs for link `%s", link)
			return errors.Wrap(err, "failed to create subs")
		}
		w.logger.Infof("finished refreshing `%s`", link)
	}

	w.logger.Infof("finished refresh submodules")
	return nil
}

func (w *Worker) createPermissions(link string) error {
	if err := w.processor.HandleGetUsersAction(data.ModulePayload{
		RequestId: "from-worker",
		Link:      link,
	}); err != nil {
		w.logger.Infof("failed to get users for link `%s`", link)
		return errors.Wrap(err, "failed to get users")
	}

	return nil
}

func (w *Worker) GetEstimatedTime() time.Duration {
	return w.estimatedTime
}
