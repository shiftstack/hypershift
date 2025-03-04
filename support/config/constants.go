package config

const (
	// NeedManagementKASAccessLabel is used by network policies
	// to prevent any pod which doesn't contain the label from accessing the management cluster KAS.
	NeedManagementKASAccessLabel = "hypershift.openshift.io/need-management-kas-access"

	// NeedMetricsServerAccessLabel is used by network policies
	// to allow egress communication to the metrics server on the management cluster.
	NeedMetricsServerAccessLabel = "hypershift.openshift.io/need-metrics-server-access"

	// EtcdPriorityClass is for etcd pods.
	EtcdPriorityClass = "hypershift-etcd"

	// APICriticalPriorityClass is for pods that are required for API calls and
	// resource admission to succeed. This includes pods like kube-apiserver,
	// aggregated API servers, and webhooks.
	APICriticalPriorityClass = "hypershift-api-critical"

	// DefaultPriorityClass is for pods in the Hypershift control plane that are
	// not API critical but still need elevated priority.
	DefaultPriorityClass = "hypershift-control-plane"

	DefaultServiceAccountIssuer  = "https://kubernetes.default.svc"
	DefaultImageRegistryHostname = "image-registry.openshift-image-registry.svc:5000"
	DefaultAdvertiseIPv4Address  = "172.20.0.1"
	DefaultAdvertiseIPv6Address  = "fd00::1"
	DefaultEtcdURL               = "https://etcd-client:2379"
	// KASSVCLBAzurePort is needed because for Azure we currently hardcode 7443 for the SVC LB as 6443 collides with public LB rule for the management cluster.
	// https://bugzilla.redhat.com/show_bug.cgi?id=2060650
	// TODO(alberto): explore exposing multiple Azure frontend IPs on the load balancer.
	KASSVCLBAzurePort           = 7443
	KASSVCPort                  = 6443
	KASPodDefaultPort           = 6443
	KASSVCIBMCloudPort          = 2040
	DefaultServiceNodePortRange = "30000-32767"
	DefaultSecurityContextUser  = 1001
	RecommendedLeaseDuration    = "137s"
	RecommendedRenewDeadline    = "107s"
	RecommendedRetryPeriod      = "26s"
	KCMRecommendedRenewDeadline = "12s"
	KCMRecommendedRetryPeriod   = "3s"

	DefaultIngressDomainEnvVar                    = "DEFAULT_INGRESS_DOMAIN"
	EnableCVOManagementClusterMetricsAccessEnvVar = "ENABLE_CVO_MANAGEMENT_CLUSTER_METRICS_ACCESS"

	EnableEtcdRecoveryEnvVar = "ENABLE_ETCD_RECOVERY"

	CPOOverridesEnvVar = "ENABLE_CPO_OVERRIDES"

	AuditWebhookService = "audit-webhook"

	// AROHCPKeyVaultManagedIdentityClientID captures the client ID of the managed identity created on an ARO HCP
	// management cluster. This managed identity is used to pull secrets and certificates out of Azure Key Vaults in the
	// management cluster's resource group in Azure.
	AROHCPKeyVaultManagedIdentityClientID = "ARO_HCP_KEY_VAULT_USER_CLIENT_ID"
)
