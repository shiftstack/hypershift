package openstack

import (
	"context"

	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"
	"github.com/openshift/hypershift/cmd/nodepool/core"

	"github.com/spf13/cobra"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type OpenStackPlatformCreateOptions struct {
	Flavor    string
	ImageName string
}

func NewCreateCommand(coreOpts *core.CreateNodePoolOptions) *cobra.Command {
	platformOpts := &OpenStackPlatformCreateOptions{}

	cmd := &cobra.Command{
		Use:          "openstack",
		Short:        "Creates basic functional NodePool resources for OpenStack platform",
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&platformOpts.Flavor, "flavor", platformOpts.Flavor, "The flavor to use for the nodepool")
	cmd.Flags().StringVar(&platformOpts.ImageName, "image-name", platformOpts.ImageName, "The image name to use for the nodepool")

	cmd.RunE = coreOpts.CreateRunFunc(platformOpts)

	return cmd
}

func (o *OpenStackPlatformCreateOptions) UpdateNodePool(_ context.Context, nodePool *hyperv1.NodePool, _ *hyperv1.HostedCluster, _ crclient.Client) error {
	nodePool.Spec.Platform.Type = hyperv1.OpenStackPlatform
	nodePool.Spec.Platform.OpenStack = &hyperv1.OpenStackNodePoolPlatform{
		Flavor:    o.Flavor,
		ImageName: o.ImageName,
	}
	return nil
}

func (o *OpenStackPlatformCreateOptions) Type() hyperv1.PlatformType {
	return hyperv1.OpenStackPlatform
}
