RUNTIME          ?= podman
IMAGE_ORG        ?= quay.io/jwmatthews
IMAGE_NAME       ?= case_watcher
IMAGE_TAG        ?= $(shell git rev-parse --short HEAD)
WATCHER_IMAGE     ?= $(IMAGE_ORG)/$(IMAGE_NAME):$(IMAGE_TAG)

build: fmt vet
	go build -o case_watcher .

# Run tests
test: build
	go test ./pkg/... -coverprofile cover.out

# Run go fmt against code
fmt:
	go fmt ./pkg/...

# Run go vet against code
vet:
	go vet -structtag=false ./pkg/...

build-image:
	$(RUNTIME) build ${CONTAINER_BUILD_PARAMS} -t $(WATCHER_IMAGE) -f Dockerfile .

push-image:
	$(RUNTIME) push $(WATCHER_IMAGE)

build-push-image: build-image push-image

.PHONY: build build-image push-image build-push-image