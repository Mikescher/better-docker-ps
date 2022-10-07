
build:
	CGO_ENABLED=0 go build -o _out/dops cmd/dops/main.go

run: build
	./_out/dops

clean:
	go clean
	rm ./_out/*