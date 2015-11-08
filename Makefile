# Makefile
# vim:ft=make

PROJECT := $(shell basename $(PWD))
SOURCES := $(shell find . -path './vendor' -prune -o -type f -name '*.go' -print)
PACKAGES=$(shell go list ./... | grep -v /vendor/)

GIT_COMMIT=`git rev-parse HEAD`
GIT_USER=`git config --get user.name`
GO_PROJECTS="/go/src/github.com/$(GIT_USER)"
BUILD_PATH=./_dist
GO_VERSION:=$(shell go version)
# ldflags does't support spaces in variables
CLEAN_GO_VERSION=$(shell echo "${GO_VERSION}" | sed -e 's/[^a-zA-Z0-9]/_/g')

# docker-compose based container name
CONTAINER="$(GIT_USER)/$(PROJECT)"

BINARY=${PROJECT}
# alternative : git describe --always --tags
# FIXME unterminating quote
# VERSION:=`cat env.go | grep "Version =" | cut -d"\"" -f 2`
VERSION:=0.2.0
BUILD_TIME=`date +%FT%T%z`

LDFLAGS="-X github.com/hackliff/sentinel.BuildTime=${BUILD_TIME} -X github.com/hackliff/sentinel.GoVersion=${CLEAN_GO_VERSION} -X github.com/hackliff/sentinel.GitCommit=${GIT_COMMIT}"

all: $(BINARY)

container:
	docker build --rm -t $(CONTAINER) -f dev.Dockerfile .
	docker run -d --name $(PROJECT) \
		-v $(PWD):$(GO_PROJECTS)/$(PROJECT) \
		-w $(GO_PROJECTS)/$(PROJECT) $(CONTAINER) sleep infinity

shell:
	docker exec -it $(PROJECT) bash

# FIXME azer/logger breaks windows/amd64 build
crossbuild: $(SOURCES)
	gox -verbose \
		-ldflags ${LDFLAGS} \
		-os="linux darwin" \
		-arch="amd64" \
		-output="$(BUILD_PATH)/$(VERSION)/{{.OS}}-{{.Arch}}/{{.Dir}}" .

#release: crossbuild
release:
	[[ -n "$(COMMENT)" ]] || $(error "WHAT")
	git tag -a v$(VERSION) -m '$(COMMENT)'

$(BINARY): $(SOURCES)
	# TODO exclude ./vendor dir
	go build -v -ldflags${LDFLAGS} -o ${BINARY}

dev.install:
	./scripts/setup.sh

install:
	go install -ldflags ${LDFLAGS}

lint:
	GO_VENDOR=1 gometalinter --deadline=25s ./...

test:
	go test $(shell go list ./... | grep -v /vendor/) $(TESTARGS)

.PHONY: godoc
godoc:
	godoc -http=0.0.0.0:6060

.PHONY: clean
clean:
	[[ -d ${BUILD_PATH} ]] && rm -rf ${BUILD_PATH}
	[[ -f ${BINARY} ]] && rm -rf ${BINARY}
