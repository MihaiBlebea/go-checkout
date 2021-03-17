setup: build up

build:
	docker build -t serbanblebea/go-checkout:v0.1 .

up:
	docker run --name go-checkout -d serbanblebea/go-checkout:v0.1