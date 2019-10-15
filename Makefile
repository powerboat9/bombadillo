BINARY     := bombadillo
man1dir    := /usr/local/share/man/man1
GOPATH     ?= ${HOME}/go
GOBIN      ?= ${GOPATH}/bin
BUILD_PATH ?= ${GOBIN}

# Use a dateformat rather than -I flag since OSX
# does not support -I. It also doesn't support 
# %:z - so settle for %z.
BUILD_TIME := ${shell date "+%Y-%m-%dT%H:%M%z"}

# If VERSION is empty or not deffined use the contents of the VERSION file
VERSION    := ${shell git describe --tags 2> /dev/null}
ifndef VERSION
	VERSION  := ${shell cat ./VERSION}
endif

LDFLAGS  := -ldflags "-s -X main.version=${VERSION} -X main.build=${BUILD_TIME}"

.PHONY: install
install: build install-bin install-man clean

.PHONY: build
build:
	go build ${LDFLAGS} -o ${BINARY}

.PHONY: install-man
install-man: bombadillo.1
	gzip -k ./bombadillo.1
	install -d ${man1dir}
	install -m 0644 ./bombadillo.1.gz ${man1dir}

.PHONY: install-bin
install-bin:
	install -d ${BUILD_PATH}
	install ./${BINARY} ${BUILD_PATH}/${BINARY}

.PHONY: clean
clean: 
	go clean
	rm -f ./bombadillo.1.gz 2> /dev/null

.PHONY: uninstall
uninstall: clean
	rm -f ${man1dir}/bombadillo.1.gz
	rm -f ${BUILD_PATH}/${BINARY}

