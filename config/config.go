/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

// Specification is the configuration specification for the extension. Configuration values can be applied
// through environment variables. Learn more through the documentation of the envconfig package.
// https://github.com/kelseyhightower/envconfig
type Specification struct {
	ClusterName string `json:"clusterName" split_words:"true" required:"false" default:"cluster-1"`

	NodeCount                     int `json:"nodeCount" split_words:"true" required:"false" default:"2"`
	DeploymentsPerNode            int `json:"deploymentsPerNode" split_words:"true" required:"false" default:"5"`
	PodsPerDeployment             int `json:"podsPerDeployment" split_words:"true" required:"false" default:"2"`
	ContainerPerPod               int `json:"containerPerPod" split_words:"true" required:"false" default:"2"`
	ChangeRateKubernetesContainer int `json:"changeRateKubernetesContainer" split_words:"true" required:"false" default:"13"` // 13% of containers will be changed
	ChangeTimeKubernetesContainer int `json:"changeRateContainer" split_words:"true" required:"false" default:"180"`          // 180 seconds between changes

	//2 containers per pod * 4 pods per deployment * 5 deployments per node * 400 nodes = 16000 containers
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
