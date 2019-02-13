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

// workerPoolReconciler reconciles WorkerPool resources of Gardener's
// `extensions.gardener.cloud` API group.
type workerPoolReconciler struct {
	logger   logr.Logger
	actuator Actuator

	ctx    context.Context
	client client.Client
}

var _ reconcile.Reconciler = &workerPoolReconciler{}

// NewReconciler creates a new reconcile.Reconciler that reconciles
// WorkerPool resources of Gardener's `extensions.gardener.cloud` API group.
func NewReconciler(logger logr.Logger, actuator Actuator) reconcile.Reconciler {
	return &workerPoolReconciler{logger: logger, actuator: actuator}
}

// InjectFunc enables dependency injection into the actuator.
func (r *workerPoolReconciler) InjectFunc(f inject.Func) error {
	return f(r.actuator)
}

// InjectClient injects the controller runtime client into the reconciler.
func (r *workerPoolReconciler) InjectClient(client client.Client) error {
	r.client = client
	return nil
}

// InjectStopChannel is an implementation for getting the respective stop channel managed by the controller-runtime.
func (r *workerPoolReconciler) InjectStopChannel(stopCh <-chan struct{}) error {
	r.ctx = controller.ContextFromStopChannel(stopCh)
	return nil
}

// Reconcile is the reconciler function that gets executed in case there are new events for the `WorkerPool`
// resources.
func (r *workerPoolReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	wp := &extensionsv1alpha1.WorkerPool{}
	if err := r.client.Get(r.ctx, request.NamespacedName, wp); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		r.logger.Error(err, "Could not fetch WorkerPool")
		return reconcile.Result{}, err
	}

	if wp.DeletionTimestamp != nil {
		return r.delete(r.ctx, wp)
	}
	return r.reconcile(r.ctx, wp)
}

func (r *workerPoolReconciler) reconcile(ctx context.Context, wp *extensionsv1alpha1.WorkerPool) (reconcile.Result, error) {
	// Add finalizer to resource if not yet done.
	if finalizers := sets.NewString(wp.Finalizers...); !finalizers.Has(FinalizerName) {
		finalizers.Insert(FinalizerName)
		wp.Finalizers = finalizers.UnsortedList()
		if err := r.client.Update(ctx, wp); err != nil {
			return reconcile.Result{}, err
		}
	}

	r.logger.Info("Reconciling worker pool triggers idempotent reconciliation.", "wp", wp.Name)

	// XXX: main logic begin
	_, err := r.actuator.GenerateWorkerPool(ctx, wp)
	if err != nil {
		r.logger.Error(err, "Unable to reconcile worker pool", "wp", wp.Name)
		return controller.ReconcileErr(err)
	}
	// XXX: main logic end

	return reconcile.Result{}, nil
}

func (r *workerPoolReconciler) delete(ctx context.Context, wp *extensionsv1alpha1.WorkerPool) (reconcile.Result, error) {
	finalizers := sets.NewString(wp.Finalizers...)
	if !finalizers.Has(FinalizerName) {
		r.logger.Info("Reconciling worker pool causes a no-op as there is no finalizer.", "wp", wp.Name)
		return reconcile.Result{}, nil
	}

	// if err := r.actuator.Delete(ctx, wp); err != nil {
	// 	r.logger.Error(err, "Error deleting worker pool", "wp", wp.Name)
	// 	return controller.ReconcileErr(err)
	// }

	r.logger.Info("worker pool deletion successful, removing finalizer.", "wp", wp.Name)
	finalizers.Delete(FinalizerName)
	wp.Finalizers = finalizers.UnsortedList()
	if err := r.client.Update(ctx, wp); err != nil {
		r.logger.Error(err, "Error removing finalizer from worker pool", "wp", wp.Name)
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}
