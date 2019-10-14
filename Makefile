BINARY     := bombadillo
man1dir    := /usr/local/share/man/man1
#PKGS       := $(shell go list ./... |grep -v /vendor)
VERSION    := $(shell git describe --tags 2> /dev/null)
BUILD		   := $(shell date)
GOPATH     ?= $(HOME)/go
GOBIN      ?= ${GOPATH}/bin
BUILD_PATH ?= ${GOBIN}

# If VERSION is empty or not deffined use the contents of the VERSION file
ifndef VERSION
	VERSION  := $(shell cat ./VERSION)
endif

LDFLAGS    := 
ifdef CONF_PATH
	LDFLAGS  := -ldflags "-s -X main.version=${VERSION} -X main.build=${BUILD} -X main.conf_path${conf_path}"
else
	LDFLAGS  := -ldflags "-s -X main.version=${VERSION} -X main.build=${BUILD}"
endif

.PHONY: test
test:
	@echo ${LDFLAGS}
	@echo ${VERSION}
	@echo ${BUILD_PATH}

.PHONY: build
build:
	@go build ${LDFLAGS} -o ${BINARY}

.PHONY: install
install: install-man
	@go build ${LDFLAGS} -o ${BUILD_PATH}/${BINARY}

.PHONY: install-man
install-man: bombadillo.1
	@gzip -k ./bombadillo.1
	@install -d ${man1dir}
	@install -m 0644 ./bombadillo.1.gz ${man1dir}

.PHONY: clean
clean: 
	@go clean -i
	@rm ./bombadillo.1.gz 2> /dev/null

.PHONY: uninstall
uninstall: clean
	@rm -f ${man1dir}/bombadillo.1.gz
	@echo Removing ${BINARY}
	@rm -f ${BUILD_PATH}/${BINARY}
