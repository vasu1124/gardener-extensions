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

package workerpool

import (
	"context"
	"fmt"

	"github.com/gardener/gardener-extensions/pkg/controller/provider/workerpool"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type actuator struct {
	client client.Client
	scheme *runtime.Scheme
	logger logr.Logger
}

// NewActuator creates a new Actuator that updates the status of the handled WorkerPoolConfigs.
func NewActuator(logger logr.Logger) workerpool.Actuator {
	return &actuator{logger: logger}
}

func (c *actuator) InjectScheme(scheme *runtime.Scheme) error {
	c.scheme = scheme
	return nil
}

func (c *actuator) InjectClient(client client.Client) error {
	c.client = client
	return nil
}

func (c *actuator) Exists(ctx context.Context, config *extensionsv1alpha1.WorkerPool) (bool, error) {
	return config.Status.LastOperation != nil, nil
}

func (c *actuator) Create(ctx context.Context, config *extensionsv1alpha1.WorkerPool) error {
	return c.reconcile(ctx, config)
}

func (c *actuator) Update(ctx context.Context, config *extensionsv1alpha1.WorkerPool) error {
	return c.reconcile(ctx, config)
}

func (c *actuator) Delete(ctx context.Context, config *extensionsv1alpha1.WorkerPool) error {
	return c.delete(ctx, config)
}

func (c *actuator) reconcile(ctx context.Context, config *extensionsv1alpha1.WorkerPool) error {
	fmt.Println("have reconciled")
	return nil
}

func (c *actuator) delete(ctx context.Context, config *extensionsv1alpha1.WorkerPool) error {
	fmt.Println("have deleted")
	return nil
}
