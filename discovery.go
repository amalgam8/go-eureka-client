// Copyright (c) 2016 IBM Corp. Licensed Materials - Property of IBM.
package go_eureka_client

type Discovery interface {
	GetApplication(appName string) (*Application, error)
	GetApplications() ([]*Application, error)
	GetInstance(appId, id string) (*Instance, error)
	GetInstancesByVip(vipAddress string) ([]*Instance, error)
	GetInstancesBySecVip(secVipAddress string) ([]*Instance, error)
}

type discovery struct {
	client client
}

// NewDiscovery creates a new client used for instances discovery without cache.
func NewDiscovery(config *Config) (Discovery, error){
	return nil,nil
}

// GetApplication returns an application instance from the registry with the appName specified as argument.
func (r *discovery) GetApplication(appName string) (*Application, error) {
	return nil,nil
}

// GetApplications retrieves all applications from the registry and returns the inside an array.
func (r *discovery) GetApplications() ([]*Application, error) {
	return nil,nil
}

// GetInstance returns from the registry an instance object with the specified appId and id given as arguments.
func (r *discovery) GetInstance(appId, id string) (*Instance, error) {
	return nil,nil
}

// GetInstancesByVip returns from the registry all the instances with the given vipAddress.
func (r *discovery) GetInstancesByVip(vipAddress string) ([]*Instance, error) {
	return nil,nil
}

// GetInstancesBySecVip return from the registry all the instances with the given secured vip address.
func (r *discovery) GetInstancesBySecVip(secVipAddress string) ([]*Instance, error) {
	return nil,nil
}