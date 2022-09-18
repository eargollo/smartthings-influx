.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-extldflags=-static" .

.PHONY: run
run:
	docker-compose up --build

.PHONY: clean
clean:
	docker-compose rm
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