package extloadtest

import (
	"context"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_sdk"
	"github.com/steadybit/extension-kit/extutil"
)

type TargetData struct {
	hosts                 []discovery_kit_api.Target
	ec2Instances          []discovery_kit_api.Target
	kubernetesClusters    []discovery_kit_api.Target
	kubernetesDeployments []discovery_kit_api.Target
	kubernetesPods        []discovery_kit_api.Target
	kubernetesContainers  []discovery_kit_api.EnrichmentData
	kubernetesNodes       []discovery_kit_api.Target
	containers            []discovery_kit_api.Target
}

func NewTargetData() *TargetData {
	hosts := createHostTargets()
	ec2Instances := createEc2InstanceTargets(hosts)
	kubernetesClusters := createKubernetesClusterTargets()
	kubernetesDeployments := createKubernetesDeploymentTargets()
	kubernetesPods := createKubernetesPodTargets(hosts, kubernetesDeployments)
	kubernetesContainers := createKubernetesContainerTargets(kubernetesPods)
	kubernetesNodes := createKubernetesNodeTargets(kubernetesContainers)
	containers := createContainerTargets(kubernetesContainers)

	return &TargetData{
		hosts:                 hosts,
		ec2Instances:          ec2Instances,
		kubernetesClusters:    kubernetesClusters,
		kubernetesDeployments: kubernetesDeployments,
		kubernetesPods:        kubernetesPods,
		kubernetesContainers:  kubernetesContainers,
		kubernetesNodes:       kubernetesNodes,
		containers:            containers,
	}
}

type ltTargetDiscovery struct {
	targets     []discovery_kit_api.Target
	description func() discovery_kit_api.DiscoveryDescription
}

func (l ltTargetDiscovery) Describe() discovery_kit_api.DiscoveryDescription {
	return l.description()
}

func (l ltTargetDiscovery) DiscoverTargets(_ context.Context) ([]discovery_kit_api.Target, error) {
	return l.targets, nil
}

type ltEdDiscovery struct {
	data        []discovery_kit_api.EnrichmentData
	description func() discovery_kit_api.DiscoveryDescription
}

func (l ltEdDiscovery) Describe() discovery_kit_api.DiscoveryDescription {
	return l.description()

}

func (l ltEdDiscovery) DiscoverEnrichmentData(_ context.Context) ([]discovery_kit_api.EnrichmentData, error) {
	return l.data, nil
}

func (t *TargetData) RegisterDiscoveries() {
	discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryHost, targets: t.hosts})
	discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryEc2Instance, targets: t.ec2Instances})
	discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryKubernetesCluster, targets: t.kubernetesClusters})
	discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryKubernetesDeployment, targets: t.kubernetesDeployments})
	discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryKubernetesPods, targets: t.kubernetesPods})
	discovery_kit_sdk.Register(&ltEdDiscovery{description: getDiscoveryKubernetesContainer, data: t.kubernetesContainers})
	discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryKubernetesNode, targets: t.kubernetesNodes})
	discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryContainer, targets: t.containers})
}

func (t *TargetData) ScheduleUpdates() {
	scheduleTargetAttributeUpdateIfNecessary(t.hosts, "com.steadybit.extension_host.host")
	scheduleTargetAttributeUpdateIfNecessary(t.ec2Instances, "com.steadybit.extension_aws.ec2-instance")
	scheduleTargetAttributeUpdateIfNecessary(t.kubernetesClusters, "com.steadybit.extension_kubernetes.kubernetes-cluster")
	scheduleTargetAttributeUpdateIfNecessary(t.kubernetesDeployments, "com.steadybit.extension_kubernetes.kubernetes-deployment")
	scheduleTargetAttributeUpdateIfNecessary(t.kubernetesPods, "com.steadybit.extension_kubernetes.kubernetes-pod")
	scheduleTargetAttributeUpdateIfNecessary(t.containers, "com.steadybit.extension_container.container")
	scheduleEnrichmentDataAttributeUpdateIfNecessary(t.kubernetesContainers, "com.steadybit.extension_kubernetes.kubernetes-container")
}

func (t *TargetData) RegisterRecreateActions() {
	action_kit_sdk.RegisterAction(NewRecreateAction(
		"com.steadybit.extension_host.host",
		action_kit_api.TargetSelectionTemplate{
			Label:       "by host name",
			Description: extutil.Ptr("Find host by host name."),
			Query:       "host.hostname=\"\"",
		},
		func(name string) {
			updateTargetId(t.hosts, name, "com.steadybit.extension_host.host")
			t.ec2Instances = createEc2InstanceTargets(t.hosts)
			t.kubernetesPods = createKubernetesPodTargets(t.hosts, t.kubernetesDeployments)
			t.kubernetesContainers = createKubernetesContainerTargets(t.kubernetesPods)
			t.kubernetesNodes = createKubernetesNodeTargets(t.kubernetesContainers)
			t.containers = createContainerTargets(t.kubernetesContainers)
		},
	))

	action_kit_sdk.RegisterAction(NewRecreateAction(
		"com.steadybit.extension_container.container",
		action_kit_api.TargetSelectionTemplate{
			Label:       "by kubernetes deployment",
			Description: extutil.Ptr("Find container by kubernetes deployment."),
			Query:       "k8s.cluster-name=\"\" and k8s.namespace=\"\" and k8s.deployment=\"\"",
		},
		func(name string) {
			updateEnrichmentDataId(t.kubernetesContainers, name, "com.steadybit.extension_kubernetes.kubernetes-container")
			t.containers = createContainerTargets(t.kubernetesContainers)
		},
	))
}
