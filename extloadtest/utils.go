// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

package extloadtest

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
)

func updateTargetId(targets []discovery_kit_api.Target, name, targetType string) []discovery_kit_api.Target {
	return updateId(targets, name, targetType, func(i discovery_kit_api.Target) string {
		return i.Id
	}, func(i *discovery_kit_api.Target, v string) {
		i.Id = v
	})
}

func updateEnrichmentDataId(dataList []discovery_kit_api.EnrichmentData, name, targetType string) []discovery_kit_api.EnrichmentData {
	return updateId(dataList, name, targetType, func(i discovery_kit_api.EnrichmentData) string {
		return i.Id
	}, func(i *discovery_kit_api.EnrichmentData, v string) {
		i.Id = v
	})
}

// updateId returns a copy of targets with the element matching name bumped to the
// next @version. A fresh slice is returned (copy-on-write) so concurrent readers
// holding the previous slice are never mutated under them. The previous in-place
// variant ranged over value copies and discarded the write, so recreate was a no-op.
func updateId[T any](targets []T, name, targetType string, getter func(i T) string, setter func(i *T, v string)) []T {
	updated := slices.Clone(targets)

	for i := range updated {
		id := getter(updated[i])
		if id != name {
			continue
		}

		version := 0
		idWithoutVersion := id
		if versionSeparator := strings.LastIndex(id, "@"); versionSeparator != -1 {
			version, _ = strconv.Atoi(id[versionSeparator+1:])
			idWithoutVersion = id[:versionSeparator]
		}
		newId := fmt.Sprintf("%s@%d", idWithoutVersion, version+1)
		setter(&updated[i], newId)

		log.Info().
			Str("type", targetType).
			Str("newId", newId).
			Str("id", id).
			Msg("recreated target")
		return updated
	}

	log.Warn().
		Str("type", targetType).
		Str("name", name).
		Msg("missing target to recreate")
	return updated
}
