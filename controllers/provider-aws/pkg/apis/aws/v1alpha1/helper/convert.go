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

package helper

import (
	awsv1alpha1 "github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/apis/aws/v1alpha1"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
)

// InfrastructureStatusFromProviderStatus converts a provider status into a awsv1alpha1 InfrastructureStatus.
func InfrastructureStatusFromProviderStatus(providerStatus *runtime.RawExtension) (*awsv1alpha1.InfrastructureStatus, error) {
	var infrastructureStatus awsv1alpha1.InfrastructureStatus
	if err := yaml.Unmarshal(providerStatus.Raw, &infrastructureStatus); err != nil {
		return nil, err
	}
	return &infrastructureStatus, nil
}

// InfrastructureStatusToProviderStatus converts a awsv1alpha1 InfrastructureStatus into a provider status.
func InfrastructureStatusToProviderStatus(in *awsv1alpha1.InfrastructureStatus) (*runtime.RawExtension, error) {
	out, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}
	return &runtime.RawExtension{Raw: out}, nil
}
