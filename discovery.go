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
	"fmt"
	"log"
)

// Discovery defines the discovery interface and actions :
type Discovery interface {
	GetApplication(appName string) (*Application, error)
	GetApplications() ([]*Application, error)
	GetInstance(appID, id string) (*Instance, error)
	GetInstancesByVip(vipAddress string) ([]*Instance, error)
	GetInstancesBySecVip(secVipAddress string) ([]*Instance, error)
}

type discovery struct {
	client *client
}

// NewDiscovery creates a new client used for instances discovery without cache.
func NewDiscovery(config *Config, handler InstanceEventHandler) (Discovery, error) {
	discoveryClient, err := newClient(config, handler)
	if err != nil {
		return nil, err
	}

	newDiscovery := &discovery{client: discoveryClient}

	return newDiscovery, nil
}

// GetApplication returns an application instance from the registry with the appName specified as argument.
// If more the one application instance with the same name exists, it will return the first one found.
func (r *discovery) GetApplication(appName string) (*Application, error) {
	apps, e := r.client.fetchApp("apps/" + appName)
	if e != nil {
		return nil, e
	}
	if len(apps.Application) > 1 {
		log.Print("Found more then one application instance with the same name, will return the first\n")
	} else if len(apps.Application) < 1 {
		return nil, fmt.Errorf("App with the name %s doesn't exist in eureka server", appName)
	}
	return apps.Application[0], e
}

// GetApplications retrieves all applications from the registry and returns the inside an array.
func (r *discovery) GetApplications() ([]*Application, error) {
	apps, e := r.client.fetchApps("apps/")
	if e != nil {
		return nil, e
	}
	if len(apps.Application) < 1 {
		return nil, fmt.Errorf("server is empty")
	}
	return apps.Application, e
}

// GetInstance returns from the registry an instance object with the specified appId and id given as arguments.
func (r *discovery) GetInstance(appID, id string) (*Instance, error) {
	inst, e := r.client.fetchInstance(appID, id)
	if e != nil {
		return nil, e
	}

	return inst, nil

}

// GetInstancesByVip returns from the registry all the instances with the given vipAddress.
func (r *discovery) GetInstancesByVip(vipAddress string) ([]*Instance, error) {
	insts, e := r.client.fetchInstancesByVip(vipAddress)
	if e != nil {
		return nil, e
	}

	return insts, nil
}

// GetInstancesBySecVip return from the registry all the instances with the given secured vip address.
func (r *discovery) GetInstancesBySecVip(secVipAddress string) ([]*Instance, error) {
	insts, e := r.client.fetchInstancesBySVip(secVipAddress)
	if e != nil {
		return nil, e
	}

	return insts, nil

}
