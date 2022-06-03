package aws

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3/dispatcher"
	"net/url"
	"testing"
)

func TestLambdaDispatcher(t *testing.T) {

	ctx := context.Background()

	dsn_str := "credentials=fixtures/credentials:example region=us-east-1"

	u := url.URL{}
	u.Scheme = "lambda"
	u.Host = "ExampleFunction"

	q := u.Query()
	q.Set("dsn", dsn_str)
	q.Set("invocation_type", "DryRun")

	u.RawQuery = q.Encode()

	dispatcher_uri := u.String()

	_, err := dispatcher.NewDispatcher(ctx, dispatcher_uri)

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}
}
