# Kubernetes resource patch controller using PromQL

[![release](https://img.shields.io/github/release/DoodleScheduling/prometheuspatch-controller/all.svg)](https://github.com/DoodleScheduling/prometheuspatch-controller/releases)
[![release](https://github.com/doodlescheduling/prometheuspatch-controller/actions/workflows/release.yaml/badge.svg)](https://github.com/doodlescheduling/prometheuspatch-controller/actions/workflows/release.yaml)
[![report](https://goreportcard.com/badge/github.com/DoodleScheduling/prometheuspatch-controller)](https://goreportcard.com/report/github.com/DoodleScheduling/prometheuspatch-controller)
[![Coverage Status](https://coveralls.io/repos/github/DoodleScheduling/prometheuspatch-controller/badge.svg?branch=master)](https://coveralls.io/github/DoodleScheduling/prometheuspatch-controller?branch=master)
[![license](https://img.shields.io/github/license/DoodleScheduling/prometheuspatch-controller.svg)](https://github.com/DoodleScheduling/prometheuspatch-controller/blob/master/LICENSE)

Apply patches to kubernetes resources based on prometheus queries.

## Example

```yaml
apiVersion: metrics.infra.doodle.com/v1beta1
kind: PrometheusPatchRule
metadata:
  name: annotate-namespace
spec:
  prometheus:
    address: http://prometheus-server.prometheus
  expr: |
    rate(nginx_ingress_controller_requests{exported_namespace="default"}[5m]) == 0
  for: 5m
  interval: 2m
  suspend: false
  json6902Patches:
  - target:
      version: v1
      resource: namespaces
      name: default
    patch:
    - op: add
      path: /metadata/annotations/has-ingress-traffic"
      value: "false"
```

## Details

### Prometheus expression
As soon as the given rule spec.expr evaluates to `true` the patches spec.patches get applied to the defined target `spec.patches[].target`.

### Pending state
You may define a window spec.for for which the rule will be in a pending condition similar to prometheus alerting rules.
As soon as the expression was `true` for the specified duration the patches get applied.

### Patches
Define a list of patches which needs a target selector as well as a list of JSON 6902 patch operations.
The target select requires at least the api version `version` as well as the resource group `resource` which is usually the kind in plural lowercase.

```yaml
json6902Patches:
- target:
    version: v1
    resource: namespaces
    name: default
  patch:
  - op: add
    path: /metadata/annotations/has-ingress-traffic"
    value: "false"
```
Instead selecting a single resource you may also select multiple ones by left out the name field.
You can filter multiple onse by specifying a comma separated label select: `labelSelector: label=value,label2=value`.

### Interval
Defines in what interval the rule is evaluated.

### Suspend
The PrometheusPatchRule may be suspended setting spec.suspend to `true`. A suspended rule does not get reconciled, meaning no patches will be applied as long as the rule is suspended.

### Remove patches
By design patches are **not** removed if the defined expression evaluates to `false` and if the patches have been added before.
The way to achieve this is to create another PrometheusPatchRule which expression does the opposite as well as reverse patches.

## Installation

### Requirements
You need a running prometheus [prometheus](https://prometheus.io/) (or any compatible prometheus service like [thanos](https://thanos.io/)).

### Permission
By default both the helm chart and the kustomize default base have a cluster rolebinding to cluster-admin.
Meaning the controller is granted full admin permission on the cluster.
This is needed as patch rules can target any kind of resources.
You may disable the binding and define fine grained cluster roles accordingly.

### Helm

Please see [chart/prometheuspatch-controller](https://github.com/DoodleScheduling/prometheuspatch-controller/tree/master/chart/prometheuspatch-controller) for the helm chart docs.

### Manifests/kustomize

Alternatively you may get the bundled manifests in each release to deploy it using kustomize or use them directly.

## Configure the controller

The controller can be configured using cmd args:
```
--concurrent int                            The number of concurrent Pod reconciles. (default 4)
--enable-leader-election                    Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.
--field-manager string                      The name of the field maanger used for server side apply https://kubernetes.io/docs/reference/using-api/server-side-apply/. (default "prometheuspatch-controller")
--graceful-shutdown-timeout duration        The duration given to the reconciler to finish before forcibly stopping. (default 10m0s)
--health-addr string                        The address the health endpoint binds to. (default ":9557")
--insecure-kubeconfig-exec                  Allow use of the user.exec section in kubeconfigs provided for remote apply.
--insecure-kubeconfig-tls                   Allow that kubeconfigs provided for remote apply can disable TLS verification.
--kube-api-burst int                        The maximum burst queries-per-second of requests sent to the Kubernetes API. (default 300)
--kube-api-qps float32                      The maximum queries-per-second of requests sent to the Kubernetes API. (default 50)
--leader-election-lease-duration duration   Interval at which non-leader candidates will wait to force acquire leadership (duration string). (default 35s)
--leader-election-release-on-cancel         Defines if the leader should step down voluntarily on controller manager shutdown. (default true)
--leader-election-renew-deadline duration   Duration that the leading controller manager will retry refreshing leadership before giving up (duration string). (default 30s)
--leader-election-retry-period duration     Duration the LeaderElector clients should wait between tries of actions (duration string). (default 5s)
--log-encoding string                       Log encoding format. Can be 'json' or 'console'. (default "json")
--log-level string                          Log verbosity level. Can be one of 'trace', 'debug', 'info', 'error'. (default "info")
--max-retry-delay duration                  The maximum amount of time for which an object being reconciled will have to wait before a retry. (default 15m0s)
--metrics-addr string                       The address the metric endpoint binds to. (default ":9556")
--min-retry-delay duration                  The minimum amount of time for which an object being reconciled will have to wait before a retry. (default 750ms)
--watch-all-namespaces                      Watch for resources in all namespaces, if set to false it will only watch the runtime namespace. (default true)
--watch-label-selector string               Watch for resources with matching labels e.g. 'sharding.fluxcd.io/shard=shard1'.

``