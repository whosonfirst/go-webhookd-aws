# go-webhookd-aws

Go package to implement the `whosonfirst/go-webhookd` interfaces for dispatching webhooks originating from GitHub to AWS services.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-webhookd-aws.svg)](https://pkg.go.dev/github.com/whosonfirst/go-webhookd-aws)

Before you begin please [read the go-webhookd documentation](https://github.com/whosonfirst/go-webhookd/blob/master/README.md) for an overview of concepts and principles.

## Usage

```
import (
	_ "github.com/go-webhookd-aws/v4"
)
```

## Dispatchers

### Lambda

The `Lambda` dispatcher will send messages to an Amazon Web Services (ASW) [Lambda function](#). It is defined as a URI string in the form of:

```
lambda://{FUNCTION}?region={AWS_REGIONS}&credentials={AWS_CREDENTIALS}&invocation_type={INVOCATION_TYPE}
```

Where `{FUNCTION}` is the name of the AWS Lambda function to invoke. Valid query parameters are:

#### Properties

| Name | Value | Description | Required |
| --- | --- | --- | --- |
| region | string | The AWS region where the Lambda function is stored. | yes |
| credentials | string | The AWS credentials string used to invoke the Lambda function. | yes |
| invocation_type | string | A valid AWS Lambda `Invocation Type` string. | no |
| halt_on_message | string | An optional regular expression that will be compared to the commit message; if it matches the transformer will return an error with code `webhookd.HaltEvent` | no |
| halt_on_author | string | An optional regular expression that will be compared to the commit author; if it matches the transformer will return an error with code `webhookd.HaltEvent` | no |


#### Credentials strings

Credentials for URIs are defined as string labels. They are:

| Label | Description |
| --- | --- |
| `anon:` | Empty or anonymous credentials. |
| `env:` | Read credentials from AWS defined environment variables. |
| `iam:` | Assume AWS IAM credentials are in effect. |
| `iam:{REGION}:{ARN}` | Assume AWS IAM credentials are in effect after assuming the IAM Role defined by `{ARN}` (in `{REGION}`). |
| `sts:{ARN}` | Assume the role defined by `{ARN}` using STS credentials. |
| `{AWS_PROFILE_NAME}` | This this profile from the default AWS credentials location. |
| `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}` | This this profile from a user-defined AWS credentials location. |

For example:

```
aws:///us-east-1?credentials=iam:
```

## See also

* https://github.com/whosonfirst/go-webhookd