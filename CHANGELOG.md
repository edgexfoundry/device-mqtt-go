
<a name="EdgeX MQTT Device Service (found in device-mqtt-go) Changelog"></a>
## EdgeX MQTT Device Service
[Github repository](https://github.com/edgexfoundry/device-mqtt-go)

### Change Logs for EdgeX Dependencies
- [device-sdk-go](https://github.com/edgexfoundry/device-sdk-go/blob/main/CHANGELOG.md)
- [go-mod-core-contracts](https://github.com/edgexfoundry/go-mod-core-contracts/blob/main/CHANGELOG.md)
- [go-mod-bootstrap](https://github.com/edgexfoundry/go-mod-bootstrap/blob/main/CHANGELOG.md)
- [go-mod-messaging](https://github.com/edgexfoundry/go-mod-messaging/blob/main/CHANGELOG.md) (indirect dependency)
- [go-mod-registry](https://github.com/edgexfoundry/go-mod-registry/blob/main/CHANGELOG.md)  (indirect dependency)
- [go-mod-secrets](https://github.com/edgexfoundry/go-mod-secrets/blob/main/CHANGELOG.md) (indirect dependency)
- [go-mod-configuration](https://github.com/edgexfoundry/go-mod-configuration/blob/main/CHANGELOG.md) (indirect dependency)

## [4.0.0] Odessa - 2025-03-12 (Only compatible with the 4.x releases)

### ✨  Features

- Enable PIE support for ASLR and full RELRO ([e0a7d63…](https://github.com/edgexfoundry/device-mqtt-go/commit/e0a7d6305e2fc591ccde04a15c4d9c78573e62c9))


### ♻ Code Refactoring

- Update module to v4 ([b35bedd…](https://github.com/edgexfoundry/device-mqtt-go/commit/b35bedd1f3478e6ee9832909a33ee3e010ee194a))
```text

BREAKING CHANGE: update go module to v4

```

### 🐛 Bug Fixes

- Only one ldflags flag is allowed ([cc886f0…](https://github.com/edgexfoundry/device-mqtt-go/commit/cc886f0163229c0230891ab3c0de8aec2f7a231b))

### 👷 Build

- Upgrade to go-1.23, Linter1.61.0 and Alpine 3.20 ([43f794b…](https://github.com/edgexfoundry/device-mqtt-go/commit/43f794b456aa3735bbca56bd2373d27e91a7aa37))


## [v3.1.0] Napa - 2023-11-15 (Only compatible with the 3.x releases)


### ✨  Features

- Remove snap packaging ([#635](https://github.com/edgexfoundry/device-mqtt-go/issues/635)) ([aecec81…](https://github.com/edgexfoundry/device-mqtt-go/commit/aecec81b0df1ccac07961d4c177179e1cb746efd))
```text

BREAKING CHANGE: Remove snap packaging ([#635](https://github.com/edgexfoundry/device-mqtt-go/issues/635))

```


### ♻ Code Refactoring

- Remove github.com/pkg/errors from Attribution.txt ([a383cfd…](https://github.com/edgexfoundry/device-mqtt-go/commit/a383cfdc3a434418355ebd57346991eaea66a3fb))


### 👷 Build

- Upgrade to go-1.21, Linter1.54.2 and Alpine 3.18 ([#624](https://github.com/edgexfoundry/device-mqtt-go/issues/624)) ([9edf56c…](https://github.com/edgexfoundry/device-mqtt-go/commit/9edf56c898603a6da7d069e5015abb4f1d1f30d1))


### 🤖 Continuous Integration

- Add automated release workflow on tag creation ([7b36b08…](https://github.com/edgexfoundry/device-mqtt-go/commit/7b36b08afcb678adb2f3795fb5abdedf4fd6d0d9))


## [v3.0.0] Minnesota - 2023-05-31 (Only compatible with the 3.x releases)

### Features ✨

- Handle receiving multiple readings in the async payload ([#]())
  ```text
  BREAKING CHANGE: The published async payload must now be json containing the resource names and values, even when a single resource is sent. 
                   - Redundant single resource Device Commands (`testrandfloat32`, `testrandfloat64`, `testping` and `testmessage`) in `mqtt.test.device.profile.yml` have be removed.
                     Use the corresponding resource name (`randfloat32`, `randfloat64`, `ping` or `message`) instead.
  ```
- Make multi-level topics the only implementation ([#15e121](https://github.com/edgexfoundry/device-mqtt-go/commit/15e121638b61dd5d2b05931572372c8312ebd12f))
  ```text
  BREAKING CHANGE: Removed single level topics so that only multi level is used 
  ```
- Consume SDK interface changes ([#24b4e6c8](https://github.com/edgexfoundry/device-mqtt-go/commit/24b4e6c8c5dfa1689a2429f2d6b13ded7328be9a))
  ```text
  BREAKING CHANGE: Consume SDK interface changes by adding Discover and ValidateDevice func on driver
  ```
- Update profile Minimum and Maximum to numeric data ([#5b00e40](https://github.com/edgexfoundry/device-mqtt-go/commit/5b00e406319c2ed6c9a36a296f26911c42639135))
  ```text
  BREAKING CHANGE: Update profile Minimum and Maximum to numeric data according to the latest core-contracts changes.
  ```
- Change configuration and device TOML files to YAML format ([#3944987](https://github.com/edgexfoundry/device-mqtt-go/commit/39449871d1528794bb93d06fd28712f06c9b0a27))
  ```text
  BREAKING CHANGE: Configuration and device files now in YAML format
  ```
- Updates for common config ([#ec4bdb7](https://github.com/edgexfoundry/device-mqtt-go/commit/ec4bdb71d113baf8f62d105c9dc125da7e597352))
  ```text
  BREAKING CHANGE: Configuration file changed to remove common config settings
  ```
- Remove ZeroMQ MessageBus capability ([#]())
  ```text
  BREAKING CHANGE: ZeroMQ no longer option for the EdgeX MessageBus
  ```

### Bug Fixes 🐛

- Update logging to avoid leaking messageBus credentials ([#51ee6b6](https://github.com/edgexfoundry/device-mqtt-go/commits/51ee6b6))
- **snap:** Refactor to avoid conflicts with readonly config provider directory ([#535](https://github.com/edgexfoundry/device-mqtt-go/issues/535)) ([#c80097c](https://github.com/edgexfoundry/device-mqtt-go/commits/c80097c))

### Code Refactoring ♻

- Updated secret path to secretName ([#dc80312a](https://github.com/edgexfoundry/device-mqtt-go/commit/dc80312a46813471bd7ca781724042db8dc8e361))
  ```text
  BREAKING CHANGE: `path` in Secret DTO renamed to `secretName`
  ```
- Consume device-sdk-go breaking changes ([#f6d0510](https://github.com/edgexfoundry/device-mqtt-go/commit/f6d05106c5497cf4be6abc7fa7efd102b354d07f))
  ```text
  BREAKING CHANGE: update ProtocolDriver implementation for the new ProtocolDriver interface changes
  ```
- Replace internal topics from config with new constants ([#90bb521](https://github.com/edgexfoundry/device-mqtt-go/commit/90bb521e5152a28ee665954924b952139261dddb))
  ```text
  BREAKING CHANGE: Internal topics no longer configurable, except the base topic.
  ```
- Rework code for refactored MessageBus Configuration ([#5c1f2ce](https://github.com/edgexfoundry/device-mqtt-go/commit/5c1f2ce46133253723bb2df5ba99f3b69c76a38a))
  ```text
  BREAKING CHANGE: MessageQueue renamed to MessageBus and fields changed. See v3 Migration guide.
  ```
- Use latest SDK for flattened config stem ([#6bdfaee](https://github.com/edgexfoundry/device-mqtt-go/commit/6bdfaee9ae606a4f8b392396fcc4c15889d23cee))
  ```text
  BREAKING CHANGE: Location of service configuration in Consul changed to edgex/v3/device-mqtt
  ```
- **snap:** Drop the support for legacy snap env options ([#bc7652b](https://github.com/edgexfoundry/device-mqtt-go/commit/bc7652beee9a90296fe0b57947c4589a421f5e02))
  ```text
  BREAKING CHANGE: Drop the support for deprecated snap options starting with `env.`
  ```
- **snap:** Update command and metadata sourcing ([#529](https://github.com/edgexfoundry/device-mqtt-go/issues/529)) ([#1f5b405](https://github.com/edgexfoundry/device-mqtt-go/commits/1f5b405))

### Build 👷

- Update to Go 1.20, Alpine 3.17 and linter v1.51.2 ([#acf7416](https://github.com/edgexfoundry/device-mqtt-go/commits/acf7416))

## [v2.3.0] Levski - 2022-11-09 (Only compatible with the 2.x releases)

### Features ✨

- Add Service Metrics configuration ([#431](https://github.com/edgexfoundry/device-mqtt-go/issues/431)) ([#b192367](https://github.com/edgexfoundry/device-mqtt-go/commits/b192367))
- Add NATS configuration and build option ([#414](https://github.com/edgexfoundry/device-mqtt-go/issues/414)) ([#626368a](https://github.com/edgexfoundry/device-mqtt-go/commits/626368a))
- Add commanding via message configuration ([#b1e10e2](https://github.com/edgexfoundry/device-mqtt-go/commits/b1e10e2))
- Add go-winio to attribution (new SPIFFE dependency) ([#383](https://github.com/edgexfoundry/device-mqtt-go/issues/383)) ([#7519713](https://github.com/edgexfoundry/device-mqtt-go/commits/7519713))
- **snap:** add config interface with unique identifier ([#6b3f84b](https://github.com/edgexfoundry/device-mqtt-go/commits/6b3f84b))

### Bug Fixes 🐛

- replace broker host config with a valid example ([#396](https://github.com/edgexfoundry/device-mqtt-go/issues/396)) ([#2c395e9](https://github.com/edgexfoundry/device-mqtt-go/commits/2c395e9))
- **snap:** Remove duplicate file copying in install hook ([#420](https://github.com/edgexfoundry/device-mqtt-go/issues/420)) ([#1a028d9](https://github.com/edgexfoundry/device-mqtt-go/commits/1a028d9))

### Code Refactoring ♻

- **snap:** edgex-snap-hooks related upgrade ([#397](https://github.com/edgexfoundry/device-mqtt-go/issues/397)) ([#8b7bfd0](https://github.com/edgexfoundry/device-mqtt-go/commits/8b7bfd0))

### Build 👷

- Upgrade to Go 1.18, fix linter issues and optimize attribution script ([#403](https://github.com/edgexfoundry/device-mqtt-go/issues/403)) ([#7e9452b](https://github.com/edgexfoundry/device-mqtt-go/commits/7e9452b))

## [v2.2.0] Kamakura - 2022-05-11  (Only compatible with the 2.x releases)

### Features ✨
- Enable security hardening ([#ec9d6ac](https://github.com/edgexfoundry/device-mqtt-go/commits/ec9d6ac))

### Bug Fixes 🐛
- **snap:** Expose parent directory in device-config plug ([#1a2d89f](https://github.com/edgexfoundry/device-mqtt-go/commits/1a2d89f))

### Code Refactoring ♻
- **snap:** Remove redundant content identifier ([#302ebdb](https://github.com/edgexfoundry/device-mqtt-go/commits/302ebdb))

### Documentation 📖
- Add deployment section to readme ([#330](https://github.com/edgexfoundry/device-mqtt-go/issues/330)) ([#7f14e9a](https://github.com/edgexfoundry/device-mqtt-go/commits/7f14e9a))

### Build 👷
- Update to latest SDK w/o ZMQ on windows ([#0aab54c](https://github.com/edgexfoundry/device-mqtt-go/commits/0aab54c))
    ```
    BREAKING CHANGE:
    ZeroMQ no longer supported on native Windows for EdgeX
    MessageBus
    ```
- **snap:** source metadata from repo ([#328](https://github.com/edgexfoundry/device-mqtt-go/issues/328)) ([#195658a](https://github.com/edgexfoundry/device-mqtt-go/commits/195658a))

### Continuous Integration 🔄
- gomod changes related for Go 1.17 ([#335](https://github.com/edgexfoundry/device-mqtt-go/issues/335)) ([#6506aa6](https://github.com/edgexfoundry/device-mqtt-go/commits/6506aa6))
- Go 1.17 related changes ([#334](https://github.com/edgexfoundry/device-mqtt-go/issues/334)) ([#620914e](https://github.com/edgexfoundry/device-mqtt-go/commits/620914e))

## [v2.1.0] Jakarta - 2021-11-18  (Only compatible with the 2.x releases)

### Features ✨
- Add support for new Object reading type ([#323](https://github.com/edgexfoundry/device-mqtt-go/issues/323)) ([#5b76d0d](https://github.com/edgexfoundry/device-mqtt-go/commits/5b76d0d))
- Update config file to include secrets file parameters ([#310](https://github.com/edgexfoundry/device-mqtt-go/issues/310)) ([#c11b57c](https://github.com/edgexfoundry/device-mqtt-go/commits/c11b57c))
- Add UseMqttLevels capability to pass in metadata ([#303](https://github.com/edgexfoundry/device-mqtt-go/issues/303)) ([#7270bca](https://github.com/edgexfoundry/device-mqtt-go/commits/7270bca))
- **snap:** bump edgex-snap-hooks to support secretsfile config ([#a36a4aa](https://github.com/edgexfoundry/device-mqtt-go/commits/a36a4aa))

### Bug Fixes 🐛
- Update all TOML to use quote and not single-quote ([#298](https://github.com/edgexfoundry/device-mqtt-go/issues/298)) ([#255d6e1](https://github.com/edgexfoundry/device-mqtt-go/commits/255d6e1))

### Code Refactoring ♻

### Documentation 📖
- Update build status badge ([#56d1e45](https://github.com/edgexfoundry/device-mqtt-go/commits/56d1e45))
- **snap:** Fix markdown and replication errors ([#8907cac](https://github.com/edgexfoundry/device-mqtt-go/commits/8907cac))

### Build 👷
- Update to use released SDK ([#322](https://github.com/edgexfoundry/device-mqtt-go/issues/322)) ([#ea5b921](https://github.com/edgexfoundry/device-mqtt-go/commits/ea5b921))
- Update alpine base to 3.14 ([#304](https://github.com/edgexfoundry/device-mqtt-go/issues/304)) ([#9a27036](https://github.com/edgexfoundry/device-mqtt-go/commits/9a27036))
- **snap:** Upgrade to snap core20 ([#2f4bb9e](https://github.com/edgexfoundry/device-mqtt-go/commits/2f4bb9e))

### Continuous Integration 🔄
- Remove need for CI specific Dockerfile ([#1a605e9](https://github.com/edgexfoundry/device-mqtt-go/commits/1a605e9))

## [v2.0.0] Ireland - 2021-06-30  (Only compatible with the 2.x releases)

### Features ✨
- Fix onConnectHandler panics and update config file ([#288](https://github.com/edgexfoundry/device-mqtt-go/pull/288))
- Using single MQTT broker config ([#277](https://github.com/edgexfoundry/device-mqtt-go/issues/277)) ([#056bd70](https://github.com/edgexfoundry/device-mqtt-go/commits/056bd70))
- Enable using MessageBus as the default ([#279](https://github.com/edgexfoundry/device-mqtt-go/issues/279)) ([#f18a6a3](https://github.com/edgexfoundry/device-mqtt-go/commits/f18a6a3))
- Extract the command response retry interval as configuration ([#62ff07c](https://github.com/edgexfoundry/device-mqtt-go/commits/62ff07c))
- Move Driver config to new custom config section ([#5b2c07b](https://github.com/edgexfoundry/device-mqtt-go/commits/5b2c07b))
- Add secure MessagBus capability ([#696b33d](https://github.com/edgexfoundry/device-mqtt-go/commits/696b33d))
- Remove Logging configuration ([#f1a7c6f](https://github.com/edgexfoundry/device-mqtt-go/commits/f1a7c6f))
- Updated Dockerfile to install dumb-init ([#bc66537](https://github.com/edgexfoundry/device-mqtt-go/commits/bc66537))
- Enable use of secret via SecretProvider for MQTT broker credentials ([#33a7955](https://github.com/edgexfoundry/device-mqtt-go/commits/33a7955))

### Bug Fixes 🐛
- Change "."s in profile name to "-"s ([#284](https://github.com/edgexfoundry/device-mqtt-go/issues/284)) ([#8213f84](https://github.com/edgexfoundry/device-mqtt-go/commits/8213f84))
- Add AuthMode settings so have ability to enable/disable Auth MQTT connections ([#269](https://github.com/edgexfoundry/device-mqtt-go/issues/269)) ([#9a33ad5](https://github.com/edgexfoundry/device-mqtt-go/commits/9a33ad5))
- Add Type='vault' to [SecretStore] config ([#7c58968](https://github.com/edgexfoundry/device-mqtt-go/commits/7c58968))
- Corrected port numbers per PR comments ([#dbf9134](https://github.com/edgexfoundry/device-mqtt-go/commits/dbf9134))
- Added missing InsecureSecrets Section and UseMessageBus = false ([#ed2040e](https://github.com/edgexfoundry/device-mqtt-go/commits/ed2040e))
### Code Refactoring ♻
- Change PublishTopicPrefix value to be 'edgex/events/device' ([#3890446](https://github.com/edgexfoundry/device-mqtt-go/commits/3890446))
- Rename the custom config name to MQTTBrokerInfo ([#d8fe7de](https://github.com/edgexfoundry/device-mqtt-go/commits/d8fe7de))
- Update configuration for change to common ServiceInfo struct ([#7ed00ab](https://github.com/edgexfoundry/device-mqtt-go/commits/7ed00ab))
    ```
    BREAKING CHANGE:
    Service configuration has changed
    ```
- Update to assign and uses new Port Assignments ([#9e27054](https://github.com/edgexfoundry/device-mqtt-go/commits/9e27054))
    ```
    BREAKING CHANGE:
    Device MQTT default port number has changed to 59982
    ```
- rename example device AutoEvent Fequency to Interval ([#3a738e3](https://github.com/edgexfoundry/device-mqtt-go/commits/3a738e3))
- Added go mod tidy to dockerfile ([#5919639](https://github.com/edgexfoundry/device-mqtt-go/commits/5919639))
- Update for new service key names and overrides for hyphen to underscore ([#356f292](https://github.com/edgexfoundry/device-mqtt-go/commits/356f292))
    ```
    BREAKING CHANGE:
    Service key names used in configuration have changed.
    ```
- use v2 device-sdk ([#5a126a9](https://github.com/edgexfoundry/device-mqtt-go/commits/5a126a9))
### Documentation 📖
- Add badges to readme ([#cfb712f](https://github.com/edgexfoundry/device-mqtt-go/commits/cfb712f))
### Build 👷
- update build files for zmq dependency ([#d53328a](https://github.com/edgexfoundry/device-mqtt-go/commits/d53328a))
- **deps:** bump github.com/eclipse/paho.mqtt.golang ([#788356c](https://github.com/edgexfoundry/device-mqtt-go/commits/788356c))
- **deps:** bump github.com/stretchr/testify from 1.5.1 to 1.7.0 ([#5dc0bc9](https://github.com/edgexfoundry/device-mqtt-go/commits/5dc0bc9))
- update Dockerfiles to use go 1.16 ([#cc189d3](https://github.com/edgexfoundry/device-mqtt-go/commits/cc189d3))
- update go.mod to go 1.16 ([#df72406](https://github.com/edgexfoundry/device-mqtt-go/commits/df72406))
- **snap:** update go to 1.16 ([#941ce85](https://github.com/edgexfoundry/device-mqtt-go/commits/941ce85))
- **snap:** update snap v2 support ([#59017e6](https://github.com/edgexfoundry/device-mqtt-go/commits/59017e6))
### Continuous Integration 🔄
- update local docker image names ([#06b6566](https://github.com/edgexfoundry/device-mqtt-go/commits/06b6566))

<a name="v1.3.1"></a>
## [v1.3.1] - 2021-02-02
### Features ✨
- **snap:** add startup-duration and startup-interval configure options ([#bad7e1b](https://github.com/edgexfoundry/device-mqtt-go/commits/bad7e1b))
### Build 👷
- **deps:** bump github.com/edgexfoundry/device-sdk-go ([#a154119](https://github.com/edgexfoundry/device-mqtt-go/commits/a154119))
### Continuous Integration 🔄
- add semantic.yml for commit linting, update PR template to latest ([#692e0b5](https://github.com/edgexfoundry/device-mqtt-go/commits/692e0b5))
- standardize dockerfiles ([#43e9764](https://github.com/edgexfoundry/device-mqtt-go/commits/43e9764))

<a name="v1.3.0"></a>
## [v1.3.0] - 2020-11-18
### Bug Fixes 🐛
- Return error instead of the panic if required config not found ([#8630507](https://github.com/edgexfoundry/device-mqtt-go/commits/8630507))
- Modify float value checking condition ([#2a661a3](https://github.com/edgexfoundry/device-mqtt-go/commits/2a661a3))
- local snap development ([#8bc9dbb](https://github.com/edgexfoundry/device-mqtt-go/commits/8bc9dbb))
### Code Refactoring ♻
- Upgrade SDK to v1.2.4-dev.34 ([#fe9eb72](https://github.com/edgexfoundry/device-mqtt-go/commits/fe9eb72))
- update dockerfile to appropriately use ENTRYPOINT and CMD, closes[#164](https://github.com/edgexfoundry/device-mqtt-go/issues/164) ([#d7447a9](https://github.com/edgexfoundry/device-mqtt-go/commits/d7447a9))
### Build 👷
- Upgrade to Go1.15 ([#b7208c3](https://github.com/edgexfoundry/device-mqtt-go/commits/b7208c3))
- add dependabot.yml ([#730afc1](https://github.com/edgexfoundry/device-mqtt-go/commits/730afc1))
- **deps:** bump github.com/edgexfoundry/device-sdk-go ([#61125b8](https://github.com/edgexfoundry/device-mqtt-go/commits/61125b8))

<a name="v1.2.2"></a>
## [v1.2.2] - 2020-08-19
### Snap
- add env override for ProfilesDir ([#48947cc](https://github.com/edgexfoundry/device-mqtt-go/commits/48947cc))
### Bug Fixes 🐛
- Optimize MQTT client creation for async value Add conn retry mechanism and use os.Exit(1) instead of the panic error ([#7f88040](https://github.com/edgexfoundry/device-mqtt-go/commits/7f88040))
### Code Refactoring ♻
- upgrade SDK to v1.2.0 ([#863c2d0](https://github.com/edgexfoundry/device-mqtt-go/commits/863c2d0))
### Documentation 📖
- Add standard PR template ([#50e33cb](https://github.com/edgexfoundry/device-mqtt-go/commits/50e33cb))
