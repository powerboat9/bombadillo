BINARY := bombadillo
PREFIX := /usr/local
MANPREFIX := ${PREFIX}/share/man

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

.PHONY: build
build:
	go build ${LDFLAGS} -o ${BINARY}

.PHONY: install
install: install-bin install-man clean

.PHONY: install-man
install-man: bombadillo.1
	gzip -k ./bombadillo.1
	install -d ${DESTDIR}${MANPREFIX}/man1
	install -m 0644 ./bombadillo.1.gz ${DESTDIR}${MANPREFIX}/man1

.PHONY: install-bin
install-bin: build
	install -d ${DESTDIR}${PREFIX}/bin
	install -m 0755 ./${BINARY} ${DESTDIR}${PREFIX}/bin/${BINARY}

.PHONY: clean
clean: 
	go clean
	rm -f ./bombadillo.1.gz 2> /dev/null

.PHONY: uninstall
uninstall: clean
	rm -f ${DESTDIR}${MANPREFIX}/man1/bombadillo.1.gz
	rm -f ${DESTDIR}${PREFIX}/bin/${BINARY}

