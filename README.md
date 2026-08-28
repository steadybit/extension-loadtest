# Steadybit extension-loadtest

Internal extension used to load test agent & platform.

## Configuration

| Environment Variable                             | Helm value                     | Meaning                                                                                                                                   | Required | Default |
|--------------------------------------------------|--------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|----------|---------|
| `STEADYBIT_EXTENSION_FAKE_TARGET_TYPE_COUNT`     | `loadtest.fakeTargetTypeCount` | Number of pseudo target types to register. Set to `0` to register none at all.                                                                 | no       | `5`     |
| `STEADYBIT_EXTENSION_FAKE_TARGETS_PER_TYPE`      | `loadtest.fakeTargetsPerType`  | Number of targets discovered per pseudo target type.                                                                                          | no       | `10`    |

Beyond the settings above, this extension supports the configuration common to all Steadybit
extensions:

- [extension-kit](https://github.com/steadybit/extension-kit#environment-variables) — HTTP and
  health ports, TLS and mutual TLS, unix domain socket, logging, and pprof.
- [Target Filtering](https://github.com/steadybit/discovery-kit/blob/main/docs/target-filtering.md) —
  stop the extension reporting targets you do not want.
- [Group Matching](https://github.com/steadybit/discovery-kit/blob/main/docs/target-enrichment.md#group-matching) —
  tag discovered targets with a group, so enrichment rules only match within it.

## Pseudo target types

Every regular discovery of this extension borrows the target type id of a real extension
(`com.steadybit.extension_host.host`, `com.steadybit.extension_kubernetes.kubernetes-deployment`, ...)
and therefore inherits its description and icon. The pseudo target types are the exception: they are
owned by this extension, carry deliberately made-up names (`com.steadybit.extension_loadtest.flux-capacitor`,
`...hoverboard`, `...rubber-duck`, ...) so they cannot be mistaken for anything real, and have no actions.

They exist to exercise the platform and the UI with target type icons that are missing or unrenderable.
The registered types rotate through the three cases, so a count of `3` or more always covers all of them:

| Index | Icon                                              | `loadtest.icon` |
|-------|---------------------------------------------------|-----------------|
| 0     | none - the `icon` field is omitted entirely       | `none`          |
| 1     | broken - a data URI whose payload is not an image | `broken`        |
| 2     | valid - a plain star SVG                          | `valid`         |
| 3     | none again, and so on                             | `none`          |

To tell a rendered fallback apart from a working icon without having to guess, the icon case is visible
in two places:

- every target carries the `loadtest.icon` attribute. It is shared by all pseudo target types (the
  `loadtest.<type>.name` and `loadtest.<type>.serial` attributes are per type), so the explorer can group
  by it and put the three cases side by side. It is also the second column of every type's target table.
- the target type label names it, e.g. *Flux Capacitor (none icon)*, *Hoverboard (broken icon)*. A target
  type carries no attributes of its own, so the label is the only place this can show where the type
  itself is rendered.

```bash
STEADYBIT_EXTENSION_FAKE_TARGET_TYPE_COUNT=5 \
STEADYBIT_EXTENSION_FAKE_TARGETS_PER_TYPE=10 \
./extension-loadtest
```

## Installation

### Kubernetes

Detailed information about agent and extension installation in kubernetes can also be found in
our [documentation](https://docs.steadybit.com/install-and-configure/install-agent/install-on-kubernetes).

#### Recommended (via agent helm chart)

All extensions provide a helm chart that is also integrated in the
[helm-chart](https://github.com/steadybit/helm-charts/tree/main/charts/steadybit-agent) of the agent.

You must provide additional values to activate this extension.

```
--set extension-loadtest.enabled=true \
```

Additional configuration options can be found in
the [helm-chart](https://github.com/steadybit/extension-loadtest/blob/main/charts/steadybit-extension-loadtest/values.yaml) of the
extension.

#### Alternative (via own helm chart)

If you need more control, you can install the extension via its
dedicated [helm-chart](https://github.com/steadybit/extension-loadtest/blob/main/charts/steadybit-extension-loadtest).

```bash
helm repo add steadybit-extension-loadtest https://steadybit.github.io/extension-loadtest
helm repo update
helm upgrade steadybit-extension-loadtest \
    --install \
    --wait \
    --timeout 5m0s \
    --create-namespace \
    --namespace steadybit-agent \
    steadybit-extension-loadtest/steadybit-extension-loadtest
```

### Linux Package

This extension is currently not available as a Linux package.

## Extension registration

Make sure that the extension is registered with the agent. In most cases this is done automatically. Please refer to
the [documentation](https://docs.steadybit.com/install-and-configure/install-agent/extension-registration) for more
information about extension registration and how to verify.

## Version and Revision

The version and revision of the extension:
- are printed during the startup of the extension
- are added as a Docker label to the image
- are available via the `version.txt`/`revision.txt` files in the root of the image
