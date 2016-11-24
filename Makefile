PWD := $(shell pwd)
PKG := github.com/ckeyer/sloth
APP := sloth

DEV_IMAGE := ckeyer/dev
DEV_UI_IMAGE := ckeyer/dev:node

VERSION := $(shell cat VERSION.txt)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

LD_FLAGS := -X $(PKG)/version.version=$(VERSION) -X $(PKG)/version.gitCommit=$(GIT_COMMIT) -w

NET := $(shell docker network inspect cknet > /dev/zero && echo "--net cknet --ip 172.16.1.8" || echo "")
UI_NET := $(shell docker network inspect cknet > /dev/zero && echo "--net cknet --ip 172.16.1.9" || echo "")

local:
	go build -a -ldflags="$(LD_FLAGS)" -o bundles/$(APP) cli/main.go

test:
	go test -ldflags="$(LD_FLAGS)" $$(go list ./... |grep -v "vendor")

dev:
	docker run --rm -it \
	 --name $(APP)-dev \
	 $(NET) \
	 -v $(PWD):/opt/gopath/src/$(PKG) \
	 -w /opt/gopath/src/$(PKG) \
	 $(DEV_IMAGE) bash

dev-ui:
	docker run --rm -it \
	 --name $(APP)-ui-dev \
	 -p 8080:8080 \
	 -v $(PWD)/ui:/opt/$(APP) \
	 -w /opt/$(APP) \
	 $(DEV_UI_IMAGE) bash