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

//import "log"

type dictionary struct {
	appNameIndex map[string]map[string]*Instance
	vipIndex     map[string]map[string]*Instance
	svipIndex    map[string]map[string]*Instance
}

func newDictionary() dictionary {
	return dictionary{appNameIndex: map[string]map[string]*Instance{}, vipIndex: map[string]map[string]*Instance{},
		svipIndex: map[string]map[string]*Instance{}}
}

func (d *dictionary) onDelete(inst *Instance, id string, app *Application) {
	if vipInsts, ok := d.vipIndex[inst.VIPAddr]; ok {
		delete(vipInsts, id)
		if len(vipInsts) == 0 {
			delete(d.vipIndex, inst.VIPAddr)
		}
	}
	if svipInsts, ok := d.svipIndex[inst.SecVIPAddr]; ok {
		delete(d.svipIndex[inst.SecVIPAddr], id)
		if len(svipInsts) == 0 {
			delete(d.svipIndex, inst.SecVIPAddr)
		}
	}
	if appInsts, ok := d.appNameIndex[app.Name]; ok {
		delete(appInsts, id)
		if len(appInsts) == 0 {
			delete(d.appNameIndex, app.Name)
		}
	}
}

func (d *dictionary) onAdd(inst *Instance, id string, app *Application) {

	if inst.VIPAddr != "" {
		insts := d.vipIndex[inst.VIPAddr]
		if insts == nil {
			insts = map[string]*Instance{}
			d.vipIndex[inst.VIPAddr] = insts
			d.vipIndex[inst.VIPAddr][inst.ID] = inst
		} else {
			d.vipIndex[inst.VIPAddr][inst.ID] = inst
		}
	}

	if inst.SecVIPAddr != "" {
		if d.svipIndex[inst.SecVIPAddr] == nil {
			d.svipIndex[inst.SecVIPAddr] = map[string]*Instance{}
			d.svipIndex[inst.SecVIPAddr][inst.ID] = inst
		} else {
			d.svipIndex[inst.SecVIPAddr][inst.ID] = inst
		}
	}

	if inst.Application != "" {
		if d.appNameIndex[app.Name] == nil {
			d.appNameIndex[app.Name] = map[string]*Instance{}
			d.appNameIndex[app.Name][inst.ID] = inst
		} else {
			d.appNameIndex[app.Name][inst.ID] = inst
		}
	}

}

func (d *dictionary) onUpdate(inst *Instance, id string, app *Application) {
	if inst.VIPAddr != "" {
		insts := d.vipIndex[inst.VIPAddr]
		if insts == nil {
			insts = map[string]*Instance{}
			d.vipIndex[inst.VIPAddr] = insts
			d.vipIndex[inst.VIPAddr][inst.ID] = inst
		} else {
			d.vipIndex[inst.VIPAddr][inst.ID] = inst
		}
	}

	if inst.SecVIPAddr != "" {
		if d.svipIndex[inst.SecVIPAddr] == nil {
			d.svipIndex[inst.SecVIPAddr] = map[string]*Instance{}
			d.svipIndex[inst.SecVIPAddr][inst.ID] = inst
		} else {
			d.svipIndex[inst.SecVIPAddr][inst.ID] = inst
		}
	}

	if inst.Application != "" {
		if d.appNameIndex[app.Name] == nil {
			d.appNameIndex[app.Name] = map[string]*Instance{}
			d.appNameIndex[app.Name][inst.ID] = inst
		} else {
			d.appNameIndex[app.Name][inst.ID] = inst
		}
	}
}

func (d *dictionary) getApplication(appName string) *Application {
	instancesMap, ok := d.appNameIndex[appName]
	if !ok {
		return nil
	}
	var instancesArray []*Instance
	for _, v := range instancesMap {
		dc := v.deepCopy()
		instancesArray = append(instancesArray, &dc)
	}
	app := &Application{Name: appName,
		Instances: instancesArray}
	return app
}

func (d *dictionary) getApplications() []*Application {
	var applicationsArray []*Application
	for appName, instancesMap := range d.appNameIndex {
		var instancesArray []*Instance
		for _, v := range instancesMap {
			dc := v.deepCopy()
			instancesArray = append(instancesArray, &dc)
		}
		app := &Application{Name: appName,
			Instances: instancesArray}
		applicationsArray = append(applicationsArray, app)
	}

	return applicationsArray
}

func (d *dictionary) GetInstancesByVip(vipAddress string) []*Instance {
	instancesMap, ok := d.vipIndex[vipAddress]
	if !ok {
		return nil
	}
	var instancesArray []*Instance
	for _, v := range instancesMap {
		dc := v.deepCopy()
		instancesArray = append(instancesArray, &dc)
	}
	return instancesArray
}

func (d *dictionary) GetInstancesBySecVip(svipAddress string) []*Instance {
	instancesMap, ok := d.svipIndex[svipAddress]
	if !ok {
		return nil
	}
	var instancesArray []*Instance
	for _, v := range instancesMap {
		dc := v.deepCopy()
		instancesArray = append(instancesArray, &dc)
	}
	return instancesArray
}

func (d *dictionary) isEmpty() bool {
	if len(d.appNameIndex) == 0 && len(d.svipIndex) == 0 && len(d.vipIndex) == 0 {
		return true
	}
	return false

}

func (d *dictionary) copyDictionary() dictionary {
	service := dictionary{appNameIndex: map[string]map[string]*Instance{}, svipIndex: map[string]map[string]*Instance{}, vipIndex: map[string]map[string]*Instance{}}
	for appIndex, insts := range d.appNameIndex {
		copyInsts := map[string]*Instance{}
		for instName, inst := range insts {
			copyInsts[instName] = inst
		}
		service.appNameIndex[appIndex] = copyInsts
	}

	for vipIndex, insts := range d.vipIndex {
		copyInsts := map[string]*Instance{}
		for instName, inst := range insts {
			copyInsts[instName] = inst
		}
		service.vipIndex[vipIndex] = copyInsts
	}

	for svipIndex, insts := range d.svipIndex {
		copyInsts := map[string]*Instance{}
		for instName, inst := range insts {
			copyInsts[instName] = inst
		}
		service.svipIndex[svipIndex] = copyInsts
	}
	return service

}
