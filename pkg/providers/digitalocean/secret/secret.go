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

package secret

import (
	"github.com/banzaicloud/pipeline/pkg/providers/digitalocean/client"
)

// DigitalOcean keys
const (
	AccessToken = "personal_access_token"
)

// DOVerify for validating DigitalOcean credentials
type DOVerify struct {
	credential *client.Credential
}

// CreateDOSecret creates a new 'DOVerify' instance
func CreateDOSecret(values map[string]string) *DOVerify {
	return &DOVerify{
		credential: CreateDOCredential(values),
	}
}

// CreateDOCredential creates a 'client.Credential' instance from secret's values
func CreateDOCredential(values map[string]string) *client.Credential {
	return &client.Credential{
		AccessToken: values[AccessToken],
	}
}

// VerifySecret validates DigitalOcean credentials
func (v *DOVerify) VerifySecret() error {
	do, _ := client.NewDO(v.credential)

	return do.Validate()
}
