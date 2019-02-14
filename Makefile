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
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/algnhsa"
	mv src/github.com/whosonfirst/go-webhookd/vendor/github.com/aws/aws-sdk-go src/github.com/aws/
	mv src/github.com/whosonfirst/go-webhookd/vendor/github.com/whosonfirst/go-whosonfirst-aws src/github.com/whosonfirst/
	mv src/github.com/whosonfirst/go-whosonfirst-aws/vendor/github.com/whosonfirst/go-whosonfirst-cli src/github.com/whosonfirst/

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go

bin: 	rmdeps self
	@GOPATH=$(shell pwd) go build -o bin/webhookd-lambda cmd/webhookd-lambda.go
	@GOPATH=$(shell pwd) go build -o bin/webhookd-lambda-task cmd/webhookd-lambda-task.go
	@GOPATH=$(shell pwd) go build -o bin/webhookd-config cmd/webhookd-config.go

lambda: lambda-webhookd lambda-task

lambda-webhookd:
	@make self
	if test -f main; then rm -f main; fi
	if test -f webhookd.zip; then rm -f webhookd.zip; fi
	@GOPATH=$(GOPATH) GOOS=linux go build -o main cmd/webhookd-lambda.go
	zip webhookd.zip main
	rm -f main

lambda-task:
	@make self
	if test -f main; then rm -f main; fi
	if test -f webhookd-task.zip; then rm -f webhookd-task.zip; fi
	@GOPATH=$(GOPATH) GOOS=linux go build -o main cmd/webhookd-lambda-task.go
	zip webhookd-task.zip main
	rm -f main
