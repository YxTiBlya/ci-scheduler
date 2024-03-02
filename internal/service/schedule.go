package service

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/YxTiBlya/ci-api/pkg/executor"
	"github.com/YxTiBlya/ci-monitor/pkg/models"
)

func (svc *Service) Schedule(msgs <-chan amqp.Delivery) {
	for b := range msgs {
		svc.log.Info().Str("body", string(b.Body)).Msg("Received message")

		var msg models.QSPipelineMsg
		if err := json.Unmarshal(b.Body, &msg); err != nil {
			svc.log.Error().Err(err).
				Str("body", string(b.Body)).
				Msg("Failed to unmarshal message")
			continue
		}

		// to nginx
		for _, task := range msg.Pipeline {
			resp, err := svc.ExecutorAPI.ExecuteTask(context.Background(), &executor.ExecuteRequest{
				Repo: msg.Repo,
				Cmd:  task.Command,
			})
			if err != nil {
				svc.log.Error().Err(err).Msg("Failed to execute task")
				continue
			}
			svc.log.Info().Str("resp", resp.String()).Msg("Executed task")
		}
	}
}
