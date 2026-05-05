.PHONY: init
init:
	@echo "> init git commit template..."
	git config --local commit.template ./.commit_template

.PHONY: certs
certs:
	@./scripts/gen-certs.sh

.PHONY: fmt
fmt:
	@echo "> format source code..."
	go mod tidy
	go fmt ./...
	find . -print | grep --regex '.*\.go' | xargs go tool goimports -w -local "github.com/gracefulm/go-template-project"

.PHONY: lint
lint:
	@echo "> linting source code..."
	go tool golangci-lint run

.PHONY: sec
sec:
	@echo "> security check for source code..."
	go tool govulncheck ./...

.PHONY: test
test: lint fmt
	@echo "> testing go files..."
	go test -race -cover ./...

.PHONY: site-deps
site-deps:
	@echo "> resolving hugo modules..."
	cd site && hugo mod get github.com/imfing/hextra@latest && hugo mod tidy

.PHONY: site-serve
site-serve: site-deps
	@echo "> hugo server (http://localhost:1313)..."
	cd site && hugo server -D

.PHONY: site-build
site-build: site-deps
	@echo "> building static site into site/public..."
	cd site && hugo --minify
