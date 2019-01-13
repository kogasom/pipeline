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

package manager

import (
	"fmt"

	"github.com/banzaicloud/pipeline/pkg/providers/digitalocean/client"
	"github.com/banzaicloud/pipeline/pkg/providers/digitalocean/model"
)

// ClusterManager for managing Cluster state
type ClusterManager struct {
	client *client.DigitalOcean
}

// NewClusterManager creates a new ClusterManager
func NewClusterManager(client *client.DigitalOcean) *ClusterManager {
	return &ClusterManager{
		client: client,
	}
}

// ValidateModel validates a DigitalOcean Cluster model
func (manager *ClusterManager) ValidateModel(model *model.Cluster) error {
	options, err := manager.client.GetOptions()
	if err != nil {
		return err
	}

	// TODO
	fmt.Println(options)

	return nil
}
