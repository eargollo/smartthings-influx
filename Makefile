VERSION=$(shell git describe --tags)

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-extldflags=-static" .

.PHONY: test
test:
	go test -cover ./...

.PHONY: run
run:
	UID=$(id -u) GID=$(id -g) docker-compose -f docker-compose-dev.yml up --build

.PHONY: clean
clean:
	docker-compose -f docker-compose-dev.yml rm
	rm -rf data

.PHONY: lint
lint: lint-code lint-security lint-vulnerability

.PHONY: lint-code 
lint-code:
	golangci-lint run

.PHONY: lint-security
lint-security:
	gosec ./...
	
.PHONY: lint-vulnerability
lint-vulnerability:
	govulncheck ./...

.PHONY: outdated
outdated:
	@go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep -F '[' || true

.PHONY: release
release:
	# Requires containerd for pulling and storing images (Settings/General in Docker Desktop)
	@echo publishing '$(VERSION)'
	docker build --platform linux/amd64,linux/arm64 --push -t eargollo/smartthings-influx:$(VERSION) .
	@echo publishing latest
	docker build --platform linux/amd64,linux/arm64 --push -t eargollo/smartthings-influx .

.PHONY: docker
docker:
	docker build . -t smartthings-influx 

.PHONY: cover
cover:
	go test -coverprofile=coverage.out -covermode=count  ./...
	@go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+'
