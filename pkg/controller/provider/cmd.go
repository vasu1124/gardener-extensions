// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"flag"

	"github.com/gardener/gardener-extensions/pkg/controller/cmd"
	"github.com/gardener/gardener-extensions/pkg/controller/version"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ExtensionsScheme is the default scheme for extensions, consisting of all Kubernetes built-in
// schemes (client-go/kubernetes/scheme) and the extensions/v1alpha1 scheme.
var ExtensionsScheme = runtime.NewScheme()

func init() {
	utilruntime.Must(scheme.AddToScheme(ExtensionsScheme))
	utilruntime.Must(extensionsv1alpha1.AddToScheme(ExtensionsScheme))
}

// CommandOptions are options used for creating an provider command.
type CommandOptions struct {
	Manager              *ManagerOptions
	WorkerPoolController *WorkerPoolConfigControllerOptions
}

// NewCommandOptions creates new CommandOptions with the given name, type name and actuator factory.
func NewCommandOptions(name, typeName string) *CommandOptions {
	return &CommandOptions{
		Manager:              NewManagerOptions(name),
		WorkerPoolController: NewWorkerPoolConfigControllerOptions(name, typeName),
	}
}

// Flags yields a NamedFlagSet with all subcomponents relevant for an provider command.
func (c *CommandOptions) Flags() cmd.NamedFlagSet {
	fss := cmd.NamedFlagSet{}

	c.WorkerPoolController.AddFlags(fss.FlagSet("controller"))

	fs := fss.FlagSet("misc")
	fs.AddGoFlagSet(flag.CommandLine)

	return fss
}

// Config produces a new CommandConfig used for creating a provider command.
func (c *CommandOptions) Config() (*CommandConfig, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	mgrConfig, err := c.Manager.Config()
	if err != nil {
		return nil, err
	}

	workerPoolCtrlConfig, err := c.WorkerPoolController.Config()
	if err != nil {
		return nil, err
	}

	return &CommandConfig{
		REST:                 restConfig,
		Manager:              mgrConfig,
		WorkerPoolController: workerPoolCtrlConfig,
	}, nil
}

// CommandConfig is the configuration for creating a provider command.
type CommandConfig struct {
	REST                 *rest.Config
	Manager              *ManagerConfig
	WorkerPoolController *ControllerConfig
}

// Complete fills in any fields not set that are required to have valid data.
func (c *CommandConfig) Complete() *CompletedConfig {
	return &CompletedConfig{&completedConfig{c}}
}

type completedConfig struct {
	*CommandConfig
}

// CompletedConfig is the completed config (all fields set) used to run an provider command.
type CompletedConfig struct {
	*completedConfig
}

// Run runs the provider command with the given completed configuration.
func Run(ctx context.Context, config *CompletedConfig) error {
	log := config.Manager.Log.WithName("entrypoint")
	log.Info("Gardener Controller Extensions", "version", version.Version)

	mgr, err := manager.New(config.REST, config.Manager.Options)
	if err != nil {
		log.Error(err, "Could not instantiate manager")
		return err
	}

	// Start worker pool config controller
	ctrl, err := controller.New(config.WorkerPoolController.Name, mgr, config.WorkerPoolController.Options)
	if err != nil {
		log.Error(err, "Could not instantiate controller")
		return err
	}

	if err := ctrl.Watch(&source.Kind{Type: &extensionsv1alpha1.WorkerPool{}}, &handler.EnqueueRequestForObject{}, config.WorkerPoolController.Predicates...); err != nil {
		log.Error(err, "Could not watch worker pool configs")
		return err
	}

	// In the future additional controllers could be started here (e.g., the Infrastructure controller).

	return mgr.Start(ctx.Done())
}
