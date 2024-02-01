FROM golang:1.18.8-alpine3.15 as BuildStage

WORKDIR /usr/src/app
COPY . .

RUN go mod download && go mod verify
RUN go build -o /dns main.go

FROM alpine:latest
WORKDIR /
COPY --from=BuildStage /dns /dns
COPY .env /

ENTRYPOINT ["/dns"]
