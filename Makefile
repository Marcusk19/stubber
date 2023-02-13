BINARY_NAME = main.out

init:
	cd stubber-ui; echo "Installing dependencies for frontend"; npm install;
	cp example.env .env
	cd ..; echo "Finished setting up frontend";
build:
	go build -o ${BINARY_NAME} main.go

run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
	rm -r pg_data

deps:
	go get github.com/joho/godotenv
	go get -u github.com/gorilla/muxs
	go get github.com/lib/pq

test:
	go test ./...
dev-up:
	docker-compose up -d
	@echo "frontend is up at http://localhost"
dev-rebuild:
	docker-compose down
	docker-compose build
	docker-compose up -d
dev-down:
	docker-compose down