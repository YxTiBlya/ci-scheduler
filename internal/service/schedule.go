package service

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/YxTiBlya/ci-api/pkg/executor"
	"github.com/YxTiBlya/ci-monitor/pkg/models"
)

func (svc *Service) Schedule(msgs <-chan amqp.Delivery) {
	for b := range msgs {
		svc.log.Infof("Received a message: %s", b.Body)

		var msg models.QSPipelineMsg
		if err := json.Unmarshal(b.Body, &msg); err != nil {
			svc.log.Error("Failed to unmarshal message", zap.String("body", string(b.Body)), zap.Error(err))
			continue
		}

		// for test
		// TODO: how i can balance and use more executors?
		for _, task := range msg.Pipeline {
			resp, err := svc.ExecutorAPI.ExecuteTask(context.Background(), &executor.ExecuteRequest{
				Repo: msg.Repo,
				Cmd:  task.Command,
			})
			if err != nil {
				svc.log.Error("Failed to execute task", zap.Error(err))
				continue
			}
			svc.log.Info("Successfully executed task", zap.Any("resp", resp))
		}
	}
}
