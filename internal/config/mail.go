package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type MailCfg struct {
	Subject string `figure:"subject,required"`
}

func (c *config) Mail() *MailCfg {
	return c.mail.Do(func() interface{} {
		var config MailCfg
		err := figure.
			Out(&config).
			From(kv.MustGetStringMap(c.getter, "mail")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out mail from config"))
		}

		return &config
	}).(*MailCfg)
}
