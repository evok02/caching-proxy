.DEFAULT_GOAL := run
.PHONY: run fmt vet build build-proxy-container run-proxy-container stop-proxy-container

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build-proxy: vet
	@go build -o ./bin/proxy ./cmd/main.go

run-proxy: build-proxy 
	@./bin/proxy

build-proxy-container: 
	@docker build -t proxy:latest .

create-brige:
	@docker network create proxy-net

run-proxy-container: build-proxy-container
	@docker run -d --name cache --network proxy-net redis:alpine
	@docker run -d --name proxy1 --network proxy-net proxy:latest
	@docker run -d --name proxy2 --network proxy-net proxy:latest
	@docker run -d --name proxy3 --network proxy-net proxy:latest
	@docker run -d --network proxy-net --name reverse_proxy -p 8080:80\
		-v $(PWD)/caddy/Caddyfile:/etc/caddy/Caddyfile caddy:alpine

stop-proxy-container:
	@docker kill reverse_proxy
	@docker kill proxy1
	@docker kill proxy2
	@docker kill proxy3
	@docker kill cache
