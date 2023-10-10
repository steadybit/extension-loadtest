package extloadtest

import (
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-kit/exthttp"
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

func (t *TargetData) RegisterDiscoveryHandlers() {
	exthttp.RegisterHttpHandler("/discovery/host", exthttp.GetterAsHandler(getDiscoveryHost))
	exthttp.RegisterHttpHandler("/discovery/host/targets", exthttp.GetterAsHandler(targets(t.hosts)))

	exthttp.RegisterHttpHandler("/discovery/ec2-instance", exthttp.GetterAsHandler(getDiscoveryEc2Instance))
	exthttp.RegisterHttpHandler("/discovery/ec2-instance/targets", exthttp.GetterAsHandler(targets(t.ec2Instances)))

	exthttp.RegisterHttpHandler("/discovery/kubernetes-cluster", exthttp.GetterAsHandler(getDiscoveryKubernetesCluster))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-cluster/targets", exthttp.GetterAsHandler(targets(t.kubernetesClusters)))

	exthttp.RegisterHttpHandler("/discovery/kubernetes-deployment", exthttp.GetterAsHandler(getDiscoveryKubernetesDeployment))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-deployment/targets", exthttp.GetterAsHandler(targets(t.kubernetesDeployments)))

	exthttp.RegisterHttpHandler("/discovery/kubernetes-pod", exthttp.GetterAsHandler(getDiscoveryKubernetesPods))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-pod/targets", exthttp.GetterAsHandler(targets(t.kubernetesPods)))

	exthttp.RegisterHttpHandler("/discovery/kubernetes-container", exthttp.GetterAsHandler(getDiscoveryKubernetesContainer))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-container/targets", exthttp.GetterAsHandler(enrichmentData(t.kubernetesContainers)))

	exthttp.RegisterHttpHandler("/discovery/kubernetes-node", exthttp.GetterAsHandler(getDiscoveryKubernetesNode))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-node/targets", exthttp.GetterAsHandler(targets(t.kubernetesNodes)))

	exthttp.RegisterHttpHandler("/discovery/container", exthttp.GetterAsHandler(getDiscoveryContainer))
	exthttp.RegisterHttpHandler("/discovery/container/targets", exthttp.GetterAsHandler(targets(t.containers)))
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
