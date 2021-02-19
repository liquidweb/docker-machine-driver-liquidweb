/*
Copyright © LiquidWeb

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package lwdriver

import (
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"

	lwclient "github.com/liquidweb/liquidweb-go/client"

	"github.com/liquidweb/docker-machine-driver-liquidweb/util"
)

const (
	DefaultTemplate   string = "DEBIAN_10_UNMANAGED"
	DefaultDockerPort int    = 2376
)

func NewDriver() *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			SSHUser: drivers.DefaultSSHUser,
			SSHPort: drivers.DefaultSSHPort,
		},
	}
}

func (self *Driver) DriverName() string {
	return "liquidweb"
}

func (self *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "LW_USERNAME",
			Name:   "lw-username",
			Usage:  "liquidweb api/account username",
			Value:  "",
		},

		mcnflag.StringFlag{
			EnvVar: "LW_PASSWORD",
			Name:   "lw-password",
			Usage:  "password associated with lw-username",
			Value:  "",
		},

		mcnflag.StringFlag{
			EnvVar: "LW_API_DOMAIN",
			Name:   "lw-api-domain",
			Usage:  "liquidweb public api domain",
			Value:  "https://api.liquidweb.com",
		},

		mcnflag.StringFlag{
			EnvVar: "LW_TEMPLATE",
			Name:   "lw-template",
			Usage:  "Name of the template to deploy on the node",
			Value:  DefaultTemplate,
		},

		mcnflag.StringFlag{
			EnvVar: "LW_NODE_ROOT_PASSWORD",
			Name:   "lw-node-root-password",
			Usage:  "root password to set on the node",
			Value:  "",
		},

		mcnflag.IntFlag{
			EnvVar: "LW_CONFIG_ID",
			Name:   "lw-config-id",
			Usage:  "config-id to deploy the node as",
		},

		mcnflag.IntFlag{
			EnvVar: "LW_ZONE_ID",
			Name:   "lw-zone-id",
			Usage:  "zone_id of the zone to deploy the node in",
			Value:  -1,
		},

		mcnflag.StringFlag{
			EnvVar: "LW_NODE_HOSTNAME",
			Name:   "lw-node-hostname",
			Usage:  "hostname to give the node",
			Value:  util.RandomHostname(),
		},

		mcnflag.IntFlag{
			EnvVar: "LW_DOCKER_PORT",
			Name:   "lw-docker-port",
			Usage:  "dockerport to use",
			Value:  DefaultDockerPort,
		},
	}
}

func (self *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	self.BaseDriver.SetSwarmConfigFromFlags(opts)

	self.LwApiUsername = opts.String("lw-username")
	self.LwApiPassword = opts.String("lw-password")
	self.LwApiDomain = opts.String("lw-api-domain")

	self.LwComputeConfigId = opts.Int("lw-config-id")
	if self.LwComputeConfigId == 0 {
		// TODO private parent child support
		return errors.New("must give a lw-config-id")
	}
	self.LwComputeZoneId = opts.Int("lw-zone-id")
	self.LwComputeNodeHostname = opts.String("lw-node-hostname")
	self.LwComputeNodeRootPassword = opts.String("lw-node-root-password")
	self.LwComputeTemplate = opts.String("lw-template")

	self.DockerPort = opts.Int("lw-docker-port")

	if self.DockerPort <= 0 {
		return errors.New("docker port must be greater than zero")
	}

	if self.LwApiUsername == "" || self.LwApiPassword == "" {
		return errors.New("lw username and password are required")
	}

	if self.LwApiDomain == "" {
		return errors.New("api domain cannot be blank")
	}

	if self.LwComputeTemplate == "" {
		return errors.New("template cannot be blank")
	}

	if strings.Contains(strings.ToUpper(self.LwComputeTemplate), "WINDOWS") {
		return errors.New("Windows is not supported")
	}

	return nil
}

// private

func (self *Driver) waitUntilMachineNodeReady(api *lwclient.API) {
	var ready bool
	for !ready {
		status, err := api.StormServer.Status(self.LwComputeNodeUniqId)
		if err != nil {
			log.Warnf("failed fetching status for [%s]: %s", self.LwComputeNodeUniqId, err)
			time.Sleep(30 * time.Second)
			continue
		}
		if strings.ToUpper(status.Status) == "RUNNING" {
			log.Infof("node [%s] has become ready", self.LwComputeNodeUniqId)
			break
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

func (self *Driver) getApiClient() (*lwclient.API, error) {
	user := self.LwApiUsername
	pass := self.LwApiPassword
	url := self.LwApiDomain
	api, err := lwclient.NewAPI(user, pass, url, 90)

	return api, err
}

func (self *Driver) createSshKey() (pubKey string, err error) {
	sshKeyPath := self.GetSSHKeyPath()
	if err = ssh.GenerateSSHKey(sshKeyPath); err != nil {
		return
	}

	var pubKeyB []byte
	pubKeyB, err = ioutil.ReadFile(sshKeyPath + ".pub")
	pubKey = string(pubKeyB)

	return
}
