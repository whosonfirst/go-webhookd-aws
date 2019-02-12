package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-webhookd/config"	
	"io/ioutil"
	"log"
	"os"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		fh, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			log.Fatal(err)
		}

		cfg := config.WebhookConfig{}
		err = json.Unmarshal(body, &cfg)
		
		if err != nil {
			log.Fatal(err)
		}
		
		enc, err := json.Marshal(cfg)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(enc))
	}
}
