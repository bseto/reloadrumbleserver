FROM golang:1.11

WORKDIR /reload_rumble_server/

RUN go get \
	gopkg.in/yaml.v2 \
	github.com/gorilla/websocket