package cluster

import (
	"context"
	"time"

	pipConfig "github.com/banzaicloud/pipeline/config"
	"github.com/banzaicloud/pipeline/internal/cluster"
	"github.com/banzaicloud/pipeline/internal/providers/digitalocean"
	"github.com/banzaicloud/pipeline/model"
	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
	pkgCommon "github.com/banzaicloud/pipeline/pkg/common"
	pkgSecret "github.com/banzaicloud/pipeline/pkg/secret"
	"github.com/banzaicloud/pipeline/secret"
	"github.com/digitalocean/godo"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// CreateDigitalOceanClusterFromRequest creates ClusterModel struct from the request
func CreateDigitalOceanClusterFromRequest(request *pkgCluster.CreateClusterRequest, orgID, userID uint) (*DigitalOceanCluster, error) {
	log.Debug("Create ClusterModel struct from the request")
	var c DigitalOceanCluster

	c.db = pipConfig.DB()

	nodePools, err := createDigitalOceanNodePoolsModelFromRequest(request.Properties.CreateClusterDigitalOcean.NodePools, userID)
	if err != nil {
		return nil, err
	}

	c.model = &digitalocean.DigitalOceanClusterModel{
		Cluster: cluster.ClusterModel{
			Name:           request.Name,
			Location:       request.Location,
			OrganizationID: orgID,
			SecretID:       request.SecretId,
			Cloud:          digitalocean.Provider,
			Distribution:   digitalocean.ClusterDistributionDigitalOcean,
			CreatedBy:      userID,
		},

		MasterVersion: request.Properties.CreateClusterDigitalOcean.VersionSlug,
		NodePools:     nodePools,
		Tags:          createDigitalOceanTagsModelFromRequest(request.Properties.CreateClusterDigitalOcean.Tags),
	}

	return &c, nil
}

type DigitalOceanCluster struct {
	db                  *gorm.DB
	model               *digitalocean.DigitalOceanClusterModel
	digitalOceanCluster *godo.KubernetesCluster //Don't use this directly
	APIEndpoint         string
	CommonClusterBase
}

// getDigitalOceanServiceClient creates service client for DigitalOcean
func (c *DigitalOceanCluster) getDigitalOceanServiceClient() (*godo.Client, error) {
	secretItem, err := c.GetSecretWithValidation()
	if err != nil {
		return nil, err
	}
	oauthClient := oauth2.NewClient(context.Background(), &tokenSource{AccessToken: secretItem.GetValue(pkgSecret.PAT)})
	return godo.NewClient(oauthClient), nil
}

//GetID returns the specified cluster id
func (c *DigitalOceanCluster) GetID() uint {
	return c.model.Cluster.ID
}

//GetUID returns the specified cluster UID
func (c *DigitalOceanCluster) GetUID() string {
	return c.model.Cluster.UID
}

// GetOrganizationId gets org where the cluster belongs
func (c *DigitalOceanCluster) GetOrganizationId() uint {
	return c.model.Cluster.OrganizationID
}

//GetName returns the name of the cluster
func (c *DigitalOceanCluster) GetName() string {
	return c.model.Cluster.Name
}

// GetCloud returns the cloud type of the cluster
func (c *DigitalOceanCluster) GetCloud() string {
	return c.model.Cluster.Cloud
}

// GetDistribution returns the distribution type of the cluster
func (c *DigitalOceanCluster) GetDistribution() string {
	return c.model.Cluster.Distribution
}

// GetLocation gets where the cluster is.
func (c *DigitalOceanCluster) GetLocation() string {
	return c.model.Cluster.Location
}

// GetCreatedBy returns cluster create userID.
func (c *DigitalOceanCluster) GetCreatedBy() uint {
	return c.model.Cluster.CreatedBy
}

// GetSecretId retrieves the secret id
func (c *DigitalOceanCluster) GetSecretId() string {
	return c.model.Cluster.SecretID
}

// GetSshSecretId retrieves the secret id
func (c *DigitalOceanCluster) GetSshSecretId() string {
	return c.model.Cluster.SSHSecretID
}

// SaveSshSecretId saves the ssh secret id to database
func (c *DigitalOceanCluster) SaveSshSecretId(sshSecretId string) error {
	c.model.Cluster.SSHSecretID = sshSecretId

	err := c.db.Save(&c.model).Error
	if err != nil {
		return errors.Wrap(err, "failed to save ssh secret id")
	}

	return nil
}

// SaveConfigSecretId saves the config secret id in database
func (c *DigitalOceanCluster) SaveConfigSecretId(configSecretId string) error {
	c.model.Cluster.ConfigSecretID = configSecretId

	err := c.db.Save(&c.model).Error
	if err != nil {
		return errors.Wrap(err, "failed to save config secret id")
	}

	return nil
}

// GetConfigSecretId return config secret id
func (c *DigitalOceanCluster) GetConfigSecretId() string {
	return c.model.Cluster.ConfigSecretID
}

// GetSecretWithValidation returns secret from vault
func (c *DigitalOceanCluster) GetSecretWithValidation() (*secret.SecretItemResponse, error) {
	return c.CommonClusterBase.getSecret(c)
}

//Persist save the cluster model
func (c *DigitalOceanCluster) Persist(status, statusMessage string) error {
	log.Infof("Model before save: %v", c.model)
	c.model.Cluster.Status = status
	c.model.Cluster.StatusMessage = statusMessage

	err := c.db.Save(&c.model).Error
	if err != nil {
		return errors.Wrap(err, "failed to persist cluster")
	}

	return nil
}

// UpdateStatus updates cluster status in database
func (c *DigitalOceanCluster) UpdateStatus(status, statusMessage string) error {
	c.model.Cluster.Status = status
	c.model.Cluster.StatusMessage = statusMessage

	err := c.db.Save(&c.model).Error
	if err != nil {
		return errors.Wrap(err, "failed to update status")
	}

	return nil
}

//DeleteFromDatabase deletes model from the database
func (c *DigitalOceanCluster) DeleteFromDatabase() error {
	if err := c.db.Delete(&c.model.Cluster).Error; err != nil {
		return err
	}

	for _, nodePool := range c.model.NodePools {
		if err := c.db.Delete(nodePool).Error; err != nil {
			return err
		}
	}

	if err := c.db.Delete(c.model).Error; err != nil {
		return err
	}

	c.model = nil

	return nil
}

//CreateCluster creates a new cluster
func (c *DigitalOceanCluster) CreateCluster() error {
	log.Info("Start create cluster (DigitalOcean)")
	log.Info("Get DigitalOcean Service Client")

	client, err := c.getDigitalOceanServiceClient()
	if err != nil {
		return err
	}

	log.Info("Get DigitalOcean Service Client succeeded")

	createNodePoolsRequest := make([]*godo.KubernetesNodePoolCreateRequest, len(c.model.NodePools))
	for nodePoolIdx, nodePool := range c.model.NodePools {
		createNodePoolRequest := &godo.KubernetesNodePoolCreateRequest{
			Name:  nodePool.Name,
			Size:  nodePool.Size,
			Count: nodePool.Count,
			Tags:  c.getStringArrayFromTags(nodePool.Tags),
		}
		createNodePoolsRequest[nodePoolIdx] = createNodePoolRequest
	}

	clusterCreateRequest := &godo.KubernetesClusterCreateRequest{
		Name:        c.model.Cluster.Name,
		RegionSlug:  c.model.Cluster.Location,
		VersionSlug: c.model.MasterVersion,
		NodePools:   createNodePoolsRequest,
		Tags:        c.getStringArrayFromTags(c.model.Tags),
	}

	ctx := context.TODO()
	cluster, _, err := client.Kubernetes.Create(ctx, clusterCreateRequest)
	if err != nil {
		return errors.Wrap(err, "failed to create cluster")
	}

	c.model.DigitalOceanID = cluster.ID
	err = c.model.Save()
	if err != nil {
		return errors.Wrap(err, "failed to save cluster")
	}

	log.Infof("Successfully created cluster %s", cluster.Name)
	c.digitalOceanCluster = cluster

	if cluster != nil {
		log.Infof("Cluster %s create is called", cluster.Name)
		log.Info("Waiting for cluster...")

		for cluster.Status.State != "running" {
			log.Infof("Cluster status: %s", cluster.Status.State)
			time.Sleep(time.Second * 5)
			cluster, _, err = client.Kubernetes.Get(ctx, cluster.ID)

			if err != nil {
				return errors.Wrap(err, "error during getting cluster status")
			}
		}

		log.Info("Cluster is running!")
	}

	// TODO: on existing cluster update it instead? like at GKE?

	return nil
}

// ValidateCreationFields validates all field
func (c *DigitalOceanCluster) ValidateCreationFields(r *pkgCluster.CreateClusterRequest) error {
	// TODO
	return nil
}

// UpdateCluster updates DigitalOcean cluster in cloud
func (c *DigitalOceanCluster) UpdateCluster(updateRequest *pkgCluster.UpdateClusterRequest, userId uint) error {
	// TODO
	return nil
}

//CheckEqualityToUpdate validates the update request
func (c *DigitalOceanCluster) CheckEqualityToUpdate(r *pkgCluster.UpdateClusterRequest) error {
	// TODO
	return nil
}

//AddDefaultsToUpdate adds defaults to update request
func (c *DigitalOceanCluster) AddDefaultsToUpdate(r *pkgCluster.UpdateClusterRequest) {
	// TODO
}

// DeleteCluster deletes cluster from digitalocean
func (c *DigitalOceanCluster) DeleteCluster() error {
	// TODO
	return nil
}

// DownloadK8sConfig downloads the kubeconfig file from cloud
func (c *DigitalOceanCluster) DownloadK8sConfig() ([]byte, error) {
	log.Info("Start to download K8s config (DigitalOcean)")
	log.Info("Get DigitalOcean Service Client")

	client, err := c.getDigitalOceanServiceClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DigitalOcean service client")
	}

	log.Info("Get DigitalOcean Service Client succeeded")

	ctx := context.TODO()
	config, _, err := client.Kubernetes.GetKubeConfig(ctx, c.model.DigitalOceanID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DigitalOcean K8s config")
	}

	return config.KubeconfigYAML, nil
}

