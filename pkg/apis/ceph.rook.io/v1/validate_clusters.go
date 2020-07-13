/*
Copyright 2020 The Rook Authors. All rights reserved.

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

package v1

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strconv"
)

// compile-time assertions ensures CephCluster implements webhook.Validator so a webhook builder
// will be registered for the validating webhook.
var _ webhook.Validator = &CephCluster{}

func (c *CephCluster) ValidateCreate() error {
	logger.Info("validate create cephcluster")
	//If external mode enabled, then check if other fields are empty
	if c.Spec.External.Enable {
		if c.Spec.Mon != (MonSpec{}) || c.Spec.Dashboard != (DashboardSpec{}) || c.Spec.Monitoring != (MonitoringSpec{}) || c.Spec.DisruptionManagement != (DisruptionManagementSpec{}) || len(c.Spec.Mgr.Modules) > 0 || len(c.Spec.Network.Provider) > 0 || len(c.Spec.Network.Selectors) > 0 {
			return errors.New("invalid create : external mode enabled cannot have mon,dashboard,monitoring,network,disruptionManagement,storage fields in CR")
		}
	}
	return nil
}

func (c *CephCluster) ValidateUpdate(old runtime.Object) error {
	logger.Info("validate update cephcluster")
	occ := old.(*CephCluster)
	return validateUpdatedCephCluster(c, occ)
}

func (c *CephCluster) ValidateDelete() error {
	return nil
}

func validateUpdatedCephCluster(updatedCephCluster *CephCluster, found *CephCluster) error {
	if updatedCephCluster.Spec.DataDirHostPath != found.Spec.DataDirHostPath {
		return errors.Errorf("invalid update: DataDirHostPath change from %q to %q is not allowed", found.Spec.DataDirHostPath, updatedCephCluster.Spec.DataDirHostPath)
	}

	if updatedCephCluster.Spec.Network.HostNetwork != found.Spec.Network.HostNetwork {
		return errors.Errorf("invalid update: HostNetwork change from %q to %q is not allowed", strconv.FormatBool(found.Spec.Network.HostNetwork), strconv.FormatBool(updatedCephCluster.Spec.Network.HostNetwork))
	}

	if updatedCephCluster.Spec.Network.Provider != found.Spec.Network.Provider {
		return errors.Errorf("invalid update: Provider change from %q to %q is not allowed", found.Spec.Network.Provider, updatedCephCluster.Spec.Network.Provider)
	}

	return nil
}
