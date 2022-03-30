# k8sprom-patch-controller helm chart

Installs the [k8sprom-patch-controller](https://github.com/DoodleScheduling/k8sprom-patch-controller).

## Installing the Chart

To install the chart with the release name `k8sprom-patch-controller`:

```console
helm repo add k8sprom-patch-controller https://doodlescheduling.github.io/k8sprom-patch-controller/
helm upgrade --install k8sprom-patch-controller k8sprom-patch-controller/k8sprom-patch-controller
```

This command deploys k8sprom-patch-controller with the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Using the Chart

The chart comes with a PodMonitor for use with the [Prometheus Operator](https://github.com/helm/charts/tree/master/stable/prometheus-operator).
You need to enable it using `podMonitor.enabled: true`.
If you're not using the Prometheus Operator, you can use `podAnnotations` as below:

```yaml
podAnnotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "metrics"
  prometheus.io/path: "/metrics"
```

## Configuration

See Customizing the Chart Before Installing. To see all configurable options with detailed comments, visit the chart's values.yaml, or run the configuration command:

```sh
$ helm show values k8sprom-patch-controller/k8sprom-patch-controller
```
