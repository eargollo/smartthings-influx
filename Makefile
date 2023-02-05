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
	docker buildx build --platform linux/amd64,linux/arm64 --push -t eargollo/smartthings-influx .
