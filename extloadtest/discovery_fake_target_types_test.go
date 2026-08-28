// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package extloadtest

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/steadybit/extension-loadtest/config"
	"github.com/stretchr/testify/require"
)

func TestFakeTargetTypeIconsRotate(t *testing.T) {
	// none -> broken -> valid, so a count of >= 3 always covers all three cases
	for index, kind := range map[int]string{
		0: iconKindNone, 1: iconKindBroken, 2: iconKindValid,
		3: iconKindNone, 4: iconKindBroken, 5: iconKindValid,
	} {
		require.Equalf(t, kind, fakeTargetTypeIconKind(index), "index %d", index)
	}

	require.Nil(t, iconForKind(iconKindNone))
	require.Equal(t, brokenIconDataUri, *iconForKind(iconKindBroken))
	require.Equal(t, validIconDataUri, *iconForKind(iconKindValid))
}

// The attribute must say what the type's icon actually is, otherwise it cannot be
// trusted to explain the fallback the UI renders.
func TestIconAttributeMatchesTheDescribedIcon(t *testing.T) {
	config.Config.FakeTargetsPerType = 1

	for i := 0; i < len(fakeTargetTypeSpecs); i++ {
		discovery := newFakeTargetDiscovery(i)
		targets, err := discovery.DiscoverTargets(context.Background())
		require.NoError(t, err)

		kind := targets[0].Attributes[fakeIconAttribute]
		icon := discovery.DescribeTarget().Icon

		switch {
		case icon == nil:
			require.Equal(t, []string{iconKindNone}, kind)
		case *icon == brokenIconDataUri:
			require.Equal(t, []string{iconKindBroken}, kind)
		default:
			require.Equal(t, []string{iconKindValid}, kind)
		}
	}
}

// The type label names the icon case too - a target type carries no attributes.
func TestFakeTargetTypeLabelNamesTheIconCase(t *testing.T) {
	config.Config.FakeTargetsPerType = 1

	require.Equal(t, "Flux Capacitor (none icon)", newFakeTargetDiscovery(0).DescribeTarget().Label.One)
	require.Equal(t, "Hoverboards (broken icon)", newFakeTargetDiscovery(1).DescribeTarget().Label.Other)
	require.Equal(t, "Rubber Duck (valid icon)", newFakeTargetDiscovery(2).DescribeTarget().Label.One)
}

// The shared attribute is described exactly once, no matter how many types run,
// otherwise the sdk logs a duplicate warning per type.
func TestSharedIconAttributeIsDescribedOnce(t *testing.T) {
	config.Config.FakeTargetsPerType = 1

	for i := 0; i < len(fakeTargetTypeSpecs); i++ {
		for _, described := range newFakeTargetDiscovery(i).DescribeAttributes() {
			require.NotEqual(t, fakeIconAttribute, described.Attribute)
		}
	}

	shared := ltFakeSharedAttributeDescriber{}.DescribeAttributes()
	require.Len(t, shared, 1)
	require.Equal(t, fakeIconAttribute, shared[0].Attribute)
}

// The valid icon must actually render, otherwise the rotation has no control case.
func TestValidIconIsAWellFormedSvgDataUri(t *testing.T) {
	payload, found := strings.CutPrefix(validIconDataUri, "data:image/svg+xml,")
	require.True(t, found)

	decoded, err := url.QueryUnescape(payload)
	require.NoError(t, err)
	require.NoError(t, xml.Unmarshal([]byte(decoded), new(struct {
		XMLName xml.Name `xml:"svg"`
	})))
}

// The description must omit the icon entirely for the 'none' slot of the rotation -
// that is the case the fake types exist to reproduce.
func TestFakeTargetDescriptionHasNoIcon(t *testing.T) {
	config.Config.FakeTargetsPerType = 2

	description := newFakeTargetDiscovery(0).DescribeTarget()
	require.Nil(t, description.Icon)
	require.Equal(t, "com.steadybit.extension_loadtest.flux-capacitor", description.Id)
	require.Equal(t, fakeTargetTypeCategory, *description.Category)
	require.Equal(t, "Flux Capacitors (none icon)", description.Label.Other)
	require.NotEmpty(t, description.Table.Columns)
}

// Every fake type shares one category, so they group together in the UI.
func TestAllFakeTargetTypesShareTheDebugCategory(t *testing.T) {
	config.Config.FakeTargetsPerType = 1

	for i := 0; i < len(fakeTargetTypeSpecs); i++ {
		require.Equal(t, "Debug", *newFakeTargetDiscovery(i).DescribeTarget().Category)
	}
}

func TestFakeTargetTypeIdsAreUniqueBeyondThePool(t *testing.T) {
	config.Config.FakeTargetsPerType = 1

	seen := make(map[string]bool)
	for i := 0; i < len(fakeTargetTypeSpecs)*2+3; i++ {
		id := newFakeTargetDiscovery(i).DescribeTarget().Id
		require.Falsef(t, seen[id], "duplicate fake target type id %s at index %d", id, i)
		seen[id] = true
	}
}

func TestFakeTargetsCarryTheirTypeAndAttributes(t *testing.T) {
	config.Config.FakeTargetsPerType = 3
	config.Config.PodUID = "PodUID1"

	discovery := newFakeTargetDiscovery(1)
	targets, err := discovery.DiscoverTargets(context.Background())
	require.NoError(t, err)
	require.Len(t, targets, 3)

	for _, target := range targets {
		require.Equal(t, "com.steadybit.extension_loadtest.hoverboard", target.TargetType)
		require.Equal(t, []string{"true"}, target.Attributes["steadybit.loadtest"])
		require.NotEmpty(t, target.Attributes["loadtest.hoverboard.name"])
	}

	// the description's columns must reference attributes the targets actually have
	for _, column := range discovery.DescribeTarget().Table.Columns {
		require.NotEmptyf(t, targets[0].Attributes[column.Attribute],
			"column %s has no matching attribute on the discovered targets", column.Attribute)
	}
}

func TestFakeTargetTypesCanBeTurnedOff(t *testing.T) {
	config.Config.FakeTargetTypeCount = 0
	require.NotPanics(t, RegisterFakeTargetTypeDiscoveries)
}

func TestFakeAttributeDescriptionsMatchTheColumns(t *testing.T) {
	config.Config.FakeTargetsPerType = 1

	discovery := newFakeTargetDiscovery(4)
	described := make(map[string]bool)
	all := append(discovery.DescribeAttributes(), ltFakeSharedAttributeDescriber{}.DescribeAttributes()...)
	for _, attribute := range all {
		described[attribute.Attribute] = true
		require.True(t, len(attribute.Label.One) > 0 && len(attribute.Label.Other) > 0)
	}
	for _, column := range discovery.DescribeTarget().Table.Columns {
		require.Truef(t, described[column.Attribute], "column %s has no attribute description", column.Attribute)
	}
	require.Equal(t, fmt.Sprintf("loadtest.%s.name", discovery.spec.key), discovery.attribute("name"))
}
