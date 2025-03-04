package openstack

import (
	"fmt"

	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"

	"k8s.io/utils/ptr"
	capiopenstackv1beta1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

func MachineTemplateSpec(hcluster *hyperv1.HostedCluster, nodePool *hyperv1.NodePool) (*capiopenstackv1beta1.OpenStackMachineTemplateSpec, error) {
	openStackMachineTemplate := &capiopenstackv1beta1.OpenStackMachineTemplateSpec{Template: capiopenstackv1beta1.OpenStackMachineTemplateResource{Spec: capiopenstackv1beta1.OpenStackMachineSpec{
		Flavor: ptr.To(nodePool.Spec.Platform.OpenStack.Flavor),
	}}}

	if nodePool.Spec.Platform.OpenStack.ImageName != "" {
		openStackMachineTemplate.Template.Spec.Image.Filter = &capiopenstackv1beta1.ImageFilter{
			Name: ptr.To(nodePool.Spec.Platform.OpenStack.ImageName),
		}
	} else {
		// TODO(emilien): Add support for using the image from the release payload.
		// This will be possible when CAPO supports managing images in the OpenStack cluster:
		// https://github.com/kubernetes-sigs/cluster-api-provider-openstack/pull/2130
		// For 4.17 we might leave this as is and let the user provide the image name as
		// we plan to deliver the OpenStack provider as a dev preview.
		return nil, fmt.Errorf("image name is required")
	}

	return openStackMachineTemplate, nil
}
