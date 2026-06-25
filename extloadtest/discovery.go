// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

package extloadtest

import (
	"context"
	"encoding/json"
	"fmt"
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
	hosts                 []discovery_kit_api.Target
	ec2Instances          []discovery_kit_api.Target
	azureInstances        []discovery_kit_api.Target
	gcpInstances          []discovery_kit_api.Target
	kubernetesClusters    []discovery_kit_api.Target
	kubernetesDeployments []discovery_kit_api.Target
	kubernetesPods        []discovery_kit_api.Target
	kubernetesContainers  []discovery_kit_api.EnrichmentData
	kubernetesNodes       []discovery_kit_api.Target
	containers            []discovery_kit_api.Target
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

	now := time.Now()
	typeId := l.description().Id

	// Simulated extension restart: the whole type is unavailable during a restart window.
	if isExtensionDown(config.Config.FindSimulateExtensionRestartSpecification(typeId), now) {
		return []discovery_kit_api.Target{}, nil
	}

	attrSpec := config.Config.FindAttributeUpdate(typeId)
	replSpec := config.Config.FindTargetReplacementsSpecification(typeId)
	total := len(*l.targets)

	host := ""
	if config.Config.ServicesEnabled {
		var key discovery_kit_sdk.HttpRequestContextKey = "httpRequest"
		if value := ctx.Value(key); value != nil {
			if httpRequest, ok := value.(*http.Request); ok && httpRequest != nil {
				host = httpRequest.Host
			}
		}
	}

	// Zero-copy fast path: with no host projection and nothing to mutate, serve the
	// shared slice directly so simulating lots of targets keeps the previous (no
	// per-request allocation) memory profile.
	hasAttr := attrSpec != nil && attrSpec.Rate > attributeUpdateDisableThreshold
	if host == "" && !hasAttr && replSpec == nil {
		return *l.targets, nil
	}

	result := make([]discovery_kit_api.Target, 0, total)
	for i := range *l.targets {
		base := &(*l.targets)[i]
		if isTargetReplaced(base.Id, total, replSpec, now) {
			continue
		}
		target := copyTargetWithHost(base, host)
		// Use the canonical base id (not the per-service host#id) so every service
		// projection of a target shares one change schedule, matching isTargetReplaced.
		applyAttributeUpdate(target.Attributes, base.Id, attrSpec, now)
		result = append(result, target)
	}
	return result, nil
}

// copyTargetWithHost deep-copies the base target so serve-time mutations never
// touch the shared slice, applying the per-service host prefix when host != "".
func copyTargetWithHost(base *discovery_kit_api.Target, host string) discovery_kit_api.Target {
	target := discovery_kit_api.Target{
		Id:         base.Id,
		TargetType: base.TargetType,
		Label:      base.Label,
		Attributes: make(map[string][]string, len(base.Attributes)),
	}
	for k, v := range base.Attributes {
		cp := make([]string, len(v))
		copy(cp, v)
		target.Attributes[k] = cp
	}
	if host != "" {
		target.Id = fmt.Sprintf("%s#%s", host, base.Id)
		target.Label = fmt.Sprintf("%s#%s", host, base.Label)
		prefixAttribute(&target, "host.hostname", host)
		prefixAttribute(&target, "aws-ec2.hostname.internal", host)
		prefixAttribute(&target, "azure-scale-set-instance.hostname", host)
		prefixAttribute(&target, "azure-vm.hostname", host)
		prefixAttribute(&target, "gcp-vm.hostname", host)
		prefixAttribute(&target, "k8s.container.id.stripped", host)
		prefixAttribute(&target, "k8s.node.name", host)
		prefixAttribute(&target, "container.host", host)
		prefixAttribute(&target, "container.host/name", host)
		prefixAttribute(&target, "container.id", host)
		prefixAttribute(&target, "container.id.stripped", host)
		prefixAttribute(&target, "container.id.stripped", host)
	}
	return target
}

// applyAttributeUpdate sets the deterministic, replica-consistent value for a
// configured update-attribute (no-op when no spec or the rate disables it).
func applyAttributeUpdate(attributes map[string][]string, id string, spec *config.AttributeUpdateSpecification, now time.Time) {
	if spec == nil || spec.Rate <= attributeUpdateDisableThreshold {
		return
	}
	attributes[spec.AttributeName] = []string{deterministicAttributeValue(id, spec, now)}
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
	now := time.Now()
	typeId := l.description().Id

	if isExtensionDown(config.Config.FindSimulateExtensionRestartSpecification(typeId), now) {
		return []discovery_kit_api.EnrichmentData{}, nil
	}

	attrSpec := config.Config.FindAttributeUpdate(typeId)
	replSpec := config.Config.FindTargetReplacementsSpecification(typeId)
	total := len(*l.data)

	// Zero-copy fast path (see DiscoverTargets): nothing to mutate, serve the shared slice.
	hasAttr := attrSpec != nil && attrSpec.Rate > attributeUpdateDisableThreshold
	if !hasAttr && replSpec == nil {
		return *l.data, nil
	}

	result := make([]discovery_kit_api.EnrichmentData, 0, total)
	for i := range *l.data {
		base := &(*l.data)[i]
		if isTargetReplaced(base.Id, total, replSpec, now) {
			continue
		}
		ed := discovery_kit_api.EnrichmentData{
			Id:                 base.Id,
			EnrichmentDataType: base.EnrichmentDataType,
			Attributes:         make(map[string][]string, len(base.Attributes)),
		}
		for k, v := range base.Attributes {
			cp := make([]string, len(v))
			copy(cp, v)
			ed.Attributes[k] = cp
		}
		applyAttributeUpdate(ed.Attributes, ed.Id, attrSpec, now)
		result = append(result, ed)
	}
	return result, nil
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

func (t *TargetData) RegisterRecreateActions() {
	action_kit_sdk.RegisterAction(NewRecreateAction(
		"com.steadybit.extension_host.host",
		action_kit_api.TargetSelectionTemplate{
			Label:       "by host name",
			Description: new("Find host by host name."),
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
			Description: new("Find container by kubernetes deployment."),
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

func (t *TargetData) updateConfigHandler(config any) exthttp.Handler {
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
	// Updates are computed deterministically at serve time, so configuration
	// changes take effect on the next discovery without any rescheduling.
	log.Info().Msg("Configuration updated; deterministic serve-time updates take effect on the next discovery")
}
