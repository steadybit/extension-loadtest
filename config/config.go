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
	"sync"
)

var (
	Config Specification

	// Mu guards the runtime-mutable Config fields (AttributeUpdates,
	// TargetReplacements, SimulateExtensionRestarts) against the data race between
	// the /config/* update handlers writing them and the discoveries reading them.
	// Readers take RLock via the Find* accessors; the update handlers take Lock
	// while unmarshalling the request body into Config.
	Mu sync.RWMutex
)

const (
	// MaxClockSkewSeconds is the cross-replica wall-clock skew the deterministic
	// serve-time updates tolerate. See extloadtest/deterministic.go.
	MaxClockSkewSeconds = 30
	// MinBucketIntervalFactor: every time-bucketed interval must be at least this
	// multiple of MaxClockSkewSeconds, bounding per-boundary disagreement to
	// <= 1/factor of the interval.
	MinBucketIntervalFactor = 20
)

// MinBucketIntervalSeconds is the smallest interval allowed for any time-bucketed
// spec (attribute updates, target replacements, extension restarts).
func MinBucketIntervalSeconds() int {
	return MaxClockSkewSeconds * MinBucketIntervalFactor
}

// Specification is the configuration specification for the extension. Configuration values can be applied
// through environment variables. Learn more through the documentation of the envconfig package.
// https://github.com/kelseyhightower/envconfig
type Specification struct {
	ClusterName string `json:"clusterName" split_words:"true" required:"false" default:"cluster-loadtest"`
	PodUID      string `json:"podUID" split_words:"true" required:"false" default:"PodUID1"`
	PodName     string `json:"podName" split_words:"true" required:"false" default:"pod-0"`

	//2 containers per pod * 4 pods per deployment * 5 deployments per node * 400 nodes = 16000 containers
	Ec2NodeCount       int `json:"ec2NodeCount" split_words:"true" required:"false" default:"10"`
	AzureNodeCount     int `json:"azureNodeCount" split_words:"true" required:"false" default:"1"`
	GcpNodeCount       int `json:"gcpNodeCount" split_words:"true" required:"false" default:"1"`
	DeploymentsPerNode int `json:"deploymentsPerNode" split_words:"true" required:"false" default:"10"`
	PodsPerDeployment  int `json:"podsPerDeployment" split_words:"true" required:"false" default:"1"`
	ContainerPerPod    int `json:"containerPerPod" split_words:"true" required:"false" default:"1"`

	AttributeUpdates          AttributeUpdateSpecifications          `split_words:"true" required:"false" default:"[]"`
	TargetReplacements        TargetReplacementsSpecifications       `split_words:"true" required:"false" default:"[]"`
	SimulateExtensionRestarts SimulateExtensionRestartSpecifications `split_words:"true" required:"false" default:"[]"`

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

	// Simulate an Event Listener
	EventListenerEnabled bool `json:"eventListenerEnabled" split_words:"true" required:"false" default:"true"`

	// discovery delay in ms
	DiscoveryDelayInMs int `json:"discoveryDelayInMs" split_words:"true" required:"false" default:"0"`
}

func IsPodZero() bool {
	return strings.HasSuffix(Config.PodName, "-0")
}

type AttributeUpdateSpecifications []AttributeUpdateSpecification

type AttributeUpdateSpecification struct {
	Type          string  `json:"type" split_words:"true"`
	AttributeName string  `json:"attributeName" split_words:"true"`
	Rate          float64 `json:"rate" split_words:"true"`
	Interval      int     `json:"interval" split_words:"true"`
}

func (s *AttributeUpdateSpecifications) Decode(value string) error {
	var specs []AttributeUpdateSpecification
	err := json.Unmarshal([]byte(value), &specs)
	if err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	*s = specs
	return nil
}

func (s *AttributeUpdateSpecification) Decode(value string) error {
	var spec AttributeUpdateSpecification
	err := json.Unmarshal([]byte(value), &spec)
	if err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	*s = spec
	return nil
}

type TargetReplacementsSpecifications []TargetReplacementsSpecification

type TargetReplacementsSpecification struct {
	Type     string `json:"type" split_words:"true"`
	Count    int    `json:"count" split_words:"true"`
	Interval int    `json:"interval" split_words:"true"`
}

func (s *TargetReplacementsSpecifications) Decode(value string) error {
	var specs []TargetReplacementsSpecification
	err := json.Unmarshal([]byte(value), &specs)
	if err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	*s = specs
	return nil
}

func (s *TargetReplacementsSpecification) Decode(value string) error {
	var spec TargetReplacementsSpecification
	err := json.Unmarshal([]byte(value), &spec)
	if err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	*s = spec
	return nil
}

type SimulateExtensionRestartSpecifications []SimulateExtensionRestartSpecification

