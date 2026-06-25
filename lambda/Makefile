# ─────────────────────────────────────────────────────────────
# Lambda Instrumentation Platform — Backend Makefile
# ─────────────────────────────────────────────────────────────

BINARY_NAME = bootstrap
BUILD_DIR   = .build
ZIP_NAME    = function.zip

# Go build flags
GOOS   = linux
GOARCH = arm64
LDFLAGS = -s -w

.PHONY: all clean deps build build-obfuscated test package deploy

# ── Default target ──
all: clean deps build package

# ── Install dependencies ──
deps:
	go mod tidy
	go mod download

# ── Run tests ──
test:
	go test ./... -v -count=1

# ── Standard build (readable, for development) ──
build:
	@echo "Building Lambda binary ($(GOOS)/$(GOARCH))..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		./cmd/lambda/
	@echo "✓ Built $(BUILD_DIR)/$(BINARY_NAME)"

# ── Obfuscated build (for production / distribution) ──
build-obfuscated:
	@echo "Building obfuscated Lambda binary..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) garble -literals -tiny -seed=random build \
		-ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		./cmd/lambda/
	@echo "✓ Built obfuscated $(BUILD_DIR)/$(BINARY_NAME)"

# ── Verify obfuscation ──
verify-obfuscation:
	@echo "Checking for leaked strings..."
	@! strings $(BUILD_DIR)/$(BINARY_NAME) | grep -i "newrelic" || echo "⚠ WARNING: Found 'newrelic' strings!"
	@! strings $(BUILD_DIR)/$(BINARY_NAME) | grep -i "handler" || echo "⚠ WARNING: Found 'handler' strings!"
	@echo "Checking for symbols..."
	@! go tool nm $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null | grep -i "handler" || echo "⚠ WARNING: Found symbols!"
	@echo "✓ Obfuscation check complete"

# ── Package for Lambda ──
package: build
	@echo "Packaging Lambda zip..."
	@cd $(BUILD_DIR) && zip -j ../$(ZIP_NAME) $(BINARY_NAME)
	@echo "✓ Created $(ZIP_NAME) ($$(du -h $(ZIP_NAME) | cut -f1))"

# ── Package obfuscated ──
package-obfuscated: build-obfuscated
	@echo "Packaging obfuscated Lambda zip..."
	@cd $(BUILD_DIR) && zip -j ../$(ZIP_NAME) $(BINARY_NAME)
	@echo "✓ Created $(ZIP_NAME) ($$(du -h $(ZIP_NAME) | cut -f1))"

# ── Deploy with SAM ──
deploy: package
	sam deploy \
		--template-file template.yaml \
		--stack-name lambda-instrumentation-platform \
		--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND \
		--resolve-s3

# ── Local testing with SAM ──
local-invoke:
	sam local invoke OrchestratorFunction \
		--template template.yaml \
		--event events/health.json

# ── Clean ──
clean:
	@rm -rf $(BUILD_DIR) $(ZIP_NAME)
	@echo "✓ Cleaned"
