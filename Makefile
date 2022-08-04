
.PHONY:  test format

test:
	docker-compose up -d && sleep 1 && go test -race ./... ; sleep 1 ; docker-compose down

format:
	go fmt ./...