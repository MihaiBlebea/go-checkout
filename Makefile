setup: build up

env-file: 
	cp ./.env.example ./.env

build:
	docker build -t serbanblebea/go-checkout:v0.1 .

up:
	docker run --rm --name go-checkout -d --env-file ./.env serbanblebea/go-checkout:v0.1

stop: 
	docker stop go-checkout

go-build:
	go build -o=./checkout .

go-test:
	go test -v ./...
