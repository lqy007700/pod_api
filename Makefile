
all: proto build docker docker-run

.PHONY: proto
proto:
	sudo docker run --rm -v $(shell pwd):$(shell pwd) -w $(shell pwd) cap1573/cap-v3 --proto_path=. --micro_out=. --go_out=:. ./proto/pod_api/pod_api.proto

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /usr/local/Cellar/go@1.19/1.19.11/bin/go build -o pod_api *.go

.PHONY: docker
docker:
	sudo docker build . -t zxnl/pod_api:latest

docker-run:
	sudo docker run -p 8082:8082 -v /Users/lqy007700/Data/code/go-application/go-paas/pod_api/micro.log:/micro.log zxnl/pod_api