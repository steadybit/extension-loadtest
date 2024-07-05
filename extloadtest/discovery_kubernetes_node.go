package extloadtest

import (
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_commons"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
)

func getDiscoveryKubernetesNode() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id: "com.steadybit.extension_kubernetes.kubernetes-node",
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

type NodeData struct {
	Deployments          map[string]bool
	Pods                 map[string]bool
	Namespaces           map[string]bool
	ContainerIds         map[string]bool
	ContainerIdsStripped map[string]bool
}

func createKubernetesNodeTargets(containerTargets []discovery_kit_api.EnrichmentData) []discovery_kit_api.Target {
	var data = map[string]NodeData{}
	for _, containerTarget := range containerTargets {
		hostName := containerTarget.Attributes["k8s.node.name"][0]

		nodeData, ok := data[hostName]
		if !ok {
			nodeData = NodeData{}
			nodeData.Deployments = map[string]bool{}
			nodeData.Pods = map[string]bool{}
			nodeData.Namespaces = map[string]bool{}
			nodeData.ContainerIds = map[string]bool{}
			nodeData.ContainerIdsStripped = map[string]bool{}
			data[hostName] = nodeData
		}
		nodeData.Deployments[containerTarget.Attributes["k8s.deployment"][0]] = true
		nodeData.Pods[containerTarget.Attributes["k8s.pod.name"][0]] = true
		nodeData.Namespaces[containerTarget.Attributes["k8s.namespace"][0]] = true
		nodeData.ContainerIds[containerTarget.Attributes["k8s.container.id"][0]] = true
		nodeData.ContainerIdsStripped[containerTarget.Attributes["k8s.container.id.stripped"][0]] = true
	}

	var result = make([]discovery_kit_api.Target, 0, len(data))
	for hostName, nodeData := range data {
		//fmt.Println(hostName, nodeData)

		target := discovery_kit_api.Target{
			Id:         hostName,
			TargetType: "com.steadybit.extension_kubernetes.kubernetes-node",
			Label:      hostName,
			Attributes: map[string][]string{
				"k8s.node.name":             {hostName},
				"host.hostname":             {hostName},
				"k8s.cluster-name":          {config.Config.ClusterName},
				"k8s.distribution":          {"kubernetes"},
				"k8s.deployment":            keys(nodeData.Deployments),
				"k8s.namespace":             keys(nodeData.Namespaces),
				"k8s.pod.name":              keys(nodeData.Pods),
				"k8s.container.id":          keys(nodeData.ContainerIds),
				"k8s.container.id.stripped": keys(nodeData.ContainerIdsStripped),
			},
		}
		result = append(result, target)

	}

	return discovery_kit_commons.ApplyAttributeExcludes(result, config.Config.DiscoveryAttributesExcludesKubernetesNode)
}

func keys(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
