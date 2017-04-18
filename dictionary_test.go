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
	"testing"

	"os"
)

var testDict dictionary
var instsApp map[string]*Instance
var instsVip map[string]*Instance
var instsSvip map[string]*Instance

func setup() {
	inst1 := Instance{
		ID:          "inst1",
		Application: "app1 ",
		VIPAddr:     "132.60.60.10",
		SecVIPAddr:  "142.20.1.3",
		Status:      "UP",
		Port:        &Port{Enabled: "True"},
		SecPort:     &Port{Enabled: "True"},
		Datacenter:  &DatacenterInfo{},
		Lease:       &LeaseInfo{},
		ActionType:  hashcodeDelimiter,
	}
	inst2 := Instance{
		ID:          "inst2",
		Application: "app1 ",
		VIPAddr:     "132.60.60.10",
		SecVIPAddr:  "142.20.1.3",
		Status:      "UP",
		Port:        &Port{Enabled: "True"},
		SecPort:     &Port{Enabled: "True"},
		Datacenter:  &DatacenterInfo{},
		Lease:       &LeaseInfo{},
		ActionType:  hashcodeDelimiter,
	}
	inst3 := Instance{
		ID:          "inst3",
		Application: "app1 ",
		VIPAddr:     "132.60.60.10",
		SecVIPAddr:  "142.20.1.3",
		Status:      "UP",
		Port:        &Port{Enabled: "True"},
		SecPort:     &Port{Enabled: "True"},
		Datacenter:  &DatacenterInfo{},
		Lease:       &LeaseInfo{},
		ActionType:  hashcodeDelimiter,
	}
	instsApp = map[string]*Instance{"inst1": &inst1, "inst2": &inst2, "inst3": &inst3}
	instsVip = map[string]*Instance{"inst1": &inst1, "inst2": &inst2, "inst3": &inst3}
	instsSvip = map[string]*Instance{"inst1": &inst1, "inst2": &inst2, "inst3": &inst3}
	appNameIndex := map[string]map[string]*Instance{"app1": instsApp}
	svipIndex := map[string]map[string]*Instance{"142.20.1.3": instsVip}
	vipIndex := map[string]map[string]*Instance{"132.60.60.10": instsSvip}
	testDict = dictionary{appNameIndex: appNameIndex, vipIndex: vipIndex, svipIndex: svipIndex}
}

func shutdown() {}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
func TestOnDelete(t *testing.T) {
	setup()
	t.Logf("len of appInd = %d, len of vipIndex = %d , "+
		"len of svipIndex = %d ", len(testDict.appNameIndex["app1"]), len(testDict.vipIndex["132.60.60.10"]), len(testDict.svipIndex["142.20.1.3"]))

	app1 := Application{Name: "app1"}
	testDict.Delete(instsApp["inst1"], "inst1", &app1)
	if len(testDict.appNameIndex["app1"]) != 2 || len(testDict.vipIndex["132.60.60.10"]) != 2 || len(testDict.svipIndex["142.20.1.3"]) != 2 {
		t.Errorf("Length of dict didn't change after delete. len of appInd = %d, len of vipIndex = %d , "+
			"len of svipIndex = %d ", len(testDict.appNameIndex["app1"]), len(testDict.vipIndex["132.60.60.10"]), len(testDict.svipIndex["142.20.1.3"]))
	}
	if _, ok := testDict.appNameIndex["app1"]["inst1"]; ok {
		t.Error("inst1 wasn't deleted from appNameIndex")
	}
	if _, ok := testDict.vipIndex["132.60.60.10"]["inst1"]; ok {
		t.Error("inst1 wasn't deleted from vipIndex")
	}
	if _, ok := testDict.svipIndex["142.20.1.3"]["inst1"]; ok {
		t.Error("inst1 wasn't deleted from svipIndex")
	}
	testDict.Delete(instsApp["inst2"], "inst2", &app1)
	if len(testDict.appNameIndex["app1"]) != 1 || len(testDict.vipIndex["132.60.60.10"]) != 1 || len(testDict.svipIndex["142.20.1.3"]) != 1 {
		t.Errorf("Length of dict didn't change after delete. len of appInd = %d, len of vipIndex = %d , "+
			"len of svipIndex = %d ", len(testDict.appNameIndex["app1"]), len(testDict.vipIndex["132.60.60.10"]), len(testDict.svipIndex["142.20.1.3"]))
	}
	if _, ok := testDict.appNameIndex["app1"]["inst2"]; ok {
		t.Error("inst2 wasn't deleted from appNameIndex")
	}
	if _, ok := testDict.vipIndex["132.60.60.10"]["inst2"]; ok {
		t.Error("inst2 wasn't deleted from vipIndex")
	}
	if _, ok := testDict.svipIndex["142.20.1.3"]["inst2"]; ok {
		t.Error("inst2 wasn't deleted from svipIndex")
	}

	testDict.Delete(instsApp["inst3"], "inst3", &app1)
	if len(testDict.appNameIndex) != 0 || len(testDict.vipIndex) != 0 || len(testDict.svipIndex) != 0 {
		t.Errorf("Length of dict didn't change after delete. len of appInd = %d, len of vipIndex = %d , "+
			"len of svipIndex = %d ", len(testDict.appNameIndex), len(testDict.vipIndex), len(testDict.svipIndex))
	}

}

