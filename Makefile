.PHONY: build docker docker-run docker-tag docker-push clean

# Configuration
APP_NAME := acs-proxy
BIN := $(APP_NAME)
REGISTRY := harbor.praction.in/i9
IMAGE := $(REGISTRY)/$(APP_NAME)
TAG ?= latest

# Build binary for Linux (used in Dockerfile stage)
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN) ./cmd/main.go

# Build Docker image with specified tag (default: latest)
docker:
	docker buildx build --platform linux/amd64 -t $(IMAGE):$(TAG) .

# Run Docker container locally for testing
docker-run:
	docker run --rm -p 3000:3000 -p 9001:9001 --name $(APP_NAME) $(IMAGE):$(TAG)

# Tag Docker image with a specific version
docker-tag:
ifndef VERSION
	$(error VERSION is not set. Use like: make docker-tag VERSION=1.0.0)
endif
	docker tag $(IMAGE):$(TAG) $(IMAGE):$(VERSION)

# Push both latest and versioned tags to registry
docker-push:
	docker push $(IMAGE):$(TAG)
ifdef VERSION
	docker push $(IMAGE):$(VERSION)
endif

# Remove local binary
clean:
	rm -f $(BIN)
