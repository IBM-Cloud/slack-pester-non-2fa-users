all: test

clean:
	rm -f slack-disable-2fa

install: prepare
	godep go install

prepare:
	go get github.com/tools/godep

build: prepare
	godep go build

test: prepare build
	echo "no tests"

.PHONY: install prepare build test
