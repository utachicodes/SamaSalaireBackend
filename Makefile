.PHONY: build run seed test tidy fmt vet lint docker clean

BIN := bin/samasalaire
SEED_BIN := bin/seed

build:
	go build -o $(BIN) ./cmd/server

run:
	go run ./cmd/server

seed:
	go run ./cmd/seed

test:
	go test ./... -race -count=1

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

vet:
	go vet ./...

lint:
	golangci-lint run ./...

docker:
	docker build -t samasalaire-backend:latest .

clean:
	rm -rf bin coverage.out coverage.html
