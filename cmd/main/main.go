package main

import (
	"flag"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"

	"github.com/YxTiBlya/ci-api/pkg/executor"
	"github.com/YxTiBlya/ci-core/rabbitmq"
	"github.com/YxTiBlya/ci-core/scheduler"

	"github.com/YxTiBlya/ci-scheduler/internal/service"
)

type Config struct {
	Service  service.Config  `yaml:"scheduler"`
	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`
	QSName   string          `yaml:"qs_name"`
}

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "cfg", "config.yaml", "app cfg path")
	flag.Parse()
}
func main() {
	logger := zap.Must(zap.NewDevelopment()).Sugar()

	yamlFile, err := os.ReadFile(cfgPath)
	if err != nil {
		logger.Fatal("failed to open config file", zap.Error(err))
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		logger.Fatal("failed to unmarshal config file", zap.Error(err))
	}
	cfg.Service.QSName = cfg.QSName

	rmq, err := rabbitmq.NewRabbitMQ(rabbitmq.WithConfig(cfg.RabbitMQ))
	if err != nil {
		logger.Fatal("failed to create rabbitmq", zap.Error(err))
	}

	grpcExecutor := executor.NewClient(
		cfg.Service.GRPCExecutorAddress,
		executor.WithClientDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)

	svc := service.New(cfg.Service, logger, service.Relations{
		QS:          rmq,
		ExecutorAPI: grpcExecutor,
	})

	sch := scheduler.NewScheduler(
		zap.Must(zap.NewDevelopment()).Sugar(),
		scheduler.NewComponent("rabbitmq", rmq),
		scheduler.NewComponent("grpc executor", grpcExecutor),
		scheduler.NewComponent("service", svc),
	)
	sch.Start()
}
