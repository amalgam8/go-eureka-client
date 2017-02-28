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

// Pre requests: a eureka server available on http://172.17.0.2:8080/eureka/v2/

import  (
	"testing"
	"time"
)

// Global test variables :
var testRegistrator Registrator
var testDiscovery Discovery
var testdiscoveryCache DiscoveryCache
var quitHeartBeats  chan struct{}
var errChan chan error
var stopdiscoveryCacheChan chan struct{}
var instances  []*Instance
// function setupTest creates discovery client, discovery cache client, registrator client, and registers apps to the server.
func setupTest(t *testing.T) error{
	quitHeartBeats = make( chan struct{})
	errChan = make(chan error)
	stopdiscoveryCacheChan = make(chan  struct{})
	e := createClients()
	if e != nil {
		return e
	}
	return registerInstances(t)

}

// function createClients creates discovery client, discovery cache client and registrator client.
func createClients() error {
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
	var e error
	testRegistrator, e = NewRegistrator(conf,&m)
	if e != nil {
		return e
	}
	testDiscovery ,e = NewDiscovery(conf,&m)
	if e != nil {
		return e
	}
	testdiscoveryCache,e = NewDiscoveryCache(conf,7*time.Second,&m)
	if e != nil {
		return e
	}
	testdiscoveryCache.Run(stopdiscoveryCacheChan)
	return e
}

// function createInstance, create anew instance with the minimal requirements in order to be registerd in the server :
func createInstance(hostName, appName,vipAddr, svipAddr string) *Instance {
	return &Instance{
		Application :  appName,
		HostName    :  hostName,
		Status:        "UP",
		Datacenter: &DatacenterInfo{Class:"class", Name:"MyOwn",Metadata:DatacenterMetadata{  }},
		VIPAddr: vipAddr,
		SecVIPAddr: svipAddr,
	}
}

// function registerInstances registers 6 instances to the server.
func registerInstances(t *testing.T) error {
	instances = []*Instance{createInstance("inst1", "app1","vip1","svip1") ,
		createInstance("inst2", "app1","vip1","svip2") ,
		createInstance("inst3", "app1","vip2","svip1") ,
		createInstance("inst1", "app2","vip2","svip2"),
		createInstance("inst2", "app2","vip2","svip1"),
		createInstance("inst3", "app2","vip1","svip2")}
	for _, inst := range instances {
		e := testRegistrator.Register(inst)
		if e!= nil {
			return e
		}
	}
	go sendHeartBeats(instances,t)
	return nil
}

func sendHeartBeats(instances []*Instance,t *testing.T) {
	for {
		select {
		case <- quitHeartBeats :
			t.Log("quit channel recieived, stop sending heartbeats")
			return
		default:
			for _, inst := range instances {
				e := testRegistrator.Heartbeat(inst)
				if e!= nil {
					errChan <- e
					return
				}
			}
			t.Log("heartbeats sent")
			time.Sleep(30*time.Second)

		}
	}
}

func tearDownTest(t *testing.T) error {
	var stop struct{}

	// Stop the go routine which sends heat-beats to instances.
	quitHeartBeats <- stop

	// Stop the discovery cache poll intervals :
	stopdiscoveryCacheChan <- stop

	// De-register all instances from server :
	for _, inst := range instances {
		e := testRegistrator.Deregister(inst)
		if e!= nil {
			t.Log("something went wrong during dregisteriation...")
			return e
		}
	}
	t.Log("Successfully de-registerd all instances from eureka server. sleep 30 seconds before return from" +
		" function")
	time.Sleep(30*time.Second)
	return nil
}


