package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type ProxyAction int

const (
	ProxyActivate ProxyAction = iota
	ProxyDeactivate
	ProxyPreserve
)

type Record struct {
	Name       string `json:"name"`
	UpdateIPv4 bool   `json:"update_ipv4"`
	UpdateIPv6 bool   `json:"update_ipv6"`
	ProxyRaw   bool   `json:"proxy"` // Do not use this field, use Proxy instead!
	Proxy      ProxyAction
	Create     bool `json:"create"`
}

func (r *Record) setDefaults(plain map[string]interface{}, defaults *defaults) error {
	if r.Name == "" {
		return errors.New(`missing "records.name" key`)
	}
	if _, ok := plain["update_ipv4"]; !ok {
		r.UpdateIPv4 = defaults.UpdateIPv4Default
	}
	if _, ok := plain["update_ipv6"]; !ok {
		r.UpdateIPv6 = defaults.UpdateIPv6Default
	}
	if _, ok := plain["proxy"]; !ok {
		r.Proxy = ProxyPreserve
	} else {
		if r.ProxyRaw {
			r.Proxy = ProxyActivate
		} else {
			r.Proxy = ProxyDeactivate
		}
	}
	return nil
}

type defaults struct {
	UpdateIPv4Default bool `json:"update_ipv4_default"`
	UpdateIPv6Default bool `json:"update_ipv6_default"`
}

type Config struct {
	APIToken      string               `json:"api_token"`
	CheckInterval int                  `json:"check_interval"`
	CheckTimeout  int                  `json:"check_timeout"`
	Zones         map[string][]*Record `json:"zones"`
	defaults      `json:",inline"`
}

func (c *Config) setDefaults(plain map[string]interface{}) error {
	if c.APIToken == "" {
		return errors.New(`missing "api_token" key`)
	}
	if c.CheckInterval == 0 {
		c.CheckInterval = 60
	}
	log.Printf(`"check_interval" set to %d seconds`, c.CheckInterval)
	if c.CheckTimeout == 0 {
		c.CheckTimeout = 5
	}
	log.Printf(`"check_timeout" set to %d seconds`, c.CheckTimeout)
	if _, ok := plain["update_ipv4_default"]; !ok {
		c.defaults.UpdateIPv4Default = true
	}
	if _, ok := plain["update_ipv6_default"]; !ok {
		c.defaults.UpdateIPv6Default = true
	}
	for zone, records := range c.Zones {
		for i, record := range records {
			// manually unmarshal the plain map to the record
			plainZones, ok := plain["zones"].(map[string]interface{})
			if !ok {
				return errors.New(`key "zones" is not a JSON dictionary`)
			}
			plainZone, ok := plainZones[zone].([]interface{})
			if !ok {
				return fmt.Errorf(`key "zones.%s" is not a JSON array`, zone)
			}
			plainRecord, ok := plainZone[i].(map[string]interface{})
			if !ok {
				return fmt.Errorf(`key "zones.%s.[%d]" is not a JSON dictionary`, zone, i)
			}
			// provide the plain record map to record.setDefaults
			if err := record.setDefaults(plainRecord, &c.defaults); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) Validate() error {
	// TODO
	return nil
}

// ConfiguresIPv6AtAll checks if any record is acutally
// configured to update their AAAA records
func (c *Config) ConfiguresIPv6AtAll() bool {
	atAll := false
outerLoop:
	for _, records := range c.Zones {
		for _, record := range records {
			if record.UpdateIPv6 {
				atAll = true
				break outerLoop
			}
		}
	}
	return atAll
}

func Load(filename string) (*Config, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var plain map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &plain); err != nil {
		return nil, err
	}
	c := Config{}
	if err := json.Unmarshal(jsonBytes, &c); err != nil {
		return nil, err
	}
	if err := c.setDefaults(plain); err != nil {
		return nil, err
	}
	return &c, nil
}
