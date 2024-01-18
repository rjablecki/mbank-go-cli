DOCKER_USERNAME ?= rjablecki
APPLICATION_NAME ?= mbank-go-cli
GIT_HASH ?= $(shell git log --format="%h" -n 1)

test:
	go mod tidy
	go test -v ./...

test-cover:
	go mod tidy
	go test -cover -v ./...

run:
	go build -o app
	./app

cloud-run-describe:
	gcloud run services describe hello

docker-build:
	docker build --no-cache --tag ${DOCKER_USERNAME}/${APPLICATION_NAME}:latest .

docker-push:
	docker push ${DOCKER_USERNAME}/${APPLICATION_NAME}:latest
	
run-docker:
	docker run --env-file .env.local --rm ${DOCKER_USERNAME}/${APPLICATION_NAME}:latest

