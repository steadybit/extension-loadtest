/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package config

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"strings"
)

// Specification is the configuration specification for the extension. Configuration values can be applied
// through environment variables. Learn more through the documentation of the envconfig package.
// https://github.com/kelseyhightower/envconfig
type Specification struct {
	ClusterName string `json:"clusterName" split_words:"true" required:"false" default:"cluster-loadtest"`
	PodUID      string `json:"podUID" split_words:"true" required:"false" default:"PodUID1"`
	PodName     string `json:"podName" split_words:"true" required:"false" default:"pod-0"`

	//2 containers per pod * 4 pods per deployment * 5 deployments per node * 400 nodes = 16000 containers
	Ec2NodeCount       int `json:"ec2NodeCount" split_words:"true" required:"false" default:"2"`
	AzureNodeCount     int `json:"azureNodeCount" split_words:"true" required:"false" default:"2"`
	GcpNodeCount       int `json:"gcpNodeCount" split_words:"true" required:"false" default:"2"`
	DeploymentsPerNode int `json:"deploymentsPerNode" split_words:"true" required:"false" default:"5"`
	PodsPerDeployment  int `json:"podsPerDeployment" split_words:"true" required:"false" default:"2"`
	ContainerPerPod    int `json:"containerPerPod" split_words:"true" required:"false" default:"2"`

	AttributeUpdates AttributeUpdateSpecifications `split_words:"true" required:"false" default:"[{\"type\": \"com.steadybit.extension_aws.ec2-instance\", \"attributeName\": \"aws-ec2.label.change-ts\", \"rate\": 0.20, \"interval\": 600},{\"type\": \"com.steadybit.extension_container.container\", \"attributeName\": \"container.label.change-ts\", \"rate\": 0.20, \"interval\": 180},{\"type\": \"com.steadybit.extension_kubernetes.kubernetes-container\", \"attributeName\": \"k8s.label.change-ts\", \"rate\": 0.20, \"interval\": 180},{\"type\": \"com.steadybit.extension_kubernetes.kubernetes-deployment\", \"attributeName\": \"k8s.label.change-ts\", \"rate\": 0.20, \"interval\": 180}]"`

	//Simulate created and deleted com.steadybit.extension_container.container targets every 180 seconds with a randomized count between 0 and the given value.
	ContainerTargetCreationsAndDeletions ContainerTargetCreationsAndDeletionsSpecification `split_words:"true" required:"false" default:"{\"count\": 10, \"interval\": 60}"`

	DiscoveryAttributesExcludesContainer            []string `json:"discoveryAttributesExcludesContainer" split_words:"true" required:"false"`
	DiscoveryAttributesExcludesEc2                  []string `json:"discoveryAttributesExcludesEc2" split_words:"true" required:"false"`
	DiscoveryAttributesExcludesHost                 []string `json:"discoveryAttributesExcludesHost" split_words:"true" required:"false"`
	DiscoveryAttributesExcludesKubernetesPod        []string `json:"discoveryAttributesExcludesKubernetesPod" split_words:"true" required:"false"`
	DiscoveryAttributesExcludesKubernetesDeployment []string `json:"discoveryAttributesExcludesKubernetesDeployment" split_words:"true" required:"false"`
	DiscoveryAttributesExcludesKubernetesContainer  []string `json:"discoveryAttributesExcludesKubernetesContainer" split_words:"true" required:"false"`
	DiscoveryAttributesExcludesKubernetesNode       []string `json:"discoveryAttributesExcludesKubernetesNode" split_words:"true" required:"false"`

	// Simulate multiple extensions
	ServicesEnabled bool `json:"servicesEnabled" split_words:"true" required:"false" default:"false"`

	// Disable Discoveries
	DisableAWSDiscovery        bool `json:"disableAWSDiscovery" split_words:"true" required:"false" default:"false"`
	DisableGCPDiscovery        bool `json:"disableGCPDiscovery" split_words:"true" required:"false" default:"false"`
	DisableAzureDiscovery      bool `json:"disableAzureDiscovery" split_words:"true" required:"false" default:"false"`
	DisableKubernetesDiscovery bool `json:"disableKubernetesDiscovery" split_words:"true" required:"false" default:"false"`
	DisableHostDiscovery       bool `json:"disableHostDiscovery" split_words:"true" required:"false" default:"false"`
	DisableContainerDiscovery  bool `json:"disableContainerDiscovery" split_words:"true" required:"false" default:"false"`

	// Simulate Enrichments
	EnrichmentHostToContainerEnabled bool `json:"enrichmentHostToContainerEnabled" split_words:"true" required:"false" default:"false"`
	EnrichmentContainerToHostEnabled bool `json:"enrichmentContainerToHostEnabled" split_words:"true" required:"false" default:"false"`
	EnrichmentEc2ToHostEnabled       bool `json:"enrichmentEc2ToHostEnabled" split_words:"true" required:"false" default:"false"`

	// discovery delay in ms
	DiscoveryDelayInMs int `json:"discoveryDelayInMs" split_words:"true" required:"false" default:"0"`
}

func IsPodZero() bool {
	return strings.HasSuffix(Config.PodName, "-0")
}

type AttributeUpdateSpecifications []AttributeUpdateSpecification

func (s *AttributeUpdateSpecifications) Decode(value string) error {
	var specs []AttributeUpdateSpecification
	err := json.Unmarshal([]byte(value), &specs)
	if err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	*s = specs
	return nil
}

type AttributeUpdateSpecification struct {
	Type          string  `json:"type" split_words:"true"`
	AttributeName string  `json:"attributeName" split_words:"true"`
	Rate          float64 `json:"rate" split_words:"true"`
	Interval      int     `json:"interval" split_words:"true"`
}

func (s *ContainerTargetCreationsAndDeletionsSpecification) Decode(value string) error {
	var spec ContainerTargetCreationsAndDeletionsSpecification
	err := json.Unmarshal([]byte(value), &spec)
	if err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	*s = spec
	return nil
}

type ContainerTargetCreationsAndDeletionsSpecification struct {
	Count    int `json:"count" split_words:"true"`
	Interval int `json:"interval" split_words:"true"`
}

var (
	Config Specification
)

func ParseConfiguration() {
	err := envconfig.Process("steadybit_extension", &Config)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to parse configuration from environment.")
	}
}

func ValidateConfiguration() {
}

func (s *Specification) FindAttributeUpdate(t string) *AttributeUpdateSpecification {
	for _, attributeUpdate := range s.AttributeUpdates {
		if attributeUpdate.Type == t {
			return &attributeUpdate
		}
	}
	return nil
}