func  TestGetInstancesByVipAddress(t *testing.T) {
	t.Log("************************************************************************************************")
	t.Log("Calling Setup test...")
	e := setupTest(t)
	if e != nil {
		// send  quit channel :
		quitHeartBeats <- struct{}{}
		stopdiscoveryCacheChan <- struct{}{}
		t.Errorf("Error during setup : %v",e)
		return
	}
	t.Log("Successfully initialized discovery, discovery cache and registrator client. sleeping 30 seconds.")
	t.Log("************************************************************************************************")
	time.Sleep(30*time.Second)

	t.Log("Testing  Get instances by vip address... ")
	dicoveryInstances, err := testDiscovery.GetInstancesByVip("vip1")

	if err != nil {
		quitHeartBeats <- struct{}{}
		t.Errorf("Error while trying to get instances from server into discovery: %v ", err)
		return
	}
	if len(dicoveryInstances) != 3 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Should have gotten 3 instances with vip addrees vip1, instead got %d",
			len(dicoveryInstances))
		return

	}
	numOfInstsFromApp1 := 0
	numOfInstsFromApp2 := 0
	numOfInstsFromInst1 := 0
	numOfInstsFromInst2 := 0
	numOfInstsFromsVip1 := 0
	numOfInstsFromsVip2:= 0
	for _, inst := range dicoveryInstances {

		if inst.Application == "APP1" {
			numOfInstsFromApp1++
		}
		if inst.Application == "APP2" {
			numOfInstsFromApp2++
		}
		if inst.HostName == "inst1" {
			numOfInstsFromInst1++
		}
		if inst.HostName == "inst2" {
			numOfInstsFromInst2++
		}
		if inst.SecVIPAddr == "svip1" {
			numOfInstsFromsVip1++
		}
		if inst.SecVIPAddr == "svip2" {
			numOfInstsFromsVip2++
		}
	}
	if numOfInstsFromApp1 !=2 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to app1 should be 2, insted : %d",numOfInstsFromApp1 )
		return
	} else if numOfInstsFromApp2 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to app2 should be 1, insted : %d",numOfInstsFromApp2 )
		return
	} else if numOfInstsFromInst1 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to id inst1 should be 1, insted : %d",numOfInstsFromInst1 )
		return
	} else if numOfInstsFromInst2 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to id inst2 should be 1, insted : %d",numOfInstsFromInst2 )
		return
	} else if numOfInstsFromsVip1 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to secure vip1 should be 1, insted : %d",numOfInstsFromsVip1 )
		return
	} else if numOfInstsFromsVip2 != 2 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to secure vip2 should be 2, insted : %d",numOfInstsFromsVip2 )
		return
	}

	discoveryCacheInstances , err := testdiscoveryCache.GetInstancesByVip("vip1")

	if err != nil {
		quitHeartBeats <- struct{}{}
		t.Errorf("Error while trying to get instances from discovery cache: %v ", err)
		return
	}

	if len(discoveryCacheInstances) != 3 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Should have gotten 3 instances with vip addrees vip1, instead got %d",
			len(discoveryCacheInstances))
		return
	}
	numOfInstsFromApp1 = 0
	numOfInstsFromApp2 = 0
	numOfInstsFromInst1 = 0
	numOfInstsFromInst2 = 0
	numOfInstsFromsVip1 = 0
	numOfInstsFromsVip2 = 0
	for _, inst := range discoveryCacheInstances {
		if inst.Application == "APP1" {
			numOfInstsFromApp1++
		}
		if inst.Application == "APP2" {
			numOfInstsFromApp2++
		}
		if inst.HostName == "inst1" {
			numOfInstsFromInst1++
		}
		if inst.HostName == "inst2" {
			numOfInstsFromInst2++
		}
		if inst.SecVIPAddr == "svip1" {
			numOfInstsFromsVip1++
		}
		if inst.SecVIPAddr == "svip2" {
			numOfInstsFromsVip2++
		}
	}
	if numOfInstsFromApp1 !=2 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to app1 should be 2, insted : %d",numOfInstsFromApp1 )
	} else if numOfInstsFromApp2 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to app2 should be 1, insted : %d",numOfInstsFromApp2 )
	} else if numOfInstsFromInst1 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to id inst1 should be 1, insted : %d",numOfInstsFromInst1 )
	} else if numOfInstsFromInst2 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to id inst2 should be 1, insted : %d",numOfInstsFromInst2 )
	} else if numOfInstsFromsVip1 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to secure vip1 should be 1, insted : %d",numOfInstsFromsVip1 )
	} else if numOfInstsFromsVip2 != 2 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to secure vip2 should be 2, insted : %d",numOfInstsFromsVip2 )
	}
	t.Log("Successfully passed test GetInstancesByVipAddress..." )

	// clean :
	t.Log("************************************************************************************************")
	t.Log("Calling tear dorwn function...")
	e = tearDownTest(t)
	if e != nil {
		t.Errorf("Error during tear down :  : %v",e)
	}
	t.Log("Tear down cpmpleted succesfully")
	t.Log("************************************************************************************************")
	time.Sleep(30*time.Second)

	t.Log("Test Ended Successfully")

}

