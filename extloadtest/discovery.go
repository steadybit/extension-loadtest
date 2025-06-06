// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

package extloadtest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/procyon-projects/chrono"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_sdk"
	"github.com/steadybit/extension-kit/exthttp"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/steadybit/extension-loadtest/config"
	"net/http"
	"time"
)

type TargetData struct {
	hosts                       []discovery_kit_api.Target
	hostsBackup                 []discovery_kit_api.Target
	ec2Instances                []discovery_kit_api.Target
	ec2InstancesBackup          []discovery_kit_api.Target
	azureInstances              []discovery_kit_api.Target
	azureInstancesBackup        []discovery_kit_api.Target
	gcpInstances                []discovery_kit_api.Target
	gcpInstancesBackup          []discovery_kit_api.Target
	kubernetesClusters          []discovery_kit_api.Target
	kubernetesClustersBackup    []discovery_kit_api.Target
	kubernetesDeployments       []discovery_kit_api.Target
	kubernetesDeploymentsBackup []discovery_kit_api.Target
	kubernetesPods              []discovery_kit_api.Target
	kubernetesPodsBackup        []discovery_kit_api.Target
	kubernetesContainers        []discovery_kit_api.EnrichmentData
	kubernetesContainersBackup  []discovery_kit_api.EnrichmentData
	kubernetesNodes             []discovery_kit_api.Target
	kubernetesNodesBackup       []discovery_kit_api.Target
	containers                  []discovery_kit_api.Target
	containersBackup            []discovery_kit_api.Target
}

