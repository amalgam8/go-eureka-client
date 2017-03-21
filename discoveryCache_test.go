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

/*
import (
	"testing"
	"time"
	//"log"

	"fmt"
	"log"
)
*/

// mockHandler implements the interface InstanceEventHandler which is invoked when there is a change in the cache.
type mockHandler string

func (m *mockHandler) OnAdd(inst *Instance) {

}

func (m *mockHandler) OnUpdate(oldInst, newInst *Instance) {

}
func (m *mockHandler) OnDelete(int *Instance) {

}

/*
func TestNewDiscoveryCache(t *testing.T){
	// Create configuration for the discovery Cache :
	conf := &Config{
		ConnectTimeoutSeconds: 10 * time.Second,
		UseDNSForServiceUrls :  false, // default false
		ServiceUrls :          map[string][]string{"eureka" : []string{"http://172.17.0.2:8080/eureka/v2/"} },
		ServerPort  :          8080, // default 8080
		PreferSameZone:        false, // default false
		RetriesCount   :       3, // default 3
		UseJSON  :             true, // default false (means XML)
	}
	var m mockHandler = "test"
	discoveryCache, e := NewDiscoveryCache(conf,7*time.Second, &m)
	if e != nil {
		t.Errorf("error = %v", e)
	}
	var stopch chan struct{}
	discoveryCache.Run(stopch)
	log.Printf("New discovery test started to run, sleeping for 15 seconds...")
	time.Sleep(15*time.Second)
	log.Printf("Invoking GetApplications method:")
	apps, _ := discoveryCache.GetApplications()
	log.Printf("Number of Apps in Discovery Cache is %v", len(apps))
	log.Printf("Name Of App in the Discovery Cache is: %s", apps[0].Name)
	log.Printf("Instances in the application details %v", apps[0].Instances)
	if (len(apps) != 1) {
		t.Errorf("1 app should be registred on server")
	}
	insts := apps[0].Instances
	if (len(insts) != 1){
		t.Error("should only have 1 instance")
	}
	inst := insts[0]
	log.Printf("Success testing GetApplication method")
	log.Printf("Invoking GetInstancesBySecVip method:")
	if inst.SecVIPAddr == "" {
		log.Printf("Insctance secVip address is empty" )
	}
	_,e = discoveryCache.GetInstancesBySecVip(inst.SecVIPAddr)
	if e.Error() != fmt.Errorf("vipAddress  %s not found",inst.SecVIPAddr ).Error() {
		t.Errorf("Failure, svip address shouldn't exist, error returned: %v\n desired error: %v",e, fmt.Errorf("vipAddress  %s not found",inst.SecVIPAddr ))
	}
	fetchedVipInst,e := discoveryCache.GetInstancesByVip(inst.VIPAddr)
	if e !=nil {
		t.Errorf("Failure, %v",e)
	}
	if len(fetchedVipInst) != 1 {
		t.Errorf("should only have 1 instance")
	}
	vipInst := fetchedVipInst[0]
	if vipInst.ID != inst.ID {
		t.Errorf("ids of instances not equal")
	}
	log.Printf("inst id  = %s", vipInst.ID)
}
*/
