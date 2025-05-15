APP_NAME=xroll-backend
PORT=8000

# Docker build
build:
	docker build -t $(APP_NAME) .

# Docker run
run:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true
	docker run -d \
		--name $(APP_NAME) \
		--env-file .env \
		-p $(PORT):$(PORT) \
		$(APP_NAME)

logs:
	docker logs -f $(APP_NAME)

stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

