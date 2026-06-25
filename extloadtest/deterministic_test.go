// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package extloadtest

import (
	"strconv"
	"testing"
	"time"

	"github.com/steadybit/extension-loadtest/config"
	"github.com/stretchr/testify/require"
)

func attrSpec(rate float64, interval int) *config.AttributeUpdateSpecification {
	return &config.AttributeUpdateSpecification{Type: "t", AttributeName: "x.no-copy-rule", Rate: rate, Interval: interval}
}

// The whole point: two replicas (any two evaluations) must agree for the same
// (id, now) — this is the regression test for the cross-replica flap.
func TestAttributeValueIsReplicaConsistent(t *testing.T) {
	spec := attrSpec(0.1, 600)
	now := time.Unix(1_900_000_123, 0)
	for _, id := range []string{"a", "b", "deployment-42", "containerd://fix-fix-d-1-ec2-Pod-1-c-1"} {
		replicaA := deterministicAttributeValue(id, spec, now)
		replicaB := deterministicAttributeValue(id, spec, now)
		require.Equal(t, replicaA, replicaB)
	}
}

// Restart-invariance: the value depends only on (id, now), never on process
// state, so a restarted pod immediately serves identical values (no init burst).
func TestAttributeValueStableUnderSkewWithinBucket(t *testing.T) {
	spec := attrSpec(0.1, 600)
	base := time.Unix(600*1000+300, 0) // ~mid-bucket, clear of boundaries
	id := "deployment-7"
	want := deterministicAttributeValue(id, spec, base)
	for _, d := range []time.Duration{0, time.Second, 10 * time.Second, 30 * time.Second, -30 * time.Second} {
		require.Equalf(t, want, deterministicAttributeValue(id, spec, base.Add(d)),
			"a clock skew of %s mid-bucket must not change the value", d)
	}
}

func TestAttributeValueChangesAcrossBuckets(t *testing.T) {
	spec := attrSpec(0.5, 600) // period 2
	id := "x"
	seen := map[string]bool{}
	for b := int64(0); b < 12; b++ {
		seen[deterministicAttributeValue(id, spec, time.Unix(b*600+300, 0))] = true
	}
	require.Greater(t, len(seen), 1, "value must roll over across buckets")
}

func TestAttributeChangeRateMatchesConfiguredRate(t *testing.T) {
	const rate, n = 0.1, 20000
	spec := attrSpec(rate, 600)
	b0 := time.Unix(1000*600+300, 0)
	b1 := time.Unix(1001*600+300, 0)
	changed := 0
	for i := 0; i < n; i++ {
		id := "id-" + strconv.Itoa(i)
		if deterministicAttributeValue(id, spec, b0) != deterministicAttributeValue(id, spec, b1) {
			changed++
		}
	}
	require.InDelta(t, rate, float64(changed)/n, 0.02, "fraction changing per bucket should match the rate")
}

func TestApplyAttributeUpdateDisabledByZeroRate(t *testing.T) {
	now := time.Unix(600*10+300, 0)
	attrs := map[string][]string{}
	applyAttributeUpdate(attrs, "id", attrSpec(0, 600), now)
	require.Empty(t, attrs, "rate 0 must emit no attribute")
	applyAttributeUpdate(attrs, "id", nil, now)
	require.Empty(t, attrs, "no spec must emit no attribute")
	applyAttributeUpdate(attrs, "id", attrSpec(0.1, 600), now)
	require.Len(t, attrs, 1, "an enabled spec must emit the attribute")
}

func TestReplacementOmitsApproximatelyCountAndRotates(t *testing.T) {
	const total = 5000
	spec := &config.TargetReplacementsSpecification{Type: "t", Count: 50, Interval: 600}
	b0 := time.Unix(3000*600+300, 0)
	b1 := time.Unix(3001*600+300, 0)

	omitted, rotated := 0, false
	for i := 0; i < total; i++ {
		id := "id-" + strconv.Itoa(i)
		o0 := isTargetReplaced(id, total, spec, b0)
		if o0 {
			omitted++
		}
		// deterministic within a bucket
		require.Equal(t, o0, isTargetReplaced(id, total, spec, b0))
		if o0 != isTargetReplaced(id, total, spec, b1) {
			rotated = true
		}
	}
	require.InDelta(t, spec.Count, omitted, float64(spec.Count), "~count targets omitted per bucket")
	require.True(t, rotated, "the omitted set must rotate across buckets")
	require.False(t, isTargetReplaced("id-1", total, nil, b0), "no spec means nothing replaced")
}

func TestExtensionDownWindow(t *testing.T) {
	spec := &config.SimulateExtensionRestartSpecification{Type: "t", Duration: 120, Interval: 600}
	guard := int64(config.MaxClockSkewSeconds)
	cycleStart := int64(5000 * 600)
	require.True(t, isExtensionDown(spec, time.Unix(cycleStart+guard+10, 0)), "should be down early in the cycle")
	require.False(t, isExtensionDown(spec, time.Unix(cycleStart+guard+300, 0)), "should be up after the duration")
	require.False(t, isExtensionDown(nil, time.Unix(cycleStart+guard+10, 0)), "no spec means never down")
}
