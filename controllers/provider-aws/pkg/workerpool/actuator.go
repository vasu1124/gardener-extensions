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

	awsinternal "github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/apis/aws"
	awsv1alpha1 "github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/apis/aws/v1alpha1"
	"github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/aws"
	"github.com/gardener/gardener-extensions/pkg/controller/provider/workerpool"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// import (
// 	"context"
// 	"fmt"

// 	awsinternal "github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/apis/aws"
// 	awsv1alpha1 "github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/apis/aws/v1alpha1"
// 	"github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/apis/aws/v1alpha1/helper"
// 	"github.com/gardener/gardener-extensions/controllers/provider-aws/pkg/aws"
// 	"github.com/gardener/gardener-extensions/pkg/controller/provider/workerpool"
// 	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
// 	"github.com/gardener/gardener/pkg/chartrenderer"
// 	"github.com/gardener/gardener/pkg/operation"
// 	"github.com/gardener/gardener/pkg/operation/common"
// 	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
// 	machinev1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"

// 	"github.com/go-logr/logr"

// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	"k8s.io/apimachinery/pkg/runtime/serializer"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/rest"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// )

type actuator struct {
	client     client.Client
	kubernetes kubernetes.Interface
	scheme     *runtime.Scheme
	codecs     serializer.CodecFactory
	logger     logr.Logger
}

// NewActuator creates a new Actuator that updates the status of the handled WorkerPoolConfigs.
func NewActuator(logger logr.Logger) workerpool.Actuator {
	return &actuator{logger: logger}
}

func (c *actuator) InjectScheme(scheme *runtime.Scheme) error {
	c.scheme = scheme

	if err := awsinternal.AddToScheme(c.scheme); err != nil {
		return err
	}
	if err := awsv1alpha1.AddToScheme(c.scheme); err != nil {
		return err
	}
	c.codecs = serializer.NewCodecFactory(c.scheme)

	return nil
}

func (c *actuator) InjectClient(client client.Client) error {
	c.client = client
	return nil
}

func (c *actuator) InjectConfig(config *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	c.kubernetes = clientset
	return nil
}

func (c *actuator) GenerateWorkerPool(ctx context.Context, wp *extensionsv1alpha1.WorkerPool) (workerpool.WorkerPool, error) {
	// WIP
	// infrastructureStatus, err := helper.InfrastructureStatusFromProviderStatus(wp.Spec.InfrastructureProviderStatus)
	// if err != nil {
	// 	return nil, err
	// }

	// secret := &corev1.Secret{}
	// if err := c.client.Get(ctx, kutil.Key(wp.Spec.SecretRef.Namespace, wp.Spec.SecretRef.Name), secret); err != nil {
	// 	return nil, err
	// }

	// chartRenderer, err := chartrenderer.New(c.kubernetes)
	// if err != nil {
	// 	return nil, err
	// }

	// var (
	// 	zones   = wp.Spec.Zones
	// 	zoneLen = len(zones)

	// 	machinePlacements []workerpool.MachinePlacement
	// )

	// ami, err := amiFromMachineImage(wp.Spec.MachineImage)
	// if err != nil {
	// 	return nil, err
	// }
	// instanceProfile, err := instanceProfileFromInfrastructureStatus(infrastructureStatus, awsv1alpha1.InfrastructurePurposeNodes)
	// if err != nil {
	// 	return nil, err
	// }
	// role, err := roleFromInfrastructureStatus(infrastructureStatus, awsv1alpha1.InfrastructurePurposeNodes)
	// if err != nil {
	// 	return nil, err
	// }

	// for zoneIndex, zoneName := range zones {
	// 	subnet, err := subnetForZoneFromInfrastructureStatus(infrastructureStatus, awsv1alpha1.InfrastructurePurposeNodes, zoneName)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	securityGroup, err := securityGroupFromInfrastructureStatus(infrastructureStatus, awsv1alpha1.InfrastructurePurposeNodes, zoneName)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	machineClassSpec := map[string]interface{}{
	// 		"ami":                ami,
	// 		"region":             wp.Spec.Region,
	// 		"machineType":        wp.Spec.MachineType,
	// 		"iamInstanceProfile": instanceProfile.Name,
	// 		"keyName":            role.Name,
	// 		"networkInterfaces": []map[string]interface{}{
	// 			{
	// 				"subnetID":         subnet.Id,
	// 				"securityGroupIDs": []string{securityGroup.Id},
	// 			},
	// 		},
	// 		"tags": map[string]string{
	// 			fmt.Sprintf("kubernetes.io/cluster/%s", wp.Namespace): "1",
	// 			"kubernetes.io/role/node":                             "1",
	// 		},
	// 		"secret": map[string]interface{}{
	// 			"cloudConfig": wp.Spec.UserData,
	// 		},
	// 		"blockDevices": []map[string]interface{}{
	// 			{
	// 				"ebs": map[string]interface{}{
	// 					"volumeSize": common.DiskSize(wp.Spec.Volume.Size),
	// 					"volumeType": wp.Spec.Volume.Type,
	// 				},
	// 			},
	// 		},
	// 	}

	// 	var (
	// 		machineClassSpecHash = common.MachineClassHash(machineClassSpec, "TODO") // k8s majorminor version
	// 		deploymentName       = fmt.Sprintf("%s-%s-z%d", wp.Namespace, wp.Name, zoneIndex+1)
	// 		className            = fmt.Sprintf("%s-%s", deploymentName, machineClassSpecHash)
	// 		secretData           = generateMachineClassSecretData()
	// 	)

	// 	machineDeploymentSpec := map[string]interface{}{
	// 		"name": deploymentName,
	// 		"replicas":

	// 	machineDeployments = append(machineDeployments, operation.MachineDeployment{
	// 		Name:           deploymentName,
	// 		ClassName:      className,
	// 		Minimum:        common.DistributeOverZones(zoneIndex, wp.Spec.Minimum, zoneLen),
	// 		Maximum:        common.DistributeOverZones(zoneIndex, wp.Spec.Maximum, zoneLen),
	// 		MaxSurge:       common.DistributePositiveIntOrPercent(zoneIndex, wp.Spec.MaxSurge, zoneLen, wp.Spec.Maximum),
	// 		MaxUnavailable: common.DistributePositiveIntOrPercent(zoneIndex, wp.Spec.MaxUnavailable, zoneLen, wp.Spec.Minimum),
	// 	})

	// 	machineClassSpec["name"] = className
	// 	machineClassSpec["secret"].(map[string]interface{})["accessKeyID"] = string(secretData[machinev1alpha1.AWSAccessKeyID])
	// 	machineClassSpec["secret"].(map[string]interface{})["secretAccessKey"] = string(secretData[machinev1alpha1.AWSSecretAccessKey])

	// 	machineClasses = append(machineClasses, machineClassSpec)
	// }

	// fmt.Println(infrastructureStatus.EC2.KeyName)

	fmt.Println("have generated wp config")
	return workerpool.WorkerPool{}, nil
}

