BINARY=siren

VERSION=1.0.0

BUILD=`date +%FT%T%z`

# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-X main.Env=production -s -w"

DEV_LDFLAGS=-ldflags "-X main.Env=dev"

TEST_LDFLAGS=-ldflags "-X main.Env=test"

default:
	go build -o ${BINARY} -v ${DEV_LDFLAGS} -tags=jsoniter

install:
	dep ensure
# Installs our project: copies binaries
dev:
	fresh -c configs/development.conf

prod:
	go build -o ${BINARY} -v ${LDFLAGS} -tags=jsoniter

beta:
	go build -o ${BINARY} -v ${DEV_LDFLAGS} -tags=jsoniter

test:
	go test -v ${TEST_LDFLAGS} -tags=jsoniter

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY:  clean run install dev prod test
