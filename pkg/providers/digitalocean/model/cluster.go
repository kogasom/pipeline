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

package model

import (
	"time"

	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
	pkgErrors "github.com/banzaicloud/pipeline/pkg/errors"
	"github.com/banzaicloud/pipeline/pkg/providers/digitalocean/cluster/request"
)

// TableName constants
const (
	ClustersTableName          = "digitalocean_doke_clusters"
	ClustersNodePoolsTableName = "digitalocean_doke_node_pools"
)

// Cluster describes the DigitalOcean cluster model
type Cluster struct {
	ID            uint `gorm:"primary_key"`
	ClusterID     string
	Name          string `gorm:"unique_index:idx_name"`
	RegionSlug    string
	VersionSlug   string
	ClusterSubnet string
	ServiceSubnet string
	IPv4          string
	Endpoint      string
	Tags          []string

	NodePools []*NodePool `gorm:"foreignkey:ClusterID"`

	CreatedBy uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Delete    bool `gorm:"-"`
}

// NodePool describes DigitalOcean node pools model of a cluster
type NodePool struct {
	ID         uint `gorm:"primary_key"`
	NodePoolID string
	Name       string `gorm:"unique_index:idx_cluster_id_name"`
	ClusterID  uint   `gorm:"unique_index:idx_cluster_id_name"`
	Size       string
	Count      int
	Tags       []string
	CreatedBy  uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Delete     bool `gorm:"-"`
	Add        bool `gorm:"-"`
}

// TableName sets the Clusters table name
func (Cluster) TableName() string {
	return ClustersTableName
}

// TableName sets the NodePools table name
func (NodePool) TableName() string {
	return ClustersNodePoolsTableName
}

// CreateModelFromCreateRequest create model from create request
func CreateModelFromCreateRequest(r *pkgCluster.CreateClusterRequest, userId uint) (cluster Cluster, err error) {

	cluster.Name = r.Name

	return CreateModelFromRequest(cluster, r.Properties.CreateClusterDO, userId)
}

// CreateModelFromRequest creates model from request
func CreateModelFromRequest(model Cluster, r *request.Cluster, userID uint) (cluster Cluster, err error) {
	model.CreatedBy = userID

	// there should be at least 1 node pool defined
	if len(r.NodePools) == 0 {
		return cluster, pkgErrors.ErrorNodePoolNotProvided
	}

	nodePools := make([]*NodePool, 0)
	for name, data := range r.NodePools {
		nodePool := model.GetNodePoolByName(name)
		if nodePool.ID == 0 {
			nodePool.Name = name
			nodePool.Size = data.Size
			nodePool.Count = data.Count
			nodePool.Add = true
		}

		nodePool.CreatedBy = userID
		nodePools = append(nodePools, nodePool)
	}

	for _, np := range model.NodePools {
		if r.NodePools[np.Name] == nil {
			np.Delete = true
			nodePools = append(nodePools, np)
		}
	}

	model.NodePools = nodePools

	return model, err
}

// GetNodePoolByName gets a NodePool from the []NodePools by name
func (c *Cluster) GetNodePoolByName(name string) *NodePool {

	for _, np := range c.NodePools {
		if np.Name == name {
			return np
		}
	}

	return &NodePool{}
}
