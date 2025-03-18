APP_NAME = hugo-frontmatter-toolbox

# Explicitly set 'all' as the default target
.DEFAULT_GOAL := all

.PHONY: build fmt lint tidy all install readme clean test cover audit

build:
	go build -o $(APP_NAME)

fmt:
	go fmt ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

install: build
	install -m 0755 $(APP_NAME) /usr/local/bin/$(APP_NAME)

# Complete 'all' target
all: tidy fmt lint build readme

audit:
	go list -m -u all
	govulncheck ./...

test:
	go test -v ./internal/... ./pkg/... ./cmd/... ./...

cover:
	go test -coverprofile=coverage.out ./internal/... ./pkg/... ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: file://$(PWD)/coverage.html"

readme: build
	@echo "📝 Generating README.md from code..."
	@mkdir -p tools
	@go run tools/readme-generator.go
	@echo "📄 README.md updated successfully"

clean:
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html