BINARY_NAME=shortlink

deps:
	@go mod tidy

build:
	@go build -o ${BINARY_NAME} main.go

run: deps build
	@./${BINARY_NAME}

clean:
	@go clean -i

test: deps
	@go test ./...

testv: deps
	@go test -v ./...

up:
	@docker compose build
	@docker compose up -d

down:
	@docker compose down --remove-orphans

