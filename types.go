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
	"encoding/json"
	"fmt"
)


const (
	UP             StatusType = "UP"
	DOWN           StatusType = "DOWN"
	STARTING       StatusType = "STARTING"
	OUT_OF_SERVICE StatusType = "OUT_OF_SERVICE"
	UNKNOWN        StatusType = "UNKNOWN"
)

const (
	defaultDurationInt uint32 = 90
	metadataTags              = "amalgam8.tags"
	extEureka                 = "eureka"
	extVIP                    = "vipAddress"
)

// Port encapsulates information needed for a port information
type Port struct {
	Enabled string      `json:"@enabled,omitempty"`
	Value   interface{} `json:"$,omitempty"`
}

// DataCenterMetadata encapsulates information needed for a datacenter metadata information
type DatacenterMetadata map[string]interface{}

// DatacenterInfo encapsulates information needed for a datacenter information
type DatacenterInfo struct {
	Class    string             `json:"@class,omitempty"`
	Name     string             `json:"name,omitempty"`
	Metadata DatacenterMetadata `json:"metadata,omitempty"`
}

// LeaseInfo encapsulates information needed for a lease information
type LeaseInfo struct {
	RenewalInt     uint32 `json:"renewalIntervalInSecs,omitempty"`
	DurationInt    uint32 `json:"durationInSecs,omitempty"`
	RegistrationTs int64  `json:"registrationTimestamp,omitempty"`
	LastRenewalTs  int64  `json:"lastRenewalTimestamp,omitempty"`
}

