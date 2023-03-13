SHELL=cmd.exe
FRONTEND_BINARY=frontendApp.exe
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker-compose
up_build: build_broker build_auth
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

## build_auth: builds the auth binary as a Linux executable
build_auth:
	@echo Building auth binary...
	chdir ../auth && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${AUTH_BINARY} ../auth/cmd/api
	@echo Done!

## build_frontend: builds the frontend binary
build_frontend:
	@echo Building frontend binary...
	chdir ../frontend && set GOOS=windows&& set GOARCH=amd64&& set CGO_ENABLED=0\
		&& go build -o ${FRONTEND_BINARY} ../frontend/cmd/web
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
