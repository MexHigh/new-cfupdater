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

	scheduler = gocron.NewScheduler(time.Now().Location())

	emergencyModev4Active = false
	emergencyModev4Set    = make(chan bool, 1)
	emergencyModev6Active = false
	emergencyModev6Set    = make(chan bool, 1)
)

func v4Schedule() {
	currentIPv4, err := externalip.GetIPv4(conf.CheckTimeout)
	if err != nil {
		log.Printf("Got error while getting IPv4: %s", err.Error())
		emergencyModev4Set <- true
		return
	}
	log.Printf("Got %s", currentIPv4)
	emergencyModev4Set <- false
	if currentIPv4 != lastIPv4 {
		log.Printf("Detected IPv4 change (was '%s', is '%s')", lastIPv4, currentIPv4)
		UpdateIPv4(currentIPv4)
	}
	lastIPv4 = currentIPv4
}

func v6Schedule() {
	currentIPv6, err := externalip.GetIPv6(conf.CheckTimeout)
	if err != nil {
		log.Printf("Got error while getting IPv6: %s", err.Error())
		emergencyModev6Set <- true
		return
	}
	log.Printf("Got %s", currentIPv6)
	emergencyModev6Set <- false
	if currentIPv6 != lastIPv6 {
		log.Printf("Detected IPv6 change (was '%s', is '%s')", lastIPv6, currentIPv6)
		UpdateIPv6(currentIPv6)
	}
	lastIPv6 = currentIPv6
}

var (
	v4CurrentJob *gocron.Job
	v6CurrentJob *gocron.Job
)

func handleEmergencies(interval int) {
	var err error
	for {
		select {
		case v4Emergency := <-emergencyModev4Set:
			if v4Emergency && !emergencyModev4Active {
				// remove normal job and add emegency job
				scheduler.RemoveByReference(v4CurrentJob)
				v4CurrentJob, err = scheduler.Every(10).Seconds().SingletonMode().WaitForSchedule().Do(v4Schedule)
				if err != nil {
					panic(err)
				}
				emergencyModev4Active = true
				log.Printf("Emergency mode activated for IPv4")
			} else if !v4Emergency && emergencyModev4Active {
				// remove emergency job and add normal job
				scheduler.RemoveByReference(v4CurrentJob)
				v4CurrentJob, err = scheduler.Every(interval).Seconds().SingletonMode().WaitForSchedule().Do(v4Schedule)
				if err != nil {
					panic(err)
				}
				emergencyModev4Active = false
				log.Printf("Emergency mode deactivated for IPv4")
			} else {
				// no change required
				continue
			}
		case v6Emergency := <-emergencyModev6Set:
			if v6Emergency && !emergencyModev6Active {
				// remove normal job and add emegency job
				scheduler.RemoveByReference(v6CurrentJob)
				v6CurrentJob, err = scheduler.Every(10).Seconds().SingletonMode().WaitForSchedule().Do(v6Schedule)
				if err != nil {
					panic(err)
				}
				emergencyModev6Active = true
				log.Printf("Emergency mode activated for IPv6")
			} else if !v6Emergency && emergencyModev6Active {
				// remove emergency job and add normal job
				scheduler.RemoveByReference(v6CurrentJob)
				v6CurrentJob, err = scheduler.Every(interval).Seconds().SingletonMode().WaitForSchedule().Do(v6Schedule)
				if err != nil {
					panic(err)
				}
				emergencyModev6Active = false
				log.Printf("Emergency mode deactivated for IPv6")
			} else {
				// no change required
				continue
			}
		}
	}
}

func StartSchedules(interval int) {

	go handleEmergencies(interval)

	var err error
	v4CurrentJob, err = scheduler.Every(interval).Seconds().SingletonMode().Do(v4Schedule)
	if err != nil {
		panic(err)
	}
	v6CurrentJob, err = scheduler.Every(interval).Seconds().SingletonMode().Do(v6Schedule)
	if err != nil {
		panic(err)
	}

	//scheduler.StartAsync()
	scheduler.StartBlocking()
}
