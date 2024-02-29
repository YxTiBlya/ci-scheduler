build-scheduler:
	@cd ./scheduler && go build -o ../bin/scheduler ./cmd/main/main.go

run-scheduler: build-scheduler
	@./bin/scheduler --cfg config.yaml