package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/YxTiBlya/ci-api/pkg/executor"
	"github.com/YxTiBlya/ci-monitor/pkg/models"
	"github.com/YxTiBlya/ci-scheduler/internal/service/dto"
)

func (svc *Service) Schedule(msgs <-chan amqp.Delivery) {
	for b := range msgs {
		svc.log.Info().Str("body", string(b.Body)).Msg("received message")

		var msg models.QSPipelineMsg
		if err := json.Unmarshal(b.Body, &msg); err != nil {
			svc.log.Error().Err(err).
				Str("body", string(b.Body)).
				Msg("failed to unmarshal message")
			continue
		}

		// to nginx
		id := uuid.New().String()
		for _, task := range msg.Pipeline {
			svc.wg.Add(1)
			go svc.execute(id, msg.Repo, task)
		}
	}
}

func (svc *Service) execute(id, repo string, task models.Pipeline) {
	defer func() {
		svc.wg.Done()
	}()

	ctx := context.Background()
	resp, err := svc.ExecutorAPI.ExecuteTask(ctx, &executor.ExecuteRequest{
		Repo: repo,
		Cmd:  task.Command,
	})
	if err != nil {
		svc.log.Error().Err(err).Msg("failed to execute task")
		return // thats ok because executor anyway returns resp
	}

	svc.log.Info().Str("resp", resp.String()).Msg("executed task")

	data := &dto.TaskData{
		Repo:       repo,
		PipelineID: id,
		Name:       task.Name,
		Command:    task.Command,
		Status:     executor.ExecuteStatus_name[int32(resp.Status)],
		ExitCode:   resp.ExitCode,
		Output:     resp.Output,
		ExecTime:   resp.Time,
	}
	if err := svc.Relations.DB.InsertPipline(ctx, data); err != nil {
		svc.log.Error().Err(err).Msg("failed to insert pipeline")
	}
}
