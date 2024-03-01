package extloadtest

import (
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryKubernetesCluster() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id: "com.steadybit.extension_kubernetes.kubernetes-cluster",
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("60m"),
		},
	}
}

func createKubernetesClusterTargets() []discovery_kit_api.Target {
	if IAmTheLeader {
		return []discovery_kit_api.Target{
			{
				Id:         config.Config.ClusterName,
				TargetType: "com.steadybit.extension_kubernetes.kubernetes-cluster",
				Label:      config.Config.ClusterName,
				Attributes: map[string][]string{
					"k8s.cluster-name": {config.Config.ClusterName},
				},
			},
		}
	}
	return []discovery_kit_api.Target{}
}