//GetAPIEndpoint returns the Kubernetes Api endpoint
func (c *DigitalOceanCluster) GetAPIEndpoint() (string, error) {
	return "", nil
}

// GetK8sConfig returns the Kubernetes config
func (c *DigitalOceanCluster) GetK8sConfig() ([]byte, error) {
	return c.CommonClusterBase.getConfig(c)
}

// RbacEnabled returns true if rbac enabled on the cluster
func (c *DigitalOceanCluster) RbacEnabled() bool {
	return c.model.Cluster.RbacEnabled
}

// NeedAdminRights returns true if rbac is enabled and need to create a cluster role binding to user
func (c *DigitalOceanCluster) NeedAdminRights() bool {
	return false
}

// GetKubernetesUserName returns the user ID which needed to create a cluster role binding which gives admin rights to the user
func (c *DigitalOceanCluster) GetKubernetesUserName() (string, error) {
	return "", nil
}

//GetStatus gets cluster status
func (c *DigitalOceanCluster) GetStatus() (*pkgCluster.GetClusterStatusResponse, error) {
	return nil, nil
}

// GetClusterDetails gets cluster details from cloud
func (c *DigitalOceanCluster) GetClusterDetails() (*pkgCluster.DetailsResponse, error) {
	return nil, nil
}

