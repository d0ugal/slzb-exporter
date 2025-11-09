# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.13.5](https://github.com/d0ugal/slzb-exporter/compare/v2.13.4...v2.13.5) (2025-11-09)


### Bug Fixes

* Update module github.com/d0ugal/promexporter to v1.11.2 ([5dfb2cb](https://github.com/d0ugal/slzb-exporter/commit/5dfb2cb7f29b223a7655b7c16c6f7c3ab1dfe96f))
* Update module github.com/quic-go/quic-go to v0.56.0 ([458cd0c](https://github.com/d0ugal/slzb-exporter/commit/458cd0ccd0f0ea2990cb63ff15cceb31fafb0d0c))
* Update module golang.org/x/arch to v0.23.0 ([b5a8299](https://github.com/d0ugal/slzb-exporter/commit/b5a82990e8416fccb2c2261843e7933856097abc))
* Update module golang.org/x/sys to v0.38.0 ([46c1372](https://github.com/d0ugal/slzb-exporter/commit/46c1372229cae6631dc2487a40ba732d87a09eac))

## [2.13.4](https://github.com/d0ugal/slzb-exporter/compare/v2.13.3...v2.13.4) (2025-11-06)


### Bug Fixes

* **ci:** fix syntax errors in check-build-needed job ([9ad7411](https://github.com/d0ugal/slzb-exporter/commit/9ad74118bcfc0ed8806689734daebf387d5d6c29))
* Update module github.com/d0ugal/promexporter to v1.11.1 ([0ef6286](https://github.com/d0ugal/slzb-exporter/commit/0ef628692f3115cd9c7817b61f98719182f11d04))

## [2.13.3](https://github.com/d0ugal/slzb-exporter/compare/v2.13.2...v2.13.3) (2025-11-05)


### Bug Fixes

* Update dependency go to v1.25.4 ([707d161](https://github.com/d0ugal/slzb-exporter/commit/707d1612742fa339928f66911b8a3c17374bb6e8))

## [2.13.2](https://github.com/d0ugal/slzb-exporter/compare/v2.13.1...v2.13.2) (2025-11-05)


### Bug Fixes

* Update module github.com/d0ugal/promexporter to v1.11.0 ([48ea144](https://github.com/d0ugal/slzb-exporter/commit/48ea1444aedb1f620c334169b7ebc71e34eab15a))

## [2.13.1](https://github.com/d0ugal/slzb-exporter/compare/v2.13.0...v2.13.1) (2025-11-04)


### Bug Fixes

* Update google.golang.org/genproto/googleapis/api digest to f26f940 ([eb5d597](https://github.com/d0ugal/slzb-exporter/commit/eb5d59772d859bb58172f360ccfd1767777dca6e))
* Update google.golang.org/genproto/googleapis/rpc digest to f26f940 ([6485e96](https://github.com/d0ugal/slzb-exporter/commit/6485e966c4efabfacea7ba1cafa4003e649bbc80))
* Update module github.com/d0ugal/promexporter to v1.9.0 ([669ebb8](https://github.com/d0ugal/slzb-exporter/commit/669ebb84c054a66cdfa6081362dde637e848018c))
* Update module go.opentelemetry.io/proto/otlp to v1.9.0 ([a4f6ae9](https://github.com/d0ugal/slzb-exporter/commit/a4f6ae941c148cdbe312d051669b9d1ca53ba566))

## [2.13.0](https://github.com/d0ugal/slzb-exporter/compare/v2.12.2...v2.13.0) (2025-11-02)


### Features

* add dev-tag Makefile target and update CI workflow ([7f81876](https://github.com/d0ugal/slzb-exporter/commit/7f8187661ea5ae12665a937ecd739e38dfbcbe4d))
* add duplication linter (dupl) to golangci configuration ([67daec6](https://github.com/d0ugal/slzb-exporter/commit/67daec69795577d32e1fc80cb9a27b120b7ef0dc))
* add OpenTelemetry HTTP client instrumentation ([931263e](https://github.com/d0ugal/slzb-exporter/commit/931263e4993b650f5b8ac7393c88f1ae3348606c))
* add tracing configuration support ([f499bb3](https://github.com/d0ugal/slzb-exporter/commit/f499bb347570b423b581752916e51c3825bf5f5d))
* **ci:** add auto-format workflow ([940cd11](https://github.com/d0ugal/slzb-exporter/commit/940cd11f455d1ce0e844411398af7f8b17f5c671))
* enhance tracing support with detailed spans ([befc016](https://github.com/d0ugal/slzb-exporter/commit/befc0165bff2dc17678e8cdc6e80b43210481d64))
* integrate OpenTelemetry tracing into collector ([80e7c7f](https://github.com/d0ugal/slzb-exporter/commit/80e7c7f8f6958771c2067d9f60e8e06e5e1ba8b0))
* trigger CI after auto-format workflow completes ([14029fe](https://github.com/d0ugal/slzb-exporter/commit/14029fea8c728486915f2598622e4fae5421eb4c))


### Bug Fixes

* correct label mapping for SLZB device metrics ([aa45595](https://github.com/d0ugal/slzb-exporter/commit/aa4559509ba78028028033ffe810ee782eca39f2))
* correct network metrics label mapping ([e54ebbb](https://github.com/d0ugal/slzb-exporter/commit/e54ebbb25734efce125b6a99a04f8ada9bad0c50))
* correct remaining device label mappings ([076de00](https://github.com/d0ugal/slzb-exporter/commit/076de004563f70955dbf3c4292f5edeb1c092619))
* correct remaining label mapping issues in SLZB collector ([a7414d1](https://github.com/d0ugal/slzb-exporter/commit/a7414d1f043c91561453d8632e88435eb802a535))
* correct SLZBFirmwareCurrentVersion label mapping ([17dedae](https://github.com/d0ugal/slzb-exporter/commit/17dedae6b51c5363b1b5e13caae74b575d1fc13b))
* correct SLZBFirmwareUpdateAvailable label mapping ([d3243b1](https://github.com/d0ugal/slzb-exporter/commit/d3243b13ce1a81e3a28514f52a07b0e88334380c))
* correct SLZBHTTPRequestsTotal label mapping ([3826ef2](https://github.com/d0ugal/slzb-exporter/commit/3826ef2240809bc8b5eeaab00dc8e8eea10b530f))
* lint ([6d57837](https://github.com/d0ugal/slzb-exporter/commit/6d57837ea982ab11914e0e29ad96a54a70446c25))
* lint ([0cd48a2](https://github.com/d0ugal/slzb-exporter/commit/0cd48a21fc15db9420deef96bfc0c32c94cc99f0))
* remove unused spanCtx variable in collectConfigurationMetrics ([8006a92](https://github.com/d0ugal/slzb-exporter/commit/8006a92c690b93dd40380b74b01e338b9c6cda7a))
* remove unused spanCtx variables from processDeviceData and collectFirmwareStatus ([0f2e390](https://github.com/d0ugal/slzb-exporter/commit/0f2e3908a10d94403c89cedeff2bc2ac6939840e))
* resolve contextcheck and whitespace linting issues ([28d2dbc](https://github.com/d0ugal/slzb-exporter/commit/28d2dbcadfb7ea58077252377ff4f779e67d6c26))
* Update google.golang.org/genproto/googleapis/api digest to ab9386a ([e46750c](https://github.com/d0ugal/slzb-exporter/commit/e46750c5e18138f6b2b1284786b7e11d1009400f))
* Update google.golang.org/genproto/googleapis/rpc digest to ab9386a ([6fa51df](https://github.com/d0ugal/slzb-exporter/commit/6fa51df65667ec102dd7b0fe2100280750fdad30))
* Update module github.com/bytedance/sonic to v1.14.2 ([05f71da](https://github.com/d0ugal/slzb-exporter/commit/05f71da4e4d090d17eac6b4fc6940bd335a45e74))
* Update module github.com/bytedance/sonic/loader to v0.4.0 ([cd6829f](https://github.com/d0ugal/slzb-exporter/commit/cd6829fdae4ee069bb2c9460d609e188df1872b5))
* Update module github.com/d0ugal/promexporter to v1.6.1 ([86c5c7f](https://github.com/d0ugal/slzb-exporter/commit/86c5c7f53b366cd566d05619e1e5cf2b5ed858ad))
* Update module github.com/d0ugal/promexporter to v1.7.1 ([d3e5fa4](https://github.com/d0ugal/slzb-exporter/commit/d3e5fa4d9debb048f5617eefd200d62d1860b7bc))
* Update module github.com/d0ugal/promexporter to v1.8.0 ([28d91a2](https://github.com/d0ugal/slzb-exporter/commit/28d91a2f9760b0556414550a0cc340851a4e9372))
* Update module github.com/gabriel-vasile/mimetype to v1.4.11 ([16b8871](https://github.com/d0ugal/slzb-exporter/commit/16b88717198e6c7adffb5a66dbcb49fd8a12b1bd))
* Update module github.com/prometheus/common to v0.67.2 ([2503713](https://github.com/d0ugal/slzb-exporter/commit/25037136a0e41b7059e5d8cdb2593cb37efd7e2c))
* Update module github.com/prometheus/procfs to v0.19.2 ([6bb2b5e](https://github.com/d0ugal/slzb-exporter/commit/6bb2b5e2c8f430389d2e2d38b0976de18db1024e))
* Update module github.com/ugorji/go/codec to v1.3.1 ([e25a33b](https://github.com/d0ugal/slzb-exporter/commit/e25a33b24677abf0257c4d943907389a77625be0))

## [2.12.2](https://github.com/d0ugal/slzb-exporter/compare/v2.12.1...v2.12.2) (2025-10-26)


### Bug Fixes

* add internal version package and update version handling ([6ecfc96](https://github.com/d0ugal/slzb-exporter/commit/6ecfc96fa53911d2815cc88d3d6f26c39f8ca408))
* Update module github.com/d0ugal/promexporter to v1.5.0 ([437d6f8](https://github.com/d0ugal/slzb-exporter/commit/437d6f8d99d46f9012cf15b5cda6cd4f133d3e7f))
* Update module github.com/prometheus/procfs to v0.19.1 ([4cd3346](https://github.com/d0ugal/slzb-exporter/commit/4cd3346f38fce8dfce6005de5428a16697d3136a))
* use WithVersionInfo to pass version info to promexporter ([5e9c04e](https://github.com/d0ugal/slzb-exporter/commit/5e9c04ed0645ee483840f34d82c6f6df2203cb61))

## [2.12.1](https://github.com/d0ugal/slzb-exporter/compare/v2.12.0...v2.12.1) (2025-10-25)


### Bug Fixes

* Update module github.com/d0ugal/promexporter to v1.4.1 ([54cf8ac](https://github.com/d0ugal/slzb-exporter/commit/54cf8acb17cf79587594a8fe4ca3379ccb1194f5))
* Update module github.com/prometheus/procfs to v0.19.0 ([a8196e3](https://github.com/d0ugal/slzb-exporter/commit/a8196e3f55b2dd64518d91296cea29780b50b639))

## [2.12.0](https://github.com/d0ugal/slzb-exporter/compare/v2.11.0...v2.12.0) (2025-10-25)


### Features

* update promexporter to v1.4.0 ([80e816f](https://github.com/d0ugal/slzb-exporter/commit/80e816f20c7bc1379b1813f10b58bbc0526f0ef6))


### Bug Fixes

* Update module github.com/d0ugal/promexporter to v1.1.0 ([4020f6b](https://github.com/d0ugal/slzb-exporter/commit/4020f6bb854e0d7ac8279dd51c89b81d1933e978))
* Update module github.com/d0ugal/promexporter to v1.3.1 ([4a5564a](https://github.com/d0ugal/slzb-exporter/commit/4a5564a9aca87099ba19b1410094d38b06857fd4))

## [2.11.0](https://github.com/d0ugal/slzb-exporter/compare/v2.10.1...v2.11.0) (2025-10-23)


### Features

* migrate slzb-exporter to promexporter library ([be9eac7](https://github.com/d0ugal/slzb-exporter/commit/be9eac7409e8e54523d845181b59b801d22043ea))
* update to promexporter v1.0.0 ([4ee43c7](https://github.com/d0ugal/slzb-exporter/commit/4ee43c70db414e6341d679b5eec2458784619ec2))


### Bug Fixes

* update config tests for promexporter v1.0.0 migration ([152f912](https://github.com/d0ugal/slzb-exporter/commit/152f912b25a5ab1efd7ed9fff9eaab9ac588066b))
* Update module github.com/d0ugal/promexporter to v1.0.1 ([ee0d64c](https://github.com/d0ugal/slzb-exporter/commit/ee0d64c4227d625c9e28cf3ef763347dfb37872c))
* Update module github.com/prometheus/procfs to v0.18.0 ([502167f](https://github.com/d0ugal/slzb-exporter/commit/502167f61d2d348d17c07c12b08ace5054d07bde))
* update to latest promexporter changes ([62e1231](https://github.com/d0ugal/slzb-exporter/commit/62e123167446f2bbc6c94ab2cefcdfb1e9b8761d))

## [2.10.1](https://github.com/d0ugal/slzb-exporter/compare/v2.10.0...v2.10.1) (2025-10-14)


### Bug Fixes

* Update module github.com/go-playground/validator/v10 to v10.28.0 ([a6479fa](https://github.com/d0ugal/slzb-exporter/commit/a6479fa17a12f8f6603262c48b9945efe944c92f))
* Update module golang.org/x/arch to v0.22.0 ([a070bf7](https://github.com/d0ugal/slzb-exporter/commit/a070bf7ff1cb691cbe4cb73da9022cb51204fc44))

## [2.10.0](https://github.com/d0ugal/slzb-exporter/compare/v2.9.1...v2.10.0) (2025-10-14)


### Features

* convert to Gin framework and set release mode unless debug logging ([03b9253](https://github.com/d0ugal/slzb-exporter/commit/03b92532f3e0e4641bfdcab46c7ddf6382342d90))


### Bug Fixes

* auto-fix import ordering and update dependencies ([d7f04af](https://github.com/d0ugal/slzb-exporter/commit/d7f04aff33a417dc570e7dd2966ad68d2bf18c34))
* correct import ordering and update dependencies ([901f1d5](https://github.com/d0ugal/slzb-exporter/commit/901f1d5a36212bc928b2de157b92a906ae64001e))
* Update module github.com/bytedance/sonic to v1.14.1 ([22c7fa7](https://github.com/d0ugal/slzb-exporter/commit/22c7fa761cea779b32f7cfb98d938d1f8c8a798c))
* Update module github.com/gabriel-vasile/mimetype to v1.4.10 ([659d52e](https://github.com/d0ugal/slzb-exporter/commit/659d52e3ba28f3d477aa11cbbe5ba93eeed2ca64))
* Update module github.com/goccy/go-json to v0.10.5 ([8ba0181](https://github.com/d0ugal/slzb-exporter/commit/8ba0181bf94d10c28354c43ed66a6b8e728af0ce))

## [2.9.1](https://github.com/d0ugal/slzb-exporter/compare/v2.9.0...v2.9.1) (2025-10-14)


### Bug Fixes

* Update dependency go to v1.25.3 ([461b87d](https://github.com/d0ugal/slzb-exporter/commit/461b87dc38f9613f76dc76bc02da878a5367984c))

## [2.9.0](https://github.com/d0ugal/slzb-exporter/compare/v2.8.0...v2.9.0) (2025-10-08)


### Features

* update dependencies to v0.37.0 ([a670a2b](https://github.com/d0ugal/slzb-exporter/commit/a670a2bf550617f345e4122feead8d74d7f490ee))


### Bug Fixes

* update gomod commitMessagePrefix from feat to fix ([cacb5fe](https://github.com/d0ugal/slzb-exporter/commit/cacb5feb7012b4e271df38489767e48027fc5874))

## [2.8.0](https://github.com/d0ugal/slzb-exporter/compare/v2.7.0...v2.8.0) (2025-10-07)


### Features

* update dependencies to v1.25.2 ([43aef2d](https://github.com/d0ugal/slzb-exporter/commit/43aef2d011a6084bd71425fd2c596e222c2108b5))

## [2.7.0](https://github.com/d0ugal/slzb-exporter/compare/v2.6.0...v2.7.0) (2025-10-07)


### Features

* **renovate:** use feat: commit messages for dependency updates ([3dec606](https://github.com/d0ugal/slzb-exporter/commit/3dec6067bcb816319eab18e4c6760e08bbd48e00))
* update dependencies to v0.67.1 ([0be4b53](https://github.com/d0ugal/slzb-exporter/commit/0be4b537c13b0c18dc9a51f127d0d64f4cd7ff42))

## [2.6.0](https://github.com/d0ugal/slzb-exporter/compare/v2.5.0...v2.6.0) (2025-10-03)


### Features

* pin versions in documentation and examples ([9931284](https://github.com/d0ugal/slzb-exporter/commit/9931284213bfc471a533ed69e3f5ad6e4624bb9d))
* **renovate:** add docs commit message format for documentation updates ([08f512c](https://github.com/d0ugal/slzb-exporter/commit/08f512cd5149676d70d40caf8ddac804d8c65adb))


### Reverts

* remove unnecessary renovate config changes ([019cdb1](https://github.com/d0ugal/slzb-exporter/commit/019cdb1f094719bce4e9c9e80c1346a2f86c59b3))

## [2.5.0](https://github.com/d0ugal/slzb-exporter/compare/v2.4.4...v2.5.0) (2025-10-02)


### Features

* **deps:** migrate to YAML v3 ([8dcd0f9](https://github.com/d0ugal/slzb-exporter/commit/8dcd0f9a245af6ed48c37488a376060b3257f52f))
* **renovate:** add gomodTidy post-update option for Go modules ([8cf2ecb](https://github.com/d0ugal/slzb-exporter/commit/8cf2ecb8d0c8a0e998176fcaaac37b9ac50d03f1))


### Reverts

* remove unnecessary renovate config changes ([d1fe352](https://github.com/d0ugal/slzb-exporter/commit/d1fe352be40c15384fd1a40bac98bd6b8bd41496))
* remove unnecessary renovate config changes ([d5d9d42](https://github.com/d0ugal/slzb-exporter/commit/d5d9d4283c385721a1326c3a7074d52e2df2096b))

## [2.4.4](https://github.com/d0ugal/slzb-exporter/compare/v2.4.3...v2.4.4) (2025-10-02)


### Reverts

* remove unnecessary renovate config changes ([b53b29b](https://github.com/d0ugal/slzb-exporter/commit/b53b29b23e045c5ee6459280c8058beac14d9aa2))

## [2.4.3](https://github.com/d0ugal/slzb-exporter/compare/v2.4.2...v2.4.3) (2025-10-02)


### Bug Fixes

* enable indirect dependency updates in renovate config ([8202695](https://github.com/d0ugal/slzb-exporter/commit/82026958fb326b61883a4ecdddf500392249b970))

## [2.4.2](https://github.com/d0ugal/slzb-exporter/compare/v2.4.1...v2.4.2) (2025-09-22)


### Bug Fixes

* resolve linting issues in config tests and server ([df57159](https://github.com/d0ugal/slzb-exporter/commit/df571597e2bff5dfbdd711b0a7a9ad61d32cb330))

## [2.4.1](https://github.com/d0ugal/slzb-exporter/compare/v2.4.0...v2.4.1) (2025-09-20)


### Bug Fixes

* **lint:** resolve gosec configuration contradiction ([fd06674](https://github.com/d0ugal/slzb-exporter/commit/fd06674ace4d89eb6287ce6faeb52310b649037d))

## [2.4.0](https://github.com/d0ugal/slzb-exporter/compare/v2.3.2...v2.4.0) (2025-09-12)


### Features

* replace latest docker tags with versioned variables for Renovate compatibility ([2dadba9](https://github.com/d0ugal/slzb-exporter/commit/2dadba9534273c4962294b32b8e2b34bdf14d050))

## [2.3.2](https://github.com/d0ugal/slzb-exporter/compare/v2.3.1...v2.3.2) (2025-09-05)


### Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([c02f4d3](https://github.com/d0ugal/slzb-exporter/commit/c02f4d32edec5a1d2fc4306ea5f70abe267cacf7))
* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([5c65bf2](https://github.com/d0ugal/slzb-exporter/commit/5c65bf2ff2c55fcb17d36145a6522a44291d45a7))

## [2.3.1](https://github.com/d0ugal/slzb-exporter/compare/v2.3.0...v2.3.1) (2025-09-04)


### Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([8d6d7d8](https://github.com/d0ugal/slzb-exporter/commit/8d6d7d8ceb448d88a206a10765df485e42918d1a))
* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([f66e810](https://github.com/d0ugal/slzb-exporter/commit/f66e810fa20dda0b6ebb2b71d8e7a10028aad574))

## [2.3.0](https://github.com/d0ugal/slzb-exporter/compare/v2.2.0...v2.3.0) (2025-09-04)


### Features

* update dev build versioning to use semver-compatible pre-release tags ([8aab899](https://github.com/d0ugal/slzb-exporter/commit/8aab8993a7f3c7741866c5658267b7c80bc96012))


### Bug Fixes

* **ci:** add v prefix to dev tags for consistent versioning ([0793607](https://github.com/d0ugal/slzb-exporter/commit/0793607d6e5cbd83f3b82b9d535824758b5ed8e9))
* use actual release version as base for dev tags instead of hardcoded 0.0.0 ([3924cf2](https://github.com/d0ugal/slzb-exporter/commit/3924cf27665a516ddfdbb533392de68b9870fe2e))
* use fetch-depth: 0 instead of fetch-tags for full git history ([bc9a10d](https://github.com/d0ugal/slzb-exporter/commit/bc9a10d74ebc3b24b10a02bce5ed1f455091ffc3))
* use fetch-tags instead of fetch-depth for GitHub Actions ([ae87d7e](https://github.com/d0ugal/slzb-exporter/commit/ae87d7e71282d5c9fdc97e748e1206124f7bcb6e))

## [2.2.0](https://github.com/d0ugal/slzb-exporter/compare/v2.1.1...v2.2.0) (2025-09-04)


### Features

* enable global automerge in Renovate config ([0cd920d](https://github.com/d0ugal/slzb-exporter/commit/0cd920de3beefdf3203e5b8c4a0f4d0f0071288d))

## [2.1.1](https://github.com/d0ugal/slzb-exporter/compare/v2.1.0...v2.1.1) (2025-09-03)


### Bug Fixes

* pin Alpine version to 3.22.1 for consistency ([2c65ff9](https://github.com/d0ugal/slzb-exporter/commit/2c65ff9f4c846648e8503efa5da29e1f51279e2f))

## [2.1.0](https://github.com/d0ugal/slzb-exporter/compare/v2.0.2...v2.1.0) (2025-08-26)


### Features

* **docker:** use an unprivileged user during runtime ([6606723](https://github.com/d0ugal/slzb-exporter/commit/6606723059aa972e2e2a8e217ce521acc859cbbb))

## [2.0.2](https://github.com/d0ugal/slzb-exporter/compare/v2.0.1...v2.0.2) (2025-08-20)


### Bug Fixes

* trigger release for server cleanup ([b09e8e8](https://github.com/d0ugal/slzb-exporter/commit/b09e8e81b96b95e6077ad95d36196d7f80cb1a7d))

## [2.0.1](https://github.com/d0ugal/slzb-exporter/compare/v2.0.0...v2.0.1) (2025-08-20)


### Bug Fixes

* remove redundant Service Information section from UI ([cac2dd3](https://github.com/d0ugal/slzb-exporter/commit/cac2dd3126e82a7550888c32788105a6dfa00334))

## [2.0.0](https://github.com/d0ugal/slzb-exporter/compare/v1.4.0...v2.0.0) (2025-08-20)


### ⚠ BREAKING CHANGES

* Metric names have been updated to comply with Prometheus naming conventions:
    - slzb_device_heap_free_kb → slzb_device_heap_free_bytes
    - slzb_device_heap_size_kb → slzb_device_heap_size_bytes
    - slzb_config_file_count → slzb_config_files

### Features

* optimize linting performance with caching and fix metric naming ([a0397b2](https://github.com/d0ugal/slzb-exporter/commit/a0397b2f9721889c50a637ef92d5a2797cbe19a7))


### Bug Fixes

* run Docker containers as current user to prevent permission issues ([59ce761](https://github.com/d0ugal/slzb-exporter/commit/59ce76138a3ff26e781a037f0b3e7e1722b6ebd6))
* temporarily disable gocyclo to allow commit to proceed ([34ae928](https://github.com/d0ugal/slzb-exporter/commit/34ae9282fccd359df879608f562eb33ae17e4493))

## [1.4.0](https://github.com/d0ugal/slzb-exporter/compare/v1.3.1...v1.4.0) (2025-08-20)


### Features

* implement template-based UI with centralized metric information ([6ad5c74](https://github.com/d0ugal/slzb-exporter/commit/6ad5c74773f71a2b4cbe5559be5108094d8ea0a2))

## [1.3.1](https://github.com/d0ugal/slzb-exporter/compare/v1.3.0...v1.3.1) (2025-08-20)


### Bug Fixes

* **tests:** remove handleMetricsInfo test and clean up unused imports ([e027723](https://github.com/d0ugal/slzb-exporter/commit/e027723c49145cbe4104d88fcf74469bf41ee1d0))

## [1.3.0](https://github.com/d0ugal/slzb-exporter/compare/v1.2.0...v1.3.0) (2025-08-19)


### Features

* **server:** add dynamic metrics information with examples ([5a48e1d](https://github.com/d0ugal/slzb-exporter/commit/5a48e1d3758e0f7affa459362e733caf40a94d96))
* **server:** make metrics list collapsible for better UX ([492d18d](https://github.com/d0ugal/slzb-exporter/commit/492d18d2172809712a0aabf192189e3e70995e9e))


### Bug Fixes

* **lint:** pre-allocate slices to resolve golangci-lint prealloc warnings ([de62003](https://github.com/d0ugal/slzb-exporter/commit/de62003aea7e2ee56731e360df8d926e2431fe76))

## [1.2.0](https://github.com/d0ugal/slzb-exporter/compare/v1.1.2...v1.2.0) (2025-08-18)


### Features

* add comprehensive monitoring features to SLZB exporter ([caa8f0b](https://github.com/d0ugal/slzb-exporter/commit/caa8f0b61c08869860b960c8e4cb829426ff1146))
* add comprehensive SLZB monitoring metrics ([da9ff58](https://github.com/d0ugal/slzb-exporter/commit/da9ff5840bdb9bb8866f7d2f625d45f515701c9d))


### Bug Fixes

* linting issue ([f5bd7dc](https://github.com/d0ugal/slzb-exporter/commit/f5bd7dcbe336c07cd6c368830a5f40be12175f9a))

## [1.1.2](https://github.com/d0ugal/slzb-exporter/compare/v1.1.1...v1.1.2) (2025-08-16)


### Bug Fixes

* add version injection to Dockerfile and CI workflow ([56fd68f](https://github.com/d0ugal/slzb-exporter/commit/56fd68f60e9f292c4cd7ffe39df7abf0470fba90))

## [1.1.1](https://github.com/d0ugal/slzb-exporter/compare/v1.1.0...v1.1.1) (2025-08-16)


### Bug Fixes

* apply golangci-lint formatting fixes for all issues ([de0954b](https://github.com/d0ugal/slzb-exporter/commit/de0954b3a31f4ea285df9e7627d91686bb2e16fa))
* correct golangci-lint config version to 2 (only valid version) ([3660290](https://github.com/d0ugal/slzb-exporter/commit/3660290c4f724ab29e1159be6b2a57cc7a3f0b16))
* update golangci-lint config to match working mqtt-exporter pattern ([a9a73df](https://github.com/d0ugal/slzb-exporter/commit/a9a73df018f4a4e4574bc0195f22befc10369c0c))

## [1.1.0](https://github.com/d0ugal/slzb-exporter/compare/v1.0.0...v1.1.0) (2025-08-16)


### Features

* add missing project files and documentation ([319e212](https://github.com/d0ugal/slzb-exporter/commit/319e2127aaf80106db72a55e37e364b03c6f0cb8))
* upgrade to Go 1.25 ([627b5d4](https://github.com/d0ugal/slzb-exporter/commit/627b5d4fd9ae091a0316f28e3d1081bb4e2ccf5d))


### Bug Fixes

* update golangci-lint config for Go 1.25 compatibility ([b14762a](https://github.com/d0ugal/slzb-exporter/commit/b14762ab761261c7e034f15131bd00bbf8b813a5))

## [Unreleased]

### Added
- Initial project setup and structure
- SLZB-06 device monitoring capabilities
- Prometheus metrics export
- Docker containerization
- Comprehensive test suite

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A

## [1.0.0] - 2025-08-16

### Added
- Initial release of SLZB-06 Prometheus Exporter
- Device status monitoring (connectivity, temperature, uptime)
- Network metrics (Ethernet and WiFi status, RSSI, connection speeds)
- System health tracking (memory usage, heap statistics)
- HTTP request monitoring (success/failure rates, response times)
- Configurable collection intervals
- Prometheus integration with metrics export on port 9110
- Health check endpoint at /health
- Graceful shutdown with proper signal handling
- Structured logging with JSON and text formats
- Environment-based configuration
- Docker containerization with multi-stage builds
- Comprehensive test suite with race detection and coverage
- golangci-lint integration for code quality
- GitHub Actions CI/CD workflow
- Semantic versioning with build-time version injection
- Detailed documentation and API reference

### Technical Details
- Go 1.22+ with modern Go modules
- Prometheus client_golang v1.22.0 for metrics handling
- Clean architecture with separate packages for collectors, config, logging, metrics, server, and version
- Alpine Linux container for minimal footprint
- Production-ready with proper error handling and resource management
