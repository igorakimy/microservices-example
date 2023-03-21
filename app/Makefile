SHELL=cmd.exe
FRONTEND_BINARY=frontendApp.exe
FRONTEND_BINARY_LINUX=frontendApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAIL_BINARY=mailerApp
LISTENER_BINARY=listenerApp

BROKER_PROTO_PATH=./../broker/logs
LOGGER_PROTO_PATH=./../logger/logs

DOCKER_USERNAME=igorakimov

BROKER_VERSION=1.0.1
AUTH_VERSION=1.0.0
LOGGER_VERSION=1.0.0
LISTENER_VERSION=1.0.0
MAILER_VERSION=1.0.0
FRONTEND_VERSION=1.0.0
CADDY_VERSION=1.0.0

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker-compose
up_build: build_broker build_auth build_logger build_mail build_listener
	@echo Stopping Docker images...
	docker-compose down
	@echo Building (when required) and starting Docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker-compose
down:
	@echo Stopping Docker compose...
	docker-compose down
	@echo Docker compose stopped!

## build_broker: builds the broker binary as a Linux executable
build_broker:
	@echo Building broker binary...
	chdir ../broker && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${BROKER_BINARY} ../broker/cmd/api
	@echo Done!

## build_logger: builds the logger binary as a Linux executable
build_logger:
	@echo Building logger binary...
	chdir ../logger && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${LOGGER_BINARY} ../logger/cmd/api
	@echo Done!

## build_listener: builds the listener binary as a Linux executable
build_listener:
	@echo Building listener binary...
	chdir ../listener && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${LISTENER_BINARY} ../listener
	@echo Done!

## build_auth: builds the auth binary as a Linux executable
build_auth:
	@echo Building auth binary...
	chdir ../auth && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${AUTH_BINARY} ../auth/cmd/api
	@echo Done!

## build_mail: builds the mail binary as a Linux executable
build_mail:
	@echo Building mail binary...
	chdir ../mail && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${MAIL_BINARY} ../mail/cmd/api
	@echo Done!

## build_frontend: builds the frontend binary
build_frontend:
	@echo Building frontend binary...
	chdir ../frontend && set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${FRONTEND_BINARY} ../frontend/cmd/web
	@echo Done!

## build_frontend_linux: builds the frontend binary for Linux system
build_frontend_linux:
	@echo Building frontend binary...
	chdir ../frontend && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${FRONTEND_BINARY_LINUX} ../frontend/cmd/web
	@echo Done!

## start: starts the frontend
start:
	@echo Starting frontend...
	chdir ../frontend && start /b ${FRONTEND_BINARY} &

## stop: stopping the frontend
stop:
	@echo Stopping frontend...
	taskkill /im ${FRONTEND_BINARY} /f
	@echo Done!

## logger_protobuf: generates go files for logger protocol buffers
logger_protobuf:
	@echo Generating logger protobuf go files...
	protoc --proto_path=${LOGGER_PROTO_PATH} --go_out=${LOGGER_PROTO_PATH} \
		--go_opt=paths=source_relative --go-grpc_out=${LOGGER_PROTO_PATH} \
		--go-grpc_opt=paths=source_relative logs.proto
	@echo Done!

## broker_protobuf: generates go files for broker protocol buffers
broker_protobuf:
	@echo Generating broker protobuf go files...
	protoc --proto_path=${BROKER_PROTO_PATH} --go_out=${BROKER_PROTO_PATH} \
		--go_opt=paths=source_relative --go-grpc_out=${BROKER_PROTO_PATH} \
		--go-grpc_opt=paths=source_relative logs.proto
	@echo Done!

## docker_build_broker: builds docker image for broker service
docker_build_broker:
	@echo Building docker image for broker service...
	docker build -f ./../broker/Dockerfile -t ${DOCKER_USERNAME}/broker-service:${BROKER_VERSION} ./../broker
	@echo Done!

## docker_build_logger: builds docker image for logger service
docker_build_logger:
	@echo Building docker image for logger service...
	docker build -f ./../logger/Dockerfile -t ${DOCKER_USERNAME}/logger-service:${LOGGER_VERSION} ./../logger
	@echo Done!

## docker_build_auth: builds docker image for auth service
docker_build_auth:
	@echo Building docker image for auth service...
	docker build -f ./../auth/Dockerfile -t ${DOCKER_USERNAME}/auth-service:${AUTH_VERSION} ./../auth
	@echo Done!

## docker_build_mail: builds docker image for mail service
docker_build_mail:
	@echo Building docker image for mail service...
	docker build -f ./../mail/Dockerfile -t ${DOCKER_USERNAME}/mail-service:${MAIL_VERSION} ./../mail
	@echo Done!

## docker_build_listener: builds docker image for listener service
docker_build_listener:
	@echo Building docker image for listener service...
	docker build -f ./../listener/Dockerfile -t ${DOCKER_USERNAME}/listener-service:${LISTENER_VERSION} ./../listener
	@echo Done!

## docker_build_frontend: builds docker image for frontend service
docker_build_frontend:
	@echo Building docker image for frontend service...
	docker build -f ./../frontend/Dockerfile -t ${DOCKER_USERNAME}/frontend-service:${FRONTEND_VERSION} ./../frontend
	@echo Done!

## docker_build_caddy: builds docker image for caddy service
docker_build_caddy:
	@echo Building docker image for caddy service...
	docker build -f caddy.Dockerfile -t ${DOCKER_USERNAME}/caddy:${CADDY_VERSION} .
	@echo Done!

## up_stack: up docker stack
up_stack:
	@echo Deploing docker stack
	docker stack deploy -c swarm.yml micro
	@echo Done!

## down_stack: down docker stack
down_stack:
	@echo Removing docker stack
	docker stack rm micro
	@echo Done!