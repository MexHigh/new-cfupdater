package main

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"

	"git.leon.wtf/leon/new-cfupdater/externalip"
)

var (
	lastIPv4 = &externalip.IPv4{}
	lastIPv6 = &externalip.IPv6{}
)

func v4Schedule() {
	currentIPv4 := externalip.GetIPv4()
	if currentIPv4.Addr != lastIPv4.Addr {
		log.Printf("Detected IPv4 change (was '%s', is '%s')", lastIPv4.Addr, currentIPv4.Addr)
		UpdateIPv4(currentIPv4)
	}
	lastIPv4 = currentIPv4
}

func v6Schedule() {
	currentIPv6 := externalip.GetIPv6()
	if currentIPv6.Addr != lastIPv6.Addr {
		log.Printf("Detected IPv6 change (was '%s', is '%s')", lastIPv6.Addr, currentIPv6.Addr)
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
