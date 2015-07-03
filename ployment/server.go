package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/anxiousmodernman/ployment/config"
	"github.com/anxiousmodernman/ployment/webhook"
)

func main() {
	var cfgPath = flag.String("config", "", "Location of config.json")
	flag.Parse()

	if *cfgPath == "" {
		fmt.Println("Error: provide a configuration file with -config")
		os.Exit(1)
	}
	cfg, err := config.FromFile(*cfgPath) //todo: use the config returned here
	ctx := &webhook.AppContext{cfg}

	h := webhook.Hook{ctx, webhook.WebhookHandler}

	// dereference the cfg string pointer
	//REDO
	if err != nil {
		fmt.Printf("Error opening config file: ", err.Error()) // fix ugly message
		os.Exit(1)
	}
	http.Handle("/hook/", h)
	http.ListenAndServe(":8080", nil)
}
