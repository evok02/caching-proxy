.DEFAULT_GOAL := run
.PHONY: run fmt vet build

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build-proxy: vet
	@go build -o ./bin/proxy .

run-proxy: build-proxy 
	@./bin/proxy

build-proxy-container: build-proxy
	docker build -t proxy:latest .

run-proxy-container: build-proxy-container
	docker run -d --name cache --network proxy-net redis:latest
	docker run -d --name proxy1 --network proxy-net proxy:latest
	docker run -d --name proxy2 --network proxy-net proxy:latest
	docker run -d --name proxy3 --network proxy-net proxy:latest
	docker run -d --network proxy-net --name reverse_proxy -p 8080:80\
		-v $(PWD)/caddy/Caddyfile:/etc/caddy/Caddyfile caddy:latest
