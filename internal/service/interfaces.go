package service

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"

	"github.com/YxTiBlya/ci-api/pkg/executor"
	"github.com/YxTiBlya/ci-core/rabbitmq"
)

type QueryService interface {
	Consume(qName, consumer string, ack, excl, nlocal, nwait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	AddMigrates(migates ...rabbitmq.Migrate) error
}

type ExecutorAPIClient interface {
	ExecuteTask(ctx context.Context, in *executor.ExecuteRequest, opts ...grpc.CallOption) (*executor.ExecuteResponse, error)
}
