COMMIT = $$(git describe --always)
BUILD_FLAGS = -ldflags "-X main.GitCommit=\"$(COMMIT)\""

deps:
	go get -d -t ./...
	go get golang.org/x/tools/cmd/cover
	go get golang.org/x/tools/cmd/vet

test: deps
	go test -v ./...
	go test -race ./...
	go vet .
