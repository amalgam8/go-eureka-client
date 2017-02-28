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
package go_eureka_client

import (
 	"time"
	"errors"
	"math/rand"
	"fmt"
)

type Config struct {
	ConnectTimeoutSeconds time.Duration // default 10s
	UseDNSForServiceUrls  bool          // default false
	DNSDiscoveryZone      string
	ServerDNSName         string
	ServiceUrls           map[string][]string // map from Zone to array of server Urls
	ServerPort            int           // default 8080
	PreferSameZone        bool          // default false
	RetriesCount          int           // default 3
	UseJSON               bool          // default false (means XML)
}

//createUrlsList creates an array of urls from the ServiceUrls map according to the following settings:
func (c *Config) createUrlsList() ( []string ,error ) {
	if c.ServiceUrls == nil {
		return nil, errors.New("$EUREKA_SERVICE_URLS must be defined")
	}
	urls := []string{}
	indMap := map[int]string{}
	i := 0
	for k, _ := range c.ServiceUrls {
		indMap[i] = k
	}

	if c.PreferSameZone == false &&  c.UseDNSForServiceUrls == false {
		zonesPerm := rand.Perm(len(c.ServiceUrls))
		for _, v := range zonesPerm {
			urlsOfZone := c.ServiceUrls[ indMap[v] ]
			urlsPerm := rand.Perm(len(urlsOfZone))

			for _, p := range urlsPerm {
				urls = append(urls,urlsOfZone[p])
			}
		}
		return urls, nil
	} else if  c.PreferSameZone == true &&  c.UseDNSForServiceUrls == false {
		urlsOfPreferredZone := c.ServiceUrls[ c.DNSDiscoveryZone ]
		urlsPerm := rand.Perm(len(urlsOfPreferredZone))

		for _, p := range urlsPerm {
			urls = append(urls,urlsOfPreferredZone[p])
		}
		zonesPerm := rand.Perm(len(c.ServiceUrls))
		for _, v := range zonesPerm {
			if indMap[v] != c.DNSDiscoveryZone {
				urlsOfZone := c.ServiceUrls[ indMap[v] ]
				urlsPerm := rand.Perm(len(urlsOfZone))
				for _, p := range urlsPerm {
					urls = append(urls,urlsOfZone[p])
				}
			}

		}
		return urls, nil
	} else {
		return nil, fmt.Errorf("DNS discovery is not supported")
	}
}