build: fmt vet
	go build -o case_watcher .
.PHONY: build

# Run tests
test: build
	go test ./pkg/... -coverprofile cover.out

# Run go fmt against code
fmt:
	go fmt ./pkg/...

# Run go vet against code
vet:
	go vet -structtag=false ./pkg/...