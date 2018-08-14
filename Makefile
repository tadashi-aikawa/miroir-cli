MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
ARGS :=
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

.PHONY: $(shell egrep -oh ^[a-zA-Z0-9][a-zA-Z0-9_-]+: $(MAKEFILE_LIST) | sed 's/://')

help: ## Print this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9][a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

version := $(shell git rev-parse --abbrev-ref HEAD)

#------

test: ## Test
	@echo Start $@
	@echo End $@

build: ## Release build
	@echo Start $@
	@cargo build --release
	@echo End $@

release: test ## Release
	@echo 'Start $@'

	@echo '1. Update a version'
	@sed -i 's/^version = ".*"$$/version = "$(version)"/g' Cargo.toml

	@echo '2. Release build'
	@make build

	@echo '3. Staging and commit'
	git add Cargo.toml
	git commit -m ':package: Version $(version)'

	@echo '4. Tags'
	git tag v$(version) -m v$(version)

	@echo '5. Push'
	git push

	@echo 'Success All!!'
	@echo 'Create a pull request and merge to master!!'
	@echo 'https://github.com/tadashi-aikawa/miroir-cli/compare/$(version)?expand=1'

	@echo 'End $@'

deploy: ## Deploy
	@echo Start $@
	@echo TODO....
	@echo End $@

