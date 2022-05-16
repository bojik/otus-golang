BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

build: build-calendar build-scheduler build-sender

run-calendar: build-calendar
	$(BIN) --config ./configs/config.yaml
run-scheduler: build-scheduler
	$(BIN_SCHEDULER) --config ./configs/config_scheduler.yaml
run-sender: build-sender
	$(BIN_SENDER) --config ./configs/config_sender.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/calendar/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

up:
	docker-compose up -d

down:
	docker-compose down

version-calendar: build-calendar
	$(BIN) --version
version-scheduler: build-scheduler
	$(BIN_SCHEDULER) --version
version-sender: build-sender
	$(BIN_SENDER) --version

dump-config: build-calendar
	$(BIN) --dump-config

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

test-integration:
	docker-compose -f docker-compose.yml -f docker-compose-test.yml up -d --scale test-integration=0; \
	docker-compose -f docker-compose.yml -f docker-compose-test.yml up --exit-code-from test-integration test-integration; if [ $$? -eq 0 ] ; then result=0; else result=1; fi;\
	docker-compose -f docker-compose.yml -f docker-compose-test.yml down; exit $$result
.PHONY: build run build-img run-img version test lint