func NewTargetData() *TargetData {
	ec2Hosts := createHostTargets(config.Config.Ec2NodeCount, "ec2")
	gcpNodeCount := 0
	if config.IsPodZero() {
		gcpNodeCount = config.Config.GcpNodeCount

	}
	gcpHosts := createHostTargets(gcpNodeCount, "gcp")
	azureNodeCount := 0
	if config.IsPodZero() {
		azureNodeCount = config.Config.AzureNodeCount
	}
	azureHosts := createHostTargets(azureNodeCount, "azure")

	ec2Instances := createEc2InstanceTargets(ec2Hosts)
	gcpInstances := createGcpInstanceTargets(gcpHosts)
	azureInstances := createAzureInstanceTargets(azureHosts)

	kubernetesClusters := createKubernetesClusterTargets()

	ec2KubernetesDeployments := createKubernetesDeploymentTargets(config.Config.Ec2NodeCount, "ec2")
	ec2KubernetesPods := createKubernetesPodTargets(ec2Hosts, ec2KubernetesDeployments)
	ec2KubernetesContainers := createKubernetesContainerTargets(ec2KubernetesPods)
	ec2KubernetesNodes := createKubernetesNodeTargets(ec2KubernetesContainers)
	ec2Containers := createContainerTargets(ec2KubernetesContainers)

	var gcpKubernetesDeployments []discovery_kit_api.Target
	var gcpKubernetesPods []discovery_kit_api.Target
	var gcpKubernetesContainers []discovery_kit_api.EnrichmentData
	var gcpKubernetesNodes []discovery_kit_api.Target
	var gcpContainers []discovery_kit_api.Target
	if config.IsPodZero() {
		gcpKubernetesDeployments = createKubernetesDeploymentTargets(config.Config.GcpNodeCount, "gcp")
		gcpKubernetesPods = createKubernetesPodTargets(gcpHosts, gcpKubernetesDeployments)
		gcpKubernetesContainers = createKubernetesContainerTargets(gcpKubernetesPods)
		gcpKubernetesNodes = createKubernetesNodeTargets(gcpKubernetesContainers)
		gcpContainers = createContainerTargets(gcpKubernetesContainers)
	}

	var azureKubernetesDeployments []discovery_kit_api.Target
	var azureKubernetesPods []discovery_kit_api.Target
	var azureKubernetesContainers []discovery_kit_api.EnrichmentData
	var azureKubernetesNodes []discovery_kit_api.Target
	var azureContainers []discovery_kit_api.Target
	if config.IsPodZero() {
		azureKubernetesDeployments = createKubernetesDeploymentTargets(config.Config.AzureNodeCount, "azure")
		azureKubernetesPods = createKubernetesPodTargets(azureHosts, azureKubernetesDeployments)
		azureKubernetesContainers = createKubernetesContainerTargets(azureKubernetesPods)
		azureKubernetesNodes = createKubernetesNodeTargets(azureKubernetesContainers)
		azureContainers = createContainerTargets(azureKubernetesContainers)
	}
	// count generated targets
	log.Info().Msgf("Generated %d hosts", len(ec2Hosts)+len(gcpHosts)+len(azureHosts))
	log.Info().Msgf("Generated %d ec2 instances", len(ec2Instances))
	log.Info().Msgf("Generated %d gcp instances", len(gcpInstances))
	log.Info().Msgf("Generated %d azure instances", len(azureInstances))
	log.Info().Msgf("Generated %d kubernetes clusters", len(kubernetesClusters))
	log.Info().Msgf("Generated %d kubernetes deployments", len(ec2KubernetesDeployments)+len(gcpKubernetesDeployments)+len(azureKubernetesDeployments))
	log.Info().Msgf("Generated %d kubernetes pods", len(ec2KubernetesPods)+len(gcpKubernetesPods)+len(azureKubernetesPods))
	log.Info().Msgf("Generated %d kubernetes containers", len(ec2KubernetesContainers)+len(gcpKubernetesContainers)+len(azureKubernetesContainers))
	log.Info().Msgf("Generated %d kubernetes nodes", len(ec2KubernetesNodes)+len(gcpKubernetesNodes)+len(azureKubernetesNodes))
	log.Info().Msgf("Generated %d containers", len(ec2Containers)+len(gcpContainers)+len(azureContainers))

	var targetsAvailable = 0
	if !config.Config.DisableHostDiscovery {
		targetsAvailable += len(ec2Hosts) + len(gcpHosts) + len(azureHosts)
	}
	if !config.Config.DisableAWSDiscovery {
		targetsAvailable += len(ec2Instances)
	}
	if !config.Config.DisableGCPDiscovery {
		targetsAvailable += len(gcpInstances)
	}
	if !config.Config.DisableAzureDiscovery {
		targetsAvailable += len(azureInstances)
	}
	if !config.Config.DisableKubernetesDiscovery {
		targetsAvailable += len(kubernetesClusters) + len(ec2KubernetesDeployments) + len(gcpKubernetesDeployments) + len(azureKubernetesDeployments) + len(ec2KubernetesPods) + len(gcpKubernetesPods) + len(azureKubernetesPods) + len(ec2KubernetesContainers) + len(gcpKubernetesContainers) + len(azureKubernetesContainers) + len(ec2KubernetesNodes) + len(gcpKubernetesNodes) + len(azureKubernetesNodes)
	}
	if !config.Config.DisableContainerDiscovery {
		targetsAvailable += len(ec2Containers) + len(gcpContainers) + len(azureContainers)
	}
	log.Info().Msgf("Total targets available: %d", targetsAvailable)

	return &TargetData{
		hosts:                 append(append(ec2Hosts, gcpHosts...), azureHosts...),
		ec2Instances:          ec2Instances,
		gcpInstances:          gcpInstances,
		azureInstances:        azureInstances,
		kubernetesClusters:    kubernetesClusters,
		kubernetesDeployments: append(append(ec2KubernetesDeployments, gcpKubernetesDeployments...), azureKubernetesDeployments...),
		kubernetesPods:        append(append(ec2KubernetesPods, gcpKubernetesPods...), azureKubernetesPods...),
		kubernetesContainers:  append(append(ec2KubernetesContainers, gcpKubernetesContainers...), azureKubernetesContainers...),
		kubernetesNodes:       append(append(ec2KubernetesNodes, gcpKubernetesNodes...), azureKubernetesNodes...),
		containers:            append(append(ec2Containers, gcpContainers...), azureContainers...),
	}
}

type ltTargetDiscovery struct {
	targets     *[]discovery_kit_api.Target
	description func() discovery_kit_api.DiscoveryDescription
}

func (l ltTargetDiscovery) Describe() discovery_kit_api.DiscoveryDescription {
	return l.description()
}

