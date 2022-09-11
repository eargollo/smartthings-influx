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
lint:
	golangci-lint run