.DEFAULT_GOAL := help
PROJ_DIR := $(shell pwd)
SHELL = /bin/bash
COMMIT = $(shell git rev-parse --short=7 HEAD)$(shell [[ $$(git status --porcelain) = "" ]] || echo -dirty)
USER=${whoami}

# adjust this if the api is incremented
PATH_WITH_GOPATH = ${PATH}:$(shell go env GOPATH)/bin

# The Following variables drive features in nc-common-tools and this repo
GO_VERSION=go1.21.3
GO_TEST_SHUFFLE := on
GO_TEST_RACE := true

# The following block of code should be after your variables and before any targets
# and it will give a message if the repo clone has not executed git submodule update --init
# it will also setup the variables and help.mk file
export COMMON_TOOLS_DIR ?= nc-common-tools
# Check to see if the nc-common-tools is initialized by git.
# NOTE: this should support submodule or full clone in pipelines
ifeq ($(wildcard $(COMMON_TOOLS_DIR)/.git),)
$(error This repo is using git submodule please run "git submodule update --init" in order to include mk files)
endif
include $(COMMON_TOOLS_DIR)/mk/help.mk

##@ Development - repository targets developing

##@@ Input validators
validate-in: ## ensure the user specified an IN=<filename>
ifndef IN
	$(error IN is undefined, please set it with e.g. make target IN=<file_path_and_name>)
endif

validate-env: ## ensure the user specified an ENV=dev
ifndef ENV
	$(error ENV is undefined, please set it with e.g. make target ENV=<string>)
endif


##@@ Build - different build process for the repository

build: build-cmds  ## Build all binaries

build-cmds: build-strgctl ## Build different command lines

build-strgctl: ## Build strgctl
    cd $(PROJ_DIR)/cmd/strgctl &&
	CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o $(PROJ_DIR)/bin

##@ Quality Assurance - Local Quality Assurance target for the repo.
validate-go:  ## Validate the go code in this repo
# TODO dschveninger need to see if we can do this with Megalinter or a Megalinter Plugin
	@[ -z "$$(find . -name "*:*")" ] || (echo error: filenames with colons are not allowed on Windows, please rename; exit 1)

validate-versionconfigs: build-vercheck build-vergen  ## Validate all version configurations in versionConfigs
	$(PROJ_DIR)/tools/validateVersionConfigs.sh $(PROJ_DIR)/tools

code-qa: go-fmt qa-lint go-test validate-common validate-go validate-versionconfigs  ## run the quality checks for the repo

# Required Common Make files for the QA Group
include $(COMMON_TOOLS_DIR)/mk/megalinter.mk # target for qa-lint
# Optional linting targets only, should be after megalinter.mk that has level 2 help title
include $(COMMON_TOOLS_DIR)/mk/megalinterhelp.mk # help targets to run custom override on megalinter targets like qa-lint
# Language specific targets
include $(COMMON_TOOLS_DIR)/mk/go.mk # linting targets for go
include $(COMMON_TOOLS_DIR)/mk/shell.mk # linting targets for shell
include $(COMMON_TOOLS_DIR)/mk/yaml.mk # linting targets for yaml
##@ QA Common Tools - nc-common-tools shared targets
include $(COMMON_TOOLS_DIR)/mk/common-tools.mk # cleaning of tracked and untrack files

.PHONY: code-qa validate-go validate-in validate-env build build-vergen build-vercheck validate-versionconfigs vercheck
