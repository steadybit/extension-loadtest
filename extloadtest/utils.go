package extloadtest

import (
	"context"
	"fmt"
	"github.com/procyon-projects/chrono"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/extension-loadtest/config"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var scheduler = chrono.NewDefaultTaskScheduler()

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

func scheduleContainerTargetChanges(containers *[]discovery_kit_api.Target, backup *[]discovery_kit_api.Target) {
	if config.Config.ContainerTargetCreationsAndDeletions.Count > 0 {
		interval := config.Config.ContainerTargetCreationsAndDeletions.Interval
		delay := time.Duration(interval) * time.Second
		_, err := scheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
			//restore previously deleted containers
			restoredCount := len(*backup)
			*containers = append(*containers, *backup...)
			*backup = []discovery_kit_api.Target{}

			//delete random containers
			deletedCount := rand.Intn(config.Config.ContainerTargetCreationsAndDeletions.Count)
			log.Debug().Int("count", deletedCount).Msgf("Deleted containers")
			for i := 0; i < deletedCount; i++ {
				index := rand.Intn(len(*containers))
				*backup = append(*backup, (*containers)[index])
				*containers = append((*containers)[:index], (*containers)[index+1:]...)
			}
			log.Info().Msgf("Deleted %d containers and (re-)added %d containers. Total count is now %d", deletedCount, restoredCount, len(*containers))
		}, delay, chrono.WithTime(time.Now().Add(delay)))
		if err != nil {
			log.Fatal().Msgf("Failed to schedule container target changes: %s", err.Error())
		}
		log.Info().
			Int("interval", interval).
			Int("maxCount", config.Config.ContainerTargetCreationsAndDeletions.Count).
			Msg("scheduled container target creation/deletion")
	}
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

func updateTargetId(targets []discovery_kit_api.Target, name, targetType string) {
	updateId(targets, name, targetType, func(i discovery_kit_api.Target) string {
		return i.Id
	}, func(i *discovery_kit_api.Target, v string) {
		i.Id = v
	})
}

func updateEnrichmentDataId(dataList []discovery_kit_api.EnrichmentData, name, targetType string) {
	updateId(dataList, name, targetType, func(i discovery_kit_api.EnrichmentData) string {
		return i.Id
	}, func(i *discovery_kit_api.EnrichmentData, v string) {
		i.Id = v
	})
}

func updateId[T any](targets []T, name, targetType string, getter func(i T) string, setter func(i *T, v string)) {
	for _, target := range targets {
		id := getter(target)

		if id == name {
			version := 0
			idWithoutVersion := id
			if versionSeparator := strings.LastIndex(id, "@"); versionSeparator != -1 {
				version, _ = strconv.Atoi(id[versionSeparator+1:])
				idWithoutVersion = id[:versionSeparator]
			}
			newId := fmt.Sprintf("%s@%d", idWithoutVersion, version+1)
			setter(&target, newId)

			log.Info().
				Str("type", targetType).
				Str("newId", newId).
				Str("id", id).
				Msg("recreated target")
			return
		}
	}

	log.Warn().
		Str("type", targetType).
		Str("name", name).
		Msg("missing target to recreate")
}
