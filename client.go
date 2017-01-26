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
"crypto/tls"
"encoding/json"
"fmt"
"io/ioutil"
"log"
"net/http"
"net/url"
"sort"
"strings"
"sync"
"time"
)

const (
	hashcodeDelimiter = "_"
	actionAdded       = "ADDED"
	actionModified    = "MODIFIED"
	actionDeleted     = "DELETED"
)

type client struct {
	sync.Mutex
	httpClient *http.Client
	eurekaURLs []string
	dictionary dictionary
	versionDelta int64
	handler InstanceEventHandler
}

func newClient(config *Config, handler InstanceEventHandler) (*client, error) {
	eurekaURLs ,err:=config. createUrlsList()
	if eurekaURLs == nil {
		return nil, err
	}

	var httpsRequired bool

	urls := make([]string, len(eurekaURLs))
	for i, eu := range eurekaURLs {
		for strings.HasSuffix(eu, "/") {
			eu = strings.TrimSuffix(eu, "/")
		}

		u, err := url.Parse(eu)
		if err != nil {
			return nil, err
		}

		if u.Scheme == "https" {
			httpsRequired = true
		}

		urls[i] = eu
	}

	hc := &http.Client{
		Timeout: config.ConnectTimeoutSeconds,
	}

	if httpsRequired {
		hc.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	cl := &client{
		httpClient: hc,
		eurekaURLs: urls,
		handler: handler,
	}
	return cl, nil
}

func (cl *client) run(pollInterval time.Duration ,stopCh chan struct{}) {
	cl.refresh(cl.handler)

	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case <-ticker.C:
			cl.refresh(cl.handler)
		case <-stopCh:
			//close(cl.)
			return
		}
	}
}

func (cl *client) refresh(handler InstanceEventHandler) {
	log.Printf("inside refresh , client dict conatins : svipindex %v\n " +
		"vipIndex : %v\n appIndex: %v\n", cl.dictionary.svipIndex,cl.dictionary.vipIndex,cl.dictionary.appNameIndex)
	for k,_ := range cl.dictionary.svipIndex {
		log.Printf("svip address: %s\n",k)
		insts := cl.dictionary.svipIndex[k]
		for _,inst := range insts {
			log.Printf("isnt svip is %s\n", inst.SecVIPAddr)
			log.Printf("inst details = %v\n", inst)
		}
	}

	var dict dictionary
	// diff is a map of key : instance_id, value: *instance
	var diff map[string]*Instance
	// If this is the 1st time then we need to retrieve the full registry,
	// otherwise a delta could be sufficient
	if  cl.dictionary.isEmpty() == false {
		// not first time :

		dict, diff = cl.fetchDelta()
	}

	if dict.appNameIndex == nil && dict.vipIndex == nil && dict.svipIndex == nil {
		// This means first time :

		fetchdDict, err  := cl.fetchAll()
		if err != nil {
			// TODO: what message to report?
			return
		}

		dict = fetchdDict
		diff = cl.populateDiff(dict)
		cl.versionDelta = 0
	}

	old_dict := cl.dictionary.copyDictionary()
	if dict.appNameIndex != nil || dict.vipIndex != nil || dict.svipIndex!= nil {
		cl.Lock()
		cl.dictionary.vipIndex = dict.vipIndex
		cl.dictionary.appNameIndex = dict.appNameIndex
		cl.dictionary.svipIndex = dict.svipIndex
		cl.Unlock()
	}

	// Send notifications
	if len(diff) > 0 {
		for name := range diff {
			if diff[name].ActionType == actionAdded {
				cl.handler.OnAdd(diff[name])
			} else if diff[name].ActionType == actionModified{
				// TODO: need to find old object
				old_obj := old_dict.vipIndex[diff[name].VIPAddr][name]
				cl.handler.OnUpdate(old_obj,diff[name])
			} else if diff[name].ActionType == actionDeleted {
				cl.handler.OnDelete(diff[name])
			}
		}
	}
}

func (cl *client) fetchAll() (dictionary, error) {
	apps, err := cl.fetchApps("apps")
	if err != nil {
		log.Printf("Faild to update full registry. %s\n", err)
		return cl.dictionary, err
	}

	dict := newDictionary()
	if apps != nil && apps.Application != nil {
		for _, app := range apps.Application {
			for _, inst := range app.Instances {
				id, err := resolveInstanceID(inst)
				if err != nil {
					log.Printf("Failed to resolve instance ID. error: %s\n", err)
					continue
				}
				inst.ID = id
				if inst.VIPAddr != "" {
					instances := dict.vipIndex[inst.VIPAddr]
					if instances == nil {
						instances = map[string]*Instance{}

						dict.vipIndex[inst.VIPAddr] = instances
						dict.vipIndex[app.Name][id] = inst
					}
				}

				if inst.SecVIPAddr != "" {
					if dict.svipIndex[inst.SecVIPAddr] == nil {
						dict.svipIndex[inst.SecVIPAddr] =  map[string]*Instance{}
						dict.svipIndex[app.Name][id] = inst
					}
				}
				if inst.Application != "" {
					if dict.appNameIndex[app.Name] == nil {
						dict.appNameIndex[app.Name] = map[string]*Instance{}
						dict.appNameIndex[app.Name][id] = inst
					}
				}
			}
		}
	}

	hashcode := calculateHashcode(dict.vipIndex)
	log.Printf("A full fetch completed. %s\n", hashcode)
	log.Printf("Returning the dict : %v", dict.appNameIndex)
	return dict, nil
}

