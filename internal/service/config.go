package service

type Config struct {
	GRPCExecutorAddress string `yaml:"grpc_executor_address"`
	QSName              string
}
