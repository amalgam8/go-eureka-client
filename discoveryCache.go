// Copyright (c) 2016 IBM Corp. Licensed Materials - Property of IBM.
package go_eureka_client

import (
	"time"
	"fmt"
)

type DiscoveryCache interface {
	Discovery
	Run(stopCh chan struct{} )
}

type discoveryCache struct {
	client client
	pollInterval time.Duration

}

// NewDiscoveryCache creates a new client used for instances discovery with internal cache.
// pollInterval defines the polling interval.
// handler is used to get notification on instances. nil indicates that no notifications are needed.
func NewDiscoveryCache(config *Config, pollInterval time.Duration, handler InstanceEventHandler) (DiscoveryCache, error){

	discoveryCacheClient, err := newClient(config, handler)
	if err != nil {
		return nil,err
	}

	newDiscoveryCache :=  &discoveryCache{client: *discoveryCacheClient,
		                             pollInterval: pollInterval,
		                             }

	return newDiscoveryCache,nil
}

// Run start running the cache.
func (d *discoveryCache) Run(stopCh chan struct{} ) {
	go d.client.run(d.pollInterval, stopCh)
}
// GetApplication returns an application instance from the cache with the appName specified as argument.
func (d *discoveryCache) GetApplication(appName string) (*Application, error) {

	app := d.client.dictionary.getApplication(appName)
	if app == nil {
		return nil, fmt.Errorf("Application Name %s not found",appName )
	}
	return app,nil
}

// GetApplications retrieves all applications from the cache and returns them inside an array.
func (d *discoveryCache) GetApplications() ([]*Application, error) {
	// TODO : check if there is a possible error
	return d.client.dictionary.getApplications(), nil
}

// GetInstance returns from the cache an instance object with the specified appId and id given as arguments.
// appId - string representing application name. id - id  string of instance
func (d *discoveryCache) GetInstance(appId, id string) (*Instance, error) {
	return d.client.dictionary.appNameIndex[appId][id],nil
}

// GetInstancesByVip returns from the cache all the instances with the given vipAddress.
func (d *discoveryCache) GetInstancesByVip(vipAddress string) ([]*Instance, error) {
	instances := d.client.dictionary.GetInstancesByVip(vipAddress)
	if instances == nil {
		return nil, fmt.Errorf("vipAddress  %s not found",vipAddress )
	}
	return instances,nil
}

// GetInstancesBySecVip return from the cache all the instances with the given secured vip address.
func (d *discoveryCache) GetInstancesBySecVip(secVipAddress string) ([]*Instance, error) {
	instances := d.client.dictionary.GetInstancesBySecVip(secVipAddress)
	if instances == nil {
		return nil, fmt.Errorf("vipAddress  %s not found",secVipAddress )
	}
	return instances,nil
}