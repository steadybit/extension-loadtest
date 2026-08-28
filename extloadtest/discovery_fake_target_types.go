// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package extloadtest

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_sdk"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-loadtest/config"
)

// The fake target types exist to exercise the platform and UI with target types
// this extension owns itself - unlike every other discovery here, which borrows
// the target type ids (and therefore the descriptions, including the icons) of
// the real extensions. Their deliberately made-up names keep them apart from
// anything a customer would ever discover. Each one comes with a no-op action, so
// it can be used in an experiment - which is where the icon (or the fallback the
// UI renders for it) shows up next to the selected targets.
//
// Their main purpose is the icon handling, so the registered types rotate through
// the three cases: no icon at all, an icon that cannot be rendered, and a perfectly
// good one. A count of >= 3 therefore always covers all of them at once.

// fakeTargetTypeIdPrefix is the namespace of all fake target types.
const fakeTargetTypeIdPrefix = "com.steadybit.extension_loadtest."

// fakeTargetTypeCategory groups every fake target type under one heading, so they
// stay together and apart from the real ones wherever the UI groups by category.
const fakeTargetTypeCategory = "Debug"

// fakeIconAttribute is shared by every fake target type - unlike the per-type
// name/serial attributes - so the explorer can group all of them at once and put
// the three icon cases side by side. Its value is one of the iconKind* constants.
const fakeIconAttribute = "loadtest.icon"

// The icon cases the registered types rotate through, as seen from a target.
const (
	iconKindNone   = "none"
	iconKindBroken = "broken"
	iconKindValid  = "valid"
)

// brokenIconDataUri is a well-formed data URI whose payload is not an image, so
// a browser using it as an image source fails to decode it and reports an error
// instead of a successful load.
const brokenIconDataUri = "data:image/png;base64,bm90LWFuLWltYWdlLWF0LWFsbA=="

// validIconDataUri is a plain star, so a type carrying it is obviously distinct
// from one that renders a fallback for a missing or unrenderable icon.
const validIconDataUri = "data:image/svg+xml,%3Csvg%20width%3D%2224%22%20height%3D%2224%22%20viewBox%3D%220%200%2024%2024%22%20fill%3D%22none%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%3E%3Cpath%20d%3D%22M12%202l2.9%206.2%206.6.9-4.8%204.7%201.2%206.7L12%2017.4%206.1%2020.5l1.2-6.7L2.5%209.1l6.6-.9L12%202z%22%20fill%3D%22%231D2632%22%2F%3E%3C%2Fsvg%3E"

type fakeTargetTypeSpec struct {
	// key is the last segment of the target type id and the attribute namespace.
	key   string
	one   string
	other string
}

// fakeTargetTypeSpecs is the pool the registered types are taken from. A count
// beyond the pool wraps around and disambiguates with a numeric suffix.
var fakeTargetTypeSpecs = []fakeTargetTypeSpec{
	{key: "flux-capacitor", one: "Flux Capacitor", other: "Flux Capacitors"},
	{key: "hoverboard", one: "Hoverboard", other: "Hoverboards"},
	{key: "rubber-duck", one: "Rubber Duck", other: "Rubber Ducks"},
	{key: "sock-drawer", one: "Sock Drawer", other: "Sock Drawers"},
	{key: "tea-kettle", one: "Tea Kettle", other: "Tea Kettles"},
	{key: "moon-cheese", one: "Moon Cheese Wheel", other: "Moon Cheese Wheels"},
	{key: "yeti-sighting", one: "Yeti Sighting", other: "Yeti Sightings"},
	{key: "banana-stand", one: "Banana Stand", other: "Banana Stands"},
}

// ltFakeTargetDiscovery adds the target type and attribute descriptions on top of
// the regular loadtest discovery, so the fake types take part in the simulated
// attribute updates, target replacements and extension restarts like any other.
type ltFakeTargetDiscovery struct {
	ltTargetDiscovery
	targetType string
	spec       fakeTargetTypeSpec
	icon       *string
	iconKind   string
}

