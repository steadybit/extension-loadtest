# Changelog

## v1.0.12

- build(deps): bump github.com/steadybit/preflight-kit/go/preflight_kit_api
- chore(chart): bump chart version to 1.0.65
- chore(deps): bump extensionlib to ^1.5.5
- chore(deps): bump golang.org/x/net to v0.61.0
- chore: remove action that causes a warning
- feat(chart): expose the pseudo target type settings
- feat: add a no-op action per pseudo target type
- feat: add pseudo target types with rotating icon cases
- feat: categorize the actions by technology and category
- feat: initialize OpenTelemetry tracing (#119)
- fix: make the recreate actions actually bump the target id (the previous in-place update was discarded, so recreate was a silent no-op)
- fix: synchronize discovery reads with the `/config/*` updates and the recreate actions so concurrent requests no longer race on the configuration and target data

## v1.0.10

- Support discovery group attribute via `STEADYBIT_EXTENSION_DISCOVERY_GROUP` env var (or `discovery.group` Helm value) — when set, the extension adds `steadybit.group=<value>` to every discovered target
- Update dependencies

## v1.0.9

- Bump Go to 1.26.3
- Update dependencies

## v1.0.8

- Bump Go to 1.25.9
- Support if-none-match for the extension list endpoint
- Update dependencies

## v1.0.7

- feat(chart): split image.name into image.registry + image.name
- Support global.priorityClassName
- Allow fixed value for poduid
- Add label with index to pods and deployment
- Update alpine packages in Docker image to address CVEs
- Update dependencies

## v1.0.6

- Add action with option parameter based on multiple target attributes
- Update dependencies

## v1.0.5

- Update dependencies

## v1.0.4

- Updated dependencies

## v1.0.3

- Added linux package

## v1.0.1

- go version 1.23

## v1.0.0

 - Initial release