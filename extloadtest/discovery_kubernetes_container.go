package extloadtest

import (
  "context"
  "fmt"
  "github.com/procyon-projects/chrono"
  "github.com/rs/zerolog/log"
  "github.com/steadybit/discovery-kit/go/discovery_kit_api"
  "github.com/steadybit/extension-kit/exthttp"
  "github.com/steadybit/extension-kit/extutil"
  "github.com/steadybit/extension-loadtest/config"
  "time"
)

func RegisterDiscoveryKubernetesContainer() {
	exthttp.RegisterHttpHandler("/discovery/kubernetes-container", exthttp.GetterAsHandler(getDiscoveryKubernetesContainer))
	exthttp.RegisterHttpHandler("/discovery/kubernetes-container/targets", exthttp.GetterAsHandler(getDiscoveryKubernetesContainerTargets))
}

func getDiscoveryKubernetesContainer() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id:         "com.steadybit.extension_container.container",
		RestrictTo: extutil.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			Method:       "GET",
			Path:         "/discovery/kubernetes-container/targets",
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

var kubernetesContainer *[]discovery_kit_api.EnrichmentData

func getDiscoveryKubernetesContainerTargets() discovery_kit_api.DiscoveryData {
	if kubernetesContainer == nil {
		kubernetesContainer = initKubernetesContainerTargets()
		taskScheduler := chrono.NewDefaultTaskScheduler()
		_, err := taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
			changeLabelsOfRandomContainer(kubernetesContainer)
		}, time.Duration(config.Config.ChangeTimeKubernetesContainer)*time.Second)
		if err != nil {
			log.Error().Err(err).Msg("Could not schedule task")
		}
	}
	return discovery_kit_api.DiscoveryData{
		EnrichmentData: kubernetesContainer,
	}
}

func changeLabelsOfRandomContainer(enrichmentData *[]discovery_kit_api.EnrichmentData) {
	enrichmentDataCopy := *enrichmentData
	count2Change := int(float64(len(enrichmentDataCopy)) / float64(100) * float64(config.Config.ChangeRateKubernetesContainer))
	log.Debug().Msgf("Changing labels of %d containers", count2Change)
	indicesToChange := getRandomIndices(count2Change, len(enrichmentDataCopy))
  for _, index := range indicesToChange {
    container := enrichmentDataCopy[index]
    container.Attributes["k8s.label.tags.changeTimestamp"] = []string{time.Now().String()}
  }
  enrichmentData = &enrichmentDataCopy
}

func initKubernetesContainerTargets() *[]discovery_kit_api.EnrichmentData {
	//example: 2 containers per pod * 4 pods per deployment * 5 deployments per node * 400 nodes = 16000 containers
	count := config.Config.NodeCount * config.Config.DeploymentsPerNode * config.Config.PodsPerDeployment * config.Config.ContainerPerPod
	result := make([]discovery_kit_api.EnrichmentData, 0, count)

	hostTargets := *getDiscoveryHostTargets().Targets
	deploymentTargets := *getDiscoveryKubernetesDeploymentTargets().Targets
	for hostIndex := 0; hostIndex < len(hostTargets); hostIndex++ {
		host := hostTargets[hostIndex]
		for deploymentIndex := 0; deploymentIndex < config.Config.DeploymentsPerNode; deploymentIndex++ {
			deployment := deploymentTargets[(hostIndex*config.Config.DeploymentsPerNode)+deploymentIndex]
			pods := deployment.Attributes["k8s.pod.name"]
			containerIdsStripped := deployment.Attributes["k8s.container.id.stripped"]
			for podIndex := 0; podIndex < len(pods); podIndex++ {
				podName := pods[podIndex]
				for containerIndex := 0; containerIndex < config.Config.ContainerPerPod; containerIndex++ {
					containerIdStripped := containerIdsStripped[(podIndex*config.Config.ContainerPerPod)+containerIndex]
					containerId := fmt.Sprintf("containerd://%s", containerIdStripped)
					target := discovery_kit_api.EnrichmentData{
						Id:                 containerId,
						EnrichmentDataType: "com.steadybit.extension_kubernetes.kubernetes-container",
						Attributes: map[string][]string{
							"k8s.cluster-name":                         {config.Config.ClusterName},
							"k8s.container.id":                         {containerId},
							"k8s.container.id.stripped":                {containerIdStripped},
							"k8s.container.image":                      {"docker.io/steadybit/loadtest-example:latest"},
							"k8s.container.name":                       {"loadtest-example"},
							"k8s.container.ready":                      {"true"},
							"k8s.deployment":                           deployment.Attributes["k8s.deployment"],
							"k8s.distribution":                         {"kubernetes"},
							"k8s.label.domain":                         {"loadtest"},
							"k8s.label.run":                            {"loadtest-example"},
							"k8s.label.service-tier":                   {"2"},
							"k8s.label.tags.datadoghq.com/service":     {"loadtest"},
							"k8s.label.tags.datadoghq.com/version":     {"1.0.0"},
							"k8s.label.tags.changeTimestamp":           {time.Now().String()},
							"k8s.namespace":                            deployment.Attributes["k8s.namespace"],
							"k8s.node.name":                            {host.Id},
							"k8s.pod.label.domain":                     {"loadtest"},
							"k8s.pod.label.run":                        {"loadtest-example"},
							"k8s.pod.label.service-tier":               {"2"},
							"k8s.pod.label.tags.datadoghq.com/service": {"loadtest"},
							"k8s.pod.label.tags.datadoghq.com/version": {"1.0.0"},
							"k8s.pod.name":                             {podName},
							"k8s.replicaset":                           deployment.Attributes["k8s.deployment"],
							"k8s.service.name":                         deployment.Attributes["k8s.deployment"],
						},
					}
					result = append(result, target)
				}
			}
		}
	}

	return &result
}
