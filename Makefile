BINARY_NAME=shortlink

build:
	@go build -o ${BINARY_NAME} main.go

run: build
	@./${BINARY_NAME}

clean:
	@go clean -i

test:
	@go test ./...

testv:
	@go test -v ./...

up:
	@docker compose build
	@docker compose up -d

down:
	@docker compose down --remove-orphans