func TestGetInstancesBySVipAddress(t *testing.T) {
	e := setupTest(t)
	if e != nil {
		// send  quit channel :
		quitHeartBeats <- struct{}{}
		stopdiscoveryCacheChan <- struct{}{}
		t.Errorf("Error during setup : %v",e)
		return
	}
	t.Log("Successfully initialized discovery, discovery cache and registrator client. sleeping 30 seconds.")
	t.Log("************************************************************************************************")
	time.Sleep(30*time.Second)

	t.Log("Testing  Get instances by secure vip address... ")
	dicoveryInstances, err := testDiscovery.GetInstancesBySecVip("svip2")

	if err != nil {
		quitHeartBeats <- struct{}{}
		t.Errorf("Error while trying to get instances from server into discovery: %v ", err)
		return
	}

	if len(dicoveryInstances) != 3 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Should have gotten 3 instances with svip addrees svip2, instead got %d",len(dicoveryInstances))
		return
	}
	numOfInstsFromApp1 := 0
	numOfInstsFromApp2 := 0
	numOfInstsFromInst1 := 0
	numOfInstsFromInst2 := 0
	numOfInstsFromVip1 := 0
	numOfInstsFromsip2:= 0
	for _, inst := range dicoveryInstances {

		if inst.Application == "APP1" {
			numOfInstsFromApp1++
		}
		if inst.Application == "APP2" {
			numOfInstsFromApp2++
		}
		if inst.HostName == "inst1" {
			numOfInstsFromInst1++
		}
		if inst.HostName == "inst2" {
			numOfInstsFromInst2++
		}
		if inst.VIPAddr == "vip1" {
			numOfInstsFromVip1++
		}
		if inst.VIPAddr == "vip2" {
			numOfInstsFromsip2++
		}
	}
	if numOfInstsFromApp1 !=1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to app1 should be 1, insted : %d",numOfInstsFromApp1 )
		return
	} else if numOfInstsFromApp2 != 2 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to app2 should be 2, insted : %d",numOfInstsFromApp2 )
		return
	} else if numOfInstsFromInst1 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to id inst1 should be 1, insted : %d",numOfInstsFromInst1 )
		return
	} else if numOfInstsFromInst2 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to id inst2 should be 1, insted : %d",numOfInstsFromInst2 )
		return
	} else if numOfInstsFromVip1 != 2 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to vip1 should be 1, insted : %d",numOfInstsFromVip1 )
		return
	} else if numOfInstsFromsip2 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to vip2 should be 2, insted : %d",numOfInstsFromsip2 )
		return
	}

	discoveryCacheInstances , err := testdiscoveryCache.GetInstancesBySecVip("svip2")

	if err != nil {
		quitHeartBeats <- struct{}{}
		t.Errorf("Error while trying to get instances from discovery cache: %v ", err)
		return
	}

	if len(discoveryCacheInstances) != 3 {
		t.Errorf("Should have gotten 3 instances with svip addrees svip2, instead got %d",len(discoveryCacheInstances))
		return
	}
	numOfInstsFromApp1 = 0
	numOfInstsFromApp2 = 0
	numOfInstsFromInst1 = 0
	numOfInstsFromInst2 = 0
	numOfInstsFromVip1 = 0
	numOfInstsFromsip2 = 0
	for _, inst := range discoveryCacheInstances {
		if inst.Application == "APP1" {
			numOfInstsFromApp1++
		}
		if inst.Application == "APP2" {
			numOfInstsFromApp2++
		}
		if inst.HostName == "inst1" {
			numOfInstsFromInst1++
		}
		if inst.HostName == "inst2" {
			numOfInstsFromInst2++
		}
		if inst.VIPAddr == "vip1" {
			numOfInstsFromVip1++
		}
		if inst.VIPAddr == "vip2" {
			numOfInstsFromsip2++
		}
	}
	if numOfInstsFromApp1 !=1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to app1 should be 1, insted : %d",numOfInstsFromApp1 )
		return
	} else if numOfInstsFromApp2 != 2 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to app2 should be 2, insted : %d",numOfInstsFromApp2 )
		return
	} else if numOfInstsFromInst1 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to id inst1 should be 1, insted : %d",numOfInstsFromInst1 )
		return
	} else if numOfInstsFromInst2 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to id inst2 should be 1, insted : %d",numOfInstsFromInst2 )
		return
	} else if numOfInstsFromVip1 != 2 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to vip1 should be 1, insted : %d",numOfInstsFromVip1 )
		return
	} else if numOfInstsFromsip2 != 1 {
		quitHeartBeats <- struct{}{}
		t.Errorf("Numer of instances belonging to vip2 should be 2, insted : %d",numOfInstsFromsip2 )
		return
	}
	t.Log("Successfully passed test GetInstancesBySVipAddress..." )


	// clean :
	t.Log("************************************************************************************************")
	t.Log("Calling tear dorwn function...")
	e = tearDownTest(t)
	if e != nil {
		t.Errorf("Error during tear down :  : %v",e)
		return
	}
	t.Log("Tear down cpmpleted succesfully")
	t.Log("************************************************************************************************")
	time.Sleep(30*time.Second)

	t.Log("Test Ended Successfully")
}
