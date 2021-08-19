.PHONY: all

all: clean build

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

compile:
	export IMAGE_TAG="$(git describe)"
	go build -ldflags "-X main.version=$(IMAGE_TAG)" ./cmd/firebase-ctl

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	@for p in $(ALL_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

build: compile vet lint

clean:
	go clean

test:
	go test ./...

ci: clean build test

generate-mocks:
	go get github.com/vektra/mockery/v2@v2.8.0
	mockery --all -r --inpackage 