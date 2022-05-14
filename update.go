package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"git.leon.wtf/leon/new-cfupdater/config"
	"git.leon.wtf/leon/new-cfupdater/externalip"
	"github.com/cloudflare/cloudflare-go"
)

type version int

const (
	v4 version = iota
	v6
)

var zoneIDs = make(map[string]string)

func GetAllZoneIDs() error {
	for zone := range conf.Zones {
		zoneID, err := cfAPI.ZoneIDByName(zone)
		if err != nil {
			return err
		}
		zoneIDs[zone] = zoneID
	}
	return nil
}

func update(v version, newIP string) error {

	for zone, records := range conf.Zones {

		zoneID := zoneIDs[zone]

		for _, record := range records {

			// evaluate record type
			var recType string
			if v == v4 && record.UpdateIPv4 {
				recType = "A"
			} else if v == v6 && record.UpdateIPv6 {
				recType = "AAAA"
			} else {
				fmt.Println("kaputt")
				continue
			}

			// get matching record(s) for the zone
			records, err := cfAPI.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{
				Type: recType,
				Name: record.Name,
			})
			if err != nil {
				return err
			}

			if l := len(records); l == 0 {

				// record does not exits
				// create if record.Create = true

				if !record.Create {
					log.Printf("Update of record %s %s (Zone: %s) skipped (does not exist and key \"create\" is false)\n", recType, record.Name, zone)
					continue
				}

				proxy := record.Proxy == config.ProxyActivate // ProxyPreserve means false in this case !!!
				resp, err := cfAPI.CreateDNSRecord(context.Background(), zoneID, cloudflare.DNSRecord{
					Type:    recType,
					Name:    record.Name,
					Content: newIP,
					TTL:     1, // automatic
					Proxied: &proxy,
				})
				if err != nil {
					return err
				}
				if !resp.Success {
					return fmt.Errorf("error while adding DNS record: %v", resp.Errors)
				}

			} else if l == 1 {

				// change existing record and push
				newRecord := records[0]
				newRecord.Content = newIP

				if err := cfAPI.UpdateDNSRecord(context.Background(), zoneID, records[0].ID, newRecord); err != nil {
					return err
				}

			} else {
				return errors.New("got more than one record")
			}

		}

	}

	return nil

}

func UpdateIPv4(newIP *externalip.IPv4) error {
	return update(v4, newIP.Addr)
}

func UpdateIPv6(newIP *externalip.IPv6) error {
	return update(v6, newIP.Addr)
}
