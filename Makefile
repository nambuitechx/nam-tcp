.PHONY: help build build-server build-client run-server run-client clean

IMAGE_TAG ?= latest
SERVER_IMAGE ?= nam-tcp-server
CLIENT_IMAGE ?= nam-tcp-client

DOCKER_BUILD_OPTS ?=

help:
	@echo "Targets:"
	@echo "  build          Build server and client images"
	@echo "  build-server   Build $(SERVER_IMAGE):$(IMAGE_TAG)"
	@echo "  build-client   Build $(CLIENT_IMAGE):$(IMAGE_TAG)"
	@echo "  run-server     Run proxy server (HTTP :8000, TCP :8888)"
	@echo "  run-client     Run client with ARGS"
	@echo "  clean          Remove built images"

build: build-server build-client

build-server:
	docker build $(DOCKER_BUILD_OPTS) -f Dockerfile.server -t $(SERVER_IMAGE):$(IMAGE_TAG) .

build-client:
	docker build $(DOCKER_BUILD_OPTS) -f Dockerfile.client -t $(CLIENT_IMAGE):$(IMAGE_TAG) .

run-server: build-server
	docker run --rm -it \
		-p 8000:8000 \
		-p 8888:8888 \
		-v nam-tcp-data:/app \
		$(SERVER_IMAGE):$(IMAGE_TAG)

run-client: build-client
	docker run --rm -it \
		--network host \
		-e NAM_TCP_PROXY \
		-e NAM_TCP_TOKEN \
		-e NAM_TCP_LOCAL \
		$(CLIENT_IMAGE):$(IMAGE_TAG) $(ARGS)

clean:
	-docker rmi $(SERVER_IMAGE):$(IMAGE_TAG) $(CLIENT_IMAGE):$(IMAGE_TAG)