var (
	_ discovery_kit_sdk.TargetDiscovery    = (*ltFakeTargetDiscovery)(nil)
	_ discovery_kit_sdk.TargetDescriber    = (*ltFakeTargetDiscovery)(nil)
	_ discovery_kit_sdk.AttributeDescriber = (*ltFakeTargetDiscovery)(nil)
)

func (l *ltFakeTargetDiscovery) DescribeTarget() discovery_kit_api.TargetDescription {
	return discovery_kit_api.TargetDescription{
		Id:       l.targetType,
		Version:  extbuild.GetSemverVersionStringOrUnknown(),
		Icon:     l.icon,
		Category: new(fakeTargetTypeCategory),
		// The label carries the icon case too: a target type has no attributes of its
		// own, so this is the only place the fallback can be named where the type
		// itself is rendered.
		Label: discovery_kit_api.PluralLabel{
			One:   fmt.Sprintf("%s (%s icon)", l.spec.one, l.iconKind),
			Other: fmt.Sprintf("%s (%s icon)", l.spec.other, l.iconKind),
		},
		Table: discovery_kit_api.Table{
			Columns: []discovery_kit_api.Column{
				{Attribute: l.attribute("name")},
				{Attribute: fakeIconAttribute},
				{Attribute: l.attribute("serial")},
			},
			OrderBy: []discovery_kit_api.OrderBy{
				{Attribute: l.attribute("name"), Direction: discovery_kit_api.ASC},
			},
		},
	}
}

// DescribeAttributes covers the per-type attributes only. fakeIconAttribute is
// shared by every type and is described once by ltFakeSharedAttributeDescriber,
// otherwise the sdk would warn about a duplicate for every registered type.
func (l *ltFakeTargetDiscovery) DescribeAttributes() []discovery_kit_api.AttributeDescription {
	return []discovery_kit_api.AttributeDescription{
		{
			Attribute: l.attribute("name"),
			Label:     discovery_kit_api.PluralLabel{One: l.spec.one + " name", Other: l.spec.one + " names"},
		},
		{
			Attribute: l.attribute("serial"),
			Label:     discovery_kit_api.PluralLabel{One: l.spec.one + " serial", Other: l.spec.one + " serials"},
		},
	}
}

// ltFakeSharedAttributeDescriber describes the attributes every fake target type
// shares, so they are registered exactly once no matter how many types there are.
type ltFakeSharedAttributeDescriber struct{}

var _ discovery_kit_sdk.AttributeDescriber = (*ltFakeSharedAttributeDescriber)(nil)

func (ltFakeSharedAttributeDescriber) DescribeAttributes() []discovery_kit_api.AttributeDescription {
	return []discovery_kit_api.AttributeDescription{
		{
			Attribute: fakeIconAttribute,
			Label:     discovery_kit_api.PluralLabel{One: "Icon case", Other: "Icon cases"},
		},
	}
}

func (l *ltFakeTargetDiscovery) attribute(name string) string {
	return fmt.Sprintf("loadtest.%s.%s", l.spec.key, name)
}

// fakeTargetTypeSpecAt returns the spec for the index-th (0-based) fake target
// type, wrapping around the pool and suffixing the repeats so ids stay unique.
func fakeTargetTypeSpecAt(index int) fakeTargetTypeSpec {
	spec := fakeTargetTypeSpecs[index%len(fakeTargetTypeSpecs)]
	if round := index / len(fakeTargetTypeSpecs); round > 0 {
		suffix := round + 1
		spec.key = fmt.Sprintf("%s-%d", spec.key, suffix)
		spec.one = fmt.Sprintf("%s %d", spec.one, suffix)
		spec.other = fmt.Sprintf("%s %d", spec.other, suffix)
	}
	return spec
}

// fakeTargetTypeIconKind returns the icon case of the index-th (0-based) fake
// target type, rotating through no icon at all, an unrenderable one and a valid one.
func fakeTargetTypeIconKind(index int) string {
	switch index % 3 {
	case 1:
		return iconKindBroken
	case 2:
		return iconKindValid
	default:
		return iconKindNone
	}
}

