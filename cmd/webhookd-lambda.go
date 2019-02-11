package main

import (
	"flag"
	"github.com/whosonfirst/algnhsa"
	"github.com/whosonfirst/go-webhookd/config"
	"github.com/whosonfirst/go-webhookd/daemon"
	"log"
	"net/http"
	_ "os"
)

func main() {

	var cfg = flag.String("config", "", "Path to a valid webhookd config file")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	wh_config, err := config.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	d, err := daemon.NewWebhookDaemonFromConfig(wh_config)

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
	// lambda_opts.BinaryContentTypes = []string{"image/png"}

	algnhsa.ListenAndServe(mux, lambda_opts)
}
