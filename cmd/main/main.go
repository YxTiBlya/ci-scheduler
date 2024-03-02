package main

import (
	"flag"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"

	"github.com/YxTiBlya/ci-api/pkg/executor"
	"github.com/YxTiBlya/ci-core/logger"
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
	logger.Init(logger.DevelopmentConfig)
	flag.StringVar(&cfgPath, "cfg", "config.yaml", "app cfg path")
	flag.Parse()
}
func main() {
	yamlFile, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatal("failed to open config file", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatal("failed to unmarshal config file", err)
	}
	cfg.Service.QSName = cfg.QSName

	rmq, err := rabbitmq.NewRabbitMQ(rabbitmq.WithConfig(cfg.RabbitMQ))
	if err != nil {
		log.Fatal("failed to create rabbitmq", err)
	}

	grpcExecutor := executor.NewClient(
		cfg.Service.GRPCExecutorAddress,
		executor.WithClientDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)

	svc := service.New(cfg.Service, service.Relations{
		QS:          rmq,
		ExecutorAPI: grpcExecutor,
	})

	sch := scheduler.NewScheduler(
		scheduler.NewComponent("rabbitmq", rmq),
		scheduler.NewComponent("grpc executor", grpcExecutor),
		scheduler.NewComponent("service", svc),
	)
	sch.Start()
}
