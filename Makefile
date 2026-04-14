.PHONY: run migrate seed migrate-seed dev build tidy lint test search-log

APP_NAME=api-admin-panel
MAIN=./cmd/main.go

## Jalankan server saja
run:
	go run $(MAIN)

## Jalankan + migrate
migrate:
	go run $(MAIN) -migrate

## Jalankan + seed
seed:
	go run $(MAIN) -seed

## Jalankan + migrate + seed (first time setup)
migrate-seed:
	go run $(MAIN) -migrate -seed

## Jalankan fresh (delete all tables) + migrate + seed (first time setup)
fresh-migrate-seed:
	go run $(MAIN) -fresh -migrate -seed

## Build binary
build:
	go build -o bin/$(APP_NAME) $(MAIN)

## Hot-reload (requires: go install github.com/air-verse/air@latest)
dev:
	air

## Tidy modules
tidy:
	go mod tidy

## Lint
lint:
	golangci-lint run ./...

## ── Testing ───────────────────────────────────────────────────────────────────
 
## Jalankan semua unit test
test:
	go test ./tests/unit/... -v
 
## Jalankan unit test dengan coverage report
test-cover:
	go test ./tests/unit/... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
 
## Jalankan integration test (butuh server berjalan)
## Contoh: make test-integration BASE_URL=http://localhost:8080
test-integration:
	TEST_BASE_URL=$(or $(BASE_URL),http://localhost:8080) \
	go test ./tests/integration/... -tags=integration -v
 
## Jalankan usecase tests saja
test-usecase:
	go test ./tests/unit/usecase/... -v
 
## Jalankan handler tests saja
test-handler:
	go test ./tests/unit/handler/... -v
 

## Cari log berdasarkan request_id di semua file log atau tanggal tertentu
## Contoh cari semua file : make search-log RID=xxxx-xxxx-xxxx-xxxx
## Contoh cari per tanggal : make search-log RID=xxxx-xxxx DATE=2026-04-11
search-log:
	@if [ -z "$(RID)" ]; then \
		echo "Usage: make search-log RID=<request-id> [DATE=2026-04-11]"; \
		exit 1; \
	fi; \
	if [ -n "$(DATE)" ]; then \
		TARGET="logs/$(APP_NAME)-$(DATE).log"; \
		if [ -f "$$TARGET" ]; then \
			echo "Searching in $$TARGET ..."; \
			grep -n "$(RID)" "$$TARGET" || echo "Not found."; \
		else \
			echo "File $$TARGET not found."; \
		fi; \
	else \
		echo "Searching in all log files ..."; \
		grep -rn "$(RID)" logs/$(APP_NAME)-*.log 2>/dev/null || echo "Not found in any log file."; \
	fi

## Docker
docker-build:
	docker build -t api-admin-panel .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app