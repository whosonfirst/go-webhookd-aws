package main

import (
	"encoding/json"
	"flag"
	"github.com/whosonfirst/algnhsa"
	"github.com/whosonfirst/go-webhookd/config"
	"github.com/whosonfirst/go-webhookd/daemon"
	"log"
	"net/http"
	"os"
)

func main() {

	flag.Parse()

	str_cfg, ok := os.LookupEnv("WEBHOOKD_CONFIG")

	if !ok {
		log.Fatal("Missing WEBHOOKD_CONFIG environment variable")
	}

	cfg := config.WebhookConfig{}
	err := json.Unmarshal([]byte(str_cfg), &cfg)

	if err != nil {
		log.Fatal(err)
	}

	d, err := daemon.NewWebhookDaemonFromConfig(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	handler, err := d.HandlerFunc()

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	lambda_opts := new(algnhsa.Options)
	algnhsa.ListenAndServe(mux, lambda_opts)
}
