# go-webhookd-aws

## Important

Work in progress.

## Dispatchers

### Lambda

The `Lambda` dispatcher will send messages to an Amazon Web Services (ASW) [Lambda function](#). It is defined as a URI string in the form of:

```
lambda://{FUNCTION}?dsn={DSN}&invocation_type={INVOCATION_TYPE}
```

#### Properties

| Name | Value | Description | Required |
| --- | --- | --- | --- |
| dsn | string | A valid `aaronland/go-aws-session` DSN string. | yes |
| function | string | The name of your Lambda function. | yes |
| invocation_type | string | A valid AWS Lambda `Invocation Type` string. | no |
