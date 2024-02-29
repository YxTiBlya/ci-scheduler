package service

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/YxTiBlya/ci-core/rabbitmq"
)

type Relations struct {
	QS          QueryService
	ExecutorAPI ExecutorAPIClient
}

func New(cfg Config, log *zap.SugaredLogger, rel Relations) *Service {
	return &Service{
		cfg:       cfg,
		log:       log,
		Relations: rel,
	}
}

type Service struct {
	cfg Config
	log *zap.SugaredLogger
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
	return nil
}
