// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package extloadtest

import (
	"testing"

	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogAction_ParameterGroups(t *testing.T) {
	desc := NewLogAction("com.steadybit.extension_loadtest.log.host", "com.steadybit.extension_host.host", action_kit_api.TargetSelectionTemplate{}).Describe()

	members := map[string][]string{}
	var ungroupedAdvanced []string
	for _, p := range desc.Parameters {
		if p.Group != nil {
			// Platforms without groups keep showing every grouped parameter in "Additional Settings".
			assert.True(t, p.Advanced != nil && *p.Advanced, "parameter %s", p.Name)
			members[*p.Group] = append(members[*p.Group], p.Name)
		} else if p.Advanced != nil && *p.Advanced {
			ungroupedAdvanced = append(ungroupedAdvanced, p.Name)
		}
	}

	assert.Equal(t, []string{"errorEndpoint", "latencyEndpoint", "latencyDuration", "targetFilter"}, members[groupFaultInjection])
	assert.Len(t, members[groupParameterTypes], 13)
	assert.Len(t, members[groupManyParameters], 15)
	assert.Len(t, members, 3)

	// The legacy fallback: an advanced parameter without a group, shown in "Additional Settings".
	assert.Equal(t, []string{"addArtifact"}, ungroupedAdvanced)

	require.NotNil(t, desc.ParameterGroups)
	var defined []string
	for _, g := range *desc.ParameterGroups {
		defined = append(defined, g.Name)
	}
	assert.Equal(t, []string{groupFaultInjection, groupParameterTypes}, defined, "the long group has no description")
}
