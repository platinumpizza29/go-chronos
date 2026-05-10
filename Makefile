IMAGE_NAME := go-chronos
PORT := 8080

.PHONY: build run stop

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d \
		--name $(IMAGE_NAME) \
		-p $(PORT):$(PORT) \
		--env-file .env \
		$(IMAGE_NAME)

stop:
	docker stop $(IMAGE_NAME) || true
	docker rm $(IMAGE_NAME) || true

logs:
	docker logs -f $(IMAGE_NAME)

restart: stop build run
