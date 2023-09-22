package extloadtest

import (
	"context"
	"github.com/procyon-projects/chrono"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-loadtest/config"
	"math/rand"
	"time"
)

var scheduler = chrono.NewDefaultTaskScheduler()

func targets(targets []discovery_kit_api.Target) func() discovery_kit_api.DiscoveryData {
	return func() discovery_kit_api.DiscoveryData {
		return discovery_kit_api.DiscoveryData{
			Targets: &targets,
		}
	}
}

func enrichmentData(enrichmentData []discovery_kit_api.EnrichmentData) func() discovery_kit_api.DiscoveryData {
	return func() discovery_kit_api.DiscoveryData {
		return discovery_kit_api.DiscoveryData{
			EnrichmentData: &enrichmentData,
		}
	}
}
func scheduleTargetAttributeUpdateIfNecessary(targets []discovery_kit_api.Target, typeId string) {
	scheduleAttributeUpdateIfNecessary(targets, typeId, func(target discovery_kit_api.Target) map[string][]string {
		return target.Attributes
	})
}

func scheduleEnrichmentDataAttributeUpdateIfNecessary(items []discovery_kit_api.EnrichmentData, typeId string) {
	scheduleAttributeUpdateIfNecessary(items, typeId, func(item discovery_kit_api.EnrichmentData) map[string][]string {
		return item.Attributes
	})
}

func scheduleAttributeUpdateIfNecessary[T any](items []T, typeId string, attributeMapAccessor func(item T) map[string][]string) {
	spec := config.Config.FindAttributeUpdate(typeId)
	if spec == nil {
		return
	}
	log.Info().
		Str("type", spec.Type).
		Str("attribute", spec.AttributeName).
		Float64("rate", spec.Rate).
		Int("interval", spec.Interval).
		Msg("scheduled attribute update")

	initAttributes(items, spec, attributeMapAccessor)

	delay := time.Duration(spec.Interval) * time.Second
	_, err := scheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		updateAttributes(items, spec, attributeMapAccessor)
	}, delay, chrono.WithTime(time.Now().Add(delay)))

	if err != nil {
		log.Fatal().Msgf("Failed to schedule attribute updates for '%s': %s", typeId, err.Error())
	}
}

func initAttributes[T any](items []T, spec *config.AttributeUpdateSpecification, getAttributeMap func(item T) map[string][]string) {
	for _, item := range items {
		getAttributeMap(item)[spec.AttributeName] = []string{time.Now().String()}
	}
}

func updateAttributes[T any](items []T, spec *config.AttributeUpdateSpecification, getAttributeMap func(item T) map[string][]string) {
	count := 0
	for _, item := range items {
		if rand.Float64() <= spec.Rate {
			count++
			getAttributeMap(item)[spec.AttributeName] = []string{time.Now().String()}
		}
	}
	log.Info().
		Str("type", spec.Type).
		Str("attribute", spec.AttributeName).
		Int("count", count).
		Msg("updated attributes")
}