func TestOnAdd(t *testing.T) {
	setup()
	t.Logf("len of appInd = %d, len of vipIndex = %d , "+
		"len of svipIndex = %d ", len(testDict.appNameIndex["app1"]), len(testDict.vipIndex["132.60.60.10"]), len(testDict.svipIndex["142.20.1.3"]))
	instToAdd := &Instance{
		ID:          "new_inst",
		Application: "app1 ",
		VIPAddr:     "132.60.60.10",
		SecVIPAddr:  "142.20.1.3",
		Status:      "UP",
		ActionType:  hashcodeDelimiter,
		Port:        &Port{Enabled: "True"},
		SecPort:     &Port{Enabled: "True"},
		Datacenter:  &DatacenterInfo{},
		Lease:       &LeaseInfo{},
	}
	app1 := Application{Name: "app1"}
	testDict.Add(instToAdd, "new_inst", &app1)
	if len(testDict.appNameIndex["app1"]) != 4 || len(testDict.vipIndex["132.60.60.10"]) != 4 || len(testDict.svipIndex["142.20.1.3"]) != 4 {
		t.Errorf("Length of dict didn't change after Add. len of appInd = %d, len of vipIndex = %d , "+
			"len of svipIndex = %d ", len(testDict.appNameIndex["app1"]), len(testDict.vipIndex["132.60.60.10"]), len(testDict.svipIndex["142.20.1.3"]))
	}
	if _, ok := testDict.appNameIndex["app1"]["new_inst"]; !ok {
		t.Error("new_inst wasn't added from appNameIndex")
	}
	if _, ok := testDict.vipIndex["132.60.60.10"]["new_inst"]; !ok {
		t.Error("new_inst wasn't added from vipIndex")
	}
	if _, ok := testDict.svipIndex["142.20.1.3"]["new_inst"]; !ok {
		t.Error("new_inst wasn't added from svipIndex")
	}

	instToAdd2 := &Instance{
		ID:          "new_inst2",
		Application: "app1 ",
		VIPAddr:     "132.60.60.55",
		SecVIPAddr:  "142.20.1.33",
		Status:      "UP",
		HostName:    "bla",
		ActionType:  hashcodeDelimiter,
		Port:        &Port{Enabled: "True"},
		SecPort:     &Port{Enabled: "True"},
		Datacenter:  &DatacenterInfo{},
		Lease:       &LeaseInfo{},
	}
	t.Logf("len of appInd = %d, len of vipIndex = %d , "+
		"len of svipIndex = %d ", len(testDict.appNameIndex["app1"]), len(testDict.vipIndex["132.60.60.10"]), len(testDict.svipIndex["142.20.1.3"]))

	testDict.Add(instToAdd2, "new_inst2", &app1)
	t.Log(testDict.vipIndex["132.60.60.10"])
	t.Logf("*******************************")
	t.Log(testDict.vipIndex["132.60.60.55"])

	if len(testDict.appNameIndex["app1"]) != 5 || len(testDict.vipIndex["132.60.60.10"]) != 4 || len(testDict.svipIndex["142.20.1.3"]) != 4 ||
		len(testDict.vipIndex["132.60.60.55"]) != 1 || len(testDict.svipIndex["142.20.1.33"]) != 1 {
		t.Errorf("Length of dict didn't change after Add. len of appInd = %d, len of vipIndex = %d , "+
			"len of svipIndex = %d   H- %d %d -H  ",
			len(testDict.appNameIndex["app1"]), len(testDict.vipIndex["132.60.60.10"]), len(testDict.svipIndex["142.20.1.3"]),
			len(testDict.vipIndex["132.60.60.55"]), len(testDict.svipIndex["142.20.1.33"]))
	}
	if _, ok := testDict.appNameIndex["app1"]["new_inst2"]; !ok {
		t.Error("new_inst2 wasn't added from appNameIndex")
	}
	if _, ok := testDict.vipIndex["132.60.60.55"]["new_inst2"]; !ok {
		t.Error("new_inst2 wasn't added from vipIndex")
	}
	if _, ok := testDict.svipIndex["142.20.1.33"]["new_inst2"]; !ok {
		t.Error("new_inst2 wasn't added from svipIndex")
	}
}

