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
)

func v4Schedule() {
	currentIPv4, err := externalip.GetIPv4()
	if err != nil {
		panic(err)
	}
	if currentIPv4 != lastIPv4 {
		log.Printf("Detected IPv4 change (was '%s', is '%s')", lastIPv4, currentIPv4)
		UpdateIPv4(currentIPv4)
	}
	lastIPv4 = currentIPv4
}

func v6Schedule() {
	currentIPv6, err := externalip.GetIPv6()
	if err != nil {
		panic(err)
	}
	if currentIPv6 != lastIPv6 {
		log.Printf("Detected IPv6 change (was '%s', is '%s')", lastIPv6, currentIPv6)
		UpdateIPv6(currentIPv6)
	}
	lastIPv6 = currentIPv6
}

func StartSchedules() {
	scheduler := gocron.NewScheduler(time.Now().Location())
	scheduler.Every(60).Seconds().SingletonMode().Do(v4Schedule)
	scheduler.Every(60).Seconds().SingletonMode().Do(v6Schedule)
	scheduler.StartAsync()
}
