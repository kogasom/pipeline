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

package cluster

import (
	"github.com/banzaicloud/pipeline/model"
	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
	pkgCommon "github.com/banzaicloud/pipeline/pkg/common"
	doClient "github.com/banzaicloud/pipeline/pkg/providers/digitalocean/client"
	"github.com/banzaicloud/pipeline/pkg/providers/digitalocean/cluster/manager"
	doModel "github.com/banzaicloud/pipeline/pkg/providers/digitalocean/model"
	doSecret "github.com/banzaicloud/pipeline/pkg/providers/digitalocean/secret"
	"github.com/banzaicloud/pipeline/secret"
)

// Cluster struct for DigitalOcean cluster
type DOCluster struct {
	modelCluster *model.ClusterModel
	CommonClusterBase
}

// Entity properties

//GetID returns the specified cluster id
func (cluster *DOCluster) GetID() uint {
	return cluster.modelCluster.ID
}

// GetUID return the cluster UID
func (cluster *DOCluster) GetUID() string {
	return cluster.modelCluster.UID
}

// GetOrganizationId gets org where the cluster belongs
func (cluster *DOCluster) GetOrganizationId() uint {
	return cluster.modelCluster.OrganizationId
}

//GetName returns the name of the cluster
func (cluster *DOCluster) GetName() string {
	return cluster.modelCluster.Name
}

// GetCloud returns the cloud type of the cluster
func (cluster *DOCluster) GetCloud() string {
	return pkgCluster.DigitalOcean
}

// GetDistribution returns the distribution type of the cluster
func (cluster *DOCluster) GetDistribution() string {
	return cluster.modelCluster.Distribution
}

// GetLocation gets where the cluster is.
func (cluster *DOCluster) GetLocation() string {
	return cluster.modelCluster.Location
}

// GetCreatedBy returns cluster create userID.
func (cluster *DOCluster) GetCreatedBy() uint {
	return cluster.modelCluster.CreatedBy
}

// Secrets

//GetSecretId retrieves the secret id
func (cluster *DOCluster) GetSecretId() string {
	return cluster.modelCluster.SecretId
}

// TODO
func (cluster *DOCluster) GetSshSecretId() string {
	return ""
}

// TODO
func (cluster *DOCluster) SaveSshSecretId(string) error {
	return nil
}

// TODO
func (cluster *DOCluster) SaveConfigSecretId(string) error {
	return nil
}

// TODO
func (cluster *DOCluster) GetConfigSecretId() string {
	return ""
}

// TODO
func (cluster *DOCluster) GetSecretWithValidation() (*secret.SecretItemResponse, error) {
	return &secret.SecretItemResponse{}, nil
}

// Persistence
// TODO
func (cluster *DOCluster) Persist(status, statusMessage string) error {
	return cluster.modelCluster.UpdateStatus(status, statusMessage)
}

// TODO
func (cluster *DOCluster) UpdateStatus(status, statusMessage string) error {
	return cluster.modelCluster.UpdateStatus(status, statusMessage)
}

// TODO
func (cluster *DOCluster) DeleteFromDatabase() error {
	err := cluster.modelCluster.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Cluster management

// TODO
func (cluster *DOCluster) CreateCluster() error {
	return nil
}

// TODO
func (cluster *DOCluster) ValidateCreationFields(r *pkgCluster.CreateClusterRequest) error {
	cm, err := cluster.GetClusterManager()
	if err != nil {
		return err
	}

	return cm.ValidateModel(&cluster.modelCluster.DOKE)
}

// TODO
func (cluster *DOCluster) UpdateCluster(*pkgCluster.UpdateClusterRequest, uint) error {
	return nil
}

// TODO
func (cluster *DOCluster) CheckEqualityToUpdate(*pkgCluster.UpdateClusterRequest) error {
	return nil
}

// TODO
func (cluster *DOCluster) AddDefaultsToUpdate(*pkgCluster.UpdateClusterRequest) {
	// ToDo
}

// TODO
func (cluster *DOCluster) DeleteCluster() error {
	return nil
}

// Kubernetes

// TODO
func (cluster *DOCluster) DownloadK8sConfig() ([]byte, error) {
	return make([]byte, 0), nil
}

// TODO
func (cluster *DOCluster) GetAPIEndpoint() (string, error) {
	return "", nil
}

// TODO
func (cluster *DOCluster) GetK8sConfig() ([]byte, error) {
	return make([]byte, 0), nil
}

// TODO
func (cluster *DOCluster) RequiresSshPublicKey() bool {
	return false
}

// TODO
func (cluster *DOCluster) RbacEnabled() bool {
	return false
}

// TODO
func (cluster *DOCluster) NeedAdminRights() bool {
	return false
}

// TODO
func (cluster *DOCluster) GetKubernetesUserName() (string, error) {
	return "", nil
}

// Cluster info

// TODO
func (cluster *DOCluster) GetStatus() (*pkgCluster.GetClusterStatusResponse, error) {
	return &pkgCluster.GetClusterStatusResponse{}, nil
}

// TODO
func (cluster *DOCluster) GetClusterDetails() (*pkgCluster.DetailsResponse, error) {
	return &pkgCluster.DetailsResponse{}, nil
}

// TODO
func (cluster *DOCluster) ListNodeNames() (pkgCommon.NodeNames, error) {
	return pkgCommon.NodeNames{}, nil
}

// TODO
func (cluster *DOCluster) NodePoolExists(nodePoolName string) bool {
	return false
}

// CreateClusterFromRequest creates a Cluster struct from the request
func CreateDOClusterFromRequest(request *pkgCluster.CreateClusterRequest, orgId, userId uint) (*DOCluster, error) {
	log.Debug("Create ClusterModel struct from the request")

	var cluster DOCluster

	cluster.modelCluster = &model.ClusterModel{
		Name:           request.Name,
		Location:       request.Location,
		Cloud:          request.Cloud,
		OrganizationId: orgId,
		SecretId:       request.SecretId,
		CreatedBy:      userId,
		Distribution:   pkgCluster.DigitalOcean,
	}

	Model, err := doModel.CreateModelFromCreateRequest(request, userId)
	if err != nil {
		return &cluster, err
	}

	cluster.modelCluster.DOKE = Model

	return &cluster, nil
}

// GetDOClient creates a new DigitalOcean client
func (cluster *DOCluster) GetDOClient() (client *doClient.DigitalOcean, err error) {
	s, err := cluster.CommonClusterBase.getSecret(cluster)
	if err != nil {
		return
	}

	return doClient.NewDO(doSecret.CreateDOCredential(s.Values))
}

// GetClusterManager creates a new ClusterManager
func (cluster *DOCluster) GetClusterManager() (m *manager.ClusterManager, err error) {
	client, err := cluster.GetDOClient()
	if err != nil {
		return
	}

	return manager.NewClusterManager(client), nil
}
