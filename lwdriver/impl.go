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
	"fmt"
	"strings"

	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/state"

	lwclient "github.com/liquidweb/liquidweb-go/client"
	"github.com/liquidweb/liquidweb-go/storm"

	"github.com/liquidweb/docker-machine-driver-liquidweb/util"
)

func (self *Driver) GetSSHKeyPath() string {
	return self.BaseDriver.GetSSHKeyPath()
}

func (self *Driver) GetMachineName() string {
	return self.BaseDriver.GetMachineName()
}

func (self *Driver) GetURL() (string, error) {
	ip, err := self.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s:%d", ip, self.DockerPort), nil
}

func (self *Driver) GetState() (state.State, error) {
	api, err := self.getApiClient()
	if err != nil {
		return state.Error, err
	}
	var status *storm.ServerStatus
	status, err = api.StormServer.Status(self.LwComputeNodeUniqId)
	if err != nil {
		return state.Error, err
	}

	switch strings.ToUpper(status.Status) {
	case "RUNNING":
		return state.Running, nil
	case "SHUTDOWN":
		return state.Stopped, nil
	case "SHUTTING DOWN":
		return state.Stopping, nil
	case "BUILDING", "BOOTING", "Updating Firewall", "Updating Network", "RESTARTING",
		"CLONING", "RESTORING IMAGE", "RE-IMAGING", "RESTORING BACKUP", "MOVING",
		"RESIZING":
		return state.Starting, nil
	}

	return state.Error, nil
}

func (self *Driver) PreCreateCheck() error {
	if self.LwComputeNodeRootPassword == "" {
		log.Info("Generating a random root password...")
		self.LwComputeNodeRootPassword = util.RandomString(30)
	}

	return nil
}

func (self *Driver) GetSSHHostname() (string, error) {
	return self.GetIP()
}

func (self *Driver) GetSSHUsername() string {
	return self.BaseDriver.GetSSHUsername()
}

func (self *Driver) GetSSHPort() (int, error) {
	return self.BaseDriver.GetSSHPort()
}

func (self *Driver) GetIP() (ip string, err error) {
	if self.IPAddress == "" {
		var api *lwclient.API
		if api, err = self.getApiClient(); err != nil {
			return
		}
		assetDetails, assetErr := api.Asset.Details(self.LwComputeNodeUniqId)
		if assetErr != nil {
			err = fmt.Errorf("error fetching primary ip: %s", assetErr)
			return
		}

		self.IPAddress = assetDetails.Ip
	}

	ip = self.IPAddress

	return
}

func (self *Driver) Create() error {
	log.Infof("Creating %s machine instance...", self.DriverName())

	pubKey, err := self.createSshKey()
	if err != nil {
		return err
	}

	api, err := self.getApiClient()
	if err != nil {
		return err
	}

	var server *storm.Server
	ca := storm.ServerParams{
		Domain:       self.LwComputeNodeHostname,
		Zone:         self.LwComputeZoneId,
		Password:     self.LwComputeNodeRootPassword,
		PublicSSHKey: pubKey,
		Template:     self.LwComputeTemplate,
		ConfigID:     self.LwComputeConfigId,
	}
	server, err = api.StormServer.Create(ca)
	if err != nil {
		return err
	}

	self.LwComputeNodeUniqId = server.UniqID

	log.Infof("Created %s instance with uniq_id [%s]", self.DriverName(), self.LwComputeNodeUniqId)

	log.Info("Waiting for machine to become ready..")
	self.waitUntilMachineNodeReady(api)
	log.Info("Machine has become ready..")

	assetDetails, assetErr := api.Asset.Details(self.LwComputeNodeUniqId)
	if assetErr != nil {
		return fmt.Errorf("error fetching primary ip: %s", assetErr)
	}

	self.IPAddress = assetDetails.Ip
	if self.IPAddress != "" {
		log.Infof("Discovered IP address [%s] for uniq_id [%s]", self.IPAddress, self.LwComputeNodeUniqId)
	}

	return nil
}

func (self *Driver) Remove() error {
	api, err := self.getApiClient()
	if err != nil {
		return err
	}

	log.Infof("removing %s compute node [%s]...", self.DriverName(), self.LwComputeNodeUniqId)
	_, err = api.StormServer.Destroy(self.LwComputeNodeUniqId)

	return err
}

func (self *Driver) Kill() error {
	api, err := self.getApiClient()
	if err != nil {
		return err
	}

	log.Infof("forcing %s compute node [%s] to shutdown...",
		self.DriverName(), self.LwComputeNodeUniqId)
	_, err = api.StormServer.Stop(self.LwComputeNodeUniqId, true)

	return err
}

func (self *Driver) Start() error {
	api, err := self.getApiClient()
	if err != nil {
		return err
	}

	log.Infof("requesting %s compute node [%s] to start...",
		self.DriverName(), self.LwComputeNodeUniqId)
	_, err = api.StormServer.Start(self.LwComputeNodeUniqId)

	return err
}

func (self *Driver) Restart() error {
	api, err := self.getApiClient()
	if err != nil {
		return err
	}

	log.Infof("rebooting %s compute node [%s]...", self.DriverName(), self.LwComputeNodeUniqId)
	_, err = api.StormServer.Reboot(self.LwComputeNodeUniqId)

	return err
}

func (self *Driver) Stop() error {
	api, err := self.getApiClient()
	if err != nil {
		return err
	}

	log.Infof("requesting %s compute node [%s] to shutdown...",
		self.DriverName(), self.LwComputeNodeUniqId)
	_, err = api.StormServer.Stop(self.LwComputeNodeUniqId)

	return err
}
