package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/YxTiBlya/ci-core/logger"
	"github.com/YxTiBlya/ci-core/rabbitmq"
)

type Relations struct {
	QS          QueryService
	ExecutorAPI ExecutorAPIClient
}

func New(cfg Config, rel Relations) *Service {
	return &Service{
		cfg:       cfg,
		log:       logger.New("service"),
		Relations: rel,
	}
}

type Service struct {
	cfg Config
	log *logger.Logger
	Relations
}

func (svc *Service) Start(ctx context.Context) error {
	err := svc.Relations.QS.AddMigrates(
		rabbitmq.WithQueue(&rabbitmq.QueueConfig{Name: svc.cfg.QSName}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to migrate rabbitmq")
	}

	msgs, err := svc.Relations.QS.Consume(svc.cfg.QSName, "", true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to consume rabbitmq")
	}

	go svc.Schedule(msgs)

	return nil
}

func (svc *Service) Stop(ctx context.Context) error {
	svc.log.Sync()
	return nil
}
