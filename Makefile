
build:
	CGO_ENABLED=0 go build -o _out/dops cmd/dops/main.go

run: build
	./_out/dops

clean:
	go clean
	rm ./_out/*

package:
	go clean
	rm -rf ./_out/*

	GOARCH=386   GOOS=linux   CGO_ENABLED=0 go build -o _out/dops_linux-386-static                     cmd/dops/main.go  # Linux - 32 bit
	GOARCH=amd64 GOOS=linux   CGO_ENABLED=0 go build -o _out/dops_linux-amd64-static                   cmd/dops/main.go  # Linux - 64 bit
	GOARCH=arm64 GOOS=linux   CGO_ENABLED=0 go build -o _out/dops_linux-arm64-static                   cmd/dops/main.go  # Linux - ARM
	GOARCH=386   GOOS=linux                 go build -o _out/dops_linux-386                            cmd/dops/main.go  # Linux - 32 bit
	GOARCH=amd64 GOOS=linux                 go build -o _out/dops_linux-amd64                          cmd/dops/main.go  # Linux - 64 bit
	GOARCH=arm64 GOOS=linux                 go build -o _out/dops_linux-arm64                          cmd/dops/main.go  # Linux - ARM
	GOARCH=amd64 GOOS=darwin                go build -o _out/dops_macos-amd64                          cmd/dops/main.go  # macOS - 32 bit
	GOARCH=amd64 GOOS=darwin                go build -o _out/dops_macos-amd64                          cmd/dops/main.go  # macOS - 64 bit
	GOARCH=amd64 GOOS=openbsd               go build -o _out/dops_openbsd-amd64                        cmd/dops/main.go  # OpenBSD - 64 bit
	GOARCH=arm64 GOOS=openbsd               go build -o _out/dops_openbsd-arm64                        cmd/dops/main.go  # OpenBSD - ARM
	GOARCH=amd64 GOOS=freebsd               go build -o _out/dops_freebsd-amd64                        cmd/dops/main.go  # FreeBSD - 64 bit
	GOARCH=arm64 GOOS=freebsd               go build -o _out/dops_freebsd-arm64                        cmd/dops/main.go  # FreeBSD - ARM

	_data/package-data/aur-git.sh
	_data/package-data/aur-bin.sh
	_data/package-data/homebrew.sh

	echo ""
	echo "[TODO]: call `make package-push-aur-git`"
	echo "[TODO]: call `make package-push-aur-bin`"
	echo "[TODO]: call `make package-push-homebrew`"
	echo "[TODO]: create github release"
	echo ""

package-push-aur-git:
	cd _out/dops-git && git push

package-push-aur-bin:
	cd _out/dops-bin && git push

package-push-homebrew:
	cd _out/homebrew-tap && git push