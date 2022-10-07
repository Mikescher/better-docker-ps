
build:
	go build -o _out/dops cmd/dops/main.go

run:
	go build -o _out/dops cmd/dops/main.go
	./_out/dops

clean:
	go clean
	rm ./_out/*