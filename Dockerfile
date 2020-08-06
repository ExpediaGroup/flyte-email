# Build image
FROM golang:1.14.6 AS build-env

ENV APP=flyte-email
ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/ExpediaGroup/$APP/

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./
RUN go test
RUN CGO_ENABLED=0 go build

# Run image
FROM alpine:3.10.2
RUN apk add --no-cache ca-certificates
ENV APP=flyte-email
COPY --from=build-env /go/src/github.com/ExpediaGroup/$APP/$APP /app/$APP
ENTRYPOINT "/app/$APP"
