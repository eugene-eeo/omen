dev:
	go build

install:
	go build -ldflags="-s -w"
	go install
