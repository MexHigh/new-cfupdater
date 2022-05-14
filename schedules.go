package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"

	"git.leon.wtf/leon/new-cfupdater/externalip"
)

var (
	lastIPv4 string
	lastIPv6 string
	//errChan  = make(chan error, 1)
)

func v4Schedule() {
	currentIPv4, err := externalip.GetIPv4()
	if err != nil {
		log.Printf("Got error while getting IPv4: %s", err.Error())
		return
	}
	log.Printf("Got %s", currentIPv4)
	if currentIPv4 != lastIPv4 {
		log.Printf("Detected IPv4 change (was '%s', is '%s')", lastIPv4, currentIPv4)
		UpdateIPv4(currentIPv4)
	}
	lastIPv4 = currentIPv4
}

func v6Schedule() {
	currentIPv6, err := externalip.GetIPv6()
	if err != nil {
		log.Printf("Got error while getting IPv6: %s", err.Error())
		return
	}
	log.Printf("Got %s", currentIPv6)
	if currentIPv6 != lastIPv6 {
		log.Printf("Detected IPv6 change (was '%s', is '%s')", lastIPv6, currentIPv6)
		UpdateIPv6(currentIPv6)
	}
	lastIPv6 = currentIPv6
}

func StartSchedules(interval int) {
	scheduler := gocron.NewScheduler(time.Now().Location())
	scheduler.Every(interval).Seconds().SingletonMode().Do(v4Schedule)
	scheduler.Every(interval).Seconds().SingletonMode().Do(v6Schedule)
	//scheduler.StartAsync()
	scheduler.StartBlocking()
}
