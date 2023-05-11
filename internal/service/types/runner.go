package types

import (
	"context"

	"github.com/acs-dl/mail-module-svc/internal/config"
)

type Runner = func(context context.Context, config config.Config)
