
run:
	@go run ./cmd/server/...

run.agent:
	@go run ./cmd/agent/...

# Обновление кода автотестов
tmpl.update:
	@git fetch template && git checkout template/main .github

build.all: build.server build.agent

build.server:
	@go build -o cmd/server/server cmd/server/*.go

build.agent:
	@go build -o cmd/agent/agent cmd/agent/*.go

test.get:
	@wget https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.9.16/metricstest
	@chmod +x metricstest

test.all: build.all
	@metricstest -test.v -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

# make test.server.i i=1
test.server.i: build.server
	@metricstest -test.v -binary-path=cmd/server/server -test.run=^TestIteration$(i)[AB]*$

# make test.agent.i i=1
test.agent.i: build.agent
	@metricstest -test.v -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -test.run=^TestIteration$(i)[AB]*$