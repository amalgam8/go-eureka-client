// Copyright 2016 IBM Corporation
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

//Package goEurekaClient Implements a go client that interacts with a eureka server
package goEurekaClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Config struct defines configurations of the eureka client in order to interact with the server.
type Config struct {
	ConnectTimeoutSeconds time.Duration       `json:"connectionTimeoutSeconds"` // default 10s
	UseDNSForServiceUrls  bool                `json:"useDNSForServiceUrls"`     // default false
	DNSDiscoveryZone      string              `json:"dnsDiscoveryZone"`
	ServerDNSName         string              `json:"serverDNSName"`
	ServiceUrls           map[string][]string `json:"serviceUrls"`    // map from Zone to array of server Urls
	ServerPort            int                 `json:"serverPort"`     // default 8080
	PreferSameZone        bool                `json:"preferSameZone"` // default false
	RetriesCount          int                 `json:"retriesCount"`   // default 3
	UseJSON               bool                `json:"useJSON"`        // default false (means XML)
}

// NewConfigFromFile reads JSON data from file and creates from it a config object.
func NewConfigFromFile(fileName string) (*Config, error) {
	configFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("opening config file %v", err.Error())
	}
	var conf *Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(conf); err != nil {
		return nil, fmt.Errorf("parsing config file %v", err.Error())
	}
	// For debugging :
	//fmt.Printf("%v %s %s", conf.ServerPort, conf.ServerPort, conf.RetriesCount)
	return conf, nil
}

//createUrlsList creates an array of urls from the ServiceUrls map according to the following settings:
func (c *Config) createUrlsList() ([]string, error) {
	if c.ServiceUrls == nil {
		return nil, errors.New("$EUREKA_SERVICE_URLS must be defined")
	}
	urls := []string{}
	indMap := map[int]string{}
	i := 0
	for k := range c.ServiceUrls {
		indMap[i] = k
	}

	if c.PreferSameZone == false && c.UseDNSForServiceUrls == false {
		zonesPerm := rand.Perm(len(c.ServiceUrls))
		for _, v := range zonesPerm {
			urlsOfZone := c.ServiceUrls[indMap[v]]
			urlsPerm := rand.Perm(len(urlsOfZone))

			for _, p := range urlsPerm {
				urls = append(urls, urlsOfZone[p])
			}
		}
		return urls, nil
	} else if c.PreferSameZone == true && c.UseDNSForServiceUrls == false {
		urlsOfPreferredZone := c.ServiceUrls[c.DNSDiscoveryZone]
		urlsPerm := rand.Perm(len(urlsOfPreferredZone))

		for _, p := range urlsPerm {
			urls = append(urls, urlsOfPreferredZone[p])
		}
		zonesPerm := rand.Perm(len(c.ServiceUrls))
		for _, v := range zonesPerm {
			if indMap[v] != c.DNSDiscoveryZone {
				urlsOfZone := c.ServiceUrls[indMap[v]]
				urlsPerm := rand.Perm(len(urlsOfZone))
				for _, p := range urlsPerm {
					urls = append(urls, urlsOfZone[p])
				}
			}

		}
		return urls, nil
	} else {
		return nil, fmt.Errorf("DNS discovery is not supported")
	}
}
