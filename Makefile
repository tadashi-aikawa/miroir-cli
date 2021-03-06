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

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "[ERROR] Required: $* !!"; \
		echo "[ERROR] Please set --> $*"; \
		exit 1; \
	fi

branch_version := $(shell git rev-parse --abbrev-ref HEAD)

#------

init:
	dep ensure

package-windows: ## Create miroir-cli for Windows.
	@mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -o dist/miroir.exe

package-linux: ## Create miroir-cli for Linux.
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o dist/miroir

clean-package:
	rm -rf dist

release: clean-package ## Build and upload packages, regarding branch name as version
	@echo '1. Update versions'
	@sed -i -r 's/const version = ".+"/const version = "$(branch_version)"/g' args.go

	@echo '2. Packaging'
	make package-linux

	@echo '2. Packaging'
	@echo 'Linux...'
	make package-linux
	tar zfc dist/miroir-$(branch_version)-x86_64-linux.tar.gz dist/miroir
	rm -rf dist/miroir
	@echo 'Windows...'
	make package-windows
	7z a dist/miroir-$(branch_version)-x86_64-windows.zip dist/miroir.exe
	rm -rf dist/miroir.exe

	@echo '3. Staging and commit'
	git add args.go
	git commit -m ':package: Version $(branch_version)'

	@echo '4. Push'
	git push

	@echo 'Success All!!'
	@echo 'Create a pull request and merge to master!!'
	@echo 'https://github.com/tadashi-aikawa/miroir-cli/compare/$(branch_version)?expand=1'
	@echo '..And deploy package!!'

deploy: guard-version ## Deploy packages [Required: $version. ex. 0.5.1]
	ghr v$(version) dist/
