build:
	@go build -o ./bin/scheduler ./cmd/main/main.go

run: build
	@./bin/scheduler --cfg config.yaml