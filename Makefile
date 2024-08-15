# Project variables
BINARY_NAME=telegram
BUILD_DIR=./bin
TMP_DIR=./tmp
CSS_SOURCE=views/css/main.css
CSS_TARGET=public/styles.css
TEMPL_DIR=.

# Go build variables
GO_BUILD_CMD=go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# TailwindCSS build command
TAILWIND_BUILD_CMD=tailwindcss -i $(CSS_SOURCE) -o $(CSS_TARGET)

# TailwindCSS build command in development mode
TAILWIND_BUILD_DEV_CMD=tailwindcss -i $(CSS_SOURCE) -o $(CSS_TARGET) --watch

# Templ commands
TEMPL_FMT_CMD=templ fmt $(TEMPL_DIR)
TEMPL_GEN_CMD=templ generate

.PHONY: all build css clean run templ

all: build css templ

build:
	@echo "Building Go binary..."
	@$(GO_BUILD_CMD)
	@echo "Build complete."

css:
	@echo "Generating CSS with TailwindCSS..."
	@$(TAILWIND_BUILD_CMD)
	@echo "CSS generation complete."

css-dev:
	@echo "Generating CSS with TailwindCSS in development mode..."
	@$(TAILWIND_BUILD_DEV_CMD)
	@echo "CSS generation complete."

templ:
	@echo "Formatting and generating templates..."
	@$(TEMPL_FMT_CMD)
	@$(TEMPL_GEN_CMD)
	@echo "Template processing complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME)
	@rm -rf $(TMP_DIR)/*
	@echo "Cleanup complete."