// ListNodeNames returns node names to label them
func (c *DigitalOceanCluster) ListNodeNames() (nodeNames pkgCommon.NodeNames, err error) {
	// nodes are labeled in create request
	return
}

// NodePoolExists returns true if node pool with nodePoolName exists
func (c *DigitalOceanCluster) NodePoolExists(nodePoolName string) bool {
	for _, np := range c.model.NodePools {
		if np != nil && np.Name == nodePoolName {
			return true
		}
	}
	return false
}

//CreateDigitalOceanClusterFromModel creates ClusterModel struct from model
func CreateDigitalOceanClusterFromModel(clusterModel *model.ClusterModel) (*DigitalOceanCluster, error) {
	log.Debug("Create ClusterModel struct from the request")
	db := pipConfig.DB()

	m := digitalocean.DigitalOceanClusterModel{
		ClusterID: clusterModel.ID,
	}

	log.Debug("Load DigitalOcean props from database")
	err := db.Where(m).Preload("Cluster").Preload("NodePools").First(&m).Error
	if err != nil {
		return nil, err
	}

	digitalOceanCluster := DigitalOceanCluster{
		db:    db,
		model: &m,
	}
	return &digitalOceanCluster, nil
}

func (c *DigitalOceanCluster) getStringArrayFromTags(tags []*digitalocean.DigitalOceanTagModel) []string {
	tagsStringArray := make([]string, len(tags))
	for tagIndex, tag := range tags {
		tagsStringArray[tagIndex] = tag.Value
	}
	return tagsStringArray
}

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
