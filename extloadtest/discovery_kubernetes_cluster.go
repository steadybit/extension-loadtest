package extloadtest

import (
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryKubernetesCluster() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id:         "com.steadybit.extension_kubernetes.kubernetes-cluster",
		RestrictTo: extutil.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			Method:       "GET",
			Path:         "/discovery/kubernetes-cluster/targets",
			CallInterval: extutil.Ptr("60m"),
		},
	}
}

func initKubernetesClusterTargets() []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, 1)
	target := discovery_kit_api.Target{
		Id:         config.Config.ClusterName,
		TargetType: "com.steadybit.extension_kubernetes.kubernetes-cluster",
		Label:      config.Config.ClusterName,
		Attributes: map[string][]string{
			"k8s.cluster-name": {config.Config.ClusterName},
		},
	}
	result = append(result, target)
	return result
}
