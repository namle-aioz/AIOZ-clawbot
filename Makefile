APP_ENV ?= app

.PHONY: run
run:
	APP_ENV=$(APP_ENV) go run ./cmd/http
