setup: build up

# build:
# 	docker build -t serbanblebea/go-lambda:0.1 .

# run:
# 	docker run -v /var/run/docker.sock:/var/run/docker.sock \
# 			-d --rm -p 8082:8082 \
# 			--name go-lambda \
# 			--env-file=.env \
# 			serbanblebea/go-lambda:0.1

root:
	docker exec -it --user root go-lambda /bin/sh

nobody:
	docker exec -it go-lambda /bin/sh

stop:
	docker stop go-lambda

refresh: stop setup nobody

go-build:
	go build -o=checkout .

# build-up: build up

build:
	docker-compose build

up:
	docker-compose up -d

migrate-docker:
	docker exec -it dashboard /bin/sh -c "./lambda migrate"

build-logs:
	docker exec -it dashboard /bin/sh -c "./lambda build-logs ceva"