func generateMachineClassSecretData(secret *corev1.Secret) map[string][]byte {
	return map[string][]byte{
		machinev1alpha1.AWSAccessKeyID:     secret.Data[aws.AccessKeyID],
		machinev1alpha1.AWSSecretAccessKey: secret.Data[aws.SecretAccessKey],
	}
}

func amiFromMachineImage(machineImage extensionsv1alpha1.MachineImage) (string, error) {
	// TODO: implement properly
	return "ami-0628e483315b5d17e", nil
}

func instanceProfileFromInfrastructureStatus(infrastructureStatus *awsv1alpha1.InfrastructureStatus, purpose string) (*awsv1alpha1.InstanceProfile, error) {
	for _, instanceProfile := range infrastructureStatus.IAM.InstanceProfiles {
		if instanceProfile.Purpose == purpose {
			return &instanceProfile, nil
		}
	}

	return nil, fmt.Errorf("could not find instance profile for purpose %s", purpose)
}

func roleFromInfrastructureStatus(infrastructureStatus *awsv1alpha1.InfrastructureStatus, purpose string) (*awsv1alpha1.Role, error) {
	for _, role := range infrastructureStatus.IAM.Roles {
		if role.Purpose == purpose {
			return &role, nil
		}
	}

	return nil, fmt.Errorf("could not find role for purpose %s", purpose)
}

func subnetForZoneFromInfrastructureStatus(infrastructureStatus *awsv1alpha1.InfrastructureStatus, purpose, zone string) (*awsv1alpha1.Subnet, error) {
	for _, subnet := range infrastructureStatus.VPC.Subnets {
		if subnet.Purpose == purpose && subnet.Zone == zone {
			return &subnet, nil
		}
	}

	return nil, fmt.Errorf("could not find subnet for zone %s and purpose %s", zone, purpose)
}

func securityGroupFromInfrastructureStatus(infrastructureStatus *awsv1alpha1.InfrastructureStatus, purpose string) (*awsv1alpha1.SecurityGroup, error) {
	for _, securityGroup := range infrastructureStatus.VPC.SecurityGroups {
		if securityGroup.Purpose == purpose {
			return &securityGroup, nil
		}
	}

	return nil, fmt.Errorf("could not find security group for purpose %s", purpose)
}
