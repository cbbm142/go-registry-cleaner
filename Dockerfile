FROM golang@sha256:55f55d3232f63391e0797acaf145ade8f6fca2ff36795dd5ae446de360724dec
# golang:1.17.1-alpine3.14

RUN addgroup -S golang && adduser -S golang -G golang
USER golang

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./clean-registry

CMD [ "./clean-registry" ]