// Instance encapsulates information needed for a service instance information
type Instance struct {
	ID            string          `json:"instanceId,omitempty"`
	HostName      string          `json:"hostName,omitempty"`
	Application   string          `json:"app,omitempty"`
	GroupName     string          `json:"appGroupName,omitempty"`
	IPAddr        string          `json:"ipAddr,omitempty"`
	VIPAddr       string          `json:"vipAddress,omitempty"`
	SecVIPAddr    string          `json:"secureVipAddress,omitempty"`
	Status        string          `json:"status,omitempty"`
	OvrStatus     string          `json:"overriddenstatus,omitempty"`
	CountryID     int             `json:"countryId,omitempty"`
	Port          *Port           `json:"port,omitempty"`
	SecPort       *Port           `json:"securePort,omitempty"`
	HomePage      string          `json:"homePageUrl,omitempty"`
	StatusPage    string          `json:"statusPageUrl,omitempty"`
	HealthCheck   string          `json:"healthCheckUrl,omitempty"`
	Datacenter    *DatacenterInfo `json:"dataCenterInfo,omitempty"`
	Lease         *LeaseInfo      `json:"leaseInfo,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	CordServer    interface{}     `json:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTs interface{}     `json:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTs   interface{}     `json:"lastDirtyTimestamp,omitempty"`
	ActionType    string          `json:"actionType,omitempty"`
}

// InstanceWrapper encapsulates information needed for a service instance registration
type InstanceWrapper struct {
	Inst *Instance `json:"instance,omitempty"`
}

// String output the structure
func (ir *Instance) String() string {
	mtlen := 0
	if ir.Metadata != nil {
		mtlen = len(ir.Metadata)
	}
	return fmt.Sprintf("vip_addres: %s, endpoint: %s:%d, hostname: %s, status: %s, metadata: %d",
		ir.VIPAddr, ir.IPAddr, ir.Port.Value, ir.HostName, ir.Status, mtlen)
}

func (ir *Instance) deepCopy() Instance {
	copyPort := Port{
		Enabled: ir.Port.Enabled,
		Value: ir.Port.Value,
	}
	copySecPort := Port{
		Enabled: ir.SecPort.Enabled,
		Value: ir.SecPort.Value,
	}

	copyDatacenter := DatacenterInfo{
		Class: ir.Datacenter.Class,
		Metadata: ir.Datacenter.Metadata,
		Name: ir.Datacenter.Name,
	}
	copyLease := LeaseInfo{
		DurationInt: ir.Lease.DurationInt,
		LastRenewalTs: ir.Lease.LastRenewalTs,
		RegistrationTs: ir.Lease.RegistrationTs,
		RenewalInt: ir.Lease.RenewalInt,

	}

	copyInst := Instance{
		ID: ir.ID,
		HostName: ir.HostName,
		Application: ir.Application,
		GroupName: ir.GroupName,
		IPAddr: ir.IPAddr,
		VIPAddr: ir.VIPAddr,
		SecVIPAddr: ir.SecVIPAddr,
		Status: ir.Status,
		OvrStatus: ir.OvrStatus,
		CountryID: ir.CountryID,
		Port : &copyPort,
		SecPort: &copySecPort,
		HomePage: ir.HomePage,
		StatusPage: ir.StatusPage,
		HealthCheck: ir.HealthCheck,
		Datacenter: &copyDatacenter,
		Lease: &copyLease,
		Metadata: ir.Metadata,
		CordServer: ir.CordServer,
		LastUpdatedTs: ir.LastDirtyTs,
		LastDirtyTs: ir.LastDirtyTs,
		ActionType: ir.ActionType,
	}
	return copyInst
}


type StatusType string


type DataCenterType string

// InstanceEventHandler can handle notifications for events that happen
// to a instance. The events are information only, so you can't return
// an error.
type InstanceEventHandler interface {
	OnAdd(inst *Instance)
	OnUpdate(oldInst, newInst *Instance)
	OnDelete(int *Instance)
}


// Application is an array of instances
type Application struct {
	Name      string      `json:"name,omitempty"`
	Instances []*Instance `json:"instance,omitempty"`
}

// UnmarshalJSON parses the JSON object of Application struct.
// We need this specific implementation because the Eureka server
// marshals differently single instance (object) and multiple instances (array).
func (app *Application) UnmarshalJSON(b []byte) error {
	type singleApplication struct {
		Name     string    `json:"name,omitempty"`
		Instance *Instance `json:"instance,omitempty"`
	}

	type multiApplication struct {
		Name      string      `json:"name,omitempty"`
		Instances []*Instance `json:"instance,omitempty"`
	}

	var mApp multiApplication
	err := json.Unmarshal(b, &mApp)
	if err != nil {
		// error probably means that we have a single instance object.
		// Thus, we try to unmarshal to a different object type
		var sApp singleApplication
		err = json.Unmarshal(b, &sApp)
		if err != nil {
			return err
		}
		app.Name = sApp.Name
		if sApp.Instance != nil {
			app.Instances = []*Instance{sApp.Instance}
		}
		return nil
	}

	app.Name = mApp.Name
	app.Instances = mApp.Instances
	return nil
}

type appVersion struct {
	VersionDelta int64  `json:"versions__delta,omitempty"`
	Hashcode     string `json:"apps__hashcode,omitempty"`
}

// Applications is an array of application objects
type Applications struct {
	appVersion
	Application []*Application `json:"application,omitempty"`
}

// UnmarshalJSON parses the JSON object of Applications struct.
// We need this specific implementation because the Eureka server
// marshals differently single application (object) and multiple applications (array).
func (apps *Applications) UnmarshalJSON(b []byte) error {
	type singleApplications struct {
		appVersion
		Application *Application `json:"application,omitempty"`
	}

	type multiApplications struct {
		appVersion
		Application []*Application `json:"application,omitempty"`
	}

	var mApps multiApplications
	err := json.Unmarshal(b, &mApps)
	if err != nil {
		// error probably means that we have a single Application object.
		// Thus, we try to unmarshal to a different object type
		var sApps singleApplications
		err = json.Unmarshal(b, &sApps)
		if err != nil {
			return err
		}
		apps.Hashcode = sApps.Hashcode
		apps.VersionDelta = sApps.VersionDelta
		if sApps.Application != nil {
			apps.Application = []*Application{sApps.Application}
		}
		return nil
	}

	apps.Hashcode = mApps.Hashcode
	apps.VersionDelta = mApps.VersionDelta
	apps.Application = mApps.Application
	return nil
}

// ApplicationsList is a list of application objects
type ApplicationsList struct {
	Applications *Applications `json:"applications,omitempty"`
}