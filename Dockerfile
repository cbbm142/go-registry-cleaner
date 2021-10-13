FROM golang@sha256:55f55d3232f63391e0797acaf145ade8f6fca2ff36795dd5ae446de360724dec as builder
# golang:1.17.1-alpine3.14

# Disable CGO so binary works after copying to 2nd stage
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./clean-registry

FROM alpine@sha256:e15947432b813e8ffa90165da919953e2ce850bef511a0ad1287d7cb86de84b5
# alpine:3.13.6

RUN addgroup -S golang && adduser -S golang -G golang
USER golang

WORKDIR /app
COPY --from=builder /build/clean-registry /app/clean-registry

CMD [ "/app/clean-registry" ]
