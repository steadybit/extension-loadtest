// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package extloadtest

import (
	"hash/fnv"
	"math"
	"strconv"
	"time"

	"github.com/steadybit/extension-loadtest/config"
)

// attributeUpdateDisableThreshold mirrors the historic "rate is effectively
// zero" check: a rate at or below this disables the attribute entirely.
const attributeUpdateDisableThreshold = 0.00000000001

// clockSkewGuard is the wall-clock skew the deterministic serve-time
// computations tolerate. Discovery is served round-robin across replicas, so
// every replica must derive the same value for the same target at the same
// instant. We quantize time into buckets and only adopt a new bucket once this
// guard has elapsed past its boundary; bucket intervals are required to be much
// larger than this (see config.ValidateConfiguration). Real EKS skew (AWS Time
// Sync) is in the millisecond range, so this is generous headroom.
var clockSkewGuard = time.Duration(config.MaxClockSkewSeconds) * time.Second

func stableHash(parts ...string) uint64 {
	h := fnv.New64a()
	for _, p := range parts {
		_, _ = h.Write([]byte(p))
		_, _ = h.Write([]byte{0})
	}
	return h.Sum64()
}

// hashUnitFloat maps the given parts deterministically into [0,1).
func hashUnitFloat(parts ...string) float64 {
	return float64(stableHash(parts...)%1_000_000) / 1_000_000.0
}

func positiveMod(a, b int64) int64 {
	m := a % b
	if m < 0 {
		m += b
	}
	return m
}

// bucketIndex returns the time-bucket index for the given interval, applying the
// clock-skew guard so a new bucket is only adopted once clockSkewGuard past its
// boundary. Computed identically on every replica from wall-clock time.
func bucketIndex(now time.Time, intervalSeconds int) int64 {
	if intervalSeconds <= 0 {
		intervalSeconds = 1
	}
	return now.Add(-clockSkewGuard).Unix() / int64(intervalSeconds)
}

// deterministicAttributeValue returns the value a target's update-attribute must
// carry at the given time. A target's value is constant for `period` buckets and
// then rolls over, where period = round(1/rate); ~rate of all targets roll on
// each bucket boundary, matching the configured change rate. It is a pure
// function of (id, now), so every replica serves the identical value.
func deterministicAttributeValue(id string, spec *config.AttributeUpdateSpecification, now time.Time) string {
	cur := bucketIndex(now, spec.Interval)
	period := int64(1)
	if spec.Rate > 0 {
		period = int64(math.Round(1.0 / spec.Rate))
		if period < 1 {
			period = 1
		}
	}
	phase := int64(stableHash(id) % uint64(period))
	lastChange := cur - positiveMod(cur-phase, period)
	return time.Unix(lastChange*int64(spec.Interval), 0).UTC().Format(time.RFC3339Nano)
}

// isTargetReplaced reports whether the target is currently in the "removed" phase
// of a target-replacement cycle. Each bucket a deterministic ~Count/total
// fraction of targets is omitted, rotating per bucket. Pure function of
// (id, now); identical on every replica.
func isTargetReplaced(id string, total int, spec *config.TargetReplacementsSpecification, now time.Time) bool {
	if spec == nil || spec.Count <= 0 || total <= 0 {
		return false
	}
	fraction := float64(spec.Count) / float64(total)
	if fraction > 1 {
		fraction = 1
	}
	bucket := bucketIndex(now, spec.Interval)
	return hashUnitFloat(id, strconv.FormatInt(bucket, 10)) < fraction
}

// isExtensionDown reports whether the extension for a type is currently
// simulating a restart (all targets unavailable). The type is down for the first
// `Duration` seconds of every `Interval`-second cycle. Pure function of
// (type, now); identical on every replica.
func isExtensionDown(spec *config.SimulateExtensionRestartSpecification, now time.Time) bool {
	if spec == nil || spec.Interval <= 0 || spec.Duration <= 0 {
		return false
	}
	secInCycle := positiveMod(now.Add(-clockSkewGuard).Unix(), int64(spec.Interval))
	return secInCycle < int64(spec.Duration)
}