type SimulateExtensionRestartSpecification struct {
	Type     string `json:"type" split_words:"true"`     // Type of the container to simulate a restart for
	Duration int    `json:"duration" split_words:"true"` // Duration how long should the targets be unavailable
	Interval int    `json:"interval" split_words:"true"` // Interval in seconds how often the restart should be simulated
}

func (s *SimulateExtensionRestartSpecifications) Decode(value string) error {
	var specs []SimulateExtensionRestartSpecification
	err := json.Unmarshal([]byte(value), &specs)
	if err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	*s = specs
	return nil
}

func (s *SimulateExtensionRestartSpecification) Decode(value string) error {
	var spec SimulateExtensionRestartSpecification
	err := json.Unmarshal([]byte(value), &spec)
	if err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	*s = spec
	return nil
}

func ParseConfiguration() {
	err := envconfig.Process("steadybit_extension", &Config)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to parse configuration from environment.")
	}
}

func ValidateConfiguration() {
	for _, tr := range Config.TargetReplacements {
		for _, sr := range Config.SimulateExtensionRestarts {
			if tr.Type == sr.Type {
				log.Fatal().Msgf("You can only use either target replacements or simulate extension restarts for type '%s', not both at the same time.", tr.Type)
			}
		}
	}

	// Serve-time updates are quantized into time buckets so every replica derives
	// the same value at the same instant. The bucket interval must be far larger
	// than the tolerated clock skew, otherwise replicas straddling a boundary
	// would disagree. Clamp too-small intervals up and warn; raise the rate/count
	// to preserve the change frequency.
	minInterval := MinBucketIntervalSeconds()
	for i := range Config.AttributeUpdates {
		if s := &Config.AttributeUpdates[i]; s.Interval > 0 && s.Interval < minInterval {
			log.Warn().Msgf("attribute-update interval %ds for '%s' is below the minimum %ds (%dx the %ds clock-skew budget); clamping to %ds. Raise 'rate' to keep the same change frequency.", s.Interval, s.Type, minInterval, MinBucketIntervalFactor, MaxClockSkewSeconds, minInterval)
			s.Interval = minInterval
		}
	}
	for i := range Config.TargetReplacements {
		if s := &Config.TargetReplacements[i]; s.Interval > 0 && s.Interval < minInterval {
			log.Warn().Msgf("target-replacement interval %ds for '%s' is below the minimum %ds; clamping to %ds. Raise 'count' to keep the same replacement frequency.", s.Interval, s.Type, minInterval, minInterval)
			s.Interval = minInterval
		}
	}
	for i := range Config.SimulateExtensionRestarts {
		if s := &Config.SimulateExtensionRestarts[i]; s.Interval > 0 && s.Interval < minInterval {
			log.Warn().Msgf("extension-restart interval %ds for '%s' is below the minimum %ds; clamping to %ds.", s.Interval, s.Type, minInterval, minInterval)
			s.Interval = minInterval
		}
		if s := &Config.SimulateExtensionRestarts[i]; s.Duration > 0 && s.Duration < MaxClockSkewSeconds {
			log.Warn().Msgf("extension-restart duration %ds for '%s' is below the %ds clock-skew budget; replicas may disagree at the window edges. Raising to %ds.", s.Duration, s.Type, MaxClockSkewSeconds, MaxClockSkewSeconds)
			s.Duration = MaxClockSkewSeconds
		}
		if s := &Config.SimulateExtensionRestarts[i]; s.Duration > 0 && s.Interval > 0 && s.Duration >= s.Interval {
			log.Fatal().Msgf("extension-restart duration %ds for '%s' must be less than interval %ds; otherwise the extension would be permanently down with no recovery window.", s.Duration, s.Type, s.Interval)
		}
	}
}

// FindAttributeUpdate returns the matching spec or nil. It returns a pointer to
// the range-loop copy (heap-allocated via escape analysis), never &s.AttributeUpdates[i],
// so callers may safely dereference it after RUnlock has fired. Any new accessor
// must keep this property — returning a slice-element pointer would let callers
// race a concurrent write through Mu. The Find* accessors below follow the same rule.
func (s *Specification) FindAttributeUpdate(t string) *AttributeUpdateSpecification {
	Mu.RLock()
	defer Mu.RUnlock()
	for _, update := range s.AttributeUpdates {
		if update.Type == t {
			return &update
		}
	}
	return nil
}

func (s *Specification) FindTargetReplacementsSpecification(t string) *TargetReplacementsSpecification {
	Mu.RLock()
	defer Mu.RUnlock()
	for _, replacement := range s.TargetReplacements {
		if replacement.Type == t {
			return &replacement
		}
	}
	return nil
}

func (s *Specification) FindSimulateExtensionRestartSpecification(t string) *SimulateExtensionRestartSpecification {
	Mu.RLock()
	defer Mu.RUnlock()
	for _, restart := range s.SimulateExtensionRestarts {
		if restart.Type == t {
			return &restart
		}
	}
	return nil
}
