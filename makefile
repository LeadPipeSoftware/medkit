VERSION=0.0.1-alpha
BUILD_DATE=`date +%FT%T%z`
COMMIT_HASH=`git describe --always`

FULL_IMPORT_PATH=github.com/LeadPipeSoftware/medkit/cmd/medkit

LDFLAGS=-ldflags "-w -s -X ${FULL_IMPORT_PATH}.Version=${VERSION} -X ${FULL_IMPORT_PATH}.BuildDate=${BUILD_DATE} -X ${FULL_IMPORT_PATH}.CommitHash=${COMMIT_HASH}"

build:
	go build ${LDFLAGS}

install:
	go install ${LDFLAGS}

clean:
	if [ -f medkit ] ; then rm medkit ; fi

.PHONY: clean install
