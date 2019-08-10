# go-webhookd-aws

What is the simplest webhook-wrangling server-daemon-thing (`go-webhookd`) that can run as an AWS Lambda function.

## Important

This is work in progress and the documentation is not complete yet. There is also a possibility this code will get merged back in to [go-webhookd](https://github.com/whosonfirst/go-webhookd) itself. I'm not sure yet.

## Install

You will need to have both `Go` (specifically version [1.12](https://golang.org/dl) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make tools
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Tools

### webhookd-lambda

This will run `webhookd` HTTP daemon as a Lambda function. In order to use it you will need to configure an AWS API Gateway endpoint.

### webhookd-lambda-task

This is a Lambda function to run an ECS task when invoked. It is principally meant to be used with the `go-webhookd` Lambda dispatcher. It should probably be renamed since it's pretty confusing, even for me.

## AWS

### Lambda

_Please write me_

### API Gateway

_Please write me_

## See also

* https://github.com/whosonfirst/go-webhookd
* https://github.com/whosonfirst/algnhsa