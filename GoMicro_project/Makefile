FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAIL_BINARY=mailApp
LISTENER_BINARY=listenerApp
FRONT_END_LINUX=frontEndLinux

## up: starts all containers in the background without forcing build
up:
	@echo "Start Docker images..."
	docker-compose up -d
	@echo "Docker images have started!"

## up_build: stops docker-compose(if running), builds all project and starts docker compose
up_build: compile_broker compile_auth compile_logger compile_mail compile_listener
	@echo "Stop docker images (if running)..."
	docker-compose down
	@echo "Build and restart docker images..."
	docker-compose up --build -d
	@echo "Docker images have built and restarted!"

## down: stop docker compose
down:
	@echo "Stop docker compose..."
	docker-compose down
	@echo "Docker images have stopped!"

## compile_broker: builds the broker binary as a linux executable
compile_broker:
	@echo "Build broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Broker service compiled!"

## compile_auth: builds the auth binary as a linux executable
compile_auth:
	@echo "Build auth binary..."
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Authentication service compiled!"

##compile_logger: builds the logger service binary
compile_logger:
	@echo "Build logger service..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Logger service compiled!"

##compile_listener: builds the listener service binary
compile_listener:
	@echo "Build listener service..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Listener service compiled!"

##compile_mail: build mail service binary
compile_mail:
	@echo "Build mail service..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Mail service compiled!"

## compile_front: builds the front-end binary
compile_front:
	@echo "Build front-end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Front-end compiled!"

## compile_front_linux: builds the front-end linux executable
compile_front_linux:
	@echo "Build front-end linux..."
	cd ../front-end && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_END_LINUX} ./cmd/web
	@echo "Front-end linux compiled!"

## start: starts the front-end
start:
		@echo "Start front-end..."
		cd ../front-end && ./${FRONT_END_BINARY} &
		@echo "Front-end started!"

## stop: stop the front-end
stop: 
	@echo "Stop front-end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Front-end stopped!"


