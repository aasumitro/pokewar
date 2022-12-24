build: swag
	@echo "Build Binary"
	mkdir ./build && mkdir ./build/db
	cp ./db/fresh.db ./build/db && cp .example.env ./build/.env
	go mod tidy -compat=1.19
	go build -o ./build/pokewar ./cmd/web/main.go

swag: tests
	@echo "Re-generate Swagger File (API Spec docs)"
	swag init --parseDependency --parseInternal \
		--parseDepth 4 -g ./cmd/web/main.go

tests: lint
	@echo "Run tests"
	gotestsum --format pkgname-and-test-fails \
		--hide-summary=skipped \
		-- -coverprofile=cover.out ./...
	rm cover.out

lint:
	@echo "Applying linter"
	golangci-lint cache clean
	golangci-lint run -c .golangci.yaml ./...

run:
	@echo "Run App"
	go mod tidy -compat=1.19
	go run ./cmd/web/main.go

prepare:
	go mod install
	cp .example.env .env