func (l ltTargetDiscovery) DiscoverTargets(ctx context.Context) ([]discovery_kit_api.Target, error) {
	if config.Config.DiscoveryDelayInMs > 0 {
		time.Sleep(time.Duration(config.Config.DiscoveryDelayInMs) * time.Millisecond)
	}
	if config.Config.ServicesEnabled {
		var key discovery_kit_sdk.HttpRequestContextKey = "httpRequest"
		value := ctx.Value(key)
		if value != nil {
			httpRequest := value.(*http.Request)
			if httpRequest != nil {
				newTargets := make([]discovery_kit_api.Target, len(*l.targets))
				//copy(newTargets, *l.targets)
				for i := range *l.targets {
					newTargets[i] = discovery_kit_api.Target{
						Id:         fmt.Sprintf("%s#%s", httpRequest.Host, (*l.targets)[i].Id),
						TargetType: (*l.targets)[i].TargetType,
						Label:      fmt.Sprintf("%s#%s", httpRequest.Host, (*l.targets)[i].Label),
						Attributes: make(map[string][]string),
					}
					for k, v := range (*l.targets)[i].Attributes {
						newTargets[i].Attributes[k] = make([]string, len(v))
						copy(newTargets[i].Attributes[k], v)
					}
					prefixAttribute(&newTargets[i], "host.hostname", httpRequest.Host)
					prefixAttribute(&newTargets[i], "aws-ec2.hostname.internal", httpRequest.Host)
					prefixAttribute(&newTargets[i], "azure-scale-set-instance.hostname", httpRequest.Host)
					prefixAttribute(&newTargets[i], "azure-vm.hostname", httpRequest.Host)
					prefixAttribute(&newTargets[i], "gcp-vm.hostname", httpRequest.Host)
					prefixAttribute(&newTargets[i], "k8s.container.id.stripped", httpRequest.Host)
					prefixAttribute(&newTargets[i], "k8s.node.name", httpRequest.Host)
					prefixAttribute(&newTargets[i], "container.host", httpRequest.Host)
					prefixAttribute(&newTargets[i], "container.host/name", httpRequest.Host)
					prefixAttribute(&newTargets[i], "container.id", httpRequest.Host)
					prefixAttribute(&newTargets[i], "container.id.stripped", httpRequest.Host)
					prefixAttribute(&newTargets[i], "container.id.stripped", httpRequest.Host)
				}
				return newTargets, nil
			}
		}
	}
	return *l.targets, nil
}

func prefixAttribute(target *discovery_kit_api.Target, attributeName string, prefix string) {
	if target.Attributes != nil && target.Attributes[attributeName] != nil && len(target.Attributes[attributeName]) > 0 {
		target.Attributes[attributeName] = []string{fmt.Sprintf("%s#%s", prefix, target.Attributes[attributeName][0])}
	}
}

type ltEdDiscovery struct {
	data        *[]discovery_kit_api.EnrichmentData
	description func() discovery_kit_api.DiscoveryDescription
}

func (l ltEdDiscovery) Describe() discovery_kit_api.DiscoveryDescription {
	return l.description()

}

func (l ltEdDiscovery) DiscoverEnrichmentData(_ context.Context) ([]discovery_kit_api.EnrichmentData, error) {
	return *l.data, nil
}

func (t *TargetData) RegisterDiscoveries() {
	if !config.Config.DisableHostDiscovery {
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryHost, targets: &t.hosts})
	}
	if !config.Config.DisableAWSDiscovery {
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryEc2Instance, targets: &t.ec2Instances})
	}
	if !config.Config.DisableGCPDiscovery {
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryGcpInstance, targets: &t.gcpInstances})
	}
	if !config.Config.DisableAzureDiscovery {
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryAzureInstance, targets: &t.azureInstances})
	}
	if !config.Config.DisableKubernetesDiscovery {
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryKubernetesCluster, targets: &t.kubernetesClusters})
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryKubernetesDeployment, targets: &t.kubernetesDeployments})
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryKubernetesPods, targets: &t.kubernetesPods})
		discovery_kit_sdk.Register(&ltEdDiscovery{description: getDiscoveryKubernetesContainer, data: &t.kubernetesContainers})
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryKubernetesNode, targets: &t.kubernetesNodes})
	}
	if !config.Config.DisableContainerDiscovery {
		discovery_kit_sdk.Register(&ltTargetDiscovery{description: getDiscoveryContainer, targets: &t.containers})
	}
}

