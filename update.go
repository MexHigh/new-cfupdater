package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"git.leon.wtf/leon/new-cfupdater/config"
	"github.com/cloudflare/cloudflare-go"
)

type version int

const (
	v4 version = iota
	v6
)

var zoneIDs = make(map[string]string)

func GetAllZoneIDs() error {
	zonesToRemove := make([]string, 0)

	for zone := range conf.Zones {
		zoneID, err := cfAPI.ZoneIDByName(zone)
		if err != nil {
			if strings.Contains(err.Error(), "zone could not be found") {
				log.Printf("Zone '%s' not found, removing it from zone update list", zone)
				zonesToRemove = append(zonesToRemove, zone)
				continue
			}
			return err
		}
		zoneIDs[zone] = zoneID
	}

	for _, zone := range zonesToRemove {
		delete(conf.Zones, zone)
	}

	return nil
}

func update(v version, newIP string) error {
	for confZone, confRecords := range conf.Zones {
		zoneID := zoneIDs[confZone]

		for _, confRecord := range confRecords {
			// evaluate record type
			var recType string
			if v == v4 && confRecord.UpdateIPv4 {
				recType = "A"
			} else if v == v6 && confRecord.UpdateIPv6 {
				recType = "AAAA"
			} else {
				// record type should not be updated
				continue
			}

			// get matching record(s) for the zone
			records, err := cfAPI.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{
				Type: recType,
				Name: confRecord.Name,
			})
			if err != nil {
				return err
			}

			if l := len(records); l == 0 {
				// record does not exits
				// create if record.Create = true

				if !confRecord.Create {
					log.Printf("Update of record %s %s (Zone: %s) skipped (does not exist and key \"create\" is false)\n", recType, confRecord.Name, confZone)
					continue
				}

				proxied := confRecord.Proxy == config.ProxyActivate // config.ProxyPreserve means false by default !!!

				resp, err := cfAPI.CreateDNSRecord(context.Background(), zoneID, cloudflare.DNSRecord{
					Type:    recType,
					Name:    confRecord.Name,
					Content: newIP,
					TTL:     1, // automatic
					Proxied: &proxied,
				})
				if err != nil {
					return err
				}
				if !resp.Success {
					return fmt.Errorf("error while adding DNS record: %v", resp.Errors)
				}

			} else if l == 1 {
				thisRecord := records[0]
				newRecord := thisRecord

				// change existing record and push
				newRecord.Content = newIP
				if confRecord.Proxy != config.ProxyPreserve {
					proxied := confRecord.Proxy == config.ProxyActivate
					newRecord.Proxied = &proxied
				} // else leave newRecord.Proxied untouched ("ProxyPreserve")

				if err := cfAPI.UpdateDNSRecord(context.Background(), zoneID, thisRecord.ID, newRecord); err != nil {
					return err
				}
			} else {
				return errors.New("got more than one record")
			}
		}

		var updatedVersion string
		if v == v4 {
			updatedVersion = "IPv4"
		} else {
			updatedVersion = "IPv6"
		}

		log.Printf("%s updates for zone %s successfull", updatedVersion, confZone)
	}

	return nil
}

func UpdateIPv4(newIP string) error {
	return update(v4, newIP)
}

func UpdateIPv6(newIP string) error {
	return update(v6, newIP)
}
