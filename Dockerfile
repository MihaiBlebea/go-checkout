FROM golang:1.16.2-buster AS build_base

RUN apt-get install git

# Set the Current Working Directory inside the container
WORKDIR /tmp/app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/lambda .

# Start fresh from a smaller image
FROM debian:buster

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

WORKDIR /app

COPY --from=build_base /tmp/app/out/lambda /app/lambda

# VOLUME ["/var/run/docker.sock"]

EXPOSE ${DASHBOARD_HTTP_PORT}
EXPOSE ${WEBHOOK_HTTP_PORT}

CMD ["./lambda", "webhook"]