package sender

import (
	"context"

	"gitlab.com/distributed_lab/acs/mail-module/internal/config"
)

func Run(ctx context.Context, cfg config.Config) {
	NewSender(cfg).Run(ctx)
}
