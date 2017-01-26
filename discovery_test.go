package go_eureka_client

import (
	"testing"
	"time"
	//"log"
	"log"
)



func TestNewDiscovery(t *testing.T){
	//func Main() {
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

	discovery, e := NewDiscovery(conf, &m)
	if e != nil {
		t.Errorf("error = %v", e)
	}

	app, e := discovery.GetApplication("HELLO-NETFLIX-OSS")
	if e != nil {
		t.Errorf("Failed to get aplication. error : %v",e)
	}
	if app.Name != "HELLO-NETFLIX-OSS" {
		t.Errorf("Unexpected app name: %s", app.Name)
	}
	insts := app.Instances
	if len(insts) != 1 {
		t.Errorf("num of instances should be 1. instead : %d", len(insts))
	}
	//inst := insts[0]

	//log.Printf("len of apps = %v", len(apps))
	//log.Printf("Name Of app: %s", apps[0].Name)
	//log.Printf("inst details %v", apps[0].Instances)

	//inst := insts[0]
	/*
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
	*/
}

func TestDiscovery_GetInstance(t *testing.T) {
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

	discovery, e := NewDiscovery(conf, &m)
	if e != nil {
		t.Errorf("error = %v", e)
	}
	inst, err := discovery.GetInstance("HELLO-NETFLIX-OSS","ea03d743d2b4")
	if err != nil {
		t.Errorf("Failed to request an instance: %v", err)
	}
	if inst.Application !=  "HELLO-NETFLIX-OSS" {
		t.Errorf("Wrong name for app of instance: %s",inst.Application)
	}
	if inst.Status != "UP" {
		t.Errorf("Wrong status: %s",inst.Status)
	}

}

func TestDiscovery_GetInstancesByVip(t *testing.T) {
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

	discovery, e := NewDiscovery(conf, &m)
	if e != nil {
		t.Errorf("error = %v", e)
	}
	insts,e := discovery.GetInstancesByVip("HELLO-NETFLIX-OSS")
	if e != nil {
		t.Errorf("Failed to request instances by vip address:  %v", e)
	}
	if (len(insts) != 1) {
		t.Errorf("Should get 1 instance from server")
	}
	inst := insts[0]
	if inst.Application !=  "HELLO-NETFLIX-OSS" {
		t.Errorf("Wrong name for app of instance: %s",inst.Application)
	}
	if inst.Status != "UP" {
		t.Errorf("Wrong status: %s",inst.Status)
	}
	log.Print(inst.HomePage)
}
