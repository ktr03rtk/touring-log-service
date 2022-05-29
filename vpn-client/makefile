SHELL:=/bin/sh
TZ?=Asia/Tokyo
HOST=master1

default: build

build: fmtcheck
	@docker compose build

run: fmtcheck
	docker compose up --build

test_up: test_down
	@COMPOSE_FILE="docker-compose.yml:docker-compose-dev.yml" \
		docker compose build
	@COMPOSE_FILE="docker-compose.yml:docker-compose-dev.yml" \
		docker compose run --rm vpn_client

test_down:
	@COMPOSE_FILE="docker-compose.yml:docker-compose-dev.yml" \
		docker compose down

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

fmtcheck:
	sh -c "'$(CURDIR)/../scripts/gofmtcheck.sh'"

sync:
	sh -c "'$(CURDIR)/../scripts/rsync.sh' $(HOST)"

.PHONY: default build run test_up test_down fmt fmtcheck sync
