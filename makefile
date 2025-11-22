# Makefile

load-secrets:
	@echo "Generating .env file..."
	@echo "GIT_PAT=$$( SOPS_AGE_KEY_FILE=$(HOME)/.config/sops/age/keys.txt sops -d secrets/git_pat.enc.txt )" > .env
	@echo "NOTION_TOKEN=$$( SOPS_AGE_KEY_FILE=$(HOME)/.config/sops/age/keys.txt sops -d secrets/notion_token.enc.txt )" >> .env
	@echo ".env generated."

decode-secrets:
	SOPS_AGE_KEY_FILE=$(HOME)/.config/sops/age/keys.txt sops -d secrets/git_pat.enc.txt > secrets/git_pat.txt
	SOPS_AGE_KEY_FILE=$(HOME)/.config/sops/age/keys.txt sops -d secrets/notion_token.enc.txt > secrets/notion_token.txt

.PHONY: test
test: load-secrets
	@echo "Running tests..."
	@export $$(cat .env | xargs) && go test -v ./...

# Run tests with coverage
.PHONY: coverage
coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Clean up test artifacts (optional)
.PHONY: clean
clean:
	@echo "Cleaning up..."
	go clean -testcache

.PHONY: run
run: load-secrets
	@echo "Running application..."
	@export $$(cat .env | xargs) && go run main.go

.PHONY: compose-up
compose-up: decode-secrets
	docker compose up --build
