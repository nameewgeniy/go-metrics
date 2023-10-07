
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
	@metricstest -test.v -server-port=8080 -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server


# make test.i i=1
test.i: build.all
	@metricstest -test.v -server-port=8080 \
		-source-path=. -file-storage-path=/home/work/go/src/go-metrics/tmp/metrics-db.json \
		-agent-binary-path=cmd/agent/agent  \
		-database-dsn='postgres://user:password@localhost:5442/db?sslmode=disable' \
		-binary-path=cmd/server/server \
		-test.run=^TestIteration$(i)[AB]*$

db.up:
	@docker-compose -f build/docker-compose.yaml up -d