func (cl *client) fetchDelta() (dictionary, map[string]*Instance) {
	apps, err := cl.fetchApps("apps/delta")
	if err != nil {
		log.Printf("Faild to update delta. %s\n", err)

		return dictionary{}, nil
	}

	if apps == nil || apps.VersionDelta == -1 {
		log.Println("Delta update is not supported")
		return dictionary{}, nil
	}

	diff := map[string]*Instance{}

	// If we have the latest version, no need to do anything
	if apps.VersionDelta == cl.versionDelta {
		log.Printf("Delta update was skipped, because we have the latest version (%d)", apps.VersionDelta)
		return cl.dictionary, diff
	}

	dict := cl.dictionary.copyDictionary()
	var updated, deleted int
	for _, app := range apps.Application {
		for _, inst := range app.Instances {
			id, err := resolveInstanceID(inst)
			if err != nil {
				log.Printf("Failed to resolve instance ID. error: %s\n", err)
				return dictionary{}, nil
			}

			inst.ID = id
			switch inst.ActionType {
			case actionDeleted:
				dict.onDelete(inst,id, app )
				deleted++
			case actionAdded:
				dict.onAdd(inst,id, app)
				updated++
			case  actionModified:
				dict.onUpdate(inst,id, app)
				updated++
			default:
				log.Printf("Unknown ActionType %s for instance %+v\n", inst.ActionType, inst)
			}

			diff[inst.ID] = inst
		}
	}

	// Calculate the new hashcode and compare it to the server
	hashcode := calculateHashcode(dict.vipIndex)
	if apps.Hashcode != hashcode {
		log.Printf("Failed to update delta (local: %s, remote %s). A full update is required\n", hashcode, apps.Hashcode)
		return dictionary{}, nil
	}

	cl.versionDelta = apps.VersionDelta
	log.Printf("Delta update completed successfully (updated: %d, deleted: %d, version: %d)\n", updated, deleted, apps.VersionDelta)

	return dict, diff
}

func (cl *client) fetchApps(path string) (*Applications, error) {
	var err error

	for _, eurl := range cl.eurekaURLs {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", eurl, path), nil)
		req.Header.Set("Accept", "application/json")

		resp, err2 := cl.httpClient.Do(req)
		if err2 != nil {

			err = err2
			continue
		}

		defer resp.Body.Close()

		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {

			err = err2
			continue
		}

		var appsList ApplicationsList
		err2 = json.Unmarshal(body, &appsList)
		if err2 != nil {

			//log.Println(body)
			err = err2
			continue
		}

		return appsList.Applications, nil
	}

	return nil, err
}

func calculateHashcode(dict map[string]map[string]*Instance) string {
	var hashcode string

	if len(dict) == 0 {
		return hashcode
	}

	hashMap := map[string]uint32{}
	for _, insts := range dict {
		for _, inst := range insts {
			if count, ok := hashMap[inst.Status]; !ok {
				hashMap[inst.Status] = 1
			} else {
				hashMap[inst.Status] = count + 1
			}
		}
	}

	var keys []string
	for status := range hashMap {
		keys = append(keys, status)
	}
	sort.Strings(keys)

	for _, status := range keys {
		count := hashMap[status]
		hashcode = hashcode + fmt.Sprintf("%s%s%d%s", status, hashcodeDelimiter, count, hashcodeDelimiter)
	}

	return hashcode
}

func (cl *client) populateDiff(dict dictionary) map[string]*Instance {
	if dict.vipIndex == nil && dict.svipIndex == nil && dict.appNameIndex == nil {
		return nil
	}

	cl.Lock()
	defer cl.Unlock()

	diff := map[string]*Instance{}
	// Scan the new dictionary and look for changes
	for vip, newInsts := range dict.vipIndex {
		if srcInsts, ok := cl.dictionary.vipIndex[vip]; ok {
			for id, newInst := range newInsts {
				if srcInst, ok := srcInsts[id]; ok {
					if newInst.Status != srcInst.Status {
						diff[id] = newInst
					}
				} else {
					diff[id] = newInst
				}
			}
		} else {
			for id, newInst := range newInsts {
				diff[id] = newInst
			}
			//diff[vip] = struct{}{}
		}
	}

	// Scan the source dictionary and look for deleted services
	for vip := range cl.dictionary.vipIndex {
		if _, ok := dict.vipIndex[vip]; !ok {

			for id, delInsts := range  cl.dictionary.vipIndex[vip] {
				diff[id] = delInsts
			}
		} else {
			for id, inst := range cl.dictionary.vipIndex[vip] {
				if _,ok := dict.vipIndex[vip][id] ; !ok {
					diff[id] = inst
				}
			}
		}
	}

	return diff
}