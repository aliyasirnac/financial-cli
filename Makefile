build:
	go build -o bin/nachboard ./cmd/cli

run: build
	@./bin/nachboard

setup:
	@zsh setup.sh