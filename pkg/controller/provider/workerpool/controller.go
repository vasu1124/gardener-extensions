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

	"github.com/gardener/gardener-extensions/pkg/controller"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	"github.com/go-logr/logr"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

const (
	// FinalizerName is the name of the finalizer written by this controller.
	FinalizerName = "extensions.gardener.cloud/workerpoolconfigs"
)

// workerPoolConfigReconciler reconciles WorkerPool resources of Gardener's
// `extensions.gardener.cloud` API group.
type workerPoolConfigReconciler struct {
	logger   logr.Logger
	actuator Actuator

	ctx    context.Context
	client client.Client
}

var _ reconcile.Reconciler = &workerPoolConfigReconciler{}

// NewReconciler creates a new reconcile.Reconciler that reconciles
// WorkerPool resources of Gardener's `extensions.gardener.cloud` API group.
func NewReconciler(logger logr.Logger, actuator Actuator) reconcile.Reconciler {
	return &workerPoolConfigReconciler{logger: logger, actuator: actuator}
}

// InjectFunc enables dependency injection into the actuator.
func (r *workerPoolConfigReconciler) InjectFunc(f inject.Func) error {
	return f(r.actuator)
}

// InjectClient injects the controller runtime client into the reconciler.
func (r *workerPoolConfigReconciler) InjectClient(client client.Client) error {
	r.client = client
	return nil
}

// InjectStopChannel is an implementation for getting the respective stop channel managed by the controller-runtime.
func (r *workerPoolConfigReconciler) InjectStopChannel(stopCh <-chan struct{}) error {
	r.ctx = controller.ContextFromStopChannel(stopCh)
	return nil
}

// Reconcile is the reconciler function that gets executed in case there are new events for the `WorkerPool`
// resources.
func (r *workerPoolConfigReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	wpc := &extensionsv1alpha1.WorkerPool{}
	if err := r.client.Get(r.ctx, request.NamespacedName, wpc); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		r.logger.Error(err, "Could not fetch WorkerPool")
		return reconcile.Result{}, err
	}

	if wpc.DeletionTimestamp != nil {
		return r.delete(r.ctx, wpc)
	}
	return r.reconcile(r.ctx, wpc)
}

func (r *workerPoolConfigReconciler) reconcile(ctx context.Context, wpc *extensionsv1alpha1.WorkerPool) (reconcile.Result, error) {
	// Add finalizer to resource if not yet done.
	if finalizers := sets.NewString(wpc.Finalizers...); !finalizers.Has(FinalizerName) {
		finalizers.Insert(FinalizerName)
		wpc.Finalizers = finalizers.UnsortedList()
		if err := r.client.Update(ctx, wpc); err != nil {
			return reconcile.Result{}, err
		}
	}

	exist, err := r.actuator.Exists(ctx, wpc)
	if err != nil {
		return reconcile.Result{}, err
	}

	if exist {
		r.logger.Info("Reconciling operating system config triggers idempotent update.", "wpc", wpc.Name)
		if err := r.actuator.Update(ctx, wpc); err != nil {
			return controller.ReconcileErr(err)
		}
		return reconcile.Result{}, nil
	}

	r.logger.Info("Reconciling operating system config triggers idempotent create.", "wpc", wpc.Name)
	if err := r.actuator.Create(ctx, wpc); err != nil {
		r.logger.Error(err, "Unable to create operating system config", "wpc", wpc.Name)
		return controller.ReconcileErr(err)
	}
	return reconcile.Result{}, nil
}

func (r *workerPoolConfigReconciler) delete(ctx context.Context, wpc *extensionsv1alpha1.WorkerPool) (reconcile.Result, error) {
	finalizers := sets.NewString(wpc.Finalizers...)
	if !finalizers.Has(FinalizerName) {
		r.logger.Info("Reconciling operating system config causes a no-op as there is no finalizer.", "wpc", wpc.Name)
		return reconcile.Result{}, nil
	}

	if err := r.actuator.Delete(ctx, wpc); err != nil {
		r.logger.Error(err, "Error deleting operating system config", "wpc", wpc.Name)
		return controller.ReconcileErr(err)
	}

	r.logger.Info("Operating system config deletion successful, removing finalizer.", "wpc", wpc.Name)
	finalizers.Delete(FinalizerName)
	wpc.Finalizers = finalizers.UnsortedList()
	if err := r.client.Update(ctx, wpc); err != nil {
		r.logger.Error(err, "Error removing finalizer from operating system config", "wpc", wpc.Name)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
