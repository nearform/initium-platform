# Changelog

## [0.2.0](https://github.com/nearform/initium-platform/compare/v0.1.0...v0.2.0) (2023-12-04)


### Features

* add certificate for ngrok tcp endpoints ([#154](https://github.com/nearform/initium-platform/issues/154)) ([0daf1fe](https://github.com/nearform/initium-platform/commit/0daf1fe9df5095decbf76533d2382a1f2e616509))
* added kubernetes replicator addon and test ([#206](https://github.com/nearform/initium-platform/issues/206)) ([a885f94](https://github.com/nearform/initium-platform/commit/a885f94ab5abd3b7b52910291297bcfeb96162ce))
* argocd custom helm chart values exposed as a release artifact ([#208](https://github.com/nearform/initium-platform/issues/208)) ([1eeb8d6](https://github.com/nearform/initium-platform/commit/1eeb8d640967c391a39669edd5a97c905ceca312))
* expose grafana and generate admin password ([#237](https://github.com/nearform/initium-platform/issues/237)) ([d31d28b](https://github.com/nearform/initium-platform/commit/d31d28b34e33a81ffd7f7f9b304bd8bbb682f2c7))
* include initium asdf plugin ([#209](https://github.com/nearform/initium-platform/issues/209)) ([b31ad11](https://github.com/nearform/initium-platform/commit/b31ad11c430463f72e8fad44d0e8f54397920b70))
* Initium rebranding ([#157](https://github.com/nearform/initium-platform/issues/157)) ([e4cc5e5](https://github.com/nearform/initium-platform/commit/e4cc5e5ee274edecb2b1c784e86ec808373523a4))
* Test the quick-start-gke document ([#197](https://github.com/nearform/initium-platform/issues/197)) ([2b551c2](https://github.com/nearform/initium-platform/commit/2b551c204c3c501dbf19e4b281da52a960206c26))
* upgrade k8s to 1.27.3, as well as knative, dex and go tests ([#156](https://github.com/nearform/initium-platform/issues/156)) ([8626970](https://github.com/nearform/initium-platform/commit/86269704b114aa785f59fd44dffd8ae7b9a5a4f9))


### Bug Fixes

* bump up jq version to support arm64 architecture ([#242](https://github.com/nearform/initium-platform/issues/242)) ([bf3571f](https://github.com/nearform/initium-platform/commit/bf3571f85e8e795798c84a523a5a521fe7228f2c))
* proper kind node image for each arch ([#203](https://github.com/nearform/initium-platform/issues/203)) ([d4fed91](https://github.com/nearform/initium-platform/commit/d4fed914f35a6fd189d041f40a8345e231c08b5c))
* reverting kind version to 0.17.0 ([#181](https://github.com/nearform/initium-platform/issues/181)) ([e175546](https://github.com/nearform/initium-platform/commit/e175546cd551a4fc4d02ee164bf8c07b4ee675fc))

## [0.1.0](https://github.com/nearform/k8s-kurated-addons/compare/v0.0.1...v0.1.0) (2023-03-23)


### Features

* Add Dex addon ([#38](https://github.com/nearform/k8s-kurated-addons/issues/38)) ([62966d4](https://github.com/nearform/k8s-kurated-addons/commit/62966d410f7d119e24ef6e012ffd1f85f4974126))
* Add Grafana Loki ([#50](https://github.com/nearform/k8s-kurated-addons/issues/50)) ([087b319](https://github.com/nearform/k8s-kurated-addons/commit/087b3195d9d24f23844716f25ecb572b7b501d90))
* Add OpenTelemetry addon ([#33](https://github.com/nearform/k8s-kurated-addons/issues/33)) ([99b2dea](https://github.com/nearform/k8s-kurated-addons/commit/99b2dea139cf377855e32feb82d4defd809b7988))
* argocd improvements ([#85](https://github.com/nearform/k8s-kurated-addons/issues/85)) ([b1c8cca](https://github.com/nearform/k8s-kurated-addons/commit/b1c8ccaa43da5bbc73a1f05cbbebee940ce2b77f))
* Configurable sync policy for ArgoCD ([#56](https://github.com/nearform/k8s-kurated-addons/issues/56)) ([f2730e8](https://github.com/nearform/k8s-kurated-addons/commit/f2730e8d2c1387374146d2f695e18eba3efd2b8f))
* create github k8s service account ([#77](https://github.com/nearform/k8s-kurated-addons/issues/77)) ([74ab411](https://github.com/nearform/k8s-kurated-addons/commit/74ab411d5872df076a8ea33fad15e75e20b61180))
* Inject helmValues from bootstrap app to addons ([#53](https://github.com/nearform/k8s-kurated-addons/issues/53)) ([9818fc7](https://github.com/nearform/k8s-kurated-addons/commit/9818fc7eb760b5c7e4a8d4ddcd5aa1a4107fe2bb))
* knative improve deployment ([#107](https://github.com/nearform/k8s-kurated-addons/issues/107)) ([2120882](https://github.com/nearform/k8s-kurated-addons/commit/212088215d1e7651bd8cf12a98921f05ecd1b848))
* Move app-of-apps manifest to helm-chart ([#31](https://github.com/nearform/k8s-kurated-addons/issues/31)) ([5a6db64](https://github.com/nearform/k8s-kurated-addons/commit/5a6db64a5bbf6d9d88086d6d4c3eddb0df278445))
* remove envsubst requirement ([#87](https://github.com/nearform/k8s-kurated-addons/issues/87)) ([8984f9c](https://github.com/nearform/k8s-kurated-addons/commit/8984f9c4b8e2fb6cbe7ada8feb3624ece9141dd1))
* remove kubernetes 1.23 from github actions matrix ([#61](https://github.com/nearform/k8s-kurated-addons/issues/61)) ([b52c5cd](https://github.com/nearform/k8s-kurated-addons/commit/b52c5cde1dfce35c1948a80133474def6c12fb54))
* Sample Grafana dashboard ([#8](https://github.com/nearform/k8s-kurated-addons/issues/8)) ([1bb35e2](https://github.com/nearform/k8s-kurated-addons/commit/1bb35e20744da739dc36b383688f65acd8c9176f))


### Bug Fixes

* change renovate commit prefix ([#60](https://github.com/nearform/k8s-kurated-addons/issues/60)) ([672e4d0](https://github.com/nearform/k8s-kurated-addons/commit/672e4d0a22294e05d6f2ae258bd2fdfc3e03b9bb))
* comment back dex default values ([#63](https://github.com/nearform/k8s-kurated-addons/issues/63)) ([b89546d](https://github.com/nearform/k8s-kurated-addons/commit/b89546decb0d602395d1ecb54ba82483601d01a7))
* Downgrade ArgoCD traffic to HTTP ([#30](https://github.com/nearform/k8s-kurated-addons/issues/30)) ([1a5b1c9](https://github.com/nearform/k8s-kurated-addons/commit/1a5b1c9ba88788354afc1f692699afbc77afb92a))
* inject values to addons from bootstrap app + revert syncPolicy templating ([#99](https://github.com/nearform/k8s-kurated-addons/issues/99)) ([c478692](https://github.com/nearform/k8s-kurated-addons/commit/c478692788fde4b72d4e785c6b1d5bca1dcddd53))
* remove the update changelog entry ([#86](https://github.com/nearform/k8s-kurated-addons/issues/86)) ([57da682](https://github.com/nearform/k8s-kurated-addons/commit/57da6822ff0eae24a52455bb431ed2e478536f81))
* replace readarray ([#92](https://github.com/nearform/k8s-kurated-addons/issues/92)) ([683bc44](https://github.com/nearform/k8s-kurated-addons/commit/683bc44582f07d4c3fc8456e3189fc0a1f46d1c7))
* updating default dex config ([#101](https://github.com/nearform/k8s-kurated-addons/issues/101)) ([26ae141](https://github.com/nearform/k8s-kurated-addons/commit/26ae1412f93a860d36a262a5d6f94feeb18f1581))
* use POSIX compatible comparison operator ([#89](https://github.com/nearform/k8s-kurated-addons/issues/89)) ([231c3cc](https://github.com/nearform/k8s-kurated-addons/commit/231c3ccc687f244903a46f681f10ba078ca2ffc0))

## Changelog
