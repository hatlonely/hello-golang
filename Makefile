NAME=hello-golang
VERSION=0.0.1
IMAGE_DEV=$(NAME)-dev

image-dev:
	docker build \
		--build-arg HTTP_PROXY=$(HTTP_PROXY) \
		--build-arg HTTPS_PROXY=$(HTTPS_PROXY) \
		-t $(IMAGE_DEV):$(VERSION) -f Dockerfile.dev .

dev-env:
	docker run -d --rm \
		--network host \
		-v $(PWD):/app \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(HOME)/.ssh:/root/.ssh \
		-w /app \
		--name $(IMAGE_DEV) \
		$(IMAGE_DEV):$(VERSION) \
		tail -f /dev/null

release:
	cargo build --release

build:
	cargo build
