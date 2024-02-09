# Steadybit extension-loadtest

Internal extension used to load test agent & platform.

## Configuration

| Environment Variable              | Helm value | Meaning                                     | Required | Default                 |
|-----------------------------------|------------|---------------------------------------------|----------|-------------------------|

The extension supports all environment variables provided by [steadybit/extension-kit](https://github.com/steadybit/extension-kit#environment-variables).

## Installation

### Using Docker

```sh
docker run \
  --rm \
  -p 8080 \
  --name steadybit-extension-loadtest \
  ghcr.io/steadybit/extension-loadtest:latest
```

### Using Helm in Kubernetes

```sh
helm repo add steadybit-extension-loadtest https://steadybit.github.io/extension-loadtest
helm repo update
helm upgrade steadybit-extension-loadtest \
    --install \
    --wait \
    --timeout 5m0s \
    --create-namespace \
    --namespace steadybit-extension \
    steadybit-extension-loadtest/steadybit-extension-loadtest
```

## Register the extension

Make sure to register the extension at the steadybit platform. Please refer to
the [documentation](https://docs.steadybit.com/integrate-with-steadybit/extensions/extension-installation) for more information.
