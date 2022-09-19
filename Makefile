build:
	go build -o sp cmd/sp/main.go

lint:
	go fmt ./...

clean:
	rm ./sp
