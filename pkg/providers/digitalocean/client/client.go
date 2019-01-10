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

package client

import (
	"context"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// DigitalOcean is for managing DigitalOcean API calls
type DigitalOcean struct {
	credential *Credential
}

// Credential describes DigitalOcean credentials for access
type Credential struct {
	AccessToken string
}

// NewDO creates a new DigitalOcean client config
func NewDO(credential *Credential) (do *DigitalOcean, err error) {
	do = &DigitalOcean{credential: credential}

	return
}

func (t Credential) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// NewClient returns a new DigitalOcean API client
func (do *DigitalOcean) NewClient() *godo.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: do.credential.AccessToken,
	})

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	return godo.NewClient(oauthClient)
}

// Validate is validates the credentials by retrieving profile information
func (do *DigitalOcean) Validate() error {
	_, _, err := do.NewClient().Account.Get(context.TODO())

	return err
}
