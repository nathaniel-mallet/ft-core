DB_NAME := ftcore
DB_PORT := 5432
DB_USER := ftcore-admin
DB_PASSWORD := ftcore-admin
DB_IMAGE := ftcore-db
DB_CONTAINER := ftcore-db
DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: db-docker-build db-up

db-docker-build:
	docker build -t ftcore-db -f Dockerfile.db .

db-up:
	docker run -d \
	--name $(DB_CONTAINER) \
	-p $(DB_PORT):5432 \
	-e POSTGRES_USER=$(DB_USER) \
	-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	-e POSTGRES_DB=$(DB_NAME) \
	$(DB_IMAGE)
