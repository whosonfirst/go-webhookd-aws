CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-webhookd-aws
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-webhookd"
	@GOPATH=$(shell pwd) go get -u "github.com/aws/aws-lambda-go/..."

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	# go fmt *.go

bin: 	rmdeps self
	# @GOPATH=$(shell pwd) go build -o bin/webhookd cmd/webhookd.go

