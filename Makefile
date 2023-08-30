-include .env.goose
export

GOCMD=go
CONFIG=config/local_config.yaml
GOBUILD=GO111MODULE=on $(GOCMD) build
CMDNOTEPATH=cmd/service_note/main.go
CMDPRODUCERPATH=cmd/service_stats_sender/main.go
BINNOTEPATH=bin/service_note
BINPRODUCERPATH=bin/service_stats_sender
CGO_ENABLED=0
GOOS=darwin
GOARCH=amd64
PROJECTNAME=note_service

.PHONY: install-goose
bin-deps: ## bin-deps - install binary dependencies
	$(info #Installing goose...)
	@if ! which goose >/dev/null ; then\
  		go install github.com/pressly/goose/v3/cmd/goose@latest; \
  	fi;

.PHONY: up-build
up-build: ## up-build - docker-compose --file docker-compose.yml up --build
	docker-compose --file docker-compose.yml up --build

.PHONY: up
up: ## up - docker-compose --file docker-compose.yml up
	docker-compose --file docker-compose.yml up

.PHONY: up-services
up-services: ## up-services - docker-compose --file docker-compose.yml up --build -d rabbitmq postgres
	docker-compose --file docker-compose.yml up --build -d rabbitmq postgres

.PHONY: down
down: ## down - docker-compose --file docker-compose.yml down -v
	docker-compose --file docker-compose.yml down -v

.PHONY: migrate-up
migrate-up: ## migrate-up
	goose -dir $(MIGRATIONS_DIR) postgres "${POSTGRESQL_URL}" up

.PHONY: migrate-down
migrate-down: ## migrate-down
	goose -dir $(MIGRATIONS_DIR) postgres "${POSTGRESQL_URL}" down

.PHONY: migrate-reset
migrate-reset: ## migrate-reset
	goose -dir $(MIGRATIONS_DIR) postgres "${POSTGRESQL_URL}" reset

define DROP_COMMAND
DO $$$$ DECLARE \
	r RECORD; \
BEGIN \
	FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP \
		EXECUTE '\''DROP TABLE IF EXISTS '\'' || quote_ident(r.tablename) || '\'' CASCADE'\''; \
	END LOOP; \
END $$$$;
endef

.PHONY: migrate-drop
export DROP_COMMAND
migrate-drop: ## migrate-drop
	psql "${POSTGRESQL_URL}" -c '${DROP_COMMAND}'

.PHONY: migrate-create
migrate-create: ## migrate-create name=create_donkey_table
ifeq ($(name),)
	@echo "You forgot to add migration name, example:\nmake create-migration name=create_users_table"
else
	goose -dir $(MIGRATIONS_DIR) create $(name) sql
endif

.PHONY: psql-cli
psql-cli: ## psql-cli
	psql "${POSTGRESQL_URL}"

.PHONY: build-note-linux
build-note-linux: ## build-note-linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o ${BINNOTEPATH} ${CMDNOTEPATH}

.PHONY: build-producer-linux
build-producer-linux: ## build-producer-linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o ${BINPRODUCERPATH} ${CMDPRODUCERPATH}

.PHONY: gen
gen: ## gen - generate proto files
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/auth.proto

.PHONY: build-note
build-note: ## build-note -  note service binary
	make gen && \
	$(GOBUILD) -o ${BINNOTEPATH} ${CMDNOTEPATH}

.PHONY: build-producer
build-producer: ## build-producer producer service binary
	make gen && \
	$(GOBUILD) -o ${BINPRODUCERPATH} ${CMDPRODUCERPATH}

.PHONY: run-note
run-note: ## run-note - service locally
	$(GOCMD) run ${CMDNOTEPATH} --config ${CONFIG}

.PHONY: run-producer
run-producer: ## run-producer - service locally
	$(GOCMD) run ${CMDPRODUCERPATH} --config ${CONFIG}

.PHONY: run-note-race
run-note-race: ## run-note-race - with race detector
	$(GOCMD) run -race ${CMDNOTEPATH} --config ${CONFIG}

.PHONY: run-producer-race
run-producer-race: ## run-rpducer-race -  with race detector
	$(GOCMD) run -race ${CMDPRODUCERPATH} -config ${CONFIG}

.PHONY: test
test: ## test - run tests
	$(GOCMD) test -v ./...


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
