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
	return r.client.register(instance)

}

// Deregister removes an instance from the registry.
func (r *registrator) Deregister(instance *Instance)  error {
	return r.client.deregister(instance)
}

// Heartbeat sends an heartbeat to the registry in order to verify its existence.
func (r *registrator) Heartbeat(instance *Instance)  error {
	return r.client.heartbeat(instance)

}

// SetStatus changes the status of an instance in the registry.
func (r *registrator) SetStatus(instance *Instance,status StatusType)  error {
	return r.client.setStatusForInstance(instance,status)
}

// SetMetadataKey sets the metaDataKey for an instance in the registry.
func (r *registrator) SetMetadataKey(inst *Instance, key string, value string)  error {
	return r.client.setMetadataKey(inst,key,value)
}

