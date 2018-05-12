# Borrowed from:
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = banyandb
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION=0
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=hanahmily
SOURCE_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
BUILD_DIR=${SOURCE_DIR}/cmd

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: clean check test linux darwin windows

install: clean linux darwin windows

linux:
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} .

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} .

windows:
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe .

test:
	cd ${BUILD_DIR}; \
	go test -v ./...

pre-commit:
	sh script/check-gofmt.sh \
	&& sh script/check-golint.sh \
	&& sh script/check-govet.sh

check: lint fmt vet

fmt:
	cd ${SOURCE_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/)

lint:
	cd ${SOURCE_DIR}; \
	golint -set_exit_status $$(go list ./... | grep -v '/vendor/')

vet:
	cd ${SOURCE_DIR}; \
	go vet $$(go list ./... | grep -v /vendor/)

clean:
	-rm -f ${BUILD_DIR}/${BINARY}-*

.PHONY: linux darwin windows test vet fmt lint clean
