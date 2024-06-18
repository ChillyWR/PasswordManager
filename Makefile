NAME=PasswordManager
PORT=5000
DOCKER_NAME=password_manager
MAIN_PATH=./cmd/server/main.go
TARGET_PATH=./bin
TARGET=${TARGET_PATH}/${NAME}
DB_CONNECTION=postgresql://admin:12345@localhost:5432/password_manager?sslmode=disable
MIGRATIONS_PATH=./internal/repo/migrations
MOCKGEN=$(shell go env GOPATH)/bin/mockgen
.DEFAULT_GOAL := help
export PM_SERVER_PORT=${PORT}

build: ## Make build of the project
	GOARCH=amd64 GOOS=darwin go build -o ${TARGET} ${MAIN_PATH}
.PHONY: build

run: ## Run the project
	${TARGET_PATH}/${NAME}
.PHONY: run

docker-build: ## Create an image in docker
	docker build -t ${DOCKER_NAME} ./
.PHONY: docker-build

docker-run: ## Run container
	docker run -p ${PORT}:${PORT} --name="${DOCKER_NAME}" ${DOCKER_NAME}
.PHONY: docker-run

docker-start: docker-build docker-run ## Create an image in docker and run a container
.PHONY: docker-start

docker-stop: ## Delete an image and container with name "password-manager"
	docker stop ${DOCKER_NAME}; docker rm ${DOCKER_NAME}; docker rmi -f ${DOCKER_NAME}
.PHONY: docker-stop

compose-up: ## Start all the services from docker-compose file in detached mode
	docker-compose up -d
.PHONY: compose-up

compose-stop: ## Drop all the services from docker-compose file
	docker-compose down
.PHONY: compose-stop

delete-docker-image: ## Deletes an image of your program
	docker image rmi ${DOCKER_NAME}
.PHONY: delete-docker-image

compose-down: compose-stop delete-docker-image
.PHONY: compose-down

up: build run ## Build and run the project
.PHONY: up

migration-up: ## Up migrates
	migrate -path ${MIGRATIONS_PATH} -database ${DB_CONNECTION} up
.PHONY: migration-up

migration-down: ## Drop migrates
	migrate -path ${MIGRATIONS_PATH} -database ${DB_CONNECTION} down
.PHONY: migration-down

connect-db: ## Open postgres container and connect to DB
	docker exec -it postgres psql ${DB_CONNECTION}
.PHONY: connect-db

mockgen-controller:
	$(MOCKGEN) -package mock -destination internal/mock/controller.go -source=internal/controller/controller.go

clean:
	go clean
	rm ${TARGET}
.PHONY: clean

help: ## Display this help screen
	@grep -E -h '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
