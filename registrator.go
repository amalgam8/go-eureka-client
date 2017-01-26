// Copyright (c) 2016 IBM Corp. Licensed Materials - Property of IBM.
package go_eureka_client


type Registrator interface {
	Register(*Instance) error
	Deregister(*Instance) error
	Heartbeat(*Instance) error
	SetStatus(inst*Instance, status StatusType) error
	SetMetadataKey(inst *Instance, key string, value string) error
}

type registrator struct {
	client client
}

// NewRegistrator creates a new client used for instance registration
func NewRegistrator(config *Config,handler InstanceEventHandler) (Registrator, error){
	registratorClient, err := newClient(config, handler)
	if err != nil {
		return nil,err
	}

	newRegistrator :=  &registrator{client: *registratorClient,
	}

	return newRegistrator,nil
}

// Register registers an instance in the registry.
func (r *registrator) Register(instance *Instance)  error {
	r.client.register(instance)
	return nil
}

// Deregister removes an instance from the registry.
func (r *registrator) Deregister(instance *Instance)  error {
	return nil
}

// Heartbeat sends an heartbeat to the registry in order to verify its existence.
func (r *registrator) Heartbeat(instance *Instance)  error {
	return nil
}

// SetStatus changes the status of an instance in the registry.
func (r *registrator) SetStatus(instance *Instance,status StatusType)  error {
	return nil
}

// SetMetadataKey sets the metaDataKey for an instance in the registry.
func (r *registrator) SetMetadataKey(inst *Instance, key string, value string)  error {
	return nil
}

