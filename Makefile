COMMIT = $$(git describe --tags --always)
BUILD_FLAGS = -ldflags "-X main.GitCommit=\"$(COMMIT)\""

clean:
	rm -f mackerel-plugin-mogilefs
	rm -fr pkg/

deps:
	go get -d -t ./...
	go get golang.org/x/tools/cmd/cover
	go get golang.org/x/tools/cmd/vet

test: deps
	go test -v ./...
	go test -race ./...
	go vet .

cover: deps
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

build: deps
	go build $(BUILD_FLAGS)

install: deps
	go install $(BUILD_FLAGS)

package: deps
	@sh -c "'$(CURDIR)/scripts/package.sh'"
