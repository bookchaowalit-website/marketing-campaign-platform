.PHONY: run build test clean proto generate

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

# Application name
APP_NAME=marketing-api-go

# Build path
BUILD_PATH=bin

all: test build

run:
	$(GORUN) cmd/main.go

build:
	$(GOBUILD) -o $(BUILD_PATH)/$(APP_NAME) -v ./cmd/main.go

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_PATH)

proto:
	@echo "Generating protobuf files..."
	cd ../../packages/grpc-protos && make generate

generate:
	@echo "Generating mocks and other files..."
	mockgen -source=internal/domain/repository/campaign.go -destination=infrastructure/repository/mocks/campaign.go -package=mocks

deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

run-docker:
	docker-compose -f ../../docker/docker-compose.yml up -d db

migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/marketing_saas?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/marketing_saas?sslmode=disable" down
