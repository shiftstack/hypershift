package openstack

import (
	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"

	utilpointer "k8s.io/utils/pointer"
	capiopenstack "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

func MachineTemplateSpec(hcluster *hyperv1.HostedCluster, nodePool *hyperv1.NodePool) (*capiopenstack.OpenStackMachineTemplateSpec, error) {
	openStackMachineTemplate := &capiopenstack.OpenStackMachineTemplateSpec{Template: capiopenstack.OpenStackMachineTemplateResource{Spec: capiopenstack.OpenStackMachineSpec{
		Flavor: nodePool.Spec.Platform.OpenStack.Flavor,
	}}}

	// TODO: figure out how to get the same image as the cluster (release payload)
	if nodePool.Spec.Platform.OpenStack.ImageName != "" {
		openStackMachineTemplate.Template.Spec.Image.Filter = &capiopenstack.ImageFilter{
			Name: utilpointer.String(nodePool.Spec.Platform.OpenStack.ImageName),
		}
	}

	return openStackMachineTemplate, nil
}