func TestGetApplication(t *testing.T) {
	setup()

	app := testDict.getApplication("app1")
	if app.Name != "app1" {
		t.Error("Unexpected appName")
	}
	//t.Log(app.Instances)
	if len(app.Instances) != 3 {
		t.Errorf("app should contain 3 instances, %d", len(app.Instances))
	}

}

func TestGetApplications(t *testing.T) {
	setup()

	apps := testDict.getApplications()
	if apps[0].Name != "app1" {
		t.Error("Unexpected appName")
	}
	//t.Log(app.Instances)
	if len(apps[0].Instances) != 3 {
		t.Errorf("app should contain 3 instances, %d", len(apps[0].Instances))
	}
	if len(apps) != 1 {
		t.Error("Unexpected number of apps")
	}

}

func TestGetInstancesByVip(t *testing.T) {
	setup()

	insts := testDict.GetInstancesByVip("142.20.1.3")
	if insts != nil {
		t.Errorf("should return nil")
	}
	insts = testDict.GetInstancesByVip("132.60.60.10")

	if len(insts) != 3 {
		t.Errorf("should contain 3 instances, %d", len(insts))
	}

}

func TestGetInstancesBySecVip(t *testing.T) {
	setup()

	insts := testDict.GetInstancesBySecVip("132.60.60.10")
	if insts != nil {
		t.Errorf("should return nil")
	}
	insts = testDict.GetInstancesBySecVip("142.20.1.3")

	if len(insts) != 3 {
		t.Errorf("should contain 3 instances, %d", len(insts))
	}

}

func TestIsEmpty(t *testing.T) {
	setup()
	if testDict.isEmpty() {
		t.Errorf("dict should not be empty")
	}
	dict := dictionary{}
	if !dict.isEmpty() {
		t.Errorf("should be Empty")
	}
	dict = dictionary{appNameIndex: nil, vipIndex: nil, svipIndex: nil}

	if !dict.isEmpty() {
		t.Errorf("should be Empty")
	}
	app1 := Application{Name: "app1"}
	testDict.Delete(instsApp["inst1"], "inst1", &app1)
	testDict.Delete(instsApp["inst2"], "inst2", &app1)
	testDict.Delete(instsApp["inst3"], "inst3", &app1)
	if !testDict.isEmpty() {
		t.Errorf("should be Empty")
	}

}

func TestCopyDictionary(t *testing.T) {
	setup()
	copyDict := testDict.copyDictionary()
	if len(copyDict.appNameIndex["app1"]) != 3 || len(copyDict.vipIndex["132.60.60.10"]) != 3 || len(copyDict.svipIndex["142.20.1.3"]) != 3 {
		t.Errorf("Length of copy dict  len of appInd = %d, len of vipIndex = %d , "+
			"len of svipIndex = %d ", len(copyDict.appNameIndex["app1"]), len(copyDict.vipIndex["132.60.60.10"]), len(copyDict.svipIndex["142.20.1.3"]))
	}
	testDict.svipIndex["142.20.1.3"]["inst1"].ActionType = actionAdded
	if copyDict.svipIndex["142.20.1.3"]["inst1"].ActionType != actionAdded {
		t.Errorf("should shallow copy the instance... ")
	}
}
