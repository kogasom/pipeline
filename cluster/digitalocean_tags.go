package cluster

import (
	"github.com/banzaicloud/pipeline/internal/providers/digitalocean"
)

// createDigitalOceanTagsModelFromRequest creates an array of DigitalOceanTagModel from the tagsData received through create/update requests
func createDigitalOceanTagsModelFromRequest(tagsData []string) []*digitalocean.DigitalOceanTagModel {
	tagsModel := make([]*digitalocean.DigitalOceanTagModel, len(tagsData))

	for tagPos, tag := range tagsData {
		tagsModel[tagPos] = &digitalocean.DigitalOceanTagModel{
			Value: tag,
		}
	}

	return tagsModel
}
