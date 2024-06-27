package openstack

import (
	"context"
	"fmt"
	"os"

	"github.com/openshift/hypershift/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"
	"github.com/openshift/hypershift/cmd/cluster/core"
	corev1 "k8s.io/api/core/v1"
)

func DefaultOptions() *CreateOptions {
	return &CreateOptions{}
}

func BindOptions(opts *CreateOptions, flags *pflag.FlagSet) {
	flags.StringVar(&opts.OpenStackCredentialsFile, "openstack-credentials-file", opts.OpenStackCACertFile, "Path to the OpenStack credentials file (required)")
	flags.StringVar(&opts.OpenStackCACertFile, "openstack-ca-cert-file", opts.OpenStackCACertFile, "Path to the OpenStack CA certificate file (optional)")
	flags.StringVar(&opts.OpenStackExternalNetworkName, "openstack-external-network-name", opts.OpenStackExternalNetworkName, "Name of the OpenStack external network (optional)")
	flags.StringVar(&opts.OpenStackNodeFlavor, "openstack-node-flavor", opts.OpenStackNodeFlavor, "The flavor to use for OpenStack nodes")
}

type CreateOptions struct {
	OpenStackCredentialsFile     string
	OpenStackCACertFile          string
	OpenStackExternalNetworkName string
	OpenStackNodeFlavor          string

	externalDNSDomain string

	name, namespace       string
	baseDomainPassthrough bool
}

func (o *CreateOptions) Validate(ctx context.Context, opts *core.CreateOptions) error {
	// Check that the OpenStack credentials file arg is set and that the file exists
	if o.OpenStackCredentialsFile == "" {
		return fmt.Errorf("OpenStack credentials file is required")
	}
	if _, err := os.Stat(o.OpenStackCredentialsFile); err != nil {
		return fmt.Errorf("OpenStack credentials file does not exist: %w", err)
	}
	// TODO(emilien) we will probably have to check the cloud name here ("openstack" by default?)

	return nil
}

func (o *CreateOptions) Complete(ctx context.Context, opts *core.CreateOptions) error {
	o.externalDNSDomain = opts.ExternalDNSDomain
	o.baseDomainPassthrough = opts.BaseDomain == ""
	o.name, o.namespace = opts.Name, opts.Namespace

	return nil
}

func (o *CreateOptions) ApplyPlatformSpecifics(cluster *hyperv1.HostedCluster) error {
	cluster.Spec.Platform = hyperv1.PlatformSpec{
		Type:      hyperv1.OpenStackPlatform,
		OpenStack: &hyperv1.OpenStackPlatformSpec{},
	}

	cluster.Spec.Platform.OpenStack.IdentityRef = hyperv1.OpenStackIdentityReference{
		Name:      credentialsSecret(cluster.Namespace, cluster.Name).Name,
		CloudName: "openstack", // TODO: make this configurable or at least check the clouds.yaml file
	}

	if o.OpenStackExternalNetworkName != "" {
		cluster.Spec.Platform.OpenStack.ExternalNetwork = &hyperv1.NetworkParam{
			Filter: &hyperv1.NetworkFilter{
				Name: o.OpenStackExternalNetworkName,
			},
		}
	}

	cluster.Spec.Services = core.GetIngressServicePublishingStrategyMapping(cluster.Spec.Networking.NetworkType, o.externalDNSDomain != "")
	if o.externalDNSDomain != "" {
		for i, svc := range cluster.Spec.Services {
			switch svc.Service {
			case hyperv1.APIServer:
				cluster.Spec.Services[i].Route = &hyperv1.RoutePublishingStrategy{
					Hostname: fmt.Sprintf("api-%s.%s", cluster.Name, o.externalDNSDomain),
				}

			case hyperv1.OAuthServer:
				cluster.Spec.Services[i].Route = &hyperv1.RoutePublishingStrategy{
					Hostname: fmt.Sprintf("oauth-%s.%s", cluster.Name, o.externalDNSDomain),
				}

			case hyperv1.Konnectivity:
				cluster.Spec.Services[i].Route = &hyperv1.RoutePublishingStrategy{
					Hostname: fmt.Sprintf("konnectivity-%s.%s", cluster.Name, o.externalDNSDomain),
				}
			case hyperv1.Ignition:
				cluster.Spec.Services[i].Route = &hyperv1.RoutePublishingStrategy{
					Hostname: fmt.Sprintf("ignition-%s.%s", cluster.Name, o.externalDNSDomain),
				}
			case hyperv1.OVNSbDb:
				cluster.Spec.Services[i].Route = &hyperv1.RoutePublishingStrategy{
					Hostname: fmt.Sprintf("ovn-sbdb-%s.%s", cluster.Name, o.externalDNSDomain),
				}
			}
		}
	}

	return nil
}

func (o *CreateOptions) GenerateNodePools(constructor core.DefaultNodePoolConstructor) []*hyperv1.NodePool {
	return nil
}

func (o *CreateOptions) GenerateResources() ([]client.Object, error) {
	credentialsContents, err := os.ReadFile(o.OpenStackCredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenStack credentials file: %w", err)
	}
	credentialsSecret := credentialsSecret(o.namespace, o.name)
	credentialsSecret.Data = map[string][]byte{
		"clouds.yaml": credentialsContents,
	}

	if o.OpenStackCACertFile != "" {
		caCertContents, err := os.ReadFile(o.OpenStackCACertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read OpenStack CA certificate file: %w", err)
		}
		credentialsSecret.Data["cacert"] = caCertContents
	}

	return []client.Object{credentialsSecret}, nil
}

func credentialsSecret(namespace, name string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-cloud-credentials",
			Namespace: namespace,
			Labels:    map[string]string{util.DeleteWithClusterLabelName: "true"},
		},
		Type: corev1.SecretTypeOpaque,
	}
}

var _ core.Platform = (*CreateOptions)(nil)

func NewCreateCommand(opts *core.CreateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "openstack",
		Short:        "Creates basic functional HostedCluster resources on OpenStack platform",
		SilenceUsage: true,
	}

	openstackOpts := DefaultOptions()
	BindOptions(openstackOpts, cmd.Flags())
	cmd.MarkPersistentFlagRequired("pull-secret")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if opts.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
			defer cancel()
		}

		if err := CreateCluster(ctx, opts, openstackOpts); err != nil {
			opts.Log.Error(err, "Failed to create cluster")
			return err
		}
		return nil
	}

	return cmd
}

func CreateCluster(ctx context.Context, opts *core.CreateOptions, kubevirtOpts *CreateOptions) error {
	return core.CreateCluster(ctx, opts, kubevirtOpts)
}
