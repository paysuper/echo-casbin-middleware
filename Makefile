ifndef VERBOSE
.SILENT:
endif

override CURRENT_DIR = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
override DOCKER_MOUNT_SUFFIX ?= consistent
override GO_PATH = $(shell echo "$(GOPATH)" | cut -d';' -f1 | sed -e "s~^\(.\):~/\1~g" -e "s~\\\~/~g" -e "s~^/go:~~g" )

TAG ?= unknown
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
DIND_UID ?= 0
DING_GUID ?= 0

ifeq ($(GO111MODULE),auto)
override GO111MODULE = on
endif

ifeq ($(OS),Windows_NT)
override ROOT_DIR = $(shell echo $(CURRENT_DIR) | sed -e "s:^/./:\U&:g")
else
override ROOT_DIR = $(CURRENT_DIR)
endif

define go_docker
	if [ "${GO_PATH}" != "" ]; then VOLUME_PKG_MOD="-v /${GO_PATH}/pkg/mod:/${GO_PATH}/pkg/mod:${DOCKER_MOUNT_SUFFIX}"; fi ;\
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	docker run --rm \
		-v /${ROOT_DIR}:/${ROOT_DIR}:${DOCKER_MOUNT_SUFFIX} \
        $${VOLUME_PKG_MOD} \
		-w /${ROOT_DIR} \
		-e GO111MODULE=on \
		-e GOPATH=/${GO_PATH}:/go \
		$${GO_IMAGE}:$${GO_IMAGE_TAG} \
		sh -c 'GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} TAG=${TAG} $(subst ",,${1}); if [ "${DIND_UID}" != "0" ]; then chown -R ${DIND_UID}:${DIND_GUID} ${ROOT_DIR}; fi'
endef

up: clean ## initialize required tools
	. ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	if [ "${DIND}" != "1" ]; then \
		go get github.com/google/wire/cmd/wire@v0.3.0 && \
		go get -u github.com/golangci/golangci-lint/cmd/golangci-lint ;\
		go get -u github.com/stretchr/testify ;\
    fi;
.PHONY: up

clean: ## remove generated files, tidy vendor dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make clean") ;\
    else \
        go mod tidy ;\
    	rm -rf *.out generated/* vendor bin ;\
    fi;
.PHONY: clean

dev-test: test lint ## test application in dev env with race and lint
.PHONY: dev-test

dind: ## useful for windows
	if [ "${GO_PATH}" != "" ]; then VOLUME_PKG_MOD="-v /${GO_PATH}/pkg/mod:/${GO_PATH}/pkg/mod:${DOCKER_MOUNT_SUFFIX}"; fi ;\
	if [ "${DIND}" = "1" ]; then \
		echo "Already in DIND" ;\
    else \
	    . ${ROOT_DIR}/scripts/common.sh ${ROOT_DIR}/scripts ;\
	    docker rm -f dind-$${PROJECT_NAME} &>/dev/null ;\
	    docker run -it --rm --name dind-$${PROJECT_NAME} --privileged \
            -v //var/run/docker.sock://var/run/docker.sock:${DOCKER_MOUNT_SUFFIX} \
            -v /${ROOT_DIR}:/${ROOT_DIR}:${DOCKER_MOUNT_SUFFIX} \
			$${VOLUME_PKG_MOD} \
            -w /${ROOT_DIR} \
			-e GOPATH=${GO_PATH} \
            $${DIND_IMAGE}:$${DIND_IMAGE_TAG} ;\
    fi;
.PHONY: dind

go-depends: ## view final versions that will be used in a build for all direct and indirect dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-depends") ;\
    else \
        cd $(ROOT_DIR) ;\
        go list -m all ;\
    fi;
.PHONY: go-depends

go-update-all: ## view available minor and patch upgrades for all direct and indirect
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-update-all") ;\
    else \
        cd $(ROOT_DIR) ;\
    	go list -u -m all ;\
    fi;
.PHONY: go-update-all

lint: ## execute linter
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make lint") ;\
    else \
        golangci-lint run --no-config --issues-exit-code=0 --deadline=30m \
          --disable-all --enable=deadcode  --enable=gocyclo --enable=golint --enable=varcheck \
          --enable=structcheck --enable=maligned --enable=errcheck --enable=dupl --enable=ineffassign \
          --enable=interfacer --enable=unconvert --enable=goconst --enable=gosec --enable=megacheck ./... ;\
    fi;
.PHONY: lint

test-with-coverage: ## test application with race and total coverage
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make test-with-coverage") ;\
    else \
		export WD=$(ROOT_DIR) ;\
         CGO_ENABLED=1 \
        go test -v -race -covermode atomic -coverprofile coverage.out ${TEST_ARGS} ./... || exit 1 ;\
        go tool cover -func=coverage.out ;\
    fi;
.PHONY: test-with-coverage

test: ## test application with race
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make test") ;\
    else \
		export WD=$(ROOT_DIR) ;\
        CGO_ENABLED=1 \
        go test -race -v ${TEST_ARGS} ./... ;\
    fi;
.PHONY: test

vendor: ## update vendor dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make vendor") ;\
    else \
        rm -rf $(ROOT_DIR)/vendor ;\
    	go mod vendor ;\
    fi;
.PHONY: vendor

go-download-deps: ## download dependencies
	if [ "${DIND}" = "1" ]; then \
		$(call go_docker,"make go-download-deps") ;\
    else \
    	go get -d ./... ;\
    fi;
.PHONY: go-download-deps

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

.DEFAULT_GOAL := help