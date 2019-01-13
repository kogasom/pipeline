// Copyright © 2018 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package request

// Cluster describes Pipeline's DigitalOcean fields of a Create/Update request
type Cluster struct {
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
	RegionSlug  string   `json:"region,omitempty" yaml:"region,omitempty"`
	VersionSlug string   `json:"version,omitempty" yaml:"version,omitempty"`
	Tags        []string `json:"tags,omitempty" yaml:"tags,omitempty"`

	NodePools map[string]*NodePool `json:"nodePools,omitempty" yaml:"nodePools,omitempty"`
}

// NodePool describes Pipeline's DigitalOcean fields of a Create/Update request
type NodePool struct {
	Name  string   `json:"name,omitempty" yaml:"name,omitempty"`
	Size  string   `json:"size,omitempty" yaml:"size,omitempty"`
	Count int      `json:"count,omitempty" yaml:"count,omitempty"`
	Tags  []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// TODO
func (c *Cluster) Validate() error {
	return nil
}
