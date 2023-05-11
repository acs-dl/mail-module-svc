package service

import (
	"context"
	"sync"

	"gitlab.com/distributed_lab/acs/mail-module/internal/service/api"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/api/handlers"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/googleApi"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/pqueue"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/processor"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/receiver"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/registrator"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/sender"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/worker"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/acs/mail-module/internal/config"
)

type svc struct {
	Name    string
	New     func(config.Config, context.Context) interface{}
	Run     func(interface{}, context.Context)
	Context func(interface{}, context.Context) context.Context
}

var services = []svc{
	{"mail", googleApi.NewGoogleAsInterface, nil, googleApi.CtxGoogleClientInstance},
	{"sender", sender.NewSenderAsInterface, sender.RunSenderAsInterface, sender.CtxSenderInstance},
	{"processor", processor.NewProcessorAsInterface, nil, processor.CtxProcessorInstance},
	{"worker", worker.NewWorkerAsInterface, worker.RunWorkerAsInterface, worker.CtxWorkerInstance},
	{"receiver", receiver.NewReceiverAsInterface, receiver.RunReceiverAsInterface, receiver.CtxReceiverInstance},
	{"registrar", registrator.NewRegistrarAsInterface, registrator.RunRegistrarAsInterface, nil},
	{"api", api.NewRouterAsInterface, api.RunRouterAsInterface, nil},
}

func Run(cfg config.Config) {
	logger := cfg.Log().WithField("service", "main")
	ctx := context.Background()
	wg := new(sync.WaitGroup)

	logger.Info("Starting all available services...")

	stopProcessQueue := make(chan struct{})
	pqueues := pqueue.NewPQueues()
	go pqueues.SuperUserPQueue.ProcessQueue(cfg.RateLimit().RequestsAmount, cfg.RateLimit().TimeLimit, stopProcessQueue)
	go pqueues.UserPQueue.ProcessQueue(cfg.RateLimit().RequestsAmount, cfg.RateLimit().TimeLimit, stopProcessQueue)
	ctx = pqueue.CtxPQueues(&pqueues, ctx)
	ctx = handlers.CtxConfig(cfg, ctx)

	for _, mySvc := range services {
		wg.Add(1)

		instance := mySvc.New(cfg, ctx)
		if instance == nil {
			logger.WithField("service", mySvc.Name).Warn("Service instance not created")
			panic(errors.Errorf("`%s` instance not created", mySvc.Name))
		}

		if mySvc.Context != nil {
			ctx = mySvc.Context(instance, ctx)
		}

		if mySvc.Run != nil {
			wg.Add(1)
			go func(structure interface{}, runner func(interface{}, context.Context)) {
				defer wg.Done()

				runner(structure, ctx)

			}(instance, mySvc.Run)
		}
		logger.WithField("service", mySvc.Name).Info("Service started")
	}

	wg.Wait()
}
