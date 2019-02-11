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

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/aws"
	awsworkerpoolconfig "github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/workerpool"
	"github.com/gardener/gardener-extensions/pkg/controller/provider"
	"github.com/gardener/gardener-extensions/pkg/controller/provider/workerpool"
	"github.com/spf13/cobra"
)

// Name is the name of the AWS controller.
const Name = "provider-aws"

// WorkerPoolConfigActuatorFactory creates a new worker pool config actuator.
func WorkerPoolConfigActuatorFactory(args *provider.ActuatorArgs) (workerpool.Actuator, error) {
	return awsworkerpoolconfig.NewActuator(args.Log), nil
}

// NewControllerCommand creates a new CoreOS controller command.
func NewControllerCommand(ctx context.Context) *cobra.Command {
	opts := provider.NewCommandOptions(Name, aws.Type)
	opts.WorkerPoolController.ActuatorFactory = WorkerPoolConfigActuatorFactory
	opts.Manager.LeaderElection = true
	opts.Manager.LeaderElectionNamespace = os.Getenv("LEADER_ELECTION_NAMESPACE")

	cmd := &cobra.Command{
		Use: "provider-aws-controller-manager",

		Run: func(cmd *cobra.Command, args []string) {
			c, err := opts.Config()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := provider.Run(ctx, c.Complete()); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()
	for _, f := range opts.Flags().FlagSets {
		fs.AddFlagSet(f)
	}

	return cmd
}
