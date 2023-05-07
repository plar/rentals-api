SVC_NAME := rentals-api

.PHONY: clean build test coverage

test:
	go test ./...

build:
	go build -o $(SVC_NAME)

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

clean:
	go clean -modcache
	rm -f $(SVC_NAME)
	rm -f coverage.out
	rm -f coverage.html

doc-build:
	docker-compose build

# Run the service in a Docker container
doc-up:
	docker-compose up -d

doc-start:
	docker-compose start

doc-stop:
	docker-compose stop

doc-down:
	docker-compose down -v --remove-orphans	

doc-logs-follow:
	docker-compose logs -f

http-test:
	http :8080/rentals/1
	http :8080/rentals
	http ':8080/rentals?offset=5&limit=3'
	http ':8080/rentals?price_min=9000&price_max=75000&sort=price&near=33.64,-117.93'