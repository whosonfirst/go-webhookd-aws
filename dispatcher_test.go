package aws

import (
	"context"
	"github.com/whosonfirst/go-webhookd/v3"
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

func TestLambdaDispatcherWithHalt(t *testing.T) {

	ctx := context.Background()

	dsn_str := "credentials=fixtures/credentials:example region=us-east-1"

	u := url.URL{}
	u.Scheme = "lambda"
	u.Host = "ExampleFunction"

	q := u.Query()
	q.Set("dsn", dsn_str)
	q.Set("invocation_type", "DryRun")

	q.Set("halt_on_message", "testing")

	u.RawQuery = q.Encode()

	dispatcher_uri := u.String()

	d, err := dispatcher.NewDispatcher(ctx, dispatcher_uri)

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}

	body := []byte(`#message I am testing things
whosonfirst-data-admin-ca`)

	err = d.Dispatch(ctx, body)

	if err == nil {
		t.Fatalf("Expected dispatcher to return an error")
	}

	if err.(*webhookd.WebhookError).Code != webhookd.HaltEvent {
		t.Fatalf("Expected HaltEvent but got %v", err)
	}
}

func TestProcessBodyWithMessage(t *testing.T) {

	ctx := context.Background()

	dsn_str := "credentials=fixtures/credentials:example region=us-east-1"

	u := url.URL{}
	u.Scheme = "lambda"
	u.Host = "ExampleFunction"

	q := u.Query()
	q.Set("dsn", dsn_str)
	q.Set("invocation_type", "DryRun")

	q.Set("halt_on_message", "testing")

	u.RawQuery = q.Encode()

	dispatcher_uri := u.String()

	d, err := dispatcher.NewDispatcher(ctx, dispatcher_uri)

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}

	body := []byte(`#message I am testing things
whosonfirst-data-admin-ca`)

	body, err = d.(*LambdaDispatcher).processBody(ctx, body)

	if err != nil && err.(*webhookd.WebhookError).Code != webhookd.HaltEvent {
		t.Fatalf("Failed to process body, %v", err)
	}

	if string(body) != "" {
		t.Fatalf("Unexpected body post-processing '%s'", string(body))
	}

}

func TestProcessBody(t *testing.T) {

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

	d, err := dispatcher.NewDispatcher(ctx, dispatcher_uri)

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}

	body := []byte(`whosonfirst-data-admin-ca`)

	body, err = d.(*LambdaDispatcher).processBody(ctx, body)

	if err != nil {
		t.Fatalf("Failed to process body, %v", err)
	}

	if string(body) != "whosonfirst-data-admin-ca" {
		t.Fatalf("Unexpected body post-processing '%s'", string(body))
	}

}

func TestPreambleRegularExpression(t *testing.T) {

	tests := []string{
		"#message hello world",
		"# message boo",
		"#author bob",
		"#bob hungey hippo",
	}

	for _, str := range tests {

		m := preamble_re.FindStringSubmatch(str)

		if len(m) != 3 {
			t.Fatalf("Preamble regular expression failed with '%s' %d", str, len(m))
		}
	}
}
