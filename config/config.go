/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package config

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

// Specification is the configuration specification for the extension. Configuration values can be applied
// through environment variables. Learn more through the documentation of the envconfig package.
// https://github.com/kelseyhightower/envconfig
type Specification struct {
	ClusterName string `json:"clusterName" split_words:"true" required:"false" default:"cluster-1"`

	//2 containers per pod * 4 pods per deployment * 5 deployments per node * 400 nodes = 16000 containers
	NodeCount          int `json:"nodeCount" split_words:"true" required:"false" default:"2"`
	DeploymentsPerNode int `json:"deploymentsPerNode" split_words:"true" required:"false" default:"5"`
	PodsPerDeployment  int `json:"podsPerDeployment" split_words:"true" required:"false" default:"2"`
	ContainerPerPod    int `json:"containerPerPod" split_words:"true" required:"false" default:"2"`

	AttributeUpdates AttributeUpdateSpecifications `split_words:"true" required:"false" default:"[{\"type\": \"com.steadybit.extension_aws.ec2-instance\", \"attributeName\": \"aws-ec2.label.change-ts\", \"rate\": 0.20, \"interval\": 600},{\"type\": \"com.steadybit.extension_container.container\", \"attributeName\": \"container.label.change-ts\", \"rate\": 0.20, \"interval\": 180},{\"type\": \"com.steadybit.extension_kubernetes.kubernetes-container\", \"attributeName\": \"k8s.label.change-ts\", \"rate\": 0.20, \"interval\": 180},{\"type\": \"com.steadybit.extension_kubernetes.kubernetes-deployment\", \"attributeName\": \"k8s.label.change-ts\", \"rate\": 0.20, \"interval\": 180}]"`

	DiscoveryAttributeExcludesContainer            []string `json:"discoveryAttributeExcludesContainer" split_words:"true" required:"false"`
	DiscoveryAttributeExcludesEc2                  []string `json:"discoveryAttributeExcludesEc2" split_words:"true" required:"false"`
	DiscoveryAttributeExcludesHost                 []string `json:"discoveryAttributeExcludesHost" split_words:"true" required:"false"`
	DiscoveryAttributeExcludesKubernetesDeployment []string `json:"discoveryAttributeExcludesKubernetesDeployment" split_words:"true" required:"false"`
	DiscoveryAttributeExcludesKubernetesContainer  []string `json:"discoveryAttributeExcludesKubernetesContainer" split_words:"true" required:"false"`
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
	Rate          float64 `json:"rate" split_words:"true" default:"0.20"`
	Interval      int     `json:"interval" split_words:"true" default:"180s"`
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