// iconForKind returns the icon a target type of that case ships, nil for none.
func iconForKind(kind string) *string {
	switch kind {
	case iconKindBroken:
		return new(brokenIconDataUri)
	case iconKindValid:
		return new(validIconDataUri)
	default:
		return nil
	}
}

func createFakeTargets(spec fakeTargetTypeSpec, targetType string, iconKind string, count int) []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, count)
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%s-%s-%d", spec.key, config.Config.PodUID, i)
		result = append(result, discovery_kit_api.Target{
			Id:         name,
			TargetType: targetType,
			Label:      name,
			Attributes: map[string][]string{
				fmt.Sprintf("loadtest.%s.name", spec.key):   {name},
				fmt.Sprintf("loadtest.%s.serial", spec.key): {fmt.Sprintf("%08d", i)},
				// what the target type's icon looks like, so the fallback the UI
				// renders can be told apart from a working icon without guessing
				fakeIconAttribute:    {iconKind},
				"steadybit.loadtest": {"true"},
			},
		})
	}
	return result
}

func newFakeTargetDiscovery(index int) *ltFakeTargetDiscovery {
	spec := fakeTargetTypeSpecAt(index)
	targetType := fakeTargetTypeIdPrefix + spec.key
	iconKind := fakeTargetTypeIconKind(index)
	targets := createFakeTargets(spec, targetType, iconKind, config.Config.FakeTargetsPerType)

	return &ltFakeTargetDiscovery{
		ltTargetDiscovery: ltTargetDiscovery{
			description: func() discovery_kit_api.DiscoveryDescription {
				return discovery_kit_api.DiscoveryDescription{
					Id: targetType,
					Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
						CallInterval: new("1m"),
					},
				}
			},
			snapshot: func() []discovery_kit_api.Target { return targets },
		},
		targetType: targetType,
		spec:       spec,
		icon:       iconForKind(iconKind),
		iconKind:   iconKind,
	}
}

// RegisterFakeTargetTypeDiscoveries registers the fake target types. Setting
// STEADYBIT_EXTENSION_FAKE_TARGET_TYPE_COUNT to 0 registers none at all.
func RegisterFakeTargetTypeDiscoveries() {
	if config.Config.FakeTargetTypeCount <= 0 {
		return
	}

	discovery_kit_sdk.Register(ltFakeSharedAttributeDescriber{})

	for i := 0; i < config.Config.FakeTargetTypeCount; i++ {
		discovery := newFakeTargetDiscovery(i)
		discovery_kit_sdk.Register(discovery)
		log.Info().
			Str("type", discovery.targetType).
			Str("icon", discovery.iconKind).
			Msgf("Registered fake target type with %d targets", config.Config.FakeTargetsPerType)
	}
	log.Info().Msgf("Registered %d fake target types", config.Config.FakeTargetTypeCount)
}

// fakeTargetTypeAction builds the no-op action of the index-th (0-based) fake
// target type. Its label names the type, since that is all the action picker shows.
func fakeTargetTypeAction(index int) action_kit_sdk.Action[DoNothingActionState] {
	spec := fakeTargetTypeSpecAt(index)
	return NewDoNothingActionWithLabel(
		fakeTargetTypeIdPrefix+spec.key,
		action_kit_api.TargetSelectionTemplate{
			Label:       "by name",
			Description: new(fmt.Sprintf("Find %s by name.", spec.one)),
			Query:       fmt.Sprintf("loadtest.%s.name=\"\"", spec.key),
		},
		fmt.Sprintf("Do Nothing (%s)", spec.one),
	)
}

// RegisterFakeTargetTypeActions registers one no-op action per fake target type.
// It loops over the same count as RegisterFakeTargetTypeDiscoveries, so a
// registered type always has an action to reach it with.
func RegisterFakeTargetTypeActions() {
	for i := 0; i < config.Config.FakeTargetTypeCount; i++ {
		action_kit_sdk.RegisterAction(fakeTargetTypeAction(i))
	}
}
