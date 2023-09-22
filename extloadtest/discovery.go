package extloadtest

import "github.com/steadybit/extension-kit/exthttp"

func RegisterAllDiscoveryHandlers() {
	hosts := initHostTargets()
	scheduleTargetAttributeUpdateIfNecessary(hosts, "com.steadybit.extension_host.host")
	exthttp.RegisterHttpHandler("/discovery/host", exthttp.GetterAsHandler(getDiscoveryHost))
	exthttp.RegisterHttpHandler("/discovery/host/targets", exthttp.GetterAsHandler(targets(hosts)))

	ec2Instances := initEc2InstanceTargets()
	scheduleTargetAttributeUpdateIfNecessary(ec2Instances, "com.steadybit.extension_aws.ec2-instance")
	exthttp.RegisterHttpHandler("/discovery/ec2-instance", exthttp.GetterAsHandler(getDiscoveryEc2Instance))
	exthttp.RegisterHttpHandler("/discovery/ec2-instance/targets", exthttp.GetterAsHandler(targets(ec2Instances)))

	kubernetesClusters := initKubernetesClusterTargets()
	scheduleTargetAttributeUpdateIfNecessary(kubernetesClusters, "com.steadybit.extension_kubernetes.kubernetes-cluster")
	exthttp.RegisterHttpHandler("/discovery/kubernetes-cluster", exthttp.GetterAsHandler(getDiscoveryKubernetesCluster))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-cluster/targets", exthttp.GetterAsHandler(targets(kubernetesClusters)))

	kubernetesDeployments := initKubernetesDeploymentTargets()
	scheduleTargetAttributeUpdateIfNecessary(kubernetesDeployments, "com.steadybit.extension_kubernetes.kubernetes-deployment")
	exthttp.RegisterHttpHandler("/discovery/kubernetes-deployment", exthttp.GetterAsHandler(getDiscoveryKubernetesDeployment))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-deployment/targets", exthttp.GetterAsHandler(targets(kubernetesDeployments)))

	kubernetesContainers := initKubernetesContainerTargets(hosts, kubernetesDeployments)
	scheduleEnrichmentDataAttributeUpdateIfNecessary(kubernetesContainers, "com.steadybit.extension_kubernetes.kubernetes-container")
	exthttp.RegisterHttpHandler("/discovery/kubernetes-container", exthttp.GetterAsHandler(getDiscoveryKubernetesContainer))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-container/targets", exthttp.GetterAsHandler(enrichmentData(kubernetesContainers)))

	containers := initContainerTargets(kubernetesContainers)
	scheduleTargetAttributeUpdateIfNecessary(containers, "com.steadybit.extension_container.container")
	exthttp.RegisterHttpHandler("/discovery/container", exthttp.GetterAsHandler(getDiscoveryContainer))
	exthttp.RegisterHttpHandler("/discovery/container/targets", exthttp.GetterAsHandler(targets(containers)))
}
