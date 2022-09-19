build:
	go build -o sp cmd/sp/*.go

lint:
	go fmt ./...

clean:
	rm ./sp
