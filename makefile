format:
	gofmt -w .

build: format
	go build -o ./bin/tfCreator .
