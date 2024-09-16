# Steadybit extension-loadtest

Internal extension used to load test agent & platform.

## Configuration

| Environment Variable              | Helm value | Meaning                                     | Required | Default                 |
|-----------------------------------|------------|---------------------------------------------|----------|-------------------------|

The extension supports all environment variables provided by [steadybit/extension-kit](https://github.com/steadybit/extension-kit#environment-variables).

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
the [documentation](https://docs.steadybit.com/install-and-configure/install-agent/extension-discovery) for more
information about extension registration and how to verify.
