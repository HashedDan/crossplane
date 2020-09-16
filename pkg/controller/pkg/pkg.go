/*
Copyright 2020 The Crossplane Authors.

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

package pkg

import (
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	"github.com/crossplane/crossplane/pkg/controller/pkg/manager"
	"github.com/crossplane/crossplane/pkg/controller/pkg/revision"
)

// Setup package controllers.
func Setup(mgr ctrl.Manager, h *rest.Config, l logging.Logger, namespace string) error {
	for _, setup := range []func(ctrl.Manager, *rest.Config, logging.Logger, string) error{
		manager.SetupConfiguration,
		manager.SetupProvider,
		revision.SetupConfigurationRevision,
		revision.SetupProviderRevision,
	} {
		if err := setup(mgr, h, l, namespace); err != nil {
			return err
		}
	}
	return nil
}
