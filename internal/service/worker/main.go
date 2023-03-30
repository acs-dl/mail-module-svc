package worker

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.com/distributed_lab/acs/mail-module/internal/config"
	"gitlab.com/distributed_lab/acs/mail-module/internal/data"
	"gitlab.com/distributed_lab/acs/mail-module/internal/data/postgres"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/processor"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const serviceName = data.ModuleName + "-worker"

type Worker interface {
	Run(ctx context.Context)
}

type worker struct {
	logger       *logan.Entry
	processor    processor.Processor
	linksQ       data.Links
	usersQ       data.Users
	permissionsQ data.Permissions
}

func NewWorker(cfg config.Config) Worker {
	return &worker{
		logger:       cfg.Log().WithField("runner", serviceName),
		processor:    processor.NewProcessor(cfg),
		linksQ:       postgres.NewLinksQ(cfg.DB()),
		usersQ:       postgres.NewUsersQ(cfg.DB()),
		permissionsQ: postgres.NewPermissionsQ(cfg.DB()),
	}
}

func (w *worker) Run(ctx context.Context) {
	running.WithBackOff(
		ctx,
		w.logger,
		serviceName,
		w.processPermissions,
		15*time.Minute,
		15*time.Minute,
		15*time.Minute,
	)
}

func (w *worker) processPermissions(_ context.Context) error {
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

func (w *worker) removeOldUsers(borderTime time.Time) error {
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

		err = w.usersQ.Delete(user.MailId)
		if err != nil {
			w.logger.Infof("failed to delete user with telegram id `%d`", user.MailId)
			return errors.Wrap(err, " failed to delete user")
		}
	}

	w.logger.Infof("finished removing old users")
	return nil
}

func (w *worker) removeOldPermissions(borderTime time.Time) error {
	w.logger.Infof("started removing old permissions")

	permissions, err := w.permissionsQ.FilterByLowerTime(borderTime).Select()
	if err != nil {
		w.logger.Infof("failed to select permissions")
		return errors.Wrap(err, " failed to select permissions")
	}

	w.logger.Infof("found `%d` permissions to delete", len(permissions))

	for _, permission := range permissions {
		err = w.permissionsQ.Delete(permission.User.MailId, permission.Link)
		if err != nil {
			w.logger.Infof("failed to delete permission")
			return errors.Wrap(err, " failed to delete permission")
		}
	}

	w.logger.Infof("finished removing old permissions")
	return nil
}

func (w *worker) createPermissions(link string) error {
	if err := w.processor.HandleNewMessage(data.ModulePayload{
		RequestId: "from-worker",
		Action:    processor.GetUsersAction,
		Link:      link,
	}); err != nil {
		w.logger.Infof("failed to get users for link `%s`", link)
		return errors.Wrap(err, "failed to get users")
	}

	return nil
}
