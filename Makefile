NAME=PasswordManager
PORT=5000
DOCKER_NAME=password_manager
MAIN_PATH=./cmd/server/main.go
TARGET_PATH=./bin
TARGET=${TARGET_PATH}/${NAME}
DB_CONNECTION=postgresql://admin:12345@localhost:5432/password_manager?sslmode=disable
MIGRATIONS_PATH=./internal/repo/migrations
.DEFAULT_GOAL := help
export PM_SERVER_PORT=${PORT}

build: ## Make build of the project
	GOARCH=amd64 GOOS=darwin go build -o ${TARGET} ${MAIN_PATH}

run: ## Run the project
	${TARGET_PATH}/${NAME}

docker-build: ## Create an image in docker
	docker build -t ${DOCKER_NAME} ./

docker-run: ## Run container
	docker run -p ${PORT}:${PORT} --name="${DOCKER_NAME}" ${DOCKER_NAME}

docker-start: docker-build docker-run ## Create an image in docker and run a container

docker-stop: ## Delete an image and container with name "password-manager"
	docker stop ${DOCKER_NAME}; docker rm ${DOCKER_NAME}; docker rmi -f ${DOCKER_NAME}

compose-up: ## Start all the services from docker-compose file in detached mode
	docker-compose up -d

compose-stop: ## Drop all the services from docker-compose file
	docker-compose down

delete-docker-image: ## Deletes an image of your program
	docker image rmi ${DOCKER_NAME}

compose-down: compose-stop delete-docker-image

up: build run ## Build and run the project

migration-up: ## Up migrates
	migrate -path ${MIGRATIONS_PATH} -database ${DB_CONNECTION} up

migration-down: ## Drop migrates
	migrate -path ${MIGRATIONS_PATH} -database ${DB_CONNECTION} down

db-connect: ## Open postgres container and connect to DB
	docker exec -it postgres psql ${DB_CONNECTION}

clean:
	go clean
	rm ${TARGET}

help: ## Display this help screen
	@grep -E -h '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build run docker-build docker-run docker-start docker-stop up migration-down migration-up db-connect
