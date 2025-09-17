.PHONY: run build clean test

run:
	go run cmd/main.go

build:
	go build -o bin/nit cmd/main.go

clean:
	rm -f bin/nit data.db

test:
	curl -X POST http://localhost:8080/users \
		-H "Content-Type: application/json" \
		-d '{"name": "John Doe", "email": "john@example.com"}'

	curl -X GET http://localhost:8080/users

install:
	go mod download

dev: install run
