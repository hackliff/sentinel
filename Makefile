# Makefile
# vim:ft=make

#  Configuration  #############################################################

BUILD_PATH       := ./_dist
RELEASE_OPTS     ?= ""

PROJECT          := $(shell basename $(PWD))
SOURCES          := $(shell find . -path './vendor' -prune -o -type f -name '*.go' -print)
PACKAGES         := $(shell go list ./... | grep -v /vendor/)

GIT_COMMIT       := $(shell git rev-parse HEAD)
GIT_USER         := $(shell git config --get user.name)
# alternative : git describe --always --tags
VERSION          := $(shell cat env.go | grep "Version =" | cut -d"\"" -f 2)
#VERSION:=0.2.0
GO_PROJECTS      := "/go/src/github.com/$(GIT_USER)"
GO_VERSION       := $(shell go version)
# ldflags does't support spaces in variables
CLEAN_GO_VERSION := $(shell echo "${GO_VERSION}" | sed -e 's/[^a-zA-Z0-9]/_/g')

# docker-compose based container name
CONTAINER        := "$(GIT_USER)/$(PROJECT)"

BUILD_TIME       := $(shell date +%FT%T%z)
BINARY           := ${PROJECT}
LDFLAGS          := "-X github.com/hackliff/sentinel.BuildTime=${BUILD_TIME} -X github.com/hackliff/sentinel.GoVersion=${CLEAN_GO_VERSION} -X github.com/hackliff/sentinel.GitCommit=${GIT_COMMIT}"

###############################################################################

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
		-output="$(BUILD_PATH)/$(VERSION)/{{.Dir}}-{{.OS}}-{{.Arch}}" .

release: crossbuild
ifndef COMMENT
	$(error no tag description provided)
endif
	git tag -a $(VERSION) -m '$(COMMENT)'
	git push --tags
	ghr $(RELEASE_OPTS) v$(VERSION) $(BUILD_PATH)/$(VERSION)/

$(BINARY): $(SOURCES)
	gom build -v -ldflags ${LDFLAGS} -o ${BINARY}

.PHONY: install.tools
install.tools:
	go get github.com/mattn/gom
	# code coverage
	go get github.com/axw/gocov/gocov
	# cross-compilation
	go get github.com/mitchellh/gox
	# github release publication
	go get github.com/tcnksm/ghr
	# code linting
	go get github.com/alecthomas/gometalinter && \
		gometalinter --install --update

install.hack: install.tools
	gom -groups=plugins,hack install

install:
	go install -ldflags ${LDFLAGS}

lint:
	GO_VENDOR=1 gometalinter --deadline=25s ./...

test:
	gom test $(shell go list ./... | grep -v /vendor/) $(TESTARGS)

.PHONY: godoc
godoc:
	godoc -http=0.0.0.0:6060

.PHONY: clean
clean:
	[[ -d ${BUILD_PATH} ]] && rm -rf ${BUILD_PATH}
	[[ -f ${BINARY} ]] && rm -rf ${BINARY}
