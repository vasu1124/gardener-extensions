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
	"github.com/go-logr/logr"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// DefaultMaxConcurrentReconciles is the default number of maximum concurrent reconciles for.
const DefaultMaxConcurrentReconciles = 5

// ControllerOptions are options used for the creation of a Controller.
type ControllerOptions struct {
	Name                    string
	Log                     logr.Logger
	Type                    string
	Predicates              []predicate.Predicate
	MaxConcurrentReconciles int
}

// ControllerConfig is the configuration for creating a provider controller.
type ControllerConfig struct {
	Name       string
	Log        logr.Logger
	Predicates []predicate.Predicate
	Options    controller.Options
}

// ActuatorArgs are arguments given to the instantiation of an Actuator.
type ActuatorArgs struct {
	Log logr.Logger
}
