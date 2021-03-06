GOCMD := go
BINARY := bombadillo
PREFIX := /usr/local
EXEC_PREFIX := ${PREFIX}
BINDIR := ${EXEC_PREFIX}/bin
DATAROOTDIR := ${PREFIX}/share
MANDIR := ${DATAROOTDIR}/man
MAN1DIR := ${MANDIR}/man1
test : GOCMD := go1.11.13

# Use a dateformat rather than -I flag since OSX
# does not support -I. It also doesn't support 
# %:z - so settle for %z.
BUILD_TIME := ${shell date "+%Y-%m-%dT%H:%M%z"}

.PHONY: build
build:
	${GOCMD} build -o ${BINARY}

.PHONY: install
install: install-bin install-man install-desktop clean

.PHONY: install-man
install-man: bombadillo.1
	gzip -k ./bombadillo.1
	install -d ${DESTDIR}${MAN1DIR}
	install -m 0644 ./bombadillo.1.gz ${DESTDIR}${MAN1DIR}

.PHONY: install-desktop
install-desktop:
ifeq ($(shell uname), Linux)
	# These steps will not work on Darwin, Plan9, or Windows
	# They would likely work on BSD systems
	install -d ${DESTDIR}${DATAROOTDIR}/applications
	install -m 0644 ./bombadillo.desktop ${DESTDIR}${DATAROOTDIR}/applications
	install -d ${DESTDIR}${DATAROOTDIR}/pixmaps
	install -m 0644 ./bombadillo-icon.png ${DESTDIR}${DATAROOTDIR}/pixmaps
	-update-desktop-database 2> /dev/null
else
	@echo "* Skipping protocol handler associations and desktop file creation for non-linux system *"
endif

.PHONY: install-bin
install-bin: build
	install -d ${DESTDIR}${BINDIR}
	install -m 0755 ./${BINARY} ${DESTDIR}${BINDIR}

.PHONY: clean
clean: 
	${GOCMD} clean
	rm -f ./bombadillo.1.gz 2> /dev/null
	rm -f ./${BINARY}_* 2> /dev/null

.PHONY: uninstall
uninstall: clean
	rm -f ${DESTDIR}${MAN1DIR}/bombadillo.1.gz
	rm -f ${DESTDIR}${BINDIR}/${BINARY}
	rm -f ${DESTDIR}${DATAROOTDIR}/applications/bombadillo.desktop
	rm -f ${DESTDIR}${DATAROOTDIR}/pixmaps/bombadillo-icon.png
	-update-desktop-database 2> /dev/null

.PHONY: release
release:
	GOOS=linux GOARCH=amd64 ${GOCMD} build ${LDFLAGS} -o ${BINARY}_linux_64
	GOOS=linux GOARCH=arm ${GOCMD} build ${LDFLAGS} -o ${BINARY}_linux_arm
	GOOS=linux GOARCH=386 ${GOCMD} build ${LDFLAGS} -o ${BINARY}_linux_32
	GOOS=darwin GOARCH=amd64 ${GOCMD} build ${LDFLAGS} -o ${BINARY}_darwin_64


.PHONY: test
test: clean build
