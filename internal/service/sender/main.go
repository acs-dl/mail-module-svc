package sender

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.com/distributed_lab/logan/v3"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/acs-dl/mail-module-svc/internal/config"
	"github.com/acs-dl/mail-module-svc/internal/data"
	"github.com/acs-dl/mail-module-svc/internal/data/postgres"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

const ServiceName = data.ModuleName + "-sender"

type Sender struct {
	publisher   *amqp.Publisher
	responsesQ  data.Responses
	log         *logan.Entry
	topic       string
	runnerDelay time.Duration
}

func NewSenderAsInterface(cfg config.Config, _ context.Context) interface{} {
	return interface{}(&Sender{
		publisher:   cfg.Amqp().Publisher,
		responsesQ:  postgres.NewResponsesQ(cfg.DB()),
		log:         logan.New().WithField("service", ServiceName),
		topic:       cfg.Amqp().Orchestrator,
		runnerDelay: cfg.Runners().Sender,
	})
}

func (s *Sender) Run(ctx context.Context) {
	go running.WithBackOff(ctx, s.log,
		ServiceName,
		s.processMessages,
		s.runnerDelay,
		s.runnerDelay,
		s.runnerDelay,
	)
}

func (s *Sender) processMessages(_ context.Context) error {
	s.log.Info("started processing responses")

	responses, err := s.responsesQ.Select()
	if err != nil {
		s.log.WithError(err).Errorf("failed to select responses")
		return errors.Wrap(err, "failed to select responses")
	}

	for _, response := range responses {
		s.log.Info("started processing response with id ", response.ID)
		err = (*s.publisher).Publish(s.topic, s.buildResponse(response))
		if err != nil {
			s.log.WithError(err).Errorf("failed to process response `%s", response.ID)
			return errors.Wrap(err, "failed to process response: "+response.ID)
		}

		err = s.responsesQ.FilterByIds(response.ID).Delete()
		if err != nil {
			s.log.WithError(err).Errorf("failed to delete processed response `%s", response.ID)
			return errors.Wrap(err, "failed to delete processed response: "+response.ID)
		}
		s.log.Info("finished processing response with id ", response.ID)
	}

	s.log.Info("finished processing responses")
	return nil
}

func (s *Sender) buildResponse(response data.Response) *message.Message {
	marshaled, err := json.Marshal(response)
	if err != nil {
		s.log.WithError(err).Errorf("failed to marshal response")
	}

	return &message.Message{
		UUID:     response.ID,
		Metadata: nil,
		Payload:  marshaled,
	}
}

func (s *Sender) SendMessageToCustomChannel(topic string, msg *message.Message) error {
	err := (*s.publisher).Publish(topic, msg)
	if err != nil {
		s.log.WithError(err).Errorf("failed to send msg `%s to `%s`", msg.UUID, topic)
		return errors.Wrap(err, "failed to send msg: "+msg.UUID)
	}

	return nil
}
