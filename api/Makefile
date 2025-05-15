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

# Docker logs
logs:
	docker logs -f $(APP_NAME)

# Docker stop & remove
stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# 🔥 Docker build cache cleanup
clean-cache:
	docker builder prune -f

# 🔄 Clean cache and then build
build-clean: clean-cache build
