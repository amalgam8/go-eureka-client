package go_eureka_client

import (
	"testing"
	"time"
	//"log"

	"fmt"
	"log"
)

//http://172.17.0.2:8080/eureka/
type mockHandler string

func (m *mockHandler) OnAdd(inst *Instance){

}

func (m *mockHandler) OnUpdate(oldInst, newInst *Instance) {

}
func (m* mockHandler )OnDelete(int *Instance){

}
func TestNewDiscoveryCache(t *testing.T){
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
		//urls, _ := conf.createUrlsList()

		//t.Errorf("urls = %v", urls)
		var m mockHandler = "test"

	discoveryCache, e := NewDiscoveryCache(conf,7*time.Second, &m)
	if e != nil {
		t.Errorf("error = %v", e)
	}
	var stopch chan struct{}
	discoveryCache.Run(stopch)
	//log.Printf("started running cache")
	time.Sleep(15*time.Second)
	apps, _ := discoveryCache.GetApplications()
	//log.Printf("len of apps = %v", len(apps))
	//log.Printf("Name Of app: %s", apps[0].Name)
	//log.Printf("inst details %v", apps[0].Instances)
	if (len(apps) != 1) {
		t.Errorf("1 app should be registred on server")
	}
	insts := apps[0].Instances
	if (len(insts) != 1){
		t.Error("shoudl only have 1 instance")
	}
	inst := insts[0]

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
