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
	"log"
	"fmt"
)
var testRegistrator Registrator
var testDiscovery Discovery
var testdiscoveryCache DiscoveryCache
var quit  chan int
var errChan chan error

// function setupTest creates discovery client, discovery cache client, registrator client, and registers apps to the server.
func setupTest() error{
	e := createClients()
	if e != nil {
		return e
	}
	return registerInstances()

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
func registerInstances() error {
	instances := []*Instance{createInstance("inst1", "app1","vip1","svip1") ,
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
	go sendHeartBeats(instances)
	return nil
}

func sendHeartBeats(instances []*Instance) {
	for {
		select {
		case i := <- quit :
			log.Printf("i: %d",i)
			log.Printf("quit channel recieived, stop sending heartbeats")
			return
		default:
			for _, inst := range instances {
				e := testRegistrator.Heartbeat(inst)
				if e!= nil {
					errChan <- e
					return
				}
			}
			log.Printf("heartbeats sent")
			time.Sleep(30*time.Second)

		}
	}
}

func TestCases(t *testing.T){
	quit = make( chan int)
	errChan = make(chan error)
	e := setupTest()
	if e != nil {
		// send  quit channel :
		quit <- 1
		t.Errorf("Error during setup : %v",e)
		return
	}
	time.Sleep(30*time.Second)
	for {
		select {
		case e = <- errChan :
			t.Errorf("Error during heartbeat sending: %v",e)
		default :
			// Test cases for discovery & discovery cache:
			// 1 . Get instances by vip address:
			log.Printf("Testing  Get instances by vip address... ")
			err := testGetInstancesByVipAddress()
			log.Printf("out")
			if err != nil {
				log.Printf("Sending q")
				quit <- 1
				log.Printf("Sending Quit signal")
				t.Errorf("Error: %v",err)
				return
			}
			log.Printf("Break... ")
			break
			// 2. Get instances by svip address :

			// 3. Get instances by application name :

			// 4. Get instances by their id :

			// 5. Get all instances :
		}
		break
	}
	quit <- 1
	log.Printf("Test Ende Successfully")
}

func  testGetInstancesByVipAddress() error {
	log.Printf("Inside test get instance by vip address: ")
	dicoveryInstances, err := testDiscovery.GetInstancesByVip("vip1")
	if err != nil {
		return err
	}
	if len(dicoveryInstances) != 3 {
		return fmt.Errorf("Should have gotten 3 instances with vip addrees vip1, instead got %d",len(dicoveryInstances))
	}
	numOfInstsFromApp1 := 0
	numOfInstsFromApp2 := 0
	numOfInstsFromInst1 := 0
	numOfInstsFromInst2 := 0
	numOfInstsFromsVip1 := 0
	numOfInstsFromsVip2:= 0
	for _, inst := range dicoveryInstances {
		log.Print(inst.Application)
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
	if numOfInstsFromApp1 !=2 || numOfInstsFromApp2 != 1 || numOfInstsFromInst1 != 1 || numOfInstsFromInst2 != 1 ||
		numOfInstsFromsVip1 != 1 || numOfInstsFromsVip2 != 2 {
		return fmt.Errorf("Unexpected content in the discovery instances. instances: %v",dicoveryInstances)
	}
	log.Printf("Should exit now...")
	return nil
}