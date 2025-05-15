APP_NAME=xroll-backend
PORT=8000

.PHONY: build run logs stop clean-cache build-clean deep-clean

# 通常のDockerビルド
build:
	docker build -t $(APP_NAME) .

# Dockerコンテナ起動
run:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true
	docker run -d \
		--name $(APP_NAME) \
		--env-file .env \
		-p $(PORT):$(PORT) \
		$(APP_NAME)

# ログ表示
logs:
	docker logs -f $(APP_NAME)

# 停止＆削除
stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# Dockerビルドキャッシュだけ削除
clean-cache:
	docker builder prune -f

# clean-cache + build
build-clean: clean-cache build

# 完全クリーン（未使用イメージ・ボリューム・ネットワーク・コンテナ）
deep-clean:
	docker system prune -a --volumes -f

