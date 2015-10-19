# Makefile
# vim:ft=make

PROJECT := $(shell basename $(PWD))
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -path './vendor' -prune -o -type f -name '*.go' -print)
PACKAGES=$(shell go list ./... | grep -v /vendor/)

GIT_COMMIT=`git rev-parse HEAD`
GIT_USER=`git config --get user.name`
GO_PROJECTS="/go/src/github.com/$(GIT_USER)"
BUILD_PATH="./_build"
GO_VERSION:=$(shell go version)
# ldflags does't support spaces in variables
CLEAN_GO_VERSION=$(shell echo "${GO_VERSION}" | sed -e 's/[^a-zA-Z0-9]/_/g')

# docker-compose based container name
CONTAINER="$(GIT_USER)/$(PROJECT)"

BINARY=${PROJECT}
# alternative : git describe --always --tags
VERSION="0.1.3"
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X github.com/hackliff/sentinel.Version=${VERSION} -X github.com/hackliff/sentinel.BuildTime=${BUILD_TIME} -X github.com/hackliff/sentinel.GoVersion=${CLEAN_GO_VERSION} -X github.com/hackliff/sentinel.GitCommit=${GIT_COMMIT}"

all: $(BINARY)

container:
	docker build --rm -t $(CONTAINER) .
	docker run -d --name $(PROJECT) \
		-v $(PWD):$(GO_PROJECTS)/$(PROJECT) \
		-w $(GO_PROJECTS)/$(PROJECT) $(CONTAINER) sleep infinity

crossbuild: $(SOURCES)
	glide install
	goxc bump
	goxc -tasks=xc -d=$(BUILD_PATH)

release: doc crossbuild
	goxc -tasks=archive -d=$(BUILD_PATH)
	goxc bintray
	goxc publish-github
	mkdocs gh-deploy --clean

$(BINARY): $(SOURCES)
	# TODO exclude ./vendor dir
	go build -v ${LDFLAGS} -o ${BINARY}

install:
	go install ${LDFLAGS}

# FIXME ./... makes gometalinter to inspect and going crazy in ./vendor
lint:
	GO_VENDOR=1 gometalinter ./...

tests:
	go test $(shell go list ./... | grep -v /vendor/) $(TESTARGS)

doc:
	cp readme.md ./docs/index.md
	cp _build/$(VERSION)/bintray.md ./docs/downloads.md
	mkdocs build --clean

.PHONY: godoc
godoc:
	godoc -http=0.0.0.0:6060

.PHONY: clean
clean:
	[[ -d ${BUILD_PATH} ]] && rm -rf ${BUILD_PATH}
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
