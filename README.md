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

This will run the `webhookd` HTTP daemon as a Lambda function. In order to use it you will need to configure an AWS API Gateway endpoint, which is discussed below.

#### Roles

Your Lambda function will need to run using a role with the following built-in AWS policies:

* `AWSLambdaBasicExecutionRole`

#### Environment variables

| Key | Value |
| --- | --- |
| WEBHOOKD_CONFIG | A valid JSON encoded `webhookd` config file | 

Including a big honking string here is not ideal, it's just how it is today. For now you'll just have to use the `webhookd-flatten-config` tool described above. Really, this should be stored in something like the AWS Secrets Manager but that will have to be "tomorrow's problem".

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

For details on `webhookd` config file please consult the [go-webhookd documentation](https://github.com/whosonfirst/go-webhookd#config-files).

And an AWS Lambda test event configured to act like an AWS API Gateway Proxy event like this:

```
{
  "body": "hello world",
  "resource": "/{proxy+}",
  "path": "/insecure",
  "httpMethod": "POST",
  "isBase64Encoded": false,
  "pathParameters": {
    "proxy": "/insecure"
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

Specifically the string `hello world` was received by the ["insecure" receiver](https://github.com/whosonfirst/go-webhookd#insecure), transformed in to `🐔 🐔` by the ["chicken" transformer](https://github.com/whosonfirst/go-webhookd#chicken) and dispatched to STDOUT (or in the case of AWS to CloudWatch) using the ["log" dispatcher](https://github.com/whosonfirst/go-webhookd#log).

For details on receivers, transformers and dispatchers please consult the main [go-webhookd documentation](https://github.com/whosonfirst/go-webhookd/blob/master/README.md).

#### API Gateway

In order to send requests to the `webhookd` Lambda function over HTTP you need to configure an API Gateway instance to sit in front of it (the Lambda function) and proxy requests and responses.

* Create a new API Gateway API
* Create a new resource for that API
* Configure as "proxy resource" (enable API Gateway CORS if you think that's necessary)
* Delete the `ANY` method and create a new `POST` method
* Configure it (the `POST` method) to use the "Lambda Function Proxy" integration type and associate it with whatever you've named your Lambda function (above)
* Deploy your API Gateway API (for the purposes of this example we're going to say you called it `STAGE`)

Now let's say we have a file called `test.svg` that looks like this:

```
<svg width="512.000000" height="512.000000" viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg"><path d="M512.000000 51.171577,502.111346 59.685313,500.130702 79.677695,505.024057 80.118699,507.237717 85.263736,512.000000 83.358551,512.000000 511.923681,0.000000 511.923681,0.000000 0.000000,512.000000 0.000000,512.000000 51.171577 Z" fill="#ffffff" fill-opacity="0.5" kind="ocean" sort_rank="200" stroke="#000000" stroke_opacity="1"/><path d="M512.000000 51.171577,502.111346 59.685313,500.130702 79.677695,505.024057 80.118699,507.237717 85.263736,512.000000 83.358551" fill="#ffffff" fill-opacity="0" kind="ocean" sort_rank="205" stroke="#000000" stroke_opacity="1"/></svg>
```

Your API gateway lives at `EXAMPLE.execute-api.us-east-1.amazonaws.com/STAGE` so you would post your SVG file like this:
 
```
$> curl -v -X POST https://EXAMPLE.execute-api.us-east-1.amazonaws.com/STAGE/insecure. -d@test.svg
...
> POST /STAGE/insecure HTTP/2
> Host: EXAMPLE.execute-api.us-east-1.amazonaws.com
> User-Agent: curl/7.54.0
> Accept: */*
> Content-Length: 682
> Content-Type: application/x-www-form-urlencoded
> 
* Connection state changed (MAX_CONCURRENT_STREAMS updated)!
* We are completely uploaded and fine
< HTTP/2 200 
< date: Sat, 10 Aug 2019 18:57:08 GMT
< content-type: application/json
< content-length: 0
< x-amzn-requestid: a714f2c0-bba0-11e9-95d9-2510c8e0110e
< x-webhookd-time-to-receive: 8.491?s
< x-webhookd-time-to-transform: 1.391499ms
< x-amz-apigw-id: EXAMPLE
< x-webhookd-time-to-process: 1.454697ms
< x-webhookd-time-to-dispatch: 54.153?s
```

And in your CloudWatch logs you'd see something like this:

```
2019/08/10 18:57:08 🐔-🐔 🐔🐔 🐔🐔 🐔"512.000000" 🐔🐔 🐔"512.000000" 🐔🐔 🐔"0 0 512 512" 🐔🐔 🐔"🐔://🐔.🐔3.🐔/2000/🐔"🐔-🐔 🐔🐔-🐔 🐔🐔 🐔🐔 🐔"🐔512.000000 51.171577,502.111346 59.685313,500.130702 79.677695,505.024057 80.118699,507.237717 85.263736,512.000000 83.358551,512.000000 511.923681,0.000000 511.923681,0.000000 0.000000,512.000000 0.000000,512.000000 51.171577 🐔" 🐔🐔 🐔"#🐔" 🐔-🐔🐔 🐔"0.5" 🐔🐔 🐔"🐔" 🐔_🐔🐔 🐔"200" 🐔🐔 🐔"#000000" 🐔_🐔🐔 🐔"1"/🐔-🐔 🐔🐔-🐔 🐔🐔 🐔🐔 🐔"🐔512.000000 51.171577,502.111346 59.685313,500.130702 79.677695,505.024057 80.118699,507.237717 85.263736,512.000000 83.358551" 🐔🐔 🐔"#🐔" 🐔-🐔🐔 🐔"0" 🐔🐔 🐔"🐔" 🐔_🐔🐔 🐔"205" 🐔🐔 🐔"#000000" 🐔_🐔🐔 🐔"1"/🐔-🐔 🐔🐔-🐔 🐔/🐔🐔-🐔 🐔
```

As of this writing it is not possible to send back the value of a transformation in the response body of a (`webhookd`) request.

#### A more concrete example

Here's another example that doesn't involve chickens (🐔). In this example we'll configure `webhookd` to listen for webhooks sent by GitHub and log the names of the repository that sent the hook.

Your config file should look like this, albeit with specific secrets and endpoints:

```
{
    "daemon": {
		"protocol": "http",
		"host": "localhost",
		"port": 8080
	},
	"receivers": {
		"github": {
			"name": "GitHub",
			"secret": "S33KRET",
			"ref": "refs/heads/master"
		}			    
	},	
	"transformations": {
		"repo": {
			"name": "GitHubRepo",
			"exclude_additions": false,
			"exclude_modifications": false,
			"exclude_deletions": true
		}				
		
	},
	"dispatchers": {
		"log": {
			"name": "Log"
		}
	},
	"webhooks": [
		{
		    "endpoint": "/ENDPOINT",
		    "receiver": "github",
		    "transformations": [ "repo" ],
		    "dispatchers": [ "log" ]
		}	    
	]
}
```

In the "Webhooks settings" page for your GitHub repository in question you'll want to plug in the following:

| Key | Value |
| --- | --- |
| Payload URL | https://EXAMPLE.execute-api.us-east-1.amazonaws.com/STAGE/ENDPOINT |
| Content type | application/json |
| Secret | S33KRET |

Where things like `EXAMPLE` and `ENDPOINT` and especially `S33KRET` are specific to your application.

In this example we have configured the [GitHub "receiver"](https://github.com/whosonfirst/go-webhookd#github) to only pay attention to things that have been committed to the `master` branch and we have configured the [GitHubRepo "transformer"](https://github.com/whosonfirst/go-webhookd#githubrepo) to ignore any deletion events (in the commit). As with the other examples we're "dispatching" everything to the ["log" dispatcher](https://github.com/whosonfirst/go-webhookd#log) so the output would like something like this:

![](docs/images/github-repo-log.png)

### webhookd-lambda-task

This is a Lambda function to run an ECS task when invoked. It is principally meant to be used with the `go-webhookd` [Lambda dispatcher](https://github.com/whosonfirst/go-webhookd/#lambda), and in particular with the [GitHubRepo transformer](https://github.com/whosonfirst/go-webhookd/#githubrepo) in a [Who's On First](https://whosonfirst.org/) context.

There is an [open ticket](https://github.com/whosonfirst/go-webhookd/issues/19) to add a dedicated "run an ECS Task" dispatcher but until that's completed this is what we're stuck with.

#### Roles

Your Lambda function will need to run using a role with the following built-in AWS policies:

* `AWSLambdaBasicExecutionRole`

Additionally you will need the following policies, or equivalents:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "ecs:RunTask",
            "Resource": "arn:aws:ecs:{AWS_REGION}:{AWS_ACCOUNT_ID}:task-definition/{ECS_TASK_NAME}:*"
        },
        {
            "Sid": "VisualEditor1",
            "Effect": "Allow",
            "Action": [
                "iam:PassRole"
            ],
            "Resource": [
                "arn:aws:iam::{AWS_ACCOUNT_ID}:role/{AWS_IAM_ROLE_FOR_THE_ECS_TASK_YOU_WANT_TO_RUN}"
            ]
        }
    ]
}
```

#### Environment variables

| Key | Value |
| --- | --- |
| WEBHOOKD_MODE | `lambda` |
| WEBHOOKD_COMMAND | `{SOME COMMAND ON YOUR CONTAINER},%s` |
| WEBHOOKD_ECS_CLUSTER | `{ECS_CLUSTER_NAME}` |
| WEBHOOKD_ECS_CONTAINER | `{ECS_CONTAINER_NAME}` |
| WEBHOOKD_ECS_DSN | `credentials=iam: region={AWS_REGION}` |
| WEBHOOKD_ECS_SECURITY_GROUP | `{AWS_SECURITY_GROUP}` |
| WEBHOOKD_ECS_SUBNET | `{AWS_SUBNET1},{AWS_SUBNET2}...` |
| WEBHOOKD_ECS_TASK | `{ECS_TASK_NAME}:{ECS_TASK_REVISION}` |

The `{SOME COMMAND ON YOUR CONTAINER},%s` string assumes the same comma-separated arguments syntax used by Docker and/or ECS container override statements.

See the way `WEBHOOKD_COMMAND` is defined as `"{SOME COMMAND ON YOUR CONTAINER},%s"` ? That's because under the hood the code was originally written to pass the payload received by the Lambda function as the second argument to the Go `fmt.Sprintf` method.

```
lambda_handler := func(ctx context.Context, payload string) (interface{}, error) {
	return launchTask(*command, payload)
}

launchTask := func(command string, args ...interface{}) (interface{}, error) {

	str_cmd := fmt.Sprintf(command, args...)
	cmd := strings.Split(str_cmd, ",")

	task_rsp, err := ecs.LaunchTask(task_opts, cmd...)
	...
}		
```

Is this awesome? No. Could it be abused? You bet. Is it a bad idea, generally? Probably.

With that in mind, the code has been updated to be strict about input when run as a Lambda function. Specifically, all payloads are validated against the following regular expression: `^[a-zA-Z0-9\-_]+$`. As in:

```
re, err := regexp.Compile(`^[a-zA-Z0-9\-_]+$`)

if err != nil {
	log.Fatal(err)
}
		
lambda_handler := func(ctx context.Context, payload string) (interface{}, error) {

	if !*command_insecure {

		if !re.MatchString(payload){
			return nil, errors.New("Invalid payload")
		}
	}
			
	return launchTask(*command, payload)
}
```

You can override this restriction by setting the `WEBHOOKD_COMMAND_INSECURE=true` environment variable. If you do it is assumed that you are confident about any input being passed to the Lambda function and then on to the ECS task.

By default, the Lambda function will launch the task and return without waiting to see whether it succeeded or not. If you want to wait and check the response of the task you need to set the following environment variables:

| Key | Value |
| --- | --- |
| WEBHOOKD_MONITOR | `true` |
| WEBHOOKD_LOGS | boolean, indicating whether to return CloudWatch logs with the response (default is false) |
| WEBHOOKD_LOGS_DSN | "credentials=iam: region={AWS_REGION}" |

## See also

* https://github.com/whosonfirst/go-webhookd
* https://github.com/whosonfirst/algnhsa