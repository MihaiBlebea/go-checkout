setup: env-file build up

setup-test: env-file build-test up-test cover-html

env-file: 
	cp ./.env.example ./.env

build:
	docker build \
		-t serbanblebea/go-checkout:v0.1 \
		.

up:
	docker run \
		--rm \
		--name go-checkout \
		-d \
		-p 8087:8087 \
		--env-file ./.env \
		serbanblebea/go-checkout:v0.1

stop: 
	docker stop go-checkout

build-test:
	docker build \
		--no-cache \
		--file ./Dockerfile.test \
		-t serbanblebea/go-checkout:test \
		.

up-test:
	docker run \
		-v ${PWD}:/app \
		--rm \
		--name go-checkout-test \
		--env-file ./.env \
		serbanblebea/go-checkout:test

cover-html:
	go tool \
		cover -html=cover.out \
		-o cover.html \
		&& open cover.html

go-build:
	go build -o=./checkout .

go-test:
	go test -v ./...
