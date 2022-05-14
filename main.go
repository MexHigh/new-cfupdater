package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"

	"git.leon.wtf/leon/new-cfupdater/config"
)

var (
	cfAPI    *cloudflare.API
	confPath *string = flag.String("config", "./config.json", "Path to the config.json file")
	conf     *config.Config
)

func main() {

	fmt.Println()
	fmt.Println("  ~ New CFUpdater by leon.wtf ~  ")
	fmt.Println()

	flag.Parse()

	log.Println("Loading config from ", *confPath)
	confTemp, err := config.Load(*confPath)
	if err != nil {
		panic(err)
	}
	if err := confTemp.Validate(); err != nil {
		panic(err)
	}
	conf = confTemp // assign to global conf

	log.Println("Logging in to Clouflare API")
	cfAPI, err = cloudflare.NewWithAPIToken(conf.APIToken)
	if err != nil {
		panic(err)
	}

	log.Println("Retrieving all zone IDs")
	if err := GetAllZoneIDs(); err != nil {
		panic(err)
	}

	//update(v4, "1.1.1.1")
	err = update(v6, "2.2.2.2")
	fmt.Println(err)
}
