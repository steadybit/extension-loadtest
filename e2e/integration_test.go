// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 Steadybit GmbH

package e2e

import (
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_test/e2e"
	"github.com/steadybit/discovery-kit/go/discovery_kit_test/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithMinikube(t *testing.T) {
	extFactory := e2e.HelmExtensionFactory{
		Name: "extension-scaffold",
		Port: 8080,
		ExtraArgs: func(m *e2e.Minikube) []string {
			return []string{"--set", "logging.level=debug"}
		},
	}

	e2e.WithDefaultMinikube(t, &extFactory, []e2e.WithMinikubeTestCase{
		{
			Name: "validate discovery",
			Test: validateDiscovery,
		},
		{
			Name: "run scaffold",
			Test: testRunscaffold,
		},
	})
}

func validateDiscovery(t *testing.T, _ *e2e.Minikube, e *e2e.Extension) {
	assert.NoError(t, validate.ValidateEndpointReferences("/", e.Client))
}

func testRunscaffold(t *testing.T, m *e2e.Minikube, e *e2e.Extension) {
	config := struct{}{}
	exec, err := e.RunAction("com.steadybit.extension_scaffold.robot.log", &action_kit_api.Target{
		Name: "robot",
	}, config, nil)
	require.NoError(t, err)
	e2e.AssertLogContains(t, m, e.Pod, "Logging in log action **start**")
	require.NoError(t, exec.Cancel())
}
