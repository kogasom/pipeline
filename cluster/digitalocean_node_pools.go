package cluster

import (
	"github.com/banzaicloud/pipeline/internal/providers/digitalocean"
	digitaloceanAPI "github.com/banzaicloud/pipeline/pkg/cluster/digitalocean"
	pkgErrors "github.com/banzaicloud/pipeline/pkg/errors"
)

// createDigitalOceanNodePoolsModelFromRequest creates an array of DigitalOceanNodePoolModel from the nodePoolsData received through create/update requests
func createDigitalOceanNodePoolsModelFromRequest(nodePoolsData map[string]*digitaloceanAPI.NodePool, userID uint) ([]*digitalocean.DigitalOceanNodePoolModel, error) {
	nodePoolsCount := len(nodePoolsData)
	if nodePoolsCount == 0 {
		return nil, pkgErrors.ErrorNodePoolNotProvided
	}
	nodePoolsModel := make([]*digitalocean.DigitalOceanNodePoolModel, nodePoolsCount)

	i := 0
	for nodePoolName, nodePoolData := range nodePoolsData {
		nodePoolsModel[i] = &digitalocean.DigitalOceanNodePoolModel{
			CreatedBy: userID,
			Name:      nodePoolName,
			Count:     nodePoolData.Count,
			Size:      nodePoolData.Size,
			Tags:      createDigitalOceanTagsModelFromRequest(nodePoolData.Tags),
		}

		i++
	}

	return nodePoolsModel, nil
}
