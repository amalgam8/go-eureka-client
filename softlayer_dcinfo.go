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
)

const (
	softlayerInstanceID = "getId"
)

// SoftLayer datacenter information
type softlayerInfo struct{}

// Returns the unique identifier of this datacenter info
func (softlayer *softlayerInfo) GetID(dcinfo *DatacenterInfo) string {
	if dcinfo == nil || dcinfo.Metadata == nil {
		return ""
	}

	uid := dcinfo.Metadata[softlayerInstanceID]

	switch v := uid.(type) {
	case float32, float64:
		return fmt.Sprintf("%.0f", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
