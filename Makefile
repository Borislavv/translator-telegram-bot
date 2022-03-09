#!make
include .env
export $(shell sed 's/=.*//' .env)

# Main app container
CONTAINER = @docker exec -it telegram_bot_app_container

# Make commands:
# 	Used vars. defined into Dockerfile as exported values.
migrate.up:
	$(CONTAINER) migrate -database "mysql://${DB_USER}:${DB_PASS}@tcp(db:3306)/${DB_NAME}" -path migrations up

migrate.down:
	$(CONTAINER) migrate -database "mysql://${DB_USER}:${DB_PASS}@tcp(db:3306)/${DB_NAME}" -path migrations down