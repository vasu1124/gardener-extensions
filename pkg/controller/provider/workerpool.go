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
	"fmt"

	"github.com/gardener/gardener-extensions/pkg/controller/provider/workerpool"
	"github.com/spf13/pflag"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type WorkerPoolConfigControllerOptions struct {
	*ControllerOptions
	ActuatorFactory WorkerPoolConfigActuatorFactory
}

// NewWorkerPoolConfigControllerOptions creates new ControllerOptions with the given name, type name and
// actuator factory.
func NewWorkerPoolConfigControllerOptions(name, typeName string) *WorkerPoolConfigControllerOptions {
	return &WorkerPoolConfigControllerOptions{
		ControllerOptions: &ControllerOptions{
			Name:                    name,
			Type:                    typeName,
			MaxConcurrentReconciles: DefaultMaxConcurrentReconciles,
		},
	}
}

// AddFlags adds all ControllerOptions relevant flags to the given FlagSet.
func (c *WorkerPoolConfigControllerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&c.MaxConcurrentReconciles, "concurrent-workerpool-syncs", c.MaxConcurrentReconciles, "The maximum number of concurrent reconciliations.")
}

// Config produces a WorkerPoolConfigControllerOptions used for instantiating a Controller.
func (c *WorkerPoolConfigControllerOptions) Config() (*ControllerConfig, error) {
	log := c.Log
	if log == nil {
		log = logf.Log
	}
	log = log.WithName(c.Name)

	if c.ActuatorFactory == nil {
		return nil, fmt.Errorf("cannot create configuration for worker pool config controller: missing ActuatorFactory")
	}

	actuator, err := c.ActuatorFactory(&ActuatorArgs{Log: log.WithName("workerpool-actuator")})
	if err != nil {
		return nil, err
	}

	predicates := c.Predicates
	if predicates == nil {
		predicates = []predicate.Predicate{workerpool.GenerationChangedPredicate()}
	}
	predicates = append(predicates, workerpool.TypePredicate(c.Type))

	return &ControllerConfig{
		Name: c.Name,
		Log:  log.WithName("controller"),
		Options: controller.Options{
			MaxConcurrentReconciles: c.MaxConcurrentReconciles,
			Reconciler:              workerpool.NewReconciler(log.WithName("workerpool-reconciler"), actuator),
		},
		Predicates: predicates,
	}, nil
}

// WorkerPoolConfigActuatorFactory is a factory used for creating Actuators.
type WorkerPoolConfigActuatorFactory func(*ActuatorArgs) (workerpool.Actuator, error)
