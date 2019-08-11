package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-webhookd/config"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	config_path := flag.String("config", "", "The path your webhookd config file")

	flag.Parse()

	abs_path, err := filepath.Abs(*config_path)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(abs_path)

	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		log.Fatal(err)
	}

	var cfg config.WebhookConfig

	err = json.Unmarshal(body, &cfg)

	if err != nil {
		log.Fatal(err)
	}

	body, err = json.Marshal(cfg)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
