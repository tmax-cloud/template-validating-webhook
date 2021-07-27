REGISTRY      ?= tmaxcloudck
VERSION       ?= 0.0.1

TVW_IMG   = $(REGISTRY)/template-validating-webhook:$(VERSION)

.PHONY: test build push

# Test webhook func
test:
	go test ./pkg/apis

# Build the docker image
build:
	docker build -f build/Dockerfile -t $(TVW_IMG) . 

# Push the docker image 
push:
	docker push $(TVW_IMG)
