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


/*
import (
	"time"
	"testing"
	"log"
	"encoding/json"
)
*/
/*
func TestNewRegistrator(t *testing.T) {
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
	registrator, e := NewRegistrator(conf, &m)
	if e != nil {
		t.Errorf("error = %v", e)
	}
	inst := &Instance{
		            Application :  "NEW_APP",
		            HostName    :  "my_host",
		            Status:        "UP",
		            Datacenter: &DatacenterInfo{Class:"class", Name:"MyOwn",Metadata:DatacenterMetadata{  }},
	}
	e = registrator.Register(inst)
	if e != nil {
		t.Errorf("Failed to register instance. error: %v", e)
	}
	// Try to get the instance from server :
	discovery, e := NewDiscovery(conf, &m)
	if e != nil {
		t.Errorf("error = %v", e)
	}
	time.Sleep(20*time.Second)
	app, e := discovery.GetApplication("NEW_APP")
	if e != nil {
		t.Errorf("Failed to get aplication. error : %v", e)
	}
	if app.Name != "NEW_APP" {
		t.Errorf("Unexpected app name: %s", app.Name)
	}
	insts := app.Instances
	if len(insts) != 1 {
		t.Errorf("Number of instances should be 1. instead : %d", len(insts))
	}
	pulled_inst := insts[0]
	if inst.HostName != "my_host" {
		t.Errorf("Host name of the instance should be my_host, instead it is : %v", pulled_inst.HostName)
	}
	//Test heartbeat :
	e = registrator.Heartbeat(pulled_inst)
	if e != nil {
		t.Errorf("Failed to send heartbeat. error : %v", e)
	}
	// Test status change :
	if pulled_inst.Status != "UP" {
		t.Errorf("Status should be up, instead it is %v", pulled_inst.Status)
	}
	e = registrator.SetStatus(pulled_inst,UNKNOWN)
	if e != nil {
		t.Errorf("Failed to send status change request. error : %v", e)
	}
	time.Sleep(30*time.Second)
	app2, _  := discovery.GetApplication("NEW_APP")
	insts2 := app2.Instances
	pulled_inst2 := insts2[0]
	log.Printf("New status is : %v",pulled_inst2.Status)
	if pulled_inst2.Status != "UNKNOWN" {
		t.Errorf("status didn't change, status is %v",app2)
	}

	// Test metadata change :
	e = registrator.SetMetadataKey(pulled_inst,"mykey", "myvalue")
	if e != nil {
		t.Errorf("Failed to send metadata change request. error : %v", e)
	}
	time.Sleep(30*time.Second)
	metaChangeIsnt,_:= discovery.GetInstance("NEW_APP","my_host")
	metadata := metaChangeIsnt.Metadata
	var metadataMap map[string]string
	json.Unmarshal(metadata, &metadataMap)
	log.Printf("metad: %v", metadataMap["mykey"] )
	if metadataMap["mykey"] != "myvalue" {
		t.Errorf("Failed to change metadata. metadata is : %v", metadataMap)
	}
	// Test deregister :
	e = registrator.Deregister(pulled_inst)
	if e != nil {
		t.Errorf("Failed to deregister. error : %v", e)
	}


}
*/