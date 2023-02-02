# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s)

# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile

install-deps: gotestsum
deps: $(GOTESTSUM)
deps:
	@ echo "Required Tools Are Available"

build: swag
	@ echo "Build Binary"
	@ mkdir ./build && mkdir ./build/db
	@ cp ./db/local-data.db ./build/db/local-data.db && cp .example.env ./build/.env
	@ go mod tidy -compat=1.19
	@ go build -o ./build/pokewar ./cmd/web/main.go
	# GOOS=windows GOARCH=amd64 go build -o ./build/pokewar.exe ./cmd/web/main.go
	@ echo "generate binary done"

swag: tests
	@ echo "Re-generate Swagger File (API Spec docs)"
	@ swag init --parseDependency --parseInternal \
		--parseDepth 4 -g ./cmd/web/main.go
	@ echo "generate swagger file done"

tests: $(GOTESTSUM) lint
	@ echo "Run tests"
	@ gotestsum --format pkgname-and-test-fails \
		--hide-summary=skipped \
		-- -coverprofile=cover.out ./...
	@ rm cover.out

lint: $(GOLANGCI)
	@ echo "Applying linter"
	@ golangci-lint cache clean
	@ golangci-lint run -c .golangci.yaml ./...

run:
	@echo "Run App"
	go mod tidy -compat=1.19
	go run ./cmd/web/main.go

prepare:
	go mod install
	cp .example.env .env
