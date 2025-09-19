.PHONY: build test run docker-build docker-run

build:
	go build -o bin/admira ./cmd/api

test:
	go test -v ./tests/...

run: build
	./bin/admira

docker-build:
	docker build -t admira-backend .

docker-run: docker-build
	docker run -p 8080:8080 --env-file .env admira-backend

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

compose-logs:
	docker-compose logs -f