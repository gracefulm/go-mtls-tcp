.PHONY: init
init:
	@echo "> init git commit template..."
	git config --local commit.template ./.commit_template

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
