// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

package extloadtest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
)

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
