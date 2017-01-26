package go_eureka_client

import (
	"time"
	"testing"
)

func TestNewRegistrator(t *testing.T) {
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

	registrator, e := NewRegistrator(conf, &m)
	if e != nil {
		t.Errorf("error = %v", e)
	}
	inst := &Instance{
		            Application :  "new_app ",
		            HostName    :  "my_host",
	}

	registrator.Register(inst)

	/*
	app, e := discovery.GetApplication("HELLO-NETFLIX-OSS")
	if e != nil {
		t.Errorf("Failed to get aplication. error : %v", e)
	}
	if app.Name != "HELLO-NETFLIX-OSS" {
		t.Errorf("Unexpected app name: %s", app.Name)
	}
	insts := app.Instances
	if len(insts) != 1 {
		t.Errorf("num of instances should be 1. instead : %d", len(insts))
	}
	*/
}