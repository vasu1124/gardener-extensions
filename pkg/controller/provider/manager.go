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
	"time"

	"github.com/go-logr/logr"
	"github.com/spf13/pflag"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// ManagerOptions are options for the creation of a Manager.
type ManagerOptions struct {
	Name                    string
	Log                     logr.Logger
	Scheme                  *runtime.Scheme
	LeaderElection          bool
	LeaderElectionID        string
	LeaderElectionNamespace string
	SyncPeriod              *time.Duration
}

// NewManagerOptions creates new ManagerOptions with the given name.
func NewManagerOptions(name string) *ManagerOptions {
	return &ManagerOptions{
		Name:                    name,
		LeaderElectionID:        fmt.Sprintf("%s-leader-election", name),
		LeaderElectionNamespace: v1.NamespaceSystem,
	}
}

// Config produces a ManagerConfig used for instantiating a Manager.
func (m *ManagerOptions) Config() (*ManagerConfig, error) {
	mgrScheme := m.Scheme
	if mgrScheme == nil {
		mgrScheme = ExtensionsScheme
	}

	opts := manager.Options{
		SyncPeriod: m.SyncPeriod,
		Scheme:     mgrScheme,
	}

	opts.LeaderElection = false //FIXME m.LeaderElection
	opts.LeaderElectionID = m.LeaderElectionID
	opts.LeaderElectionNamespace = m.LeaderElectionNamespace

	log := m.Log
	if log == nil {
		log = logf.Log
	}
	log = log.WithName(m.Name)

	return &ManagerConfig{
		Options: opts,
		Log:     log,
	}, nil
}

// AddFlags adds all ManagerOptions relevant flags to the given FlagSet.
func (m *ManagerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&m.LeaderElection, "leader-election", m.LeaderElection, "Whether to use leader election or not when running this controller manager.")
	fs.StringVar(&m.LeaderElectionID, "leader-election-id", m.LeaderElectionID, "The leader election id to use.")
	fs.StringVar(&m.LeaderElectionNamespace, "leader-election-namespace", m.LeaderElectionNamespace, "The namespace to do leader election in.")
}

// ManagerConfig is the configuration for creating a provider manager.
type ManagerConfig struct {
	Options manager.Options
	Log     logr.Logger
}
