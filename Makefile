build:
	CGO_ENABLED=0 go build .

install:
	CGO_ENABLED=0 go install 

run-messanger:
	go run . messanger --port 8080

run-reporter:
	go run . accounting --port 8081

run-migrate:
	go run . migrate -p $(ROOT)/migrations