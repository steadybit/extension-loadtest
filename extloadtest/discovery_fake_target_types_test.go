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

	"github.com/steadybit/action-kit/go/action_kit_api/v2"

	"github.com/steadybit/extension-loadtest/config"
	"github.com/stretchr/testify/require"
)

// withFakeTargetsPerType sets the config for one test and restores it afterwards,
// so a later test in this package cannot inherit a mutated package-level Config.
func withFakeTargetsPerType(t *testing.T, perType int) {
	t.Helper()
	previous := config.Config
	t.Cleanup(func() { config.Config = previous })
	config.Config.FakeTargetsPerType = perType
}

func withFakeTargetTypeCount(t *testing.T, count int) {
	t.Helper()
	previous := config.Config
	t.Cleanup(func() { config.Config = previous })
	config.Config.FakeTargetTypeCount = count
}

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
	withFakeTargetsPerType(t, 1)

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
	withFakeTargetsPerType(t, 1)

	require.Equal(t, "Flux Capacitor (none icon)", newFakeTargetDiscovery(0).DescribeTarget().Label.One)
	require.Equal(t, "Hoverboards (broken icon)", newFakeTargetDiscovery(1).DescribeTarget().Label.Other)
	require.Equal(t, "Rubber Duck (valid icon)", newFakeTargetDiscovery(2).DescribeTarget().Label.One)
}

// The shared attribute is described exactly once, no matter how many types run,
// otherwise the sdk logs a duplicate warning per type.
func TestSharedIconAttributeIsDescribedOnce(t *testing.T) {
	withFakeTargetsPerType(t, 1)

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
	withFakeTargetsPerType(t, 2)

	description := newFakeTargetDiscovery(0).DescribeTarget()
	require.Nil(t, description.Icon)
	require.Equal(t, "com.steadybit.extension_loadtest.flux-capacitor", description.Id)
	require.Equal(t, fakeTargetTypeCategory, *description.Category)
	require.Equal(t, "Flux Capacitors (none icon)", description.Label.Other)
	require.NotEmpty(t, description.Table.Columns)
}

// Every fake type shares one category, so they group together in the UI.
func TestAllFakeTargetTypesShareTheDebugCategory(t *testing.T) {
	withFakeTargetsPerType(t, 1)

	for i := 0; i < len(fakeTargetTypeSpecs); i++ {
		require.Equal(t, "Debug", *newFakeTargetDiscovery(i).DescribeTarget().Category)
	}
}

func TestFakeTargetTypeIdsAreUniqueBeyondThePool(t *testing.T) {
	withFakeTargetsPerType(t, 1)

	seen := make(map[string]bool)
	for i := 0; i < len(fakeTargetTypeSpecs)*2+3; i++ {
		id := newFakeTargetDiscovery(i).DescribeTarget().Id
		require.Falsef(t, seen[id], "duplicate fake target type id %s at index %d", id, i)
		seen[id] = true
	}
}

func TestFakeTargetsCarryTheirTypeAndAttributes(t *testing.T) {
	withFakeTargetsPerType(t, 3)
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
	withFakeTargetTypeCount(t, 0)
	require.NotPanics(t, RegisterFakeTargetTypeDiscoveries)
}

func TestFakeAttributeDescriptionsMatchTheColumns(t *testing.T) {
	withFakeTargetsPerType(t, 1)

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

// Each fake target type needs its own action, otherwise it cannot be reached from
// an experiment - which is where its icon (or the fallback) is rendered.
func TestEveryFakeTargetTypeHasItsOwnAction(t *testing.T) {
	withFakeTargetsPerType(t, 1)

	ids := make(map[string]bool)
	for i := 0; i < len(fakeTargetTypeSpecs); i++ {
		targetType := newFakeTargetDiscovery(i).DescribeTarget().Id
		description := fakeTargetTypeAction(i).Describe()

		require.NotNil(t, description.TargetSelection)
		require.Equal(t, targetType, description.TargetSelection.TargetType)
		require.Falsef(t, ids[description.Id], "duplicate action id %s at index %d", description.Id, i)
		ids[description.Id] = true
	}
	require.Len(t, ids, len(fakeTargetTypeSpecs))
}

// The action picker shows the label and nothing else, so a shared "Do Nothing"
// across all of them would be unusable.
func TestFakeTargetTypeActionLabelsNameTheirType(t *testing.T) {
	withFakeTargetsPerType(t, 1)

	require.Equal(t, "Do Nothing (Flux Capacitor)", fakeTargetTypeAction(0).Describe().Label)
	require.Equal(t, "Do Nothing (Hoverboard)", fakeTargetTypeAction(1).Describe().Label)

	// the shared constructor keeps its own label for the pre-existing registrations
	require.Equal(t, "Do Nothing", NewDoNothingAction("com.steadybit.extension_container.container",
		action_kit_api.TargetSelectionTemplate{}).Describe().Label)
}

// Technology is the umbrella, category the sub-grouping inside it: the pseudo
// target type actions must not land among the long-standing loadtest ones.
func TestFakeTargetTypeActionsHaveTheirOwnCategory(t *testing.T) {
	withFakeTargetsPerType(t, 1)

	for i := 0; i < len(fakeTargetTypeSpecs); i++ {
		description := fakeTargetTypeAction(i).Describe()
		require.Equal(t, actionTechnology, *description.Technology)
		require.Equal(t, actionCategoryBrokenTarget, *description.Category)
	}

	// the pre-existing registrations stay in the plain loadtest category
	existing := NewDoNothingAction("com.steadybit.extension_container.container",
		action_kit_api.TargetSelectionTemplate{}).Describe()
	require.Equal(t, actionTechnology, *existing.Technology)
	require.Equal(t, actionCategoryLoadtest, *existing.Category)
}

// The selection template must query an attribute the targets actually carry.
func TestFakeTargetTypeActionSelectionTemplateMatchesTheTargets(t *testing.T) {
	withFakeTargetsPerType(t, 1)

	for i := 0; i < len(fakeTargetTypeSpecs); i++ {
		discovery := newFakeTargetDiscovery(i)
		targets, err := discovery.DiscoverTargets(context.Background())
		require.NoError(t, err)

		templates := fakeTargetTypeAction(i).Describe().TargetSelection.SelectionTemplates
		require.NotNil(t, templates)
		require.Len(t, *templates, 1)

		attribute, _, found := strings.Cut((*templates)[0].Query, "=")
		require.True(t, found)
		require.NotEmptyf(t, targets[0].Attributes[attribute],
			"selection template queries %s, which the targets do not carry", attribute)
	}
}
