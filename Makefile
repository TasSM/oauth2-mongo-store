
.PHONY:  test format

test:
	docker-compose up -d && sleep 5 && go test -race ./... || docker-compose down

format:
	go fmt ./...