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

#### Roles

Your Lambda function will need to run using a role with the following built-in AWS policies:

* `AWSLambdaBasicExecutionRole`

#### Environment variables

| Key | Value |
| --- | --- | --- |
| WEBHOOKD_CONFIG | A valid JSON encoded `webhookd` config file | 

For details on `webhookd` config file please consult the [go-webhookd documentation](https://github.com/whosonfirst/go-webhookd#config-files).

Including a big honking string here is not ideal, it's just how it is today. Really this should be stored in something like the AWS Secrets Manager but that will have to be "tomorrow's problem".

For example, let's start with a config file that looks like this:

```
{
    	"daemon": {
		"protocol": "http",
		"host": "localhost",
		"port": 8080
	},
	"receivers": {
		"insecure": {
			"name": "Insecure"
		}		
	},	
	"transformations": {
		"chicken": {
			"name": "Chicken",
			"language": "zxx",
			"clucking": false
		}				
		
	},
	"dispatchers": {
		"log": {
			"name": "Log"
		}
	},
	"webhooks": [
		{
			"endpoint": "/insecure",
		 	"receiver": "insecure",
			"transformations": [ "chicken" ],
			"dispatchers": [ "log" ]
		}
	]
}
```

And an AWS Lambda test event configured to act like an AWS API Gateway Proxy event like this:

```
{
  "body": "hello world",
  "resource": "/{proxy+}",
  "path": "/insecure",
  "httpMethod": "POST",
  "isBase64Encoded": false,
  "pathParameters": {
    "proxy": "/cnSkGhLKpLqTWYOKpEmdLeljtYMeIlTPmFONZcSNGSGotFBPMSGBfgjzOwlTKaMM"
  },
  "headers": {
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
    "Accept-Encoding": "gzip, deflate, sdch",
    "Accept-Language": "en-US,en;q=0.8",
    "Cache-Control": "max-age=0",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "US",
    "Host": "1234567890.execute-api.us-east-1.amazonaws.com",
    "Upgrade-Insecure-Requests": "1",
    "User-Agent": "Custom User Agent String",
    "Via": "1.1 08f323deadbeefa7af34d5feb414ce27.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "cDehVQoZnx43VYQb9j2-nvCh-9z396Uhbp027Y2JvkCPNLmGJHqlaA==",
    "X-Forwarded-For": "127.0.0.1, 127.0.0.2",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  }
}
```

When you run the test you should see something like this:

```
START RequestId: 5c237ed9-f03f-407f-acaa-433625aa6950 Version: $LATEST
2019/08/10 18:08:09 🐔 🐔
END RequestId: 5c237ed9-f03f-407f-acaa-433625aa6950
REPORT RequestId: 5c237ed9-f03f-407f-acaa-433625aa6950	Duration: 3.14 ms	Billed Duration: 100 ms 	Memory Size: 512 MB	Max Memory Used: 55 MB	
```

Specifically the string `hello world` was received by the "insecure" receiver, transformed in to `🐔 🐔` by the "chicken" transformer and distpatched to STDOUT (or in the case of AWS to CloudWatch) using the "log" dispatcher.

For details on receivers, transformers and dispatchers please consult the main [go-webhookd documentation](https://github.com/whosonfirst/go-webhookd/blob/master/README.md).

#### API Gateway

_Please write me_

### webhookd-lambda-task

This is a Lambda function to run an ECS task when invoked. It is principally meant to be used with the `go-webhookd` Lambda dispatcher. It should probably be renamed since it's pretty confusing, even for me.

## See also

* https://github.com/whosonfirst/go-webhookd
* https://github.com/whosonfirst/algnhsa