func (t *TargetData) ScheduleUpdates() {
	if !config.Config.DisableHostDiscovery {
		scheduleTargetAttributeUpdateIfNecessary(t.hosts, "com.steadybit.extension_host.host")
		scheduleTargetReplacementIfNecessary(&t.hosts, &t.hostsBackup, "com.steadybit.extension_host.host")
		scheduleTargetExtensionRestartIfNecessary(&t.hosts, &t.hostsBackup, "com.steadybit.extension_host.host")
	}
	if !config.Config.DisableAWSDiscovery {
		scheduleTargetAttributeUpdateIfNecessary(t.ec2Instances, "com.steadybit.extension_aws.ec2-instance")
		scheduleTargetReplacementIfNecessary(&t.ec2Instances, &t.ec2InstancesBackup, "com.steadybit.extension_aws.ec2-instance")
		scheduleTargetExtensionRestartIfNecessary(&t.ec2Instances, &t.ec2InstancesBackup, "com.steadybit.extension_aws.ec2-instance")

	}
	if !config.Config.DisableGCPDiscovery {
		scheduleTargetAttributeUpdateIfNecessary(t.gcpInstances, "com.steadybit.extension_gcp.vm")
		scheduleTargetReplacementIfNecessary(&t.gcpInstances, &t.gcpInstancesBackup, "com.steadybit.extension_gcp.vm")
		scheduleTargetExtensionRestartIfNecessary(&t.gcpInstances, &t.gcpInstancesBackup, "com.steadybit.extension_gcp.vm")

	}
	if !config.Config.DisableAzureDiscovery {
		scheduleTargetAttributeUpdateIfNecessary(t.azureInstances, "com.steadybit.extension_azure.scale_set.instance")
		scheduleTargetReplacementIfNecessary(&t.azureInstances, &t.azureInstancesBackup, "com.steadybit.extension_azure.scale_set.instance")
		scheduleTargetExtensionRestartIfNecessary(&t.azureInstances, &t.azureInstancesBackup, "com.steadybit.extension_azure.scale_set.instance")
	}

	if !config.Config.DisableKubernetesDiscovery {
		scheduleTargetAttributeUpdateIfNecessary(t.kubernetesClusters, "com.steadybit.extension_kubernetes.kubernetes-cluster")
		scheduleTargetReplacementIfNecessary(&t.kubernetesClusters, &t.kubernetesClustersBackup, "com.steadybit.extension_kubernetes.kubernetes-cluster")
		scheduleTargetExtensionRestartIfNecessary(&t.kubernetesClusters, &t.kubernetesClustersBackup, "com.steadybit.extension_kubernetes.kubernetes-cluster")

		scheduleTargetAttributeUpdateIfNecessary(t.kubernetesDeployments, "com.steadybit.extension_kubernetes.kubernetes-deployment")
		scheduleTargetReplacementIfNecessary(&t.kubernetesDeployments, &t.kubernetesDeploymentsBackup, "com.steadybit.extension_kubernetes.kubernetes-deployment")
		scheduleTargetExtensionRestartIfNecessary(&t.kubernetesDeployments, &t.kubernetesDeploymentsBackup, "com.steadybit.extension_kubernetes.kubernetes-deployment")

		scheduleTargetAttributeUpdateIfNecessary(t.kubernetesPods, "com.steadybit.extension_kubernetes.kubernetes-pod")
		scheduleTargetReplacementIfNecessary(&t.kubernetesPods, &t.kubernetesPodsBackup, "com.steadybit.extension_kubernetes.kubernetes-pod")
		scheduleTargetExtensionRestartIfNecessary(&t.kubernetesPods, &t.kubernetesPodsBackup, "com.steadybit.extension_kubernetes.kubernetes-pod")
	}
	if !config.Config.DisableContainerDiscovery {
		scheduleTargetAttributeUpdateIfNecessary(t.containers, "com.steadybit.extension_container.container")
		scheduleTargetReplacementIfNecessary(&t.containers, &t.containersBackup, "com.steadybit.extension_container.container")
		scheduleTargetExtensionRestartIfNecessary(&t.containers, &t.containersBackup, "com.steadybit.extension_container.container")
	}
	if !config.Config.DisableKubernetesDiscovery {
		scheduleEnrichmentDataAttributeUpdateIfNecessary(t.kubernetesContainers, "com.steadybit.extension_kubernetes.kubernetes-container")
		scheduleEnrichmentDataReplacementIfNecessary(&t.kubernetesContainers, &t.kubernetesContainersBackup, "com.steadybit.extension_kubernetes.kubernetes-container")
	}
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

func (t *TargetData) RegisterConfigUpdateHandlers() {
	exthttp.RegisterHttpHandler("/config/targetReplacements", t.updateConfigHandler(&config.Config.TargetReplacements))
	exthttp.RegisterHttpHandler("/config/attributeUpdates", t.updateConfigHandler(&config.Config.AttributeUpdates))
}

func (t *TargetData) updateConfigHandler(config interface{}) exthttp.Handler {
	return func(w http.ResponseWriter, r *http.Request, body []byte) {
		if r.Method == http.MethodPost {
			clone := extutil.JsonMangle(config)
			err := json.Unmarshal(body, clone)
			if err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(err.Error()))
			}
			config = clone
			t.rescheduleUpdates()
			exthttp.WriteBody(w, config)
		} else if r.Method == http.MethodGet {
			exthttp.WriteBody(w, config)
		} else {
			w.WriteHeader(405)
		}
	}
}

func (t *TargetData) rescheduleUpdates() {
	log.Info().Msg("Stopping all scheduled updates")
	scheduler.Shutdown()
	scheduler = chrono.NewDefaultTaskScheduler()
	t.ScheduleUpdates()
}
