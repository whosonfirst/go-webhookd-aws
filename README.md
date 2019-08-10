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

### webhookd-flatten-config

A helper utility for encoding a valid webhookd config file in to a string that can be copy-paste-ed as a `webhookd-lambda` environment variable.

```
./bin/webhookd-flatten-config -config config.json
{"daemon":{"host":"localhost","port":8080,"allow_debug":false},"receivers":{"github":{"name":"GitHub","secret":"s33kret","Ref":"refs/heads/master"},"insecure":{"name":"Insecure","Ref":""}},"dispatchers":{"log":{"name":"Log"},"null":{"name":"Null"}},"transformations":{"chicken":{"name":"Chicken","language":"zxx"},"clucking":{"name":"Chicken","language":"eng","clucking":true},"commits":{"name":"GitHubCommits","exclude_additions":true,"exclude_modifications":true,"exclude_deletions":true},"null":{"name":"Null"}},"webhooks":[{"endpoint":"/github-test","receiver":"github","transformations":["chicken"],"dispatchers":["log"]},{"endpoint":"/insecure-test","receiver":"insecure","transformations":["chicken"],"dispatchers":["log"]}]}
```

### webhookd-lambda

This will run `webhookd` HTTP daemon as a Lambda function. In order to use it you will need to configure an AWS API Gateway endpoint.

### webhookd-lambda-task

This is a Lambda function to run an ECS task when invoked. It is principally meant to be used with the `go-webhookd` Lambda dispatcher. It should probably be renamed since it's pretty confusing, even for me.

## AWS

### Lambda

#### Roles

Your Lambda function will need to run using a role with the following built-in AWS policies:

* `AWSLambdaBasicExecutionRole`

#### Environment variables

| Key | Value | Notes |
| --- | --- | --- |
| WEBHOOKD_CONFIG | A valid JSON encoded `webhookd` config file | Including a big honking string here is not ideal, it's just how it is today |

For details on `webhookd` config file please consult the [go-webhookd documentation](https://github.com/whosonfirst/go-webhookd#config-files).

### API Gateway

_Please write me_

## See also

* https://github.com/whosonfirst/go-webhookd
* https://github.com/whosonfirst/algnhsa