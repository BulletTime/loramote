GOARCH = amd64

SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=loramote

VERSION = 1.0
COMMIT = `git rev-parse HEAD`

LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.build=${COMMIT}"

# .DEFAULT_GOAL: $(BINARY)
#
# $(BINARY): $(SOURCES)
# 	go build ${LDFLAGS} -o ${BINARY} main.go

all: $(SOURCES) linux darwin windows

linux:
	env GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o $(BINARY)-linux-$(GOARCH) main.go

darwin:
	env GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o $(BINARY)-darwin-$(GOARCH) main.go

windows:
	env GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o $(BINARY)-windows-$(GOARCH).exe main.go

install:
	go install ${LDFLAGS} ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f $(BINARY)-linux-$(GOARCH) ] ; then rm $(BINARY)-linux-$(GOARCH) ; fi
	if [ -f $(BINARY)-darwin-$(GOARCH) ] ; then rm $(BINARY)-darwin-$(GOARCH) ; fi
	if [ -f $(BINARY)-windows-$(GOARCH).exe ] ; then rm $(BINARY)-windows-$(GOARCH).exe ; fi


.PHONY: all-platforms linux darwin windows install